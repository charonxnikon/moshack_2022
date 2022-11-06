<template>
    <div class="map__list">
        <!-- <div class="map__list-header">
            Эталон | {{ data.Address }}
        </div> -->

        <div v-for="(item, index) in data.Analogs" :key="index" class="map__list-item" :class="excludedLines.includes(item.ID) ? 'map__list-item-excluded' : ''">
            <span>{{ item.Address }}</span>
            <button class="map__button" type="button" @click="exclude(item.ID)"></button>
        </div>
    </div>
</template>

<script>
export default {
    props: {
        data: {
            type: Object,
        },
        excludeLine: {
            type: Function,
        },
    },

    data: () => ({
        excludedLines: [],
    }),

    // computed: {
    //     checkExcluded() {
    //         return this.isExcluded;
    //     }
    // },

    methods: {
        exclude(id) {
            const index = this.excludedLines.indexOf(id);
            if (index > -1) {
                this.excludedLines.splice(index, 1);
            } else {
                this.excludedLines.push(id);
            }
            this.excludeLine(this.data.ID, id);
        }
    },
}
</script>
