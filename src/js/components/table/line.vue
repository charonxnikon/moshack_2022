<template>
    <div class="table__line" :class="lineClass" @click="handleClick">
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

        <div v-if="isActive" class="table__line-navigation">
            <button class="table__button" type="button" @click="calculateMethod">Рассчитать</button>
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
            calculate: {
                type: Function,
            },
        },

        computed: {
            lineClass: function() {
                var currClass = this.isActive ? 'table__line-active ' : '';
                currClass += this.hideFields ? 'table__line-collapsed ' : '';
                return currClass;
            },

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
            }
        },

        methods: {
            handleClick: function() {
                this.makeActive(this.id);
            },

            expandLine: function() {
                this.hideFields = !this.hideFields;
            },

            calculateMethod: function(id) {
                this.calculate(this.data.ID);
            },
        }
    }
</script>
