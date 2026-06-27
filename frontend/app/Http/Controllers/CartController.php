<?php

namespace App\Http\Controllers;

use App\Services\BackendApiClient;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Session;

class CartController extends Controller
{
    public function __construct(protected BackendApiClient $backend)
    {
    }

    public function index()
    {
        $cart = Session::get('cart', []);
        $total = array_sum(array_map(
            fn ($item) => $item['list_price'] * $item['quantity'],
            $cart
        ));

        return view('cart.index', [
            'cart' => $cart,
            'total' => $total,
            'customers' => $this->backend->customers() ?? [],
        ]);
    }

    public function add(Request $request)
    {
        $data = $request->validate([
            'product_id' => ['required', 'integer'],
            'quantity' => ['required', 'numeric', 'min:1'],
        ]);

        $product = $this->backend->product($data['product_id']);

        $cart = Session::get('cart', []);
        $productId = $data['product_id'];
        $quantity = (float) $data['quantity'];

        if (isset($cart[$productId])) {
            $cart[$productId]['quantity'] += $quantity;
        } else {
            $cart[$productId] = [
                'product_id' => $productId,
                'product_name' => $product['product_name'] ?? "Produk #{$productId}",
                'list_price' => (float) ($product['list_price'] ?? 0),
                'quantity' => $quantity,
            ];
        }

        Session::put('cart', $cart);

        return redirect()->route('catalog.index')->with('status', 'Produk ditambahkan ke keranjang.');
    }

    public function remove(Request $request)
    {
        $data = $request->validate(['product_id' => ['required', 'integer']]);

        $cart = Session::get('cart', []);
        unset($cart[$data['product_id']]);
        Session::put('cart', $cart);

        return redirect()->route('cart.index');
    }
}
