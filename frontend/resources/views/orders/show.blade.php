@extends('layouts.app')

@section('content')
    <div class="d-flex align-items-center justify-content-between mb-4">
        <h1 class="mb-0"><i class="bi bi-receipt me-2"></i>Pesanan #{{ $order['id'] }}</h1>
        <a href="{{ route('orders.index') }}" class="btn btn-outline-secondary">
            <i class="bi bi-arrow-left me-1"></i>Kembali
        </a>
    </div>

    <div class="surface-card bg-white p-4 mb-4">
        <div class="row">
            <div class="col-md-4">
                <div class="text-muted small">Tanggal</div>
                <div class="fw-semibold">{{ $order['order_date'] ? \Illuminate\Support\Carbon::parse($order['order_date'])->format('d M Y H:i') : '-' }}</div>
            </div>
            <div class="col-md-4">
                <div class="text-muted small">Customer ID</div>
                <div class="fw-semibold">{{ $order['customer_id'] ?? '-' }}</div>
            </div>
            <div class="col-md-4">
                <div class="text-muted small">Catatan</div>
                <div class="fw-semibold">{{ $order['notes'] ?? '-' }}</div>
            </div>
        </div>
    </div>

    <div class="surface-card bg-white p-0">
        <table class="table table-clean align-middle mb-0">
            <thead>
                <tr>
                    <th>Produk</th>
                    <th>Qty</th>
                    <th>Harga</th>
                    <th>Subtotal</th>
                </tr>
            </thead>
            <tbody>
                @foreach ($order['items'] ?? [] as $item)
                    <tr>
                        <td class="fw-semibold">{{ $item['product_name'] ?? 'Produk #' . $item['product_id'] }}</td>
                        <td>{{ $item['quantity'] }}</td>
                        <td>Rp {{ number_format($item['unit_price'], 0, ',', '.') }}</td>
                        <td class="fw-semibold">Rp {{ number_format($item['quantity'] * $item['unit_price'] * (1 - $item['discount']), 0, ',', '.') }}</td>
                    </tr>
                @endforeach
            </tbody>
            <tfoot>
                <tr>
                    <th colspan="3" class="text-end">Total</th>
                    <th class="product-price">Rp {{ number_format($order['total'], 0, ',', '.') }}</th>
                </tr>
            </tfoot>
        </table>
    </div>
@endsection
