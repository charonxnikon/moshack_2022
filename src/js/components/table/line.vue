<template>
    <div class="table__line" @click="handleClick">
        <div class="table__line-fields">
            <div class="table__field column-2in15"> {{ data.address }} </div>
            <div class="table__field column-1in15"> {{ data.rooms }} </div>
            <div class="table__field column-2in15"> {{ data.type }} </div>
            <div class="table__field column-1in15"> {{ data.height }} </div>
            <div class="table__field column-1in15"> {{ data.material }} </div>
            <div class="table__field column-1in15"> {{ data.floor }} </div>
            <div class="table__field column-1in15"> {{ data.area }} </div>
            <div class="table__field column-1in15"> {{ data.kitchen }} </div>
            <div class="table__field column-1in15"> {{ data.balcony }} </div>
            <div class="table__field column-2in15"> {{ data.metro }} </div>
            <div class="table__field column-2in15"> {{ data.condition }} </div>
        </div>

        <div v-if="isActive" class="table__line-navigation">
            <button class="table__button" type="button">Рассчитать</button>
        </div>
    </div>
</template>

<script>
    export default {
        props : {
            data: {
                type: Array,
            },
            isActive: {
                type: Boolean,
                default: false,
            },
            id: {
                type: Number,
            },
            makeActive: {
                type: Function,
            },
        },

        methods: {
            handleClick: function() {
                this.makeActive(this.id);
            }
        }
    }
</script>
