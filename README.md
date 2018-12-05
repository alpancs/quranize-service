# Quranize

Quranize transforms Alquran's transliteration into its arabic form, like "bismillah" into "بسم الله".
https://quranize.herokuapp.com

Telegram channel https://t.me/quranize

## API

### GET `/api/encode`

- Query
  - keyword: string

- Example
  - Request
    ```sh
    $ curl 'https://quranize.herokuapp.com/api/encode?keyword=bismillah'
    ```
  - Response
    ```json
    [
      "بسم الله",
      "بشماله"
    ]
    ```

### GET `/api/locate`

- Query
  - keyword: string

- Example
  - Request
    ```sh
    $ curl 'https://quranize.herokuapp.com/api/locate?keyword=بسم+الله'
    ```
  - Response
    ```json
    [
      {
        "suraNumber": 1,
        "suraName": "الفاتحة",
        "ayaNumber": 1,
        "ayaText": "بِسْمِ اللَّهِ الرَّحْمَنِ الرَّحِيمِ",
        "beginHighlight": 0,
        "endHighlight": 14
      },
      {
        "suraNumber": 11,
        "suraName": "هود",
        "ayaNumber": 41,
        "ayaText": "وَقَالَ ارْكَبُوا فِيهَا بِسْمِ اللَّهِ مَجْرَاهَا وَمُرْسَاهَا إِنَّ رَبِّي لَغَفُورٌ رَّحِيمٌ",
        "beginHighlight": 25,
        "endHighlight": 39
      },
      {
        "suraNumber": 27,
        "suraName": "النمل",
        "ayaNumber": 30,
        "ayaText": "إِنَّهُ مِن سُلَيْمَانَ وَإِنَّهُ بِسْمِ اللَّهِ الرَّحْمَنِ الرَّحِيمِ",
        "beginHighlight": 34,
        "endHighlight": 48
      }
    ]
    ```

### GET `/api/aya/{sura: int}/{aya: int}`

  - Example
    - Request
      ```sh
      $ curl 'https://quranize.herokuapp.com/api/aya/1/2'
      ```
    - Response
      ```json
      "الْحَمْدُ لِلَّهِ رَبِّ الْعَالَمِينَ"
      ```

### GET `/api/translation/{sura: int}/{aya: int}`

  - Example
    - Request
    ```sh
    $ curl 'https://quranize.herokuapp.com/api/translation/1/2'
    ```
    - Response
    ```json
    "Segala puji bagi Allah, Tuhan semesta alam."
    ```

### GET `/api/tafsir/{sura: int}/{aya: int}`

  - Example
    - Request
      ```sh
      $ curl 'https://quranize.herokuapp.com/api/tafsir/1/2'
      ```
    - Response
      ```json
      "Segala puja dan puji kita persembahkan kepada Allah semata, karena Dialah Yang menciptakan dan memelihara seluruh makhluk."
      ```

### GET `/api/trending_keywords`

  - Query:
    - limit: int (optional, default: 6)

  - Example
    - Request
    ```sh
    $ curl 'https://quranize.herokuapp.com/api/trending_keywords'
    ```
    - Response
    ```json
    [
      "bismillah",
      "akho hum syu'aiba",
      "alhamdu",
      "ardi",
      "qurban",
      "rubama"
    ]
    ```

## Related Project

https://github.com/alpancs/quranize
