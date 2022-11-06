<template>
    <div class="settings">
        <div class="settings__header">
            <div v-for="item in tabs" :key="item.key" class="settings__header-item" :class="activeTab == item.key ? 'settings__header-item-active' : ''" @click="selectTab(item.key)">
                {{ item.value }}
            </div>
        </div>

        <div v-for="key in dataKeys" :key="key" class="settings__tab" :class="activeTab == key ? 'settings__tab-active' : ''">
            <div v-if="key == 'tender'" class="settings__tab-content">
                <div class="settings__line">
                    <div class="settings__label settings__label-horizontal">
                        {{ labels.tender[1] }}
                    </div>
                </div>
                <div class="settings__line">
                    <div class="settings__label settings__label-vertical">
                        {{ labels[key][0] }}
                    </div>
                    <input
                        type="text"
                        class="settings__field"
                        :value="data[key][0]"
                        :data-tab="key"
                        @input="handleInput"
                    >
                </div>
            </div>

            <div v-else class="settings__tab-content">
                <div class="settings__line">
                    <div v-for="(item, index) in labels[key]" :key="index" class="settings__label settings__label-horizontal">
                        {{ item }}
                    </div>
                </div>
                <div v-for="(line, i) in data[key]" :key="i" class="settings__line">
                    <div class="settings__label settings__label-vertical">
                        {{ labels[key][i] }}
                    </div>
                    <input
                        v-for="(item, j) in line" :key="j"
                        type="text"
                        class="settings__field"
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
        activeTab: 'tender',
        tabs: [
            {
                key: 'tender',
                value: 'Торг',
            },
            {
                key: 'floor',
                value: 'Этаж',
            },
            {
                key: 'area',
                value: 'Общая площадь',
            },
            {
                key: 'kitchen',
                value: 'Площадь кухни',
            },
            {
                key: 'balcony',
                value: 'Наличие балкона',
            },
            {
                key: 'metro',
                value: 'Расстояние до метро',
            },
            {
                key: 'condition',
                value: 'Состояние',
            },
        ],
        labels: {
            tender: ['Корректировка на торг',  'Значение'],
            floor: ['Первый', 'Средние', 'Последний'],
            area: ['<30', '30-50', '50-65', '65-90', '90-120', '210+'],
            kitchen: ['До 7', '7-10', '10-15'],
            balcony: ['Нет', 'Есть'],
            metro: ['До 5', '5-10', '10-15', '15-30', '30-60', '60-90'],
            condition: ['Без отделки', 'Эконом', 'Улучшенный'],
        }
    }),

    computed: {
        dataKeys() {
            return Object.keys(this.data);
        }
    },

    methods: {
        selectTab(tab) {
            this.activeTab = tab;
        },

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
