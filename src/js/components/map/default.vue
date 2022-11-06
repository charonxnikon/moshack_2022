<template>
    <div class="maps">
        <controllers-view
            :data="settingsValues"
            :changeSettings="changeSettings">
        </controllers-view>

        <div v-for="(item, index) in data" :key="mapKey[item.ID]" class="map">
            <div class="map__header" @click="expandMap(item.ID)">
                <span>{{ item.Address }}</span>
                <button type="button" name="button" @click="recalculate(item.ID)">Перерассчитать</button>
            </div>

            <map-frame
                v-if="loadedData[item.ID] !== null"
                :key=""
                :data="item"
                :loaded="loadedData[item.ID]"
                :excludeAnalog="excludeAnalog">
            </map-frame>
        </div>
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
        loadedData: {},
        mapKey: {},
        excludedAnalogs: {},
        settingsValues: {
            tender: [-4.5],
            floor: [
                [0, -7, -3.1],
                [7.5, 0, 4.2],
                [3.2, 4, 0],
            ],
            area: [
                [0, 6, 14, 21, 28, 31],
                [-6, 0, 7, 14, 21, 24],
                [-12, -7, 0, 6, 13, 6],
                [-17, -12, -6, 0, 6, 9],
                [-22, -17, -11, -6, 0, 3],
                [-24, -19, -13, -8, -3, 0],
            ],
            kitchen: [
                [0, -2.9, -8.3],
                [3, 0, 5.5],
                [9, 5.8, 0],
            ],
            balcony: [
                [0, -5],
                [5.3, 0],
            ],
            metro: [
                [0, 7, 12, 17, 24, 29],
                [-7, 0, 4, 9, 15, 20],
                [-11, -4, 0, 5, 11, 15],
                [-15, -8, -5, 0, 6, 10],
                [-19, -13, -10, -6, 0, 4],
                [-22, -17, -13, -9, -4, 0],
            ],
            condition: [
                [0, -13400, -20100],
                [13400, 0, 6700],
                [20100, 6700, 0],
            ],
        }
    }),

    methods: {
        expandMap(id) {
            if (this.loadedData[id] === null) {
                const json = JSON.stringify({ id: id });
                axios.post('/estimation', json, {
                    headers: {
                        'Content-Type': 'application/json'
                    }
                })
                .then(response => {
                    if (response.status == 200) {
                        this.loadedData[id] = response.data != "" ? response.data : null;
                    }

                    this.forceRerender();

                    console.log(response)
                })
                .catch(error => {
                    console.log(error.response)
                });
            } else {
                console.log('skip');
            }
        },

        recalculate(id) {
            var analogues = this.loadedData[id].Analogs.filter((item) => {
                const index = this.excludedAnalogs[id].indexOf(item.ID);
                return index == -1;
            }).map((item) => {
                return item.ID;
            });

            var requestData = {
                id: id,
                analogs: analogues,
                ...this.settingsValues,
            }

            const json = JSON.stringify(requestData);
            axios.post('/reestimation', json, {
                headers: {
                    'Content-Type': 'application/json'
                }
            })
            .then(response => {
                console.log(response)
            })
            .catch(error => {
                console.log(error.response)
            });
        },

        excludeAnalog(parentID, childID) {
            const index = this.excludedAnalogs[parentID].indexOf(childID);
            if (index > -1) {
                this.excludedAnalogs[parentID].splice(index, 1);
            } else {
                this.excludedAnalogs[parentID].push(childID);
            }
        },

        changeSettings(field, value, row, col) {
            if (field == 'tender') {
                this.settingsValues[field][row] = parseFloat(value);
            } else {
                this.settingsValues[field][row][col] = parseFloat(value);
            }
        },

        forceRerender() {
            this.mapKey += 1;
        },
    },

    created() {
        var IDs = this.data.map((item) => {
            return item.ID;
        });
        IDs.forEach((item) => {
            this.loadedData[item] = null;
            this.mapKey[item] = 0;
            this.excludedAnalogs[item] = [];
        });
    },
}
</script>
