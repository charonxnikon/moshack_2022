<template>
    <div class="table">
        <div class="table__content">
            <line-view
                v-for="(item, ind) in data" :key="ind"
                :data="item"
                :id="ind"
                :selectLine="selectLine"
                :calculate="calculateMethod"
            ></line-view>
        </div>

        <div class="table__footer">
            <button class="table__submit" type="button" @click="calculateMethod">Рассчитать</button>
        </div>
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
                selectedLines: {},
            }
        },

        computed: {
            dataMapped: function() {
                return Object.assign({}, this.data);
            },
        },

        methods: {
            selectLine: function(ind) {
                if (this.selectedLines[ind] === undefined) {
                    this.selectedLines[ind] = true;
                } else {
                    this.selectedLines[ind] = !this.selectedLines[ind];
                }
            },

            calculateMethod: function() {
                this.calculate(this.selectedLines);
            },
        },
    }
</script>
