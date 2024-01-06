<template>
    <form class="d-flex" @submit.prevent>
        <input type="text" class="form-control" placeholder="Name" v-model="customOption.name">
        <button @click="emit('addButtonClicked', customOption)" class="btn btn-success mx-2">
            Add
        </button>
        <button @click="emit('cancelButtonClicked')" class="btn btn-danger">
            Cancel
        </button>
    </form>
</template>

<script setup>
import { ref } from 'vue'
const emit = defineEmits(['addButtonClicked', 'cancelButtonClicked'])
const customOption = ref({
    name: ""
})
</script>

<style scoped></style>