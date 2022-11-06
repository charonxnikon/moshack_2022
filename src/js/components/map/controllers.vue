<template>
    <div class="settings">
        <div v-for="key in dataKeys" :key="key" class="settings__tab" :class="activeTab ? 'settings__tab-active' : ''">
            <div v-if="key == 'tender'">
                <div class="">
                    {{ key }}
                </div>

                <input
                    type="text"
                    class="settings__tab-field"
                    :value="data[key][0]"
                    :data-tab="key"
                    @input="handleInput"
                >
            </div>

            <div v-else class="">
                <div class="">
                    {{ key }}
                </div>

                <div v-for="(line, i) in data[key]" :key="i" class="">
                    <input
                        v-for="(item, j) in line" :key="j"
                        type="text"
                        class="settings__tab-field"
                        :value="item"
                        :data-tab="key"
                        :data-row="i"
                        :data-col="j"
                        @input="handleInput"
                    >
                </div>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    props : {
        data: {
            type: Object,
        },
        changeSettings: {
            type: Function,
        },
    },

    data: () => ({
        activeTab: -1,
    }),

    computed: {
        dataKeys() {
            return Object.keys(this.data);
        }
    },

    methods: {
        handleInput(e) {
            var target = e.target;

            if (target.dataset.tab == 'tender') {
                this.changeSettings(target.dataset.tab, target.value, 0, 0);
            } else {
                this.changeSettings(target.dataset.tab, target.value, target.dataset.row, target.dataset.col);
            }

        }
    },
}
</script>
