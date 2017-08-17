let app = new Vue({
  el: '#app',

  data: {
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
      this.$set(location, 'loadingShift', true)
      let dataPromise = location.ayaNumber+n === location.original.ayaNumber ?
        Promise.resolve(location.original) :
        axios.get(`/api/aya/${location.suraNumber}/${location.ayaNumber+n}`)
        .then((response)=> ({beforeHighlightedAya: response.data}))
      
      dataPromise.then((data) => {
        this.$set(location, 'ayaNumber', location.ayaNumber+n)
        this.$set(location, 'beforeHighlightedAya', data.beforeHighlightedAya)
        this.$set(location, 'highlightedAya', data.highlightedAya)
        this.$set(location, 'afterHighlightedAya', data.afterHighlightedAya)

        this.$set(location, 'translation', undefined)
        this.$set(location, 'tafsir', undefined)
        this.$set(location, 'showTranslation', false)
        this.$set(location, 'showTafsir', false)
      })
      .catch((error) => error.response.status !== 400 ? this.notify('connection problem') : undefined)
      .then(() => this.$set(location, 'loadingShift', false))
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
