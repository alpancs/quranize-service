let app = new Vue({
  el: '#app',

  data: {
    ayaCount: [7,286,200,176,120,165,206,75,129,109,123,111,43,52,99,128,111,110,98,135,112,78,118,64,77,227,93,88,69,60,34,30,73,54,45,83,182,88,75,85,54,53,89,59,37,35,38,29,18,45,60,49,62,55,78,96,29,22,24,13,14,11,11,18,12,12,30,52,52,44,28,28,20,56,40,31,50,40,46,42,29,19,36,25,22,17,19,26,30,20,15,21,11,8,8,19,5,8,8,11,11,8,3,9,5,4,7,3,6,3,5,4,5,6],
    keyword: '',
    encodeds: [],
    loading: 0,
    trendingKeywords: [],
    shareLink: '',

    logged: false,
    lastRequestTime: 0,
    willRequest: false,
  },

  computed: {
    trimmedKeyword() {
      return this.keyword.trim()
    },
    noResults() {
      return !this.willRequest && this.trimmedKeyword !== '' && this.encodeds.length === 0
    },
    transliteration() {
      return this.encodeds.length ? this.trimmedKeyword : 'transliteration'
    },
    quran() {
      return this.encodeds.length ? this.encodeds[0].text : "Alquran"
    },
  },

  watch: {
    keyword() {
      this.willRequest = true
      this.updateResult()
    },
  },

  methods: {
    updateResult: _.debounce(function() {
      sessionStorage.setItem('keyword', this.trimmedKeyword)
      this.logged = false
      ++this.loading
      let currentRequestTime = Date.now()
      axios.get('/api/encode', {params: {keyword: this.trimmedKeyword}})
        .then((response) => {
          if (this.lastRequestTime < currentRequestTime) {
            this.lastRequestTime = currentRequestTime
            this.encodeds = response.data.map((text) => ({text}))
            this.shareLink = location.origin + '/' + this.trimmedKeyword.replace(/ /g,'').toLowerCase()
          }
        })
        .then(() => {
          if (this.$refs.encodeds)
            componentHandler.upgradeElements(this.$refs.encodeds)
        })
        .catch(() => {this.encodeds = []; this.notify('connection problem')})
        .then(() => {--this.loading; this.willRequest = this.loading > 0})
    }, 500),

    locate(encoded) {
      this.log()
      this.$set(encoded, 'expanded', !encoded.expanded)
      if (encoded.locations) return
      this.$set(encoded, 'loading', true)
      axios.get('/api/locate', {params: {keyword: encoded.text}})
      .then((response) => {
        let locations = response.data
        locations.forEach((location) => {
          location.beforeHighlightedAya = location.ayaText.substring(0, location.beginHighlight)
          location.highlightedAya = location.ayaText.substring(location.beginHighlight, location.endHighlight)
          location.afterHighlightedAya = location.ayaText.substring(location.endHighlight)
          location.original = {
            ayaNumber: location.ayaNumber,
            beforeHighlightedAya: location.beforeHighlightedAya,
            highlightedAya: location.highlightedAya,
            afterHighlightedAya: location.afterHighlightedAya,
          }
        })
        this.$set(encoded, 'locations', locations)
      })
      .then(() => componentHandler.upgradeElements(this.$refs[encoded.text]))
      .catch(() => {this.$set(encoded, 'expanded', false); this.notify('connection problem')})
      .then(() => this.$set(encoded, 'loading', false))
    },

    translate(location) {
      this.$set(location, 'showTranslation', !location.showTranslation)
      if (location.translation) return
      this.$set(location, 'loadingTranslation', true)
      axios.get(`/api/translate/${location.suraNumber}/${location.ayaNumber}`)
      .then((response) => this.$set(location, 'translation', response.data))
      .catch(() => {this.$set(location, 'showTranslation', false); this.notify('connection problem')})
      .then(() => this.$set(location, 'loadingTranslation', false))
    },

    tafsir(location) {
      this.$set(location, 'showTafsir', !location.showTafsir)
      if (location.tafsir) return
      this.$set(location, 'loadingTafsir', true)
      axios.get(`/api/tafsir/${location.suraNumber}/${location.ayaNumber}`)
      .then((response) => this.$set(location, 'tafsir', response.data))
      .catch(() => {this.$set(location, 'showTafsir', false); this.notify('connection problem')})
      .then(() => this.$set(location, 'loadingTafsir', false))
    },

    shift(location, n) {
      this.$set(location, 'shiftButtonDisabled', true)
      this.$set(location, 'loadingAya', true)
      let ayaPromise = location.ayaNumber+n === location.original.ayaNumber ?
        Promise.resolve(location.original) :
        axios.get(`/api/aya/${location.suraNumber}/${location.ayaNumber+n}`)
        .then((response)=> ({beforeHighlightedAya: response.data}))
      ayaPromise = ayaPromise.then((aya) => {
        this.$set(location, 'beforeHighlightedAya', aya.beforeHighlightedAya)
        this.$set(location, 'highlightedAya', aya.highlightedAya)
        this.$set(location, 'afterHighlightedAya', aya.afterHighlightedAya)
        this.$set(location, 'ayaNumber', location.ayaNumber+n)
      })
      .then(() => this.$set(location, 'loadingAya', false))

      this.$set(location, 'loadingTranslation', true)
      let translationPromise = location.showTranslation ?
        axios.get(`/api/translate/${location.suraNumber}/${location.ayaNumber+n}`)
        .then((response)=> response.data) :
        Promise.resolve()
      translationPromise = translationPromise.then((translation) => this.$set(location, 'translation', translation))
      .then(() => this.$set(location, 'loadingTranslation', false))

      this.$set(location, 'loadingTafsir', true)
      let tafsirPromise = location.showTafsir ?
        axios.get(`/api/tafsir/${location.suraNumber}/${location.ayaNumber+n}`)
        .then((response)=> response.data) :
        Promise.resolve()
      tafsirPromise = tafsirPromise.then((tafsir) => this.$set(location, 'tafsir', tafsir))
      .then(() => this.$set(location, 'loadingTafsir', false))

      Promise.all([ayaPromise, translationPromise, tafsirPromise])
      .catch(() => this.notify('connection problem'))
      .then(() => {
        this.$set(location, 'loadingAya', false)
        this.$set(location, 'loadingTranslation', false)
        this.$set(location, 'loadingTafsir', false)
        this.$set(location, 'shiftButtonDisabled', false)
      })
    },

    log() {
      if (!this.logged) {
        this.logged = true
        axios.post('/api/log', this.trimmedKeyword)
        .catch(() => this.logged = false)
      }
    },

    notify(message) {
      this.$refs['notification'].MaterialSnackbar.showSnackbar({message})
    },
  },
})

axios.get('/api/trending_keywords')
.then((response) => app.trendingKeywords = response.data)
.catch(() => {})

let clipboard = new Clipboard('#share-link')
clipboard.on('success', () => app.notify('share link copied to clipboard'))
