<template>
    <div class="map">
        <div class="map__header" @click="expandMap">{{ data.Address }}</div>

        <map-frame
        v-if="mapLoadOnce"
        :data="mapData">
    </map-frame>
</div>
</template>

<script>
import MapFrame from './map.vue';
import Controllers from './controllers.vue';

import axios from 'axios'

export default {
    components: {
        'map-frame' : MapFrame,
        'controllers-view' : Controllers,
    },

    props : {
        data: {
            type: Array,
        },
    },

    data: () => ({
        mapLoadOnce: false,
        mapData: [],
    }),

    methods: {
        expandMap() {
            console.log(this.data.ID);
            const json = JSON.stringify({ id: this.data.ID });
            axios.post('/estimation', json, {
                headers: {
                    'Content-Type': 'application/json'
                }
            })
            .then(response => {
                // if (response.status == 200) {
                //     this.mapData = response.data != "" ? response.data : null;
                //
                //     if (this.mapData) {
                //         this.showMap = true;
                //     }
                // }

                console.log(response)
            })
            .catch(error => {
                console.log(error.response)
            });
        },
    }
}
</script>
