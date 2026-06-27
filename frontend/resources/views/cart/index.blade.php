@extends('layouts.app')

@section('content')
    <h1 class="mb-4"><i class="bi bi-cart3 me-2"></i>Keranjang</h1>

    @if (empty($cart))
        <div class="empty-state surface-card bg-white">
            <i class="bi bi-cart-x"></i>
            Keranjang masih kosong.
            <div class="mt-3">
                <a href="{{ route('catalog.index') }}" class="btn btn-brand">
                    <i class="bi bi-grid me-1"></i>Lihat Katalog
                </a>
            </div>
        </div>
    @else
        <div class="row g-4">
            <div class="col-lg-8">
                <div class="surface-card bg-white p-0">
                    <table class="table table-clean align-middle mb-0">
                        <thead>
                            <tr>
                                <th>Produk</th>
                                <th>Harga</th>
                                <th>Qty</th>
                                <th>Subtotal</th>
                                <th></th>
                            </tr>
                        </thead>
                        <tbody>
                            @foreach ($cart as $item)
                                <tr>
                                    <td class="fw-semibold">{{ $item['product_name'] }}</td>
                                    <td>Rp {{ number_format($item['list_price'], 0, ',', '.') }}</td>
                                    <td>{{ $item['quantity'] }}</td>
                                    <td class="fw-semibold">Rp {{ number_format($item['list_price'] * $item['quantity'], 0, ',', '.') }}</td>
                                    <td class="text-end">
                                        <form method="POST" action="{{ route('cart.remove') }}">
                                            @csrf
                                            <input type="hidden" name="product_id" value="{{ $item['product_id'] }}">
                                            <button type="submit" class="btn btn-sm btn-outline-danger">
                                                <i class="bi bi-trash3"></i>
                                            </button>
                                        </form>
                                    </td>
                                </tr>
                            @endforeach
                        </tbody>
                        <tfoot>
                            <tr>
                                <th colspan="3" class="text-end">Total</th>
                                <th class="product-price">Rp {{ number_format($total, 0, ',', '.') }}</th>
                                <th></th>
                            </tr>
                        </tfoot>
                    </table>
                </div>
            </div>

            <div class="col-lg-4">
                <div class="card surface-card">
                    <div class="card-body">
                        <h5 class="card-title mb-3"><i class="bi bi-bag-check me-2"></i>Checkout</h5>
                        <form method="POST" action="{{ route('orders.store') }}">
                            @csrf
                            <div class="mb-3">
                                <label class="form-label">Pemesan</label>
                                <select name="customer_id" class="form-select" required>
                                    <option value="">-- pilih customer --</option>
                                    @foreach ($customers as $customer)
                                        <option value="{{ $customer['id'] }}">
                                            {{ $customer['company'] ?? ($customer['first_name'] . ' ' . $customer['last_name']) }}
                                        </option>
                                    @endforeach
                                </select>
                            </div>
                            <div class="mb-3">
                                <label class="form-label">Catatan (opsional)</label>
                                <textarea name="notes" class="form-control" rows="2"></textarea>
                            </div>
                            <button type="submit" class="btn btn-brand w-100">
                                <i class="bi bi-check-circle me-1"></i>Buat Pesanan
                            </button>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    @endif
@endsection
