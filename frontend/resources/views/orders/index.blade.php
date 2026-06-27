@extends('layouts.app')

@section('content')
    <h1 class="mb-4">Riwayat Pesanan</h1>

    @if (empty($orders))
        <p class="text-muted">Belum ada pesanan.</p>
    @else
        <table class="table bg-white">
            <thead>
                <tr>
                    <th>#</th>
                    <th>Tanggal</th>
                    <th>Customer</th>
                    <th>Total</th>
                    <th></th>
                </tr>
            </thead>
            <tbody>
                @foreach ($orders as $order)
                    <tr>
                        <td>{{ $order['id'] }}</td>
                        <td>{{ $order['order_date'] ? \Illuminate\Support\Carbon::parse($order['order_date'])->format('d M Y H:i') : '-' }}</td>
                        <td>{{ $order['customer_id'] ?? '-' }}</td>
                        <td>Rp {{ number_format($order['total'], 0, ',', '.') }}</td>
                        <td><a href="{{ route('orders.show', $order['id']) }}" class="btn btn-sm btn-outline-secondary">Detail</a></td>
                    </tr>
                @endforeach
            </tbody>
        </table>
    @endif
@endsection
