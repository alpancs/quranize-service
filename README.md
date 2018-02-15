# Quranize

Quranize transforms alquran into القرآن.

https://quranize.herokuapp.com


# API

- `GET /api/encode`
  - Query
    - `keyword: string`
  - Example
    - Request: `GET /api/encode?keyword=bismillah`
    - Response: `["بسم الله","بشماله"]`

- `GET /api/locate`
  - Query
    - `keyword: string`
  - Example
    - Request: `GET /api/locate?keyword=بسم الله`
    - Response: `[{"suraNumber":1,"suraName":"الفاتحة","ayaNumber":1,"ayaText":"بِسْمِ اللَّهِ الرَّحْمَنِ الرَّحِيمِ","beginHighlight":0,"endHighlight":14},{"suraNumber":11,"suraName":"هود","ayaNumber":41,"ayaText":"وَقَالَ ارْكَبُوا فِيهَا بِسْمِ اللَّهِ مَجْرَاهَا وَمُرْسَاهَا إِنَّ رَبِّي لَغَفُورٌ رَّحِيمٌ","beginHighlight":25,"endHighlight":39},{"suraNumber":27,"suraName":"النمل","ayaNumber":30,"ayaText":"إِنَّهُ مِن سُلَيْمَانَ وَإِنَّهُ بِسْمِ اللَّهِ الرَّحْمَنِ الرَّحِيمِ","beginHighlight":34,"endHighlight":48}]`

- `GET /api/aya/{sura: int}/{aya: int}`
  - Example
    - Request: `GET /api/aya/1/2`
    - Response: `"الْحَمْدُ لِلَّهِ رَبِّ الْعَالَمِينَ"`

- `GET /api/translation/{sura: int}/{aya: int}`
  - Example
    - Request: `GET /api/translation/1/2`
    - Response: `Segala puji bagi Allah, Tuhan semesta alam.`

- `GET /api/tafsir/{sura: int}/{aya: int}`
  - Example
    - Request: `GET /api/tafsir/1/2`
    - Response: `Segala puja dan puji kita persembahkan kepada Allah semata, karena Dialah Yang menciptakan dan memelihara seluruh makhluk.`

- `GET /api/trending_keywords`
  - Query:
    - `limit: int` (optional, default: `6`)
  - Example
    - Request: `GET /api/trending_keywords`
    - Response: `["bismillah","akho hum syu'aiba","alhamdu","ardi","qurban","rubama"]`

- `POST /api/log`
  - Example
    - Request: `POST /api/log` with post data `bismillah`
    - Response: empty response with status code `200`
