<template>
    <form class="d-flex" @submit.prevent>
        <input type="text" class="form-control" placeholder="Name" v-model="comparison.name">
        <button @click="emit('addButtonClicked', comparison)" class="btn btn-success mx-2">
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
const comparison = ref({
    name: ""
})
</script>

<style scoped></style>