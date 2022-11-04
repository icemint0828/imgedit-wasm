"use strict";

registerSW();

function registerSW() {

    // Service Worker 対応ブラウザの場合、スコープに基づいてService Worker を登録する

    if ('serviceWorker' in navigator) {
        window.addEventListener('load', function () {
            navigator.serviceWorker.register('./sw.js', { scope: './' }).then(function (registration) {
                console.log('ServiceWorker registration successful with scope: ', registration.scope);
            }, function (err) {
                console.log('ServiceWorker registration failed: ', err);
            });
        });
    }
}