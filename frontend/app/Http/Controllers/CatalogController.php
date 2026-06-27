<?php

namespace App\Http\Controllers;

use App\Services\BackendApiClient;
use Illuminate\Http\Request;

class CatalogController extends Controller
{
    public function __construct(protected BackendApiClient $backend)
    {
    }

    public function index(Request $request)
    {
        $category = $request->query('category');
        $products = $this->backend->products($category) ?? [];
        $categories = collect($products)->pluck('category')->filter()->unique()->sort()->values();

        return view('catalog.index', [
            'products' => $products,
            'categories' => $categories,
            'selectedCategory' => $category,
        ]);
    }
}
