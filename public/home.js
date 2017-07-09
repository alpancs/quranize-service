let app = new Vue({
  el: '#app',

  data: {
    keyword: '',
    encodeds: [],
    loading: false,
    logged: false,
    trendingKeywords: [],
  },

  computed: {
    trimmedKeyword() {
      return this.keyword.trim()
    },
    noResults() {
      return !this.loading && this.trimmedKeyword !== '' && this.encodeds.length === 0
    },
    alphabet() {
      return this.encodeds.length ? this.trimmedKeyword : 'alphabet'
    },
    quran() {
      return this.encodeds.length ? this.encodeds[0].text : "Alquran"
    },
  },

  watch: {
    keyword() {
      this.updateResult()
    },
  },

  methods: {
    updateResult: _.debounce(function() {
      this.logged = false
      this.loading = true
      axios.get('/api/encode/' + this.trimmedKeyword)
      .then((response) => this.encodeds = response.data.map((text) => ({text})))
      .then(() => componentHandler.upgradeElements(this.$refs.encodeds))
      .catch(() => this.encodeds = [])
      .then(() => this.loading = false)
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
      .catch(() => this.$set(encoded, 'locations', undefined))
      .then(() => this.$set(encoded, 'loading', false))
    },

    log() {
      if (!this.logged) {
        this.logged = true
        axios.post('/api/log/' + this.trimmedKeyword)
        .catch(() => this.logged = false)
      }
    },
  },
})

axios.get('/api/trending-keywords')
.then((response) => app.trendingKeywords = response.data)
.catch(() => {})
