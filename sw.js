var dataCacheName = 'Gonewz_data';
var cacheName = 'Gonewz_cache';
var filesToCache = [
  '/',
  '/index.html',
  '/assests/js/scripts.js',
  '/assests/js/alpine.min.js',
  '/assests/css/style.css',
  '/assests/css/tailwind.min.css',
  'assets/img/favicon.ico',
];

self.addEventListener('install', function(e) {
  console.log('[ServiceWorker] Install');
  e.waitUntil(
    caches.open(cacheName).then(function(cache) {
      console.log('[ServiceWorker] Caching app shell');
      return cache.addAll(filesToCache);
    })
  );
});

self.addEventListener('activate', function(e) {
  console.log('[ServiceWorker] Activate');
  e.waitUntil(
    caches.keys().then(function(keyList) {
      return Promise.all(keyList.map(function(key) {
        if (key !== cacheName && key !== dataCacheName) {
          console.log('[ServiceWorker] Removing old cache', key);
          return caches.delete(key);
        }
      }));
    })
  );
  return self.clients.claim();
});

self.addEventListener('fetch', function(e) {
  console.log('[Service Worker] Fetch', e.request.url);
  var dataUrl = 'https://saurav.tech/NewsAPI/top-headlines/category/general/in.json';
  if (e.request.url.indexOf(dataUrl) > -1) {
    e.respondWith(
      caches.open(dataCacheName).then(async function(cache) {
        const response = await fetch(e.request);
        cache.put(e.request.url, response.clone());
        return response;
      })
    );
  } else {
    e.respondWith(
      caches.match(e.request).then(function(response) {
        return response || fetch(e.request);
      })
    );
  }
});
