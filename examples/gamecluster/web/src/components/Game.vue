<template>
    <div class="container">
        <div class="gamecanvas" ref="game">
            <div class="player">
                <div class="score"></div>
                <canvas class="tetris" width="240" height="400"></canvas>
            </div>

            <div class="teammate">
                <div id="teamate1">
                    <div class="score"></div>
                    <canvas class="tetris" width="60" height="100"></canvas>
                </div>
                <div id="teammate2">
                    <div class="score"></div>
                    <canvas class="tetris" width="60" height="100"></canvas>
                </div>
            </div>


            <div class="enemy">
                <div id="enemy0">
                    <div class="score"></div>
                    <canvas class="tetris" width="60" height="100"></canvas>
                </div>
                <div id="enemy1">
                    <div class="score"></div>
                    <canvas class="tetris" width="60" height="100"></canvas>
                </div>
                <div id="enemy2">
                    <div class="score"></div>
                    <canvas class="tetris" width="60" height="100"></canvas>
                </div>
            </div>


        </div>

        <button type="button" v-on:click="over()">结束</button>
        <button type="button" v-on:click="moveLeft()">左</button>
        <button type="button" v-on:click="moveRight()">右</button>
        <button type="button" v-on:click="drop()">下</button>
        <button type="button" v-on:click="rotate()">变</button>
        <button type="button" v-on:click="dropDown()">速降</button>
    </div>
</template>

<script>
    import {mapGetters} from 'vuex'
    import TetrisManager from '../lib/tetris-manager'
    import ConnectionManager from '../lib/connection-manager'
    export default {
        name: "Game",
        data () {
            return {
                countdown: 0,
                player: null,
                tetrisManager: null,
                connectionManager: null,
            }
        },
        computed: {
            ...mapGetters({
                starx: 'starx',
                user: 'user',
                blue: 'blue',
                red: 'red',
            })
        },
        mounted() {
            console.log("Game start ......");
            let game = this.$refs.game;
            this.tetrisManager = new TetrisManager();
            const tetrisLocal = this.tetrisManager.createPlayer(game.querySelector(".player"));
            tetrisLocal.element.classList.add('local');
            this.player = tetrisLocal.player;
            this.connectionManager = new ConnectionManager(this.starx, this.tetrisManager, this.user.uId);

            let teammates = [];
            let enemys = [];
            if(this.red.includes(this.user.uId)){
                this.red.forEach((uId, index)=>{
                    if(uId !== this.user.uId){
                        let name = "#teammate"+index;
                        let element = game.querySelector(name);
                        console.log("name:", name, "element:", element);
                        teammates.push({
                            uId,
                            element
                        })
                    }
                })
                this.blue.forEach((uId, index)=>{
                    let name = "#enemy"+index;
                    let element = game.querySelector(name)
                    console.log("enemy name:", name, "element:", element);
                    enemys.push({
                        uId,
                        element
                    })
                })
            } else {
                this.blue.forEach((uId, index)=>{
                    if(uId !== this.user.uId){
                        let name = "#teammate"+index;
                        let element = game.querySelector(name);
                        console.log("name:", name, "element:", element);
                        teammates.push({
                            uId,
                            element
                        })
                    }
                })
                this.red.forEach((uId, index)=>{
                    let name = "#enemy"+index;
                    let element = game.querySelector(name)
                    console.log("enemy name:", name, "element:", element);
                    enemys.push({
                        uId,
                        element
                    })
                })
            }
            this.connectionManager?.updateManager(teammates);
            this.connectionManager?.updateManager(enemys);

            this.tetrisManager?.start();


            this.starx.on("onOver", this.onOver);
            this.starx.on("onStopAndSettle", this.onStopAndSettle);
        },
        destroyed:function () {
            console.log("Game destoryed ......");
            this.starx.off("onOver", this.onOver);
            this.starx.off("onStopAndSettle", this.onStopAndSettle);
        },
        methods: {
            moveLeft: function() {
                this.player.move(-1);
            },
            moveRight: function() {
                this.player.move(1);
            },
            drop: function() {

            },
            dropDown: function() {
                this.player.dropDown()
            },
            rotate: function() {
                this.player.rotate(1)
            },
            over: function () {
                this.starx.request("TableService.Over", {}, (data) => {
                    console.log("TableService.Over:", data);
                })
            },
            onOver: function (data) {
                console.log("TableService.onOver:", data);
            },
            onStopAndSettle: function (data) {
                console.log("Game over ......")
                console.log("TableService.onStopAndSettle:", data);
                this.$router.push('settle');
            }
        }
    }
</script>

<style scoped>
    .container {
        max-width: 100%;
        height: 76vh;

        margin-top: 80px;
    }
    /*.gamecanvas {*/
    /*    margin-top: 50px;*/
    /*}*/
    .player {
        float: left;
        width: 60%;
        padding: 5px;
    }
    .teammate {
        float: left;
        width: 20%;
        padding: 5px;
    }
    .enemy {
        float: left;
        width: 20%;
        padding: 5px;
    }

</style>