<template>
    <div className="container">
        <div v-if="countdown !== 0">
            <h3>倒计时...{{countdown}}</h3>
        </div>

        <div v-if="countdown === 0">
            <p>
                已加入
            </p>

            <button type="button" v-on:click="ready()">准备</button>
            <button type="button" v-on:click="cancelReady()">取消</button>
            <button type="button" v-on:click="leave()">离开桌子</button>

            <h3></h3>
        </div>
    </div>
</template>

<script>
    import {mapGetters, mapActions} from 'vuex'

    export default {
        name: "Ready",
        data () {
            return {
                countdown: 0
            }
        },
        computed: {
            ...mapGetters({
                starx: 'starx',
            })
        },
        created:function () {
            this.starx.on("onJoinTable", this.onJoinTable);
            this.starx.on("onReady", this.onReady);
            this.starx.on("onCountdown", this.onCountdown);
            this.starx.on("onStart", this.onStart);
            this.starx.on("onCancelReady", this.onCancelReady);
        },
        destroyed: function(){
            this.starx.off("onJoinTable", this.onJoinTable);
            this.starx.off("onReady", this.onReady);
            this.starx.off("onCountdown", this.onCountdown);
            this.starx.off("onStart", this.onStart);
            this.starx.off("onCancelReady", this.onCancelReady);
        },
        methods: {
            ...mapActions([
                "setBlue",
                "setRed",
            ]),
            ready: function () {
                this.starx.request("TableService.Ready", {}, (data) => {
                    console.log("TableService.Ready:", data);
                })
            },
            cancelReady: function () {
                this.starx.request("TableService.CancelReady", {}, (data) => {
                    console.log("TableService.CancelReady:", data);
                })
            },
            leave: function() {
                this.starx.request("TableService.Leave", {}, (data)=>{
                    console.log("TableService.Leave:", data);
                    this.$router.push('tablelist');
                })
            },
            onJoinTable: function (data) {
                console.log("OnJoinTable:", data);
            },
            onReady: function (data) {
                console.log("OnReady:", data);
            },
            onCountdown: function (data) {
                console.log("OnCountdown.........:", data);
                this.countdown = data.countdown;
            },
            onStart: function (data) {
                console.log("OnStart:", data);
                this.setBlue(data.blue);
                this.setRed(data.red);
                this.$router.push('game');
            },
            onCancelReady: function (data) {
                console.log("onCancelReady:", data);
            }
        }
    }
</script>

<style scoped>
    .container {
        max-width: 100%;
        height: 76vh;

        display: flex;
        flex-direction: column;
        margin-bottom: 20px;
    }

</style>