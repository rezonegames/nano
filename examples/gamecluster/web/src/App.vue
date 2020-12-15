<template>
  <div id="app">
    <div class="header-bar" v-if="hasconnect">
      <ProgressBar/>
      <SoundToggle/>
      <User/>
    </div>
    <div class="content" v-else>
      <h3>正在重连.....</h3>
    </div>

    <RouterView/>

  </div>
</template>

<script>
  import {mapGetters} from "vuex";
  import ProgressBar from "./components/ProgressBar";
  import SoundToggle from "./components/SoundToggle";
  import User from "./components/User";
  export default {
    components: {
      User,
      SoundToggle,
      ProgressBar
    },
    computed: {
      ...mapGetters({
        hasconnect:'hasconnect',
        routePath:'routePath'
      })
    },
}
</script>

<style>

  <style>
  html {
    height: 100%;
  }

  body {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    background-color: grey;
    margin: 0;
  }

  #app {
    color: #2c3e50;
    margin-top: -50px;
    max-width: 600px;
    font-family: Source Sans Pro, Helvetica, sans-serif;
    text-align: center;
    display: flex;
    flex-direction: column;
    min-height: 100vh;
    justify-content: center;
    align-items: center;
  }

  _:-ms-fullscreen, :root #app {
    margin-top: 0px !important;
  }
</style>
