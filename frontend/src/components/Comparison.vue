<template>
    <div class="p-2 rounded border border-3 d-flex justify-content-between" @mouseover="hover" @mouseleave="unhover">
        {{ comparison.name }}
        <button v-if="isHovered" @click.prevent="emit('deleteButtonClicked', comparison.id)"
            class="btn btn-danger px-1 py-0">X</button>
    </div>
</template>

<script setup>
import { ref } from 'vue'
const { comparison } = defineProps(['comparison'])
const emit = defineEmits(['deleteButtonClicked'])
const isHovered = ref(false)
function hover() {
    isHovered.value = true
}
function unhover() {
    isHovered.value = false
}

</script>

<style scoped></style>