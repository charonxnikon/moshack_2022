<template>
    <div class="map" style="width: 768px; height: 432px">
        <yandex-map :coords="data[0].coords">
            <ymap-marker
                v-for="(item, index) in data" :key="index"
                :marker-id="item.id"
                :coords="item.coords"
                :balloon-template="balloonTemplate"
                @click="onMarkerClick"
            />
        </yandex-map>
    </div>
</template>

<script>
import { yandexMap, ymapMarker } from 'vue-yandex-maps'

export default {
    components: { yandexMap, ymapMarker },

    props : {
        data: {
            type: Array,
        },
    },

    data: () => ({
        selectedCoords: [],
        selectedId: -1,
    }),

    computed: {
        balloonTemplate() {
            return `
            <h1 class="red">I have id ${this.selectedId}</h1>
            <p>I am here: ${this.selectedCoords}</p>
            `
        },
    },

    methods: {
        onMarkerClick(e) {
            var marker = e.get('target');
            this.selectedCoords = e.get('coords');
            this.selectedId = marker.properties.get('markerId');
        },
    }
}
</script>
