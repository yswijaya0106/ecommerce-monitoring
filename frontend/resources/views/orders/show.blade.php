@extends('layouts.app')

@section('content')
    <h1 class="mb-4">Pesanan #{{ $order['id'] }}</h1>

    <p>
        Tanggal: {{ $order['order_date'] ? \Illuminate\Support\Carbon::parse($order['order_date'])->format('d M Y H:i') : '-' }}<br>
        Customer ID: {{ $order['customer_id'] ?? '-' }}<br>
        Catatan: {{ $order['notes'] ?? '-' }}
    </p>

    <table class="table bg-white">
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
                    <td>{{ $item['product_name'] ?? 'Produk #' . $item['product_id'] }}</td>
                    <td>{{ $item['quantity'] }}</td>
                    <td>Rp {{ number_format($item['unit_price'], 0, ',', '.') }}</td>
                    <td>Rp {{ number_format($item['quantity'] * $item['unit_price'] * (1 - $item['discount']), 0, ',', '.') }}</td>
                </tr>
            @endforeach
        </tbody>
        <tfoot>
            <tr>
                <th colspan="3" class="text-end">Total</th>
                <th>Rp {{ number_format($order['total'], 0, ',', '.') }}</th>
            </tr>
        </tfoot>
    </table>

    <a href="{{ route('orders.index') }}" class="btn btn-outline-secondary">Kembali</a>
@endsection
