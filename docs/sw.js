// Define Service Worker version and App Shell to cache
const NAME = 'imgedit-wasm-';
const VERSION = '001';
const CACHE_NAME = NAME + VERSION;
const urlsToCache = [
    './index.html',
    './assets/js/main.js',
    './assets/js/drop-zone.js',
    './assets/js/wasm.js',
    './assets/js/wasm_exec.js',
    './assets/css/wasm.css',
    './assets/image/img_bg.gif',
    './assets/image/img_download.png',
    './assets/image/img_logo.png',
    './assets/image/img_reset.png',
];

// Install files to Service Worker
self.addEventListener('install', function (event) {
    event.waitUntil(
        caches.open(CACHE_NAME)
            .then(function (cache) {
                console.log('Opened cache');
                return cache.addAll(urlsToCache);
            })
    );
});

// If the requested file is cached in the Service Worker
// Return response from cache
self.addEventListener('fetch', function (event) {
    if (event.request.cache === 'only-if-cached' && event.request.mode !== 'same-origin')
        return;
    event.respondWith(
        caches.match(event.request)
            .then(function (response) {
                if (response) {
                    return response;
                }
                return fetch(event.request);
            })
    );
});

// If there is a change in the service worker keys cached in Cache Storage,
// delete the cache of the old version after installing the new version.
// This file considers CACHE_NAME as the value of key and detects changes.
self.addEventListener('activate', event => {
    event.waitUntil(
        caches.keys().then(keys => Promise.all(
            keys.map(key => {
                if (!CACHE_NAME.includes(key)) {
                    return caches.delete(key);
                }
            })
        )).then(() => {
            console.log(CACHE_NAME + "activated");
        })
    );
});