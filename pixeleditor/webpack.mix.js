const mix = require('laravel-mix');

mix.ts('./src/app.ts', './public')
mix.sass('./src/styles.scss', './public');
