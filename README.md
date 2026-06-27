# Ecommerce Monitoring — Pemesanan Barang Northwind

Aplikasi demo pemesanan barang yang membaca data dari database **Northwind**, terdiri dari:

- **Backend API** (Go) — REST API untuk produk, customer, dan pemesanan (orders), terinstrumentasi penuh dengan OpenTelemetry.
- **Frontend** (Laravel) — katalog produk, keranjang belanja, checkout, dan riwayat pesanan. Tanpa database sendiri — semua data diambil lewat REST API ke backend.
- **Monitoring stack** ringan — OpenTelemetry Collector + Jaeger (traces), Prometheus + Grafana (metrics), Loki + Promtail (logs), dan Portainer (manajemen container).

Semua komponen berjalan di Docker Compose, dirancang agar tetap ringan untuk dijalankan di laptop.

## Arsitektur

```
                     ┌─────────────┐
   Browser  ───────► │  Laravel    │  (frontend, :8000)
                     │  frontend   │
                     └──────┬──────┘
                            │ REST (JSON)
                            ▼
                     ┌─────────────┐        OTLP gRPC       ┌──────────────────┐
                     │  Go backend │ ─────────────────────► │  OTel Collector  │
                     │  orders-api │                        │     (:4327/4328) │
                     │   (:8081)   │                        └─────────┬────────┘
                     └──────┬──────┘                       traces │  metrics
                            │ SQL                                 ▼        ▼
                            ▼                              ┌─────────┐ ┌────────────┐
                  ┌──────────────────┐                     │ Jaeger  │ │ Prometheus │
                  │  MariaDB/MySQL   │                     │ (:16687)│ │  (:9090)   │
                  │  Northwind (DB   │                     └─────────┘ └─────┬──────┘
                  │  sudah berjalan, │                                       │
                  │  bukan di compose)│         ┌─────────┐                  │
                  └──────────────────┘          │  Loki   │◄── Promtail      │
                                                 │ (:3100) │   (docker logs)  │
                                                 └────┬────┘                  │
                                                      │                       │
                                                      ▼                       ▼
                                                 ┌──────────────────────────────┐
                                                 │           Grafana            │
                                                 │            (:3000)           │
                                                 └──────────────────────────────┘

                                                 ┌──────────────────────────────┐
                                                 │   Portainer (:9000)          │
                                                 │   Kelola semua container     │
                                                 └──────────────────────────────┘
```

Database Northwind **tidak** dijalankan oleh docker-compose ini — backend connect langsung ke server MariaDB/MySQL Northwind yang sudah berjalan di jaringan Anda (lihat [Konfigurasi](#konfigurasi)).

## Struktur Proyek

```
.
├── docker-compose.yml          # Seluruh stack: backend, frontend, monitoring
├── .env.example                # Template koneksi database
├── database/init/01-northwind.sql   # Referensi skema Northwind (tidak dipakai compose, untuk dokumentasi)
├── backend/                     # Go REST API
│   ├── cmd/api/main.go
│   └── internal/
│       ├── config/              # Load env vars
│       ├── telemetry/           # Setup OpenTelemetry SDK (trace + metric)
│       ├── db/                  # Koneksi DB ter-instrumentasi (otelsql)
│       ├── httpserver/          # Router (otelhttp)
│       ├── product/             # GET /products, /products/{id}
│       ├── customer/            # GET /customers, /customers/{id}
│       └── order/                # GET/POST /orders — checkout transactional
├── frontend/                    # Laravel app (tanpa DB sendiri)
│   ├── app/Services/BackendApiClient.php   # Wrapper HTTP client ke backend
│   ├── app/Http/Controllers/    # CatalogController, CartController, OrderController
│   └── resources/views/         # Blade views (Bootstrap via CDN, tanpa build step)
└── monitoring/
    ├── otel-collector/config.yaml
    ├── prometheus/prometheus.yml
    ├── loki/loki-config.yml
    ├── promtail/promtail-config.yml
    └── grafana/provisioning/datasources/   # Auto-provision Prometheus + Loki + Jaeger
```

## Prasyarat

- Docker & Docker Compose v2 (`docker compose version`)
- Server MariaDB/MySQL dengan database **Northwind** yang sudah berjalan dan bisa diakses dari mesin Anda (lihat skema referensi di `database/init/01-northwind.sql` — tabel `customers`, `products`, `orders`, `order_details`, dll).

## Konfigurasi

1. Copy `.env.example` menjadi `.env` di root proyek, lalu isi sesuai server database Anda:

   ```bash
   cp .env.example .env
   ```

   ```env
   DB_HOST=192.168.1.100   # ganti dengan IP server MariaDB/MySQL Anda
   DB_PORT=3306
   DB_USER=root
   DB_PASS=
   DB_NAME=northwind
   ```

2. (Opsional) Kalau pakai MCP server (misal untuk eksplorasi database lewat AI assistant), copy juga:

   ```bash
   cp .vscode/mcp.json.example .vscode/mcp.json
   ```

   lalu isi `MYSQL_HOST`/`MYSQL_USER`/`MYSQL_PASS` sesuai server Anda. File ini **tidak** dibutuhkan untuk menjalankan aplikasi — hanya alat bantu development.

## Menjalankan Aplikasi

```bash
docker compose up -d --build
```

Tunggu semua container `Up`, lalu cek statusnya:

```bash
docker compose ps
```

### Daftar Service & URL

| Service        | URL                            | Keterangan                              |
|----------------|---------------------------------|------------------------------------------|
| Frontend       | http://localhost:8000           | Aplikasi pemesanan barang (Laravel)       |
| Backend API    | http://localhost:8081           | REST API (`/products`, `/customers`, `/orders`) |
| Jaeger UI      | http://localhost:16687          | Distributed tracing                       |
| Prometheus     | http://localhost:9090           | Metrics & query                           |
| Grafana        | http://localhost:3000           | Dashboard (anonymous login, role Admin)   |
| Loki           | http://localhost:3100           | Log storage (diakses lewat Grafana)       |
| Portainer      | http://localhost:9000           | Manajemen container Docker                |

> Catatan port: backend/Jaeger/OTel Collector di-mapping ke port non-default (8081, 16687, 4327/4328) untuk menghindari konflik kalau ada stack Docker lain yang sudah memakai port standarnya (8080, 16686, 4317/4318) di mesin yang sama. Sesuaikan di `docker-compose.yml` kalau tidak ada konflik di mesin Anda.

## Cara Pakai Aplikasi

1. Buka http://localhost:8000 — otomatis redirect ke **Katalog**.
2. Pilih produk, atur jumlah, klik **Tambah** untuk masuk ke keranjang.
3. Buka **Keranjang**, pilih **Pemesan** (customer) dari dropdown, isi catatan (opsional), klik **Buat Pesanan**.
4. Pesanan langsung tersimpan ke database Northwind via backend, dan Anda diarahkan ke halaman detail pesanan.
5. Riwayat semua pesanan bisa dilihat di menu **Riwayat Pesanan**.

## Observability

### Traces (Jaeger)

1. Buka http://localhost:16687
2. Pilih service `orders-api` di dropdown **Service**, klik **Find Traces**.
3. Setiap request ke backend (termasuk saat checkout) muncul sebagai trace dengan span HTTP + span query SQL.

### Metrics (Prometheus / Grafana)

1. Buka http://localhost:3000 (login otomatis sebagai Admin, anonymous access aktif).
2. Datasource **Prometheus**, **Loki**, dan **Jaeger** sudah otomatis ter-konfigurasi (lihat menu Connections → Data sources).
3. Contoh query metric yang tersedia (bisa langsung dicoba di Explore):
   - `http_server_request_duration_seconds_count` — jumlah request HTTP per route
   - `db_sql_connection_open` — koneksi database aktif
   - `db_sql_latency_milliseconds_sum` — total waktu eksekusi query

### Logs (Loki)

Di Grafana → Explore → pilih datasource **Loki**, lalu query misalnya:

```
{service="backend"}
{service="frontend"}
```

Promtail otomatis menarik log dari semua container Docker di host (lewat `docker_sd_configs`), jadi log container baru juga ikut tertarik tanpa konfigurasi tambahan.

## Portainer — Setup Admin & Setup Token

Portainer CE mewajibkan pembuatan akun admin pertama dalam waktu terbatas (±5 menit) setelah container start, sebagai pengaman supaya tidak ada pihak lain yang bisa mengklaim akun admin lebih dulu. Kalau dibiarkan tanpa setup melewati batas waktu itu, instance akan **timeout dan terkunci** — muncul pesan *"Your Portainer instance timed out for security purposes"*.

### Setup Pertama Kali

1. Buka http://localhost:9000 segera setelah `docker compose up`.
2. Portainer akan menampilkan form **Create the initial administrator user**, dengan field:
   - **Username** — bebas, misalnya `admin`.
   - **Password** — minimal 12 karakter.
   - **Setup token** — token yang di-generate otomatis oleh Portainer setiap kali container start, dicetak ke log container.
3. Ambil setup token dari log:

   ```bash
   docker compose logs portainer | grep setup_token
   ```

   Contoh output:

   ```
   ... setup_token=9227502b2b7152a3e5f9565131711b78168dbd6a983f244c9d9c509c4cd4e7f6
   ```

   Copy string setelah `setup_token=` (tanpa spasi) ke field **Setup token** di form.
4. Isi username + password, submit. Setup selesai — login berikutnya hanya butuh username/password, tidak perlu token lagi.

### Kalau Sudah Timeout / Terkunci

Restart container Portainer untuk membuka kembali wizard setup-nya, lalu ambil token yang baru (token berubah setiap restart):

```bash
docker compose restart portainer
docker compose logs portainer | grep setup_token | tail -1
```

Segera buka http://localhost:9000 dan selesaikan setup admin sebelum timeout lagi.

## Menghentikan / Membersihkan

```bash
# Hentikan semua container, volume (data Grafana/Prometheus/Loki/Portainer) tetap ada
docker compose down

# Hentikan + hapus seluruh data volume (reset total)
docker compose down -v
```

## Troubleshooting

- **Backend gagal connect DB** — pastikan `.env` punya `DB_HOST` yang benar dan server database bisa dijangkau dari host Docker (`mysql -h <DB_HOST> -u <DB_USER> northwind`).
- **Port sudah terpakai saat `docker compose up`** — cek `docker ps` untuk container lain yang memakai port yang sama, lalu ubah mapping port di `docker-compose.yml`.
- **Frontend error 500 setelah ubah env di `docker-compose.yml`** — frontend menjalankan `php artisan config:cache` saat start supaya environment variable dari Docker terbaca dengan benar (`php artisan serve` secara default hanya forward sebagian env var ke proses child-nya). Kalau habis ubah `environment:` di compose, jalankan ulang `docker compose up -d --build frontend`.
