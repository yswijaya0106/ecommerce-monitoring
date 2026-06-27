@extends('layouts.app')

@section('content')
    <h1 class="mb-4"><i class="bi bi-receipt me-2"></i>Riwayat Pesanan</h1>

    @if (empty($orders))
        <div class="empty-state surface-card bg-white">
            <i class="bi bi-receipt-cutoff"></i>
            Belum ada pesanan.
            <div class="mt-3">
                <a href="{{ route('catalog.index') }}" class="btn btn-brand">
                    <i class="bi bi-grid me-1"></i>Lihat Katalog
                </a>
            </div>
        </div>
    @else
        <div class="surface-card bg-white p-0">
            <table class="table table-clean align-middle mb-0">
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
                            <td class="fw-semibold">#{{ $order['id'] }}</td>
                            <td>{{ $order['order_date'] ? \Illuminate\Support\Carbon::parse($order['order_date'])->format('d M Y H:i') : '-' }}</td>
                            <td>{{ $order['customer_id'] ?? '-' }}</td>
                            <td class="fw-semibold">Rp {{ number_format($order['total'], 0, ',', '.') }}</td>
                            <td class="text-end">
                                <a href="{{ route('orders.show', $order['id']) }}" class="btn btn-sm btn-outline-secondary">
                                    Detail <i class="bi bi-arrow-right ms-1"></i>
                                </a>
                            </td>
                        </tr>
                    @endforeach
                </tbody>
            </table>
        </div>
    @endif
@endsection
