<?php

namespace App\Services;

use Illuminate\Http\Client\PendingRequest;
use Illuminate\Support\Facades\Http;

/**
 * Thin wrapper around the Go orders-api. The Laravel app has no database of
 * its own — every read/write of products, customers and orders goes through
 * this client.
 */
class BackendApiClient
{
    protected function client(): PendingRequest
    {
        return Http::baseUrl(config('backend.url'))->acceptJson();
    }

    public function products(?string $category = null): array
    {
        return $this->client()
            ->get('/products', $category ? ['category' => $category] : [])
            ->throw()
            ->json();
    }

    public function product(int $id): array
    {
        return $this->client()->get("/products/{$id}")->throw()->json();
    }

    public function customers(): array
    {
        return $this->client()->get('/customers')->throw()->json();
    }

    public function orders(): array
    {
        return $this->client()->get('/orders')->throw()->json();
    }

    public function order(int $id): array
    {
        return $this->client()->get("/orders/{$id}")->throw()->json();
    }

    public function createOrder(array $payload): array
    {
        return $this->client()->post('/orders', $payload)->throw()->json();
    }
}
