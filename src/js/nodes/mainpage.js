require("./../../css/modules/mainpage.scss");
// require('vue-slider-component/theme/default.css');

import Vue from 'vue';

import Mainpage from './../modules/mainpage.vue';

new Vue({
    el: '#mainpage',
    render(h) {
        return h(Mainpage);
    }
});
