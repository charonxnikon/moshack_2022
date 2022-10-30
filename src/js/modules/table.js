require("./../../css/modules/table.scss");

import Vue from 'vue';

import Table from './../components/table/default.vue';

new Vue({
    el: '#table',
    render(h) {
        return h(Table);
    }
});
