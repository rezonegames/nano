<template>
    <div className="container">
        <div className="score-row" v-for="(table, index) in talblelist" v-bind:key="index">
            <button type="button" v-on:click="joinTable(table)">ID：{{table.tableId}}-名字：{{table.name}}-拥有者：{{table.ownerName}}-描述：{{table.desc}}</button>
        </div>

        <button type="button" v-on:click="createTable()">创建桌子</button>
        <button type="button" v-on:click="leave()">退出房间</button>
    </div>
</template>

<script>
    import {mapGetters} from 'vuex'

    export default {
        name: "TableList",
        computed: {
            ...mapGetters({
                starx: 'starx',
                talblelist:'tablelist'
            })
        },
        methods: {
            joinTable: function (table) {
                this.starx.request("TableService.Join", {
                    tableId: table.tableId
                }, (data) => {
                    console.log("TableService.Join:", data);
                    this.$router.push('ready');
                })
            },
            createTable: function () {
                let service = 'TableService';
                this.starx.request(service + ".Create", {
                    name:"我的第一个桌子",
                    desc:"欢迎大家加入！！！"
                }, (data) => {
                    console.log("TableService.Join:", data);
                    this.$router.push('ready');
                })
            },
            leave: function () {
                let service = 'RoomService';
                this.starx.request(service + ".Leave", {
                }, (data) => {
                    console.log("RoomService.Leave:", data);
                    this.$router.push('roomlist');
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
        margin-bottom: 20px;
    }

</style>