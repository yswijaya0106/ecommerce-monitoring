@extends('layouts.app')

@section('content')
    <h1 class="mb-4">Katalog Produk</h1>

    <div class="btn-group mb-4 flex-wrap">
        <a href="{{ route('catalog.index') }}" class="btn btn-sm {{ $selectedCategory ? 'btn-outline-secondary' : 'btn-secondary' }}">Semua</a>
        @foreach ($categories as $category)
            <a href="{{ route('catalog.index', ['category' => $category]) }}"
               class="btn btn-sm {{ $selectedCategory === $category ? 'btn-secondary' : 'btn-outline-secondary' }}">
                {{ $category }}
            </a>
        @endforeach
    </div>

    <div class="row">
        @forelse ($products as $product)
            <div class="col-md-4 mb-4">
                <div class="card h-100">
                    <div class="card-body d-flex flex-column">
                        <h5 class="card-title">{{ $product['product_name'] ?? 'Produk #' . $product['id'] }}</h5>
                        <p class="card-text text-muted small">{{ $product['category'] ?? '-' }}</p>
                        <p class="card-text">{{ $product['description'] ?? '' }}</p>
                        <p class="fw-bold mt-auto">Rp {{ number_format($product['list_price'], 0, ',', '.') }}</p>
                        @if ($product['discontinued'])
                            <span class="badge bg-secondary">Tidak tersedia</span>
                        @else
                            <form method="POST" action="{{ route('cart.add') }}" class="d-flex gap-2">
                                @csrf
                                <input type="hidden" name="product_id" value="{{ $product['id'] }}">
                                <input type="number" name="quantity" value="1" min="1" class="form-control form-control-sm" style="width: 5rem">
                                <button type="submit" class="btn btn-sm btn-primary">Tambah</button>
                            </form>
                        @endif
                    </div>
                </div>
            </div>
        @empty
            <p class="text-muted">Tidak ada produk.</p>
        @endforelse
    </div>
@endsection
