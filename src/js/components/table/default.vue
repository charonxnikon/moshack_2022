<template>
    <div class="table">
        <line-view
            v-for="(item, ind) in data" :key="ind"
            :data="item"
            :isActive="activeLine == ind ? true : false"
            :id = "ind"
            :makeActive = "setActiveLine"
            :calculate="calculateMethod"
        ></line-view>
    </div>
</template>

<script>
    import LineView from './line.vue';

    export default {
        components: {
            'line-view' : LineView,
        },

        props : {
            data: {
                type: Array,
            },
            calculate: {
                type: Function,
            },
        },

        data() {
            return {
                activeLine: -1,
            }
        },

        methods: {
            setActiveLine: function(ind) {
                this.activeLine = ind;
            },

            calculateMethod: function(id) {
                this.calculate(id);
            },
        }
    }
</script>
