let app = new Vue({
  el: '#app',

  data: {
    input: '',
    encodeds: [],
    loading: false,
  },

  computed: {
    trimmedInput() {
      return this.input.trim()
    },
    noResult() {
      return !this.loading && this.trimmedInput !== '' && this.encodeds.length === 0
    },
    alphabet() {
      return this.encodeds.length ? this.trimmedInput : 'alphabet'
    },
    quran() {
      return this.encodeds.length ? this.encodeds[0].text : "Al-Qu'ran"
    },
  },

  watch: {
    input: function(newInput) {
      this.updateResult()
    },
  },

  methods: {
    updateResult: _.debounce(function() {
      this.loading = true
      axios.get('/api/encode/' + this.trimmedInput)
      .then((response) => this.encodeds = response.data.map((text) => ({text})))
      .catch(() => this.encodeds = [])
      .then(() => this.loading = false)
    }, 500),

    locate(encoded) {
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
  },
})