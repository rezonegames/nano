<template>
    <div class="container">
        <div class="score-row" v-for="(room, index) in roomlist" v-bind:key="index">
            <button type="button" v-on:click="joinRoom(room.service)">{{room.name}}-人数{{room.num}}</button>
        </div>
        <div v-if="roomlist.length === 0">
            <h3>Loading...</h3>
        </div>
    </div>
</template>

<script>
    import {mapActions,mapGetters} from 'vuex'
    export default {
        name: "RoomList",
        computed:{
            ...mapGetters({
                starx:'starx',
                roomlist:'roomlist'
            })
        },
        methods:{
            ...mapActions([
                'setTableList',
            ]),
            joinRoom: function (service) {
                this.starx.request(service+".Join", {}, (data)=>{
                    console.log("tablelist:", data.tables);
                    this.setTableList(data.tables);
                    this.$router.push('tablelist');
                })
            }
        }
    }
</script>

<style scoped>
    .container {
        max-width: 100%;
        display: flex;
        flex-direction: column;
    }
</style>