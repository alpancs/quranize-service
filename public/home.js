new Vue({
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
      return !this.loading && this.trimmedInput && !this.encodeds.length
    },
    alphabet() {
      return this.encodeds.length ? this.trimmedInput : ''
    },
    quran() {
      return this.encodeds.length ? this.encodeds[0].text : ''
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
      axios.get('/api/encode/' + this.input.trim())
      .then((response) => this.encodeds = response.data.map((e) => ({text: e})))
      .catch(() => this.encodeds = [])
      .then(() => this.loading = false)
    }, 500),

    locate(encoded) {
      this.$set(encoded, 'loading', true)
      axios.get('/api/locate/' + encoded.text)
      .then((response) => {
        let data = response.data
        data.forEach((location) => {
          let aya = location.Aya
          let begin = location.Location.Index
          let end = location.Location.Index + encoded.text.length
          location.AyaBeforeHL = aya.substring(0, begin)
          location.AyaHL = aya.substring(begin, end)
          location.AyaAfterHL = aya.substring(end)
        })
        this.$set(encoded, 'locations', data)
      })
      .catch(() => this.$set(encoded, 'locations', undefined))
      .then(() => this.$set(encoded, 'loading', false))
    },
  },
})