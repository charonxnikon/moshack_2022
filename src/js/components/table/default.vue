<template>
    <div class="table">
        <line-view
            v-for="(item, ind) in data" :key="ind"
            :data="item"
            :isActive="activeLine == ind ? true : false"
            :id = "ind"
            :makeActive = "setActiveLine"
        ></line-view>
    </div>
</template>

<script>
    import LineView from './line.vue';

    export default {
        components: {
            'line-view' : LineView,
        },

        data() {
            return {
                activeLine: -1,
                data: [{
                    address: "улица Пупкина д. 1",
                    rooms: "3",
                    type: "Новостройка",
                    height: "10",
                    material: "кирпич",
                    floor: "7",
                    area: "100",
                    kitchen: "20",
                    balcony: "есть",
                    metro: "Авиамоторная",
                    condition: "под чистовую отделку",
                },
                {
                    address: "улица Пупкина д. 1",
                    rooms: "3",
                    type: "Новостройка",
                    height: "10",
                    material: "кирпич",
                    floor: "7",
                    area: "100",
                    kitchen: "20",
                    balcony: "есть",
                    metro: "Авиамоторная",
                    condition: "под чистовую отделку",
                },{
                    address: "улица Пупкина д. 1",
                    rooms: "3",
                    type: "Новостройка",
                    height: "10",
                    material: "кирпич",
                    floor: "7",
                    area: "100",
                    kitchen: "20",
                    balcony: "есть",
                    metro: "Авиамоторная",
                    condition: "под чистовую отделку",
                },{
                    address: "улица Пупкина д. 1",
                    rooms: "3",
                    type: "Новостройка",
                    height: "10",
                    material: "кирпич",
                    floor: "7",
                    area: "100",
                    kitchen: "20",
                    balcony: "есть",
                    metro: "Авиамоторная",
                    condition: "под чистовую отделку",
                },],
            }
        },

        methods: {
            setActiveLine: function(ind) {
                this.activeLine = ind;
            }
        }
    }
</script>
