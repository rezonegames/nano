import Vue from 'vue'
import VueRouter from 'vue-router'
import {sync} from 'vuex-router-sync'
import {mapGetters,mapActions} from 'vuex'
import UUID from "vue-uuid";
const axios = require('axios');

import store from "./store";
import App from './App.vue'
import LoginView from "./components/LoginView";
import Ranking from "./components/Ranking";
import RoomList from "./components/RoomList";
import TableList from "./components/TableList";
import Game from "./components/Game";
import Ready from "./components/Ready";
import Settle from "./components/Settle";

Vue.config.productionTip = false;

Vue.use(VueRouter);
Vue.use(UUID);

const routes = [
  {
    path: '/',
    component: App,
    children: [
      {
        path: 'loginview',
        component: LoginView,
      },
      {
        path: 'ranking',
        component: Ranking,
      },
      {
        path: 'roomlist',
        component: RoomList,
      },
      {
        path: 'tablelist',
        component: TableList,
      },
      {
        path: 'ranking',
        component: Ranking,
      },
      {
        path: 'game',
        component: Game,
      },
      {
        path: 'ready',
        component: Ready,
      },
      {
        path: 'settle',
        component: Settle,
      },
    ]
  }
];

const router = new VueRouter({
  routes // short for routes: routes
});

router.beforeEach((to, from, next) => {
  // console.log("beforeEach to:", to, "from:", from);
  if (to.path === '/') {
    next({path: '/loginview'})
    return
  }
  if (!store.getters.user && to.path !== "/loginview") {
    next({path: '/loginview'})
    return
  }
  next()
});

sync(store, router);

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  render: c=>c(App),
  computed:{
    ...mapGetters({
      starx: 'starx',
    })
  },
  methods:{
    ...mapActions([
        "setStarx",
        "setConnected"
    ]),
  },
  created: function () {
    axios.get("http://127.0.0.1:8090/addr").then((response, data)=>{
      console.log("advise gate addr:", response, data);
      let result = response.data;
      let host = result.addr.split(":")[0];
      let port = result.addr.split(":")[1];
      let path = result.path;
      console.log("ws gate addr:", host, port, path);
      //store starx
      this.setStarx(window.starx);
      this.starx.init({host, port, path}, () => {
        console.log("initialized");
        this.setConnected(true);
        if(!this.starx.hasListeners('close')){
          this.starx.on("close", ()=>{
            this.setConnected(false);
          });
        }
      })
    })
  },
  destroyed: function () {
    this.starx.removeAllListeners();
  }
});
