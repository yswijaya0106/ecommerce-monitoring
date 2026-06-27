@extends('layouts.app')

@section('content')
    <h1 class="mb-4">Keranjang</h1>

    @if (empty($cart))
        <p class="text-muted">Keranjang masih kosong. <a href="{{ route('catalog.index') }}">Lihat katalog</a>.</p>
    @else
        <table class="table bg-white">
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
                        <td>{{ $item['product_name'] }}</td>
                        <td>Rp {{ number_format($item['list_price'], 0, ',', '.') }}</td>
                        <td>{{ $item['quantity'] }}</td>
                        <td>Rp {{ number_format($item['list_price'] * $item['quantity'], 0, ',', '.') }}</td>
                        <td>
                            <form method="POST" action="{{ route('cart.remove') }}">
                                @csrf
                                <input type="hidden" name="product_id" value="{{ $item['product_id'] }}">
                                <button type="submit" class="btn btn-sm btn-outline-danger">Hapus</button>
                            </form>
                        </td>
                    </tr>
                @endforeach
            </tbody>
            <tfoot>
                <tr>
                    <th colspan="3" class="text-end">Total</th>
                    <th>Rp {{ number_format($total, 0, ',', '.') }}</th>
                    <th></th>
                </tr>
            </tfoot>
        </table>

        <div class="card mt-4" style="max-width: 30rem">
            <div class="card-body">
                <h5 class="card-title">Checkout</h5>
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
                    <button type="submit" class="btn btn-primary">Buat Pesanan</button>
                </form>
            </div>
        </div>
    @endif
@endsection
