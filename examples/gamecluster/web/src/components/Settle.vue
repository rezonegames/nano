<template>
    <div className="container">
        <p>
            game
        </p>

        <button type="button" v-on:click="nextRound()">继续下一局</button>
        <button type="button" v-on:click="leaveTable()">退出桌子</button>
    </div>
</template>

<script>
    import {mapGetters} from 'vuex'

    export default {
        name: "Settle",
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
            console.log("Over ......")
            // this.starx.on("OnOver", this.OnOver);
        },
        methods: {
            nextRound: function () {
                this.starx.request("TableService.Ready", {}, (data) => {
                    console.log("TableService.Ready next round:", data);
                    this.$router.push('ready');
                })
            },
            leaveTable: function () {
                this.starx.request("TableService.Leave", {}, (data) => {
                    console.log("TableService.Leave:", data);
                    this.$router.push('tablelist');
                })
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