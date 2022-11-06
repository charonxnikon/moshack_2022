<template>
    <div class="map__container">
        <div class="map__wrapper" style="width: 50%; height: 400px">
            <yandex-map :coords="[mapData.Latitude, mapData.Longitude]">
                <ymap-marker
                    :marker-id="mapData.ID"
                    :coords="[mapData.Latitude, mapData.Longitude]"
                    :balloon-template="balloonTemplate"
                    @click="onMarkerClick"
                />
                <ymap-marker
                    v-for="(item, index) in mapData.Analogs" :key="index"
                    :marker-id="item.ID"
                    :coords="[item.Latitude, item.Longitude]"
                    :balloon-template="balloonTemplate"
                    @click="onMarkerClick"
                />
            </yandex-map>
        </div>

        <list-view
            :data="mapData"
            :showSettings="showSettings"
            :excludeLine="excludeLine">
        </list-view>
    </div>
</template>

<script>
import { yandexMap, ymapMarker } from 'vue-yandex-maps'
import List from './list.vue';

import axios from 'axios'

export default {
    components: {
        yandexMap,
        ymapMarker,
        'list-view' : List,
    },

    props : {
        data: {
            type: Object,
        },
        loaded: {
            type: Object,
        },
        excludeAnalog: {
            type: Function,
        },
    },

    data: () => ({
        selectedCoords: [],
        selectedId: -1,
        showControllers: false,
        settingsData: [],
    }),

    computed: {
        mapData() {
            this.data.PriceM2 = this.loaded.PriceM2;
            this.data.TotalPrice = this.loaded.TotalPrice;

            return {
                ...this.loaded,
                ...this.data,
            }
        },
        
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

        excludeLine(parentID, childID) {
            this.excludeAnalog(parentID, childID);
        }
    },
}
</script>
