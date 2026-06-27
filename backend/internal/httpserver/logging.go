package httpserver

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
)

// maxLoggedBodyBytes caps how much of a request/response body is captured
// for the OpenTelemetry log record, so large payloads don't blow up exports.
const maxLoggedBodyBytes = 4096

var requestLogger = global.Logger("ecommerce-backend/httpserver")

// requestResponseLogging captures the request and response bodies and emits
// them as an OpenTelemetry log record. Emitting with r.Context() lets the SDK
// attach the trace/span IDs of the active otelhttp span, so the log line is
// correlated with the corresponding trace in Jaeger.
func requestResponseLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		reqBody, _ := io.ReadAll(io.LimitReader(r.Body, maxLoggedBodyBytes+1))
		r.Body.Close()
		r.Body = io.NopCloser(io.MultiReader(bytes.NewReader(reqBody), r.Body))

		rec := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		duration := time.Since(start)

		var record otellog.Record
		record.SetTimestamp(time.Now())
		record.SetSeverity(severityFor(rec.status))
		record.SetBody(otellog.StringValue(fmt.Sprintf("%s %s -> %d (%s)", r.Method, r.URL.Path, rec.status, duration)))
		record.AddAttributes(
			otellog.String("http.method", r.Method),
			otellog.String("http.target", r.URL.Path),
			otellog.Int("http.status_code", rec.status),
			otellog.Float64("duration_ms", float64(duration.Microseconds())/1000),
			otellog.String("http.request.body", truncatedBody(reqBody)),
			otellog.String("http.response.body", truncatedBody(rec.body.Bytes())),
		)
		requestLogger.Emit(r.Context(), record)
	})
}

func severityFor(status int) otellog.Severity {
	switch {
	case status >= http.StatusInternalServerError:
		return otellog.SeverityError
	case status >= http.StatusBadRequest:
		return otellog.SeverityWarn
	default:
		return otellog.SeverityInfo
	}
}

func truncatedBody(b []byte) string {
	if len(b) > maxLoggedBodyBytes {
		return string(b[:maxLoggedBodyBytes]) + "...(truncated)"
	}
	return string(b)
}

// responseRecorder wraps http.ResponseWriter to capture the status code and
// up to maxLoggedBodyBytes of the response body for logging.
type responseRecorder struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (rr *responseRecorder) WriteHeader(status int) {
	rr.status = status
	rr.ResponseWriter.WriteHeader(status)
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	if remaining := maxLoggedBodyBytes - rr.body.Len(); remaining > 0 {
		if remaining > len(b) {
			remaining = len(b)
		}
		rr.body.Write(b[:remaining])
	}
	return rr.ResponseWriter.Write(b)
}
