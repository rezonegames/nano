<template>
    <div class="container">
        <button v-on:click="loginFacebook">facebook</button>

        <!--        <button v-on:click="loginGuest">guest</button>-->
    </div>
</template>

<script>
    import {mapGetters, mapActions} from 'vuex'
    import { uuid } from 'vue-uuid';
    export default {
        name: "LoginView",
        computed:{
          ...mapGetters({
              starx:'starx',
          })
        },
        methods:{
            ...mapActions([
                'setUser',
                'setRoomList',
            ]),
            loginFacebook: function () {
                this.starx.request("GateService.Login", {deviceId:uuid.v1()}, (data)=>{
                    console.log("login data:", data);
                    this.setUser({"uId":data.uId,"displayName":data.name,"photoURL":data.pic,"email":data.email,"emailVerified":false,"identifierNumber":null,"isAnonymous":false, "redirectEventId":null});
                    this.setRoomList([{name:'普通房间', num:100, service:"RoomService"}]);
                    //to roomlist
                    this.$router.push('roomlist');
                })
            },
            loginGuest: function () {
                // this.setUser({"uid":2,"displayName":"rezongame","photoURL":"https://avatars3.githubusercontent.com/u/41884353?v=4","email":"duhe@live.hk","emailVerified":false,"identifierNumber":null,"isAnonymous":false,"providerData":[{"uid":"41884353","displayName":null,"photoURL":"https://avatars3.githubusercontent.com/u/41884353?v=4","email":"duhe@live.hk","providerId":"github.com"}],"apiKey":"AIzaSyDzKhNG-bafYPDTbAFFGqSscGJ7a-rdB9M","appName":"[DEFAULT]","authDomain":"javascript-game.firebaseapp.com","stsTokenManager":{"apiKey":"AIzaSyDzKhNG-bafYPDTbAFFGqSscGJ7a-rdB9M","refreshToken":"AG8BCnefPeplD3gZc6PAswVlkWhFVTM3Y56Vppm0yZ4X1fpSbRkBeN1Kry_sX2uJvX1NKiJi-QvfTEpr-gdD7eGxFoJdMeWo4Rv8J2o9ttiyM7rfawsmSLYlCeT7JEh3Xpj7VyTGEHG_qIjpvcF3Vy-NXtmzG175P-LJ-G_LMuz0x-mufbX4H5C2isJ2UeYFIOZobAzhlPmTcfiGeg-ZkUKXR2gIES0XuLMSZ_RfjpbLaE6MiUbClTTOFrFuwdzksv59cZUguS5_aDV3PeFGkinyuKfzlh_yVLqaVdlwQxKHPRqCNZFuwG0TioCBCuN8pfiejNm2fzGH","accessToken":"eyJhbGciOiJSUzI1NiIsImtpZCI6ImI5ODI2ZDA5Mzc3N2NlMDA1ZTQzYTMyN2ZmMjAyNjUyMTQ1ZTk2MDQiLCJ0eXAiOiJKV1QifQ.eyJwaWN0dXJlIjoiaHR0cHM6Ly9hdmF0YXJzMy5naXRodWJ1c2VyY29udGVudC5jb20vdS80MTg4NDM1Mz92PTQiLCJpc3MiOiJodHRwczovL3NlY3VyZXRva2VuLmdvb2dsZS5jb20vamF2YXNjcmlwdC1nYW1lIiwiYXVkIjoiamF2YXNjcmlwdC1nYW1lIiwiYXV0aF90aW1lIjoxNjA3NTA3NTYwLCJ1c2VyX2lkIjoiREdMOWthYUdLZ1FvWEF2bFo0Yk9IbGJCcVc5MiIsInN1YiI6IkRHTDlrYWFHS2dRb1hBdmxaNGJPSGxiQnFXOTIiLCJpYXQiOjE2MDc1MTEyNDIsImV4cCI6MTYwNzUxNDg0MiwiZW1haWwiOiJkdWhlQGxpdmUuaGsiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImZpcmViYXNlIjp7ImlkZW50aXRpZXMiOnsiZ2l0aHViLmNvbSI6WyI0MTg4NDM1MyJdLCJlbWFpbCI6WyJkdWhlQGxpdmUuaGsiXX0sInNpZ25faW5fcHJvdmlkZXIiOiJnaXRodWIuY29tIn19.SOTBE5lkmOQV7TYbXYp6QfDkPd9qEkYiz9W38kkjfVpTJ_-conXJmrXfoK-azHgESUyz4ND93H6lqIdRj5Gxcokypa1gdwk_WfGShbLg3XfgABZ_7Z6bWA0J_rCxpfL0GNYypWiUUtia4YKFpTbk_ogSSwXFOhWz7Xz4g4b-hwiji-LeQOH-rG_zwXkjMaK5X2QCNg0CyXanhzHMgEHIF4qkXV0cGVLF6GRE5UeZQ1J9s-iqFk-HvyDQjBY-TwzfKvDvQyxIqK96B4pB1Y_EQSpr18hZIyxsbKOqHu3m0GSFyow8FcLrKF9cgDPkzhKLJfVRaeXbYHlEYtYFKGfV2Q","expirationTime":1607514843047},"redirectEventId":null})
                // this.$router.push('roomlist')
            }
        }
    }
</script>

<style scoped>

</style>