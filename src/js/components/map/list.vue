<template>
    <div class="map__list" style="width: 50%; height: 400px">
        <div class="map__list-header">
            Эталон | {{ data.Address }}
        </div>

        <div v-for="(item, index) in data.Analogs" :key="index" class="map__list-item">
            <span>{{ item.Address }}</span>
            <button class="map__button" :class="isExcluded ? 'map__button-remove' : ''" type="button" @click="exclude(item.ID)"></button>
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
        isExcluded: false,
    }),

    computed: {
        checkExcluded() {
            return this.isExcluded;
        }
    },

    methods: {
        exclude(id) {
            this.isExcluded = !this.isExcluded;
            this.excludeLine(this.data.ID, id);
        }
    },
}
</script>
