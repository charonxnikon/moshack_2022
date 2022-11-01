require("./../../css/modules/mainpage.scss");

import Vue from 'vue';

import Mainpage from './../modules/mainpage.vue';

new Vue({
    el: '#mainpage',
    render(h) {
        return h(Mainpage);
    }
});
