const cacheName = 'v1::static';

self.addEventListener('install', e => {
  e.waitUntil(
    caches.open(cacheName).then(async cache => {
      return cache.addAll([
        '/',
        '/index.html',
        '/assests/js/scripts.js',
        '/assests/js/alpine.min.js',
        '/assests/css/style.css',
        '/assests/css/tailwind.min.css',
        'assets/img/favicon.ico',
      ]).then(() => self.skipWaiting());
    })
  );
});

self.addEventListener('fetch', event => {
  event.respondWith(
    caches.open(cacheName).then(async cache => {
      return cache.match(event.request).then(res => {
        return res || fetch(event.request)
      });
    })
  );
});