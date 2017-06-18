let app = new Vue({
  el: '#app',

  data: {
    input: '',
    results: [],
    loading: false,
  },

  computed: {
    trimmedInput() {
      return this.input.trim()
    },
    noResult() {
      return this.trimmedInput && !this.results.length
    },
    alphabet() {
      return this.noResult ? '' : this.trimmedInput
    },
    quran() {
      return this.noResult ? '' : this.results[0]
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
      .then((response) => this.results = response.data)
      .catch(() => this.results = [])
      .then(() => this.loading = false)
    }, 500)
  },
})