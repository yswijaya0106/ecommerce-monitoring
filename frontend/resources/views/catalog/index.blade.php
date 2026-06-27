@extends('layouts.app')

@php
    $categoryIcons = [
        'Beverages' => 'bi-cup-straw',
        'Condiments' => 'bi-droplet',
        'Confections' => 'bi-egg-fried',
        'Dairy products' => 'bi-egg',
        'Grains/Cereals' => 'bi-basket',
        'Meat/Poultry' => 'bi-egg-fried',
        'Produce' => 'bi-flower1',
        'Seafood' => 'bi-water',
    ];
    $iconFor = fn ($category) => $categoryIcons[$category] ?? 'bi-box-seam';
@endphp

@section('content')
    <div class="page-hero">
        <h1 class="mb-2"><i class="bi bi-shop me-2"></i>Katalog Produk</h1>
        <p>Pilih produk Northwind Traders favorit Anda dan tambahkan ke keranjang.</p>
    </div>

    <ul class="nav category-pills gap-2 mb-4 flex-wrap">
        <li class="nav-item">
            <a href="{{ route('catalog.index') }}" class="nav-link {{ $selectedCategory ? '' : 'active' }}">
                Semua
            </a>
        </li>
        @foreach ($categories as $category)
            <li class="nav-item">
                <a href="{{ route('catalog.index', ['category' => $category]) }}"
                   class="nav-link {{ $selectedCategory === $category ? 'active' : '' }}">
                    <i class="bi {{ $iconFor($category) }} me-1"></i>{{ $category }}
                </a>
            </li>
        @endforeach
    </ul>

    <div class="row">
        @forelse ($products as $product)
            <div class="col-md-4 mb-4">
                <div class="card product-card h-100">
                    <div class="card-body d-flex flex-column">
                        <div class="d-flex align-items-start justify-content-between mb-3">
                            <div class="product-icon">
                                <i class="bi {{ $iconFor($product['category'] ?? null) }}"></i>
                            </div>
                            @if ($product['discontinued'])
                                <span class="badge bg-secondary">Tidak tersedia</span>
                            @endif
                        </div>
                        <h5 class="card-title mb-1">{{ $product['product_name'] ?? 'Produk #' . $product['id'] }}</h5>
                        <p class="card-text text-muted small mb-2">{{ $product['category'] ?? '-' }}</p>
                        <p class="card-text text-muted small flex-grow-1">{{ $product['description'] ?? '' }}</p>
                        <p class="product-price mb-3">Rp {{ number_format($product['list_price'], 0, ',', '.') }}</p>
                        @unless ($product['discontinued'])
                            <form method="POST" action="{{ route('cart.add') }}" class="d-flex gap-2">
                                @csrf
                                <input type="hidden" name="product_id" value="{{ $product['id'] }}">
                                <input type="number" name="quantity" value="1" min="1" class="form-control form-control-sm" style="width: 5rem">
                                <button type="submit" class="btn btn-sm btn-brand flex-grow-1">
                                    <i class="bi bi-cart-plus me-1"></i>Tambah
                                </button>
                            </form>
                        @endunless
                    </div>
                </div>
            </div>
        @empty
            <div class="col-12">
                <div class="empty-state surface-card">
                    <i class="bi bi-inboxes"></i>
                    Tidak ada produk pada kategori ini.
                </div>
            </div>
        @endforelse
    </div>
@endsection
