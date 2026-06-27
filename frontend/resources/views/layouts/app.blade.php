<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{ config('app.name') }}</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&display=swap" rel="stylesheet">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.css" rel="stylesheet">
    <link href="{{ asset('css/site.css') }}" rel="stylesheet">
</head>
<body>
    <nav class="navbar navbar-expand navbar-dark bg-brand shadow-sm sticky-top">
        <div class="container">
            <a class="navbar-brand d-flex align-items-center gap-2" href="{{ route('catalog.index') }}">
                <i class="bi bi-shop"></i> {{ config('app.name') }}
            </a>
            <div class="navbar-nav">
                <a class="nav-link {{ request()->routeIs('catalog.*') ? 'active' : '' }}" href="{{ route('catalog.index') }}">
                    <i class="bi bi-grid me-1"></i>Katalog
                </a>
                <a class="nav-link {{ request()->routeIs('cart.*') ? 'active' : '' }}" href="{{ route('cart.index') }}">
                    <i class="bi bi-cart3 me-1"></i>Keranjang
                    @if (($itemCount = collect(session('cart', []))->sum('quantity')) > 0)
                        <span class="badge bg-white text-dark rounded-pill cart-badge">{{ $itemCount }}</span>
                    @endif
                </a>
                <a class="nav-link {{ request()->routeIs('orders.*') ? 'active' : '' }}" href="{{ route('orders.index') }}">
                    <i class="bi bi-receipt me-1"></i>Riwayat Pesanan
                </a>
            </div>
        </div>
    </nav>

    <main class="container pb-5 pt-4">
        @if (session('status'))
            <div class="alert alert-success d-flex align-items-center gap-2">
                <i class="bi bi-check-circle-fill"></i> {{ session('status') }}
            </div>
        @endif
        @if ($errors->any())
            <div class="alert alert-danger d-flex align-items-start gap-2">
                <i class="bi bi-exclamation-triangle-fill mt-1"></i>
                <ul class="mb-0">
                    @foreach ($errors->all() as $error)
                        <li>{{ $error }}</li>
                    @endforeach
                </ul>
            </div>
        @endif

        @yield('content')
    </main>

    <footer class="site-footer text-center border-top">
        <div class="container">
            &copy; {{ now()->year }} {{ config('app.name') }} &middot; Northwind Traders
        </div>
    </footer>
</body>
</html>
