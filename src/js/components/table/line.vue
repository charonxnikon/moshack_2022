<template>
    <div class="table__line" :class="hideFields ? 'table__line-collapsed' : ''">
        <div class="table__line-fields">
            <div class="table__field table__field-wide" @click="expandLine">
                <span class="table__field-label">Адрес</span>
                <span class="table__field-info">{{ data.Address }}</span>
            </div>
            <div v-for="(field, ind) in fields" :key="ind" class="table__field table__field-hidden">
                <span class="table__field-label">{{ field.fieldLabel }}</span>
                <span class="table__field-info">{{ field.fieldData }}</span>
            </div>
        </div>

        <div class="table__line-navigation">
            <button class="table__button" :class="!isSelected ? '' : 'table__button-remove'" type="button" @click="handleClick"></button>
        </div>
    </div>
</template>

<script>
    export default {
        props : {
            data: {
                type: Array,
            },
            id: {
                type: Number,
            },
            selectLine: {
                type: Function,
            },
        },

        computed: {
            fields: function() {
                var keys = Object.keys(this.data).filter(key => !(['ID', 'Latitude', 'Longitude', 'Address'].includes(key)));

                return keys.map((key) => {
                    var label = '';
                    switch (key) {
                        case 'Rooms':
                            label = 'Кол-во комнат';
                            break;
                        case 'Height':
                            label = 'Кол-во этажей';
                            break;
                        case 'Floor':
                            label = 'Этаж';
                            break;
                        case 'Area':
                            label = 'Площадь';
                                break;
                        case 'Kitchen':
                            label = 'Площадь кухни';
                            break;
                        case 'Balcony':
                            label = 'Балкон';
                                break;
                        case 'Type':
                            label = 'Тип';
                            break;
                        case 'Material':
                            label = 'Материал';
                            break;
                        case 'Metro':
                            label = 'Станция';
                            break;
                        case 'Condition':
                            label = 'Состояние';
                                break;
                        case 'PriceM2':
                            label = 'Цена за квадратный метр';
                                break;
                        case 'TotalPrice':
                            label = 'Цена';
                                break;
                    }

                    return {
                        fieldLabel: label,
                        fieldData: this.data[key]
                    }
                });
            },
        },

        data() {
            return {
                hideFields: true,
                isSelected: false,
            }
        },

        methods: {
            handleClick: function() {
                this.isSelected = !this.isSelected;
                this.selectLine(this.id);
            },

            expandLine: function(e) {
                this.hideFields = !this.hideFields;
            },
        }
    }
</script>
