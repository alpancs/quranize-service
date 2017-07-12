let willRequest = false
let logged = false
let lastRequestTime = 0

let app = new Vue({
  el: '#app',

  data: {
    keyword: '',
    encodeds: [],
    loading: false,
    trendingKeywords: [],
  },

  computed: {
    trimmedKeyword() {
      return this.keyword.trim()
    },
    noResults() {
      return !willRequest && this.trimmedKeyword !== '' && this.encodeds.length === 0
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
      willRequest = true
      this.updateResult()
    },
  },

  methods: {
    updateResult: _.debounce(function() {
      logged = false
      this.loading = true
      let currentRequestTime = Date.now()
      axios.get('/api/encode/' + this.trimmedKeyword)
      .then((response) => {
        if (lastRequestTime < currentRequestTime) {
          lastRequestTime = currentRequestTime
          this.encodeds = response.data.map((text) => ({text}))
        }
      })
      .then(() => {
        if (this.$refs.encodeds)
          componentHandler.upgradeElements(this.$refs.encodeds)
      })
      .catch(() => {this.encodeds = []; this.showError()})
      .then(() => {this.loading = false; willRequest = false})
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
      .catch(() => {this.$set(encoded, 'expanded', false); this.showError()})
      .then(() => this.$set(encoded, 'loading', false))
    },

    translate(location) {
      this.$set(location, 'loading', true)
      axios.get(`/api/translate/${location.Sura+1}-${location.Aya+1}`)
      .then((response) => this.$set(location, 'Translation', response.data))
      .catch(() => {this.showError()})
      .then(() => this.$set(location, 'loading', false))
    },

    log() {
      if (!logged) {
        logged = true
        axios.post('/api/log/' + this.trimmedKeyword)
        .catch(() => logged = false)
      }
    },

    showError() {
      this.$refs['network-error'].MaterialSnackbar.showSnackbar({message: 'network error'})
    },
  },
})

axios.get('/api/trending-keywords')
.then((response) => app.trendingKeywords = response.data)
.catch(() => {})
