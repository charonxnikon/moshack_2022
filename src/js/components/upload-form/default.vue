<template>
    <form class="form" @submit="uploadFile">
        <label for="upload-file" class="form__label">Excel-файл:</label>
        <input id="upload-file" class="form__input" type="file" name="xls_file">
        <div v-if="errorMessage" class="form__warning">Что то пошло не так! Попробуйте еще раз.</div>
        <input class="form__button" type="submit" value="Загрузить">
    </form>
</template>

<script>
    export default {
        props : {
            errorMessage: {
                type: Boolean,
            },
            onSubmit: {
                type: Function,
            },
        },

        methods: {
            uploadFile: function(event) {
                event.preventDefault();

                this.onSubmit(event);
            }
        }
    }
</script>
