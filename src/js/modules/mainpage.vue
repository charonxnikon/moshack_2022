<template>
    <div class="mainpage">
        <form-view
            v-if="!hideForm"
            :errorMessage="errorMessage"
            :onSubmit="uploadFile">
        </form-view>

        <div v-if="hideForm" class="mainpage__success">
            <div class="mainpage__success-message">
                Файл успешно загружен!
            </div>
            <div class="mainpage__success-button" @click="loadData">
                Отобразить результат
            </div>
        </div>

        <table-view
            v-if="showTable && loaded !== null"
            :data="loaded">
        </table-view>
    </div>
</template>

<script>
import UploadForm from '../components/upload-form/default.vue';
import Table from '../components/table/default.vue';

import axios from 'axios'

export default {
    components: {
        'form-view' : UploadForm,
        'table-view' : Table,
    },

    data() {
        return {
            hideForm: false,
            errorMessage: false,
            showTable: false,
            loaded: null,
        }
    },

    methods: {
        uploadFile: function(event) {
            this.hideForm = false;
            this.errorMessage = false;

            var formData = new FormData();
            var file = event.target.querySelector('#upload-file');
            formData.append("xls_file", file.files[0]);

            axios.post('/loadxls', formData, {
                headers: {
                    'Content-Type': 'multipart/form-data'
                }
            })
            .then(response => {
                if (response.status == 200) {
                    this.hideForm = true;
                } else {
                    this.errorMessage = true;
                }
                console.log(response)
            })
            .catch(error => {
                this.errorMessage = true;

                console.log(error.response)
            });
        },

        loadData: function() {
            axios.get('/estimation')
            .then(response => {
                if (response.status == 200) {
                    this.loaded = response.data != "" ? response.data : null;

                    if (this.loaded) {
                        this.showTable = true;
                    }
                }

                console.log(response)
            })
            .catch(error => {
                console.log(error.response)
            });
        },
    }
}
</script>
