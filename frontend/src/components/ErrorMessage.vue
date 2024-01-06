<template>
    <div class="my-2 text-danger text-center">
        <p>Error occured:</p>
        <p>{{ error }}</p>
    </div>
</template>

<script setup>
const { error } = defineProps(['error'])
</script>

<style scoped></style>