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
      if (this.trimmedKeyword === '') {
        this.encodeds = []
        this.willRequest = false
      } else {
        this.logged = false
        ++this.loading
        let currentRequestTime = Date.now()
        axios.get('/api/encode/' + this.trimmedKeyword)
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
      }
    }, 500),

    locate(encoded) {
      this.log()
      this.$set(encoded, 'expanded', !encoded.expanded)
      if (encoded.locations) return
      this.$set(encoded, 'loading', true)
      axios.get('/api/locate/' + encoded.text)
      .then((response) => {
        let locations = response.data
        locations.forEach((loc) => {
          loc.AyaBeforeHL = loc.AyaText.substring(0, loc.Begin)
          loc.AyaHL = loc.AyaText.substring(loc.Begin, loc.End)
          loc.AyaAfterHL = loc.AyaText.substring(loc.End)
        })
        this.$set(encoded, 'locations', locations)
      })
      .then(() => componentHandler.upgradeElements(this.$refs[encoded.text]))
      .catch(() => {this.$set(encoded, 'expanded', false); this.notify('connection problem')})
      .then(() => this.$set(encoded, 'loading', false))
    },

    translate(location) {
      this.$set(location, 'showTranslation', !location.showTranslation)
      if (location.Translation) return
      this.$set(location, 'loadingTranslation', true)
      axios.get(`/api/translate/${location.Sura+1}-${location.Aya+1}`)
      .then((response) => this.$set(location, 'Translation', response.data))
      .catch(() => {this.$set(location, 'showTranslation', false); this.notify('connection problem')})
      .then(() => this.$set(location, 'loadingTranslation', false))
    },

    tafsir(location) {
      this.$set(location, 'showTafsir', !location.showTafsir)
      if (location.Tafsir) return
      this.$set(location, 'loadingTafsir', true)
      axios.get(`/api/tafsir/${location.Sura+1}-${location.Aya+1}`)
      .then((response) => this.$set(location, 'Tafsir', response.data))
      .catch(() => {this.$set(location, 'showTafsir', false); this.notify('connection problem')})
      .then(() => this.$set(location, 'loadingTafsir', false))
    },

    log() {
      if (!this.logged) {
        this.logged = true
        axios.post('/api/log/' + this.trimmedKeyword)
        .catch(() => this.logged = false)
      }
    },

    notify(message) {
      this.$refs['notification'].MaterialSnackbar.showSnackbar({message})
    },
  },
})

axios.get('/api/trending-keywords')
.then((response) => app.trendingKeywords = response.data)
.catch(() => {})

let clipboard = new Clipboard('#share-link')
clipboard.on('success', () => app.notify('share link copied to clipboard'))
