<template>
    <div class="table">
        <div class="table__container">
            <div class="table__header">
                <div class="table__header-wrapper">
                    <div class="table__field column-2in15">Address</div>
                    <div class="table__field column-1in15">Rooms</div>
                    <div class="table__field column-2in15">Type</div>
                    <div class="table__field column-1in15">Height</div>
                    <div class="table__field column-2in15">Material</div>
                    <div class="table__field column-1in15">Floor</div>
                    <div class="table__field column-1in15">Area</div>
                    <div class="table__field column-1in15">Kitchen</div>
                    <div class="table__field column-1in15">Balcony</div>
                    <div class="table__field column-1in15">Subway</div>
                    <div class="table__field column-2in15">Condition</div>
                </div>
            </div>
            <div class="table__table">
                <line-view
                    v-for="(item, ind) in data" :key="ind"
                    :data="item"
                    :isActive="activeLine == ind ? true : false"
                    :id = "ind"
                    :makeActive = "setActiveLine"
                ></line-view>
            </div>
        </div>
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
                // fields: ["Address", "Rooms", "Type", "Height", "Material", "Floor", "Area", "Kitchen", "Balcony", "Subway", "Condition",]
            }
        },

        methods: {
            setActiveLine: function(ind) {
                this.activeLine = ind;
            }
        }
    }
</script>
