<?php

namespace App\Http\Controllers;

use App\Services\BackendApiClient;
use Illuminate\Http\Client\RequestException;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Session;

class OrderController extends Controller
{
    public function __construct(protected BackendApiClient $backend)
    {
    }

    public function index()
    {
        return view('orders.index', ['orders' => $this->backend->orders() ?? []]);
    }

    public function show(int $id)
    {
        return view('orders.show', ['order' => $this->backend->order($id)]);
    }

    public function store(Request $request)
    {
        $data = $request->validate([
            'customer_id' => ['required', 'integer'],
            'notes' => ['nullable', 'string', 'max:1000'],
        ]);

        $cart = Session::get('cart', []);
        if (empty($cart)) {
            return redirect()->route('cart.index')->withErrors(['cart' => 'Keranjang masih kosong.']);
        }

        try {
            $order = $this->backend->createOrder([
                'customer_id' => (int) $data['customer_id'],
                'notes' => $data['notes'] ?? null,
                'items' => array_map(fn ($item) => [
                    'product_id' => (int) $item['product_id'],
                    'quantity' => (float) $item['quantity'],
                ], array_values($cart)),
            ]);
        } catch (RequestException $e) {
            $message = $e->response->json('error', 'Gagal membuat pesanan.');

            return redirect()->route('cart.index')->withErrors(['cart' => $message]);
        }

        Session::forget('cart');

        return redirect()->route('orders.show', $order['id'])->with('status', 'Pesanan berhasil dibuat.');
    }
}
