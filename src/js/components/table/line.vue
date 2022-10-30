<template>
    <div class="table__line" :class="isActive ? 'table__line-active' : ''" @click="handleClick">
        <div class="table__line-fields">
            <div v-for="(field, ind) in fields" :key="ind" class="table__field" :class="ind == 0 ? 'table__field-wide' : ''">
                <span class="table__field-label">{{ field.fieldLabel }}</span>
                <span class="table__field-info">{{ field.fieldData }}</span>
            </div>
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

        computed: {
            fields: function () {
                return Object.keys(this.data).map((key) => {
                    var label = '';
                    switch (key) {
                        case 'address':
                            label = 'Адрес';
                            break;
                        case 'rooms':
                            label = 'Кол-во комнат';
                            break;
                        case 'height':
                            label = 'Кол-во этажей';
                            break;
                        case 'floor':
                            label = 'Этаж';
                            break;
                        case 'area':
                            label = 'Площадь';
                                break;
                        case 'kitchen':
                            label = 'Площадь кухни';
                            break;
                        case 'balcony':
                            label = 'Балкон';
                                break;
                        case 'type':
                            label = 'Тип';
                            break;
                        case 'material':
                            label = 'Материал';
                            break;
                        case 'metro':
                            label = 'Станция';
                            break;
                        case 'condition':
                            label = 'Состояние';
                                break;
                    }

                    return {
                        fieldLabel: label,
                        fieldData: this.data[key]
                    }
                });
            },
        },

        methods: {
            handleClick: function() {
                this.makeActive(this.id);
            }
        }
    }
</script>
