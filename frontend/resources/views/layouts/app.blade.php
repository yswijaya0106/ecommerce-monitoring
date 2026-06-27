<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{ config('app.name') }}</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body class="bg-light">
    <nav class="navbar navbar-expand navbar-dark bg-dark mb-4">
        <div class="container">
            <a class="navbar-brand" href="{{ route('catalog.index') }}">{{ config('app.name') }}</a>
            <div class="navbar-nav">
                <a class="nav-link" href="{{ route('catalog.index') }}">Katalog</a>
                <a class="nav-link" href="{{ route('cart.index') }}">
                    Keranjang
                    @if (($itemCount = collect(session('cart', []))->sum('quantity')) > 0)
                        <span class="badge bg-primary">{{ $itemCount }}</span>
                    @endif
                </a>
                <a class="nav-link" href="{{ route('orders.index') }}">Riwayat Pesanan</a>
            </div>
        </div>
    </nav>

    <main class="container pb-5">
        @if (session('status'))
            <div class="alert alert-success">{{ session('status') }}</div>
        @endif
        @if ($errors->any())
            <div class="alert alert-danger">
                <ul class="mb-0">
                    @foreach ($errors->all() as $error)
                        <li>{{ $error }}</li>
                    @endforeach
                </ul>
            </div>
        @endif

        @yield('content')
    </main>
</body>
</html>
