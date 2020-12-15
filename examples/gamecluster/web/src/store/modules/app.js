import api from '../../lib/app'
import * as types from '../mutation-types'

// initial state
const state = {
  logos: [],
  tempLogos: [],
  answerCount: 0,
  amount: 0,
  currentLogo: {},
  previousLogo: {},
  options: [],
  gameFinished: false,
  startTime: new Date().getTime(),
  endTime: 0,
  highScores: [],
  roomlist: [],
  tablelist:[],
  starx: null,
  hasconnect: false,
  red:[],
  blue:[],
}

// getters
const getters = {
  logos: state => state.logos,
  answerCount: state => state.answerCount,
  currentLogo: state => state.currentLogo,
  previousLogo: state => state.previousLogo,
  options: state => state.options,
  gameFinished: state => state.gameFinished,
  amount: state => state.amount,
  startTime: state => state.startTime,
  endTime: state => state.endTime,
  highScores: state => state.highScores,
  roomlist: state => state.roomlist,
  starx: state => state.starx,
  tablelist: state=> state.tablelist,
  hasconnect: state=>state.hasconnect,
  red: state=>state.red,
  blue: state=>state.blue,
}

// actions
const actions = {
  initializeLogos ({ commit }, callback) {
    api.getJSON('static/logos.json', (error, tempLogos) => {
      if (typeof tempLogos === 'string') { // IE11 fix
        tempLogos = JSON.parse(tempLogos)
      }

      if (error) {
        // Fetch from localStorage
        tempLogos = JSON.parse(window.localStorage.getItem('logos'))
      } else {
        // Update localStorage
        // window.jsTools = JSON.parse(JSON.stringify(tempLogos)) try to get rid of this
        window.localStorage.setItem('logos', JSON.stringify(tempLogos))
      }

      commit(types.SET_TEMP_LOGOS, { tempLogos })

      callback()
    })
  },
  setCurrentLogo ({ commit, state }, currentLogo) {
    const previousLogo = state.currentLogo
    commit(types.SET_PREVIOUS_LOGO, { previousLogo })
    commit(types.SET_CURRENT_LOGO, { currentLogo })
  },
  setOptions ({ commit, state }) {
    const options = api.getAnswerOptions(state.logos, state.amount, state.currentLogo.id, state.previousLogo.id)
    commit(types.SET_OPTIONS, { options })
  },
  increaseAnswerCount ({ commit }) {
    commit(types.INCREASE_ANSWER_COUNT)
  },
  finishGame ({ commit }) {
    commit(types.FINISH_GAME)
  },
  restartGame ({ commit, state }) {
    const logos = api.generateIDs(api.shuffle(JSON.parse(JSON.stringify(state.tempLogos))))
    const count = 0
    const amount = logos.length
    const currentLogo = logos[count]
    const previousLogo = {}
    const options = api.getAnswerOptions(logos, amount, currentLogo.id)

    commit(types.SET_LOGOS, { logos })
    commit(types.SET_AMOUNT, { amount })
    commit(types.SET_ANSWER_COUNT, { count })
    commit(types.SET_CURRENT_LOGO, { currentLogo })
    commit(types.SET_PREVIOUS_LOGO, { previousLogo })
    commit(types.SET_OPTIONS, { options })
    commit(types.START_GAME)
  },
  setHighScores ({ commit }, scores) {
    commit(types.SET_HIGH_SCORES, { scores })
  },
  setRoomList ({ commit }, roomlist) {
    commit(types.ROOM_LIST, { roomlist })
  },
  setStarx ({ commit }, starx) {
    commit(types.STARX, { starx })
  },
  setTableList ({ commit }, tablelist) {
    commit(types.TABLE_LIST, { tablelist })
  },
  setConnected ({ commit }, hasconnect) {
    commit(types.HAS_CONNECT, { hasconnect })
  },
  setRed ({ commit }, red) {
    commit(types.RED, { red })
  },
  setBlue ({ commit }, blue) {
    commit(types.BLUE, { blue })
  },
}

// mutations
const mutations = {
  [types.SET_TEMP_LOGOS] (state, { tempLogos }) {
    state.tempLogos = tempLogos
  },
  [types.SET_AMOUNT] (state, { amount }) {
    state.amount = amount
  },
  [types.SET_LOGOS] (state, { logos }) {
    state.logos = logos
  },
  [types.SET_CURRENT_LOGO] (state, { currentLogo }) {
    state.currentLogo = currentLogo
  },
  [types.SET_PREVIOUS_LOGO] (state, { previousLogo }) {
    state.previousLogo = previousLogo
  },
  [types.SET_OPTIONS] (state, { options }) {
    state.options = options
  },
  [types.INCREASE_ANSWER_COUNT] (state) {
    state.answerCount++
  },
  [types.SET_ANSWER_COUNT] (state, { count }) {
    state.answerCount = count
  },
  [types.START_GAME] (state) {
    state.gameFinished = false
    state.startTime = new Date().getTime()
    state.endTime = 0
  },
  [types.FINISH_GAME] (state) {
    state.gameFinished = true
    state.endTime = new Date().getTime()
  },
  [types.SET_HIGH_SCORES] (state, { scores }) {
    state.highScores = scores
  },
  [types.ROOM_LIST] (state, { roomlist }) {
    state.roomlist = roomlist
  },
  [types.TABLE_LIST] (state, { tablelist }) {
    state.tablelist = tablelist
  },
  [types.STARX] (state, { starx }) {
    state.starx = starx
  },
  [types.HAS_CONNECT] (state, { hasconnect }) {
    state.hasconnect = hasconnect
  },
  [types.RED] (state, { red }) {
    state.red = red
  },
  [types.BLUE] (state, { blue }) {
    state.blue = blue
  },
}

export default {
  state,
  getters,
  actions,
  mutations
}
