<template>
    <div class="table">
        <div class="table__content">
            <line-view
                v-for="(item, ind) in data" :key="ind"
                :data="item"
                :id="ind"
                :selectLine="selectLine"
            ></line-view>
        </div>

        <div class="table__footer">
            <button class="table__submit" type="button" @click="getSelected">Рассчитать</button>
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
            getSelectedLines: {
                type: Function,
            },
        },

        data() {
            return {
                selectedLines: {},
                selectedItems: [],
            }
        },

        methods: {
            selectLine: function(ind) {
                if (this.selectedLines[ind] === undefined) {
                    this.selectedLines[ind] = true;
                } else {
                    this.selectedLines[ind] = !this.selectedLines[ind];
                }
            },

            getSelected: function() {
                var selectedIndexes = Object.keys(this.selectedLines).filter(key => this.selectedLines[key]);
                if (selectedIndexes.length == 0) {
                    return;
                }
                this.getSelectedLines(selectedIndexes);
            },
        },
    }
</script>
