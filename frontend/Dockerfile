FROM php:8.2-cli-alpine

RUN apk add --no-cache git unzip libzip-dev libxml2-dev oniguruma-dev sqlite-dev \
    && docker-php-ext-install pdo pdo_sqlite mbstring xml dom

COPY --from=composer:2 /usr/bin/composer /usr/bin/composer

WORKDIR /app
COPY composer.json composer.lock ./
RUN composer install --no-dev --no-scripts --no-interaction --optimize-autoloader

COPY . .
RUN composer dump-autoload --optimize \
    && php -r "file_exists('.env') || copy('.env.example', '.env');" \
    && php artisan key:generate --force \
    && php -r "file_exists('database/database.sqlite') || touch('database/database.sqlite');"

EXPOSE 8000

# `php artisan serve` only forwards a small whitelist of env vars to the PHP
# subprocess it spawns (see Illuminate\Foundation\Console\ServeCommand), so
# docker-compose `environment:` overrides (BACKEND_API_URL, SESSION_DRIVER,
# etc.) would otherwise be silently ignored. Caching config first bakes the
# real runtime env into bootstrap/cache/config.php before that subprocess
# starts, so it never needs to re-read env() itself.
CMD ["sh", "-c", "php artisan config:cache && php artisan serve --host=0.0.0.0 --port=8000"]
