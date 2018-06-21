let app = new Vue({
  el: '#app',

  delimiters: ['${', '}'],

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
      return this.encodeds.length ? this.trimmedKeyword : 'alquran'
    },
    quran() {
      return this.encodeds.length ? this.encodeds[0].text : 'القرآن'
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
      if (this.trimmedKeyword != history.state)
        history.pushState(this.trimmedKeyword, "Quranize", "/"+this.trimmedKeyword)

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
        .then(() => this.$refs.encodeds ? componentHandler.upgradeElements(this.$refs.encodeds) : undefined)
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
          location.audioSource = this.getAudioSource(location)
        })
        this.$set(encoded, 'locations', locations)
      })
      .then(() => componentHandler.upgradeElements(this.$refs[encoded.text]))
      .catch(() => {this.$set(encoded, 'expanded', false); this.notify('connection problem')})
      .then(() => this.$set(encoded, 'loading', false))
    },

    getAudioSource(location) {
      return '//verses.quran.com/AbdulBaset/Mujawwad/mp3/' + this.zeroLead(location.suraNumber, 3) + this.zeroLead(location.ayaNumber, 3) + '.mp3'
    },

    zeroLead(x, n) {
      let result = x.toString()
      while (result.length < n) result = '0' + result
      return result
    },

    setLocation(location, command) {
      this.$set(location, command+'Shown', !location[command+'Shown'])
      if (location[command]) return
      this.$set(location, command+'Loading', true)
      axios.get(`/api/${command}/${location.suraNumber}/${location.ayaNumber}`)
      .then((response) => this.$set(location, command, response.data))
      .catch(() => {this.$set(location, command+'Shown', false); this.notify('connection problem')})
      .then(() => this.$set(location, command+'Loading', false))
    },

    toggle(obj, attr) {
      this.$set(obj, attr, !obj[attr])
    },

    shift(location, n) {
      let keys = ['shiftButtonDisabled', 'ayaLoading', 'translationLoading', 'tafsirLoading']
      keys.forEach((key) => this.$set(location, key, true))

      let ayaPromise = location.ayaNumber+n === location.original.ayaNumber ?
        Promise.resolve(location.original) :
        axios.get(`/api/aya/${location.suraNumber}/${location.ayaNumber+n}`)
        .then((response)=> ({beforeHighlightedAya: response.data}))
      let translationPromise = nextTranslation(location, n, 'translation')
      let tafsirPromise = nextTranslation(location, n, 'tafsir')

      Promise.all([ayaPromise, translationPromise, tafsirPromise])
      .then(([aya, translation, tafsir]) => {
        this.$set(location, 'ayaNumber', location.ayaNumber+n)
        this.$set(location, 'audioSource', this.getAudioSource(location))
        this.$set(location, 'beforeHighlightedAya', aya.beforeHighlightedAya)
        this.$set(location, 'highlightedAya', aya.highlightedAya)
        this.$set(location, 'afterHighlightedAya', aya.afterHighlightedAya)
        this.$set(location, 'translation', translation)
        this.$set(location, 'tafsir', tafsir)
      })
      .catch(() => this.notify('connection problem'))
      .then(() => keys.forEach((key) => this.$set(location, key, false)))
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

let nextTranslation = (location, n, command) =>
  location[command+'Shown'] ?
  axios.get(`/api/${command}/${location.suraNumber}/${location.ayaNumber+n}`)
  .then((response)=> response.data) :
  Promise.resolve()

let updateTrendingKeywords = () => axios.get('/api/trending_keywords')
  .then((response) => app.trendingKeywords = response.data, undefined)
updateTrendingKeywords()
setInterval(updateTrendingKeywords, 30*1000)

let shareLinkClipboard = new Clipboard('#share-link', {text: () => app.shareLink})
shareLinkClipboard.on('success', () => app.notify('share link copied to clipboard'))

let quranTextClipboard = new Clipboard('.clipboard', {text: (trigger) => trigger.innerText})
quranTextClipboard.on('success', (e) => {
  select(e.trigger)
  let splittedText = e.text.split(' ')
  let text = splittedText.length <= 5 ? e.text : '...' + splittedText.slice(0, 5).join(' ')
  app.notify(text + ' copied to clipboard')
})

let rangeObj = document.createRange()
function select(element) {
  rangeObj.selectNodeContents(element)
  let selection = window.getSelection()
  selection.removeAllRanges()
  selection.addRange(rangeObj)
}

window.onpopstate = function(event) {
  app.keyword = event.state || ''
}
