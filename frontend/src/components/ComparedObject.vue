<template>
    <div class="d-flex" @mouseover="hover" @mouseleave="unhover">
        <div class="p-3 rounded border border-4 col-11">
            <div class="row">
                <div class="col-2 col-md-2">
                    <img class="img-fluid" :src="photoSrc" @error="defaultPhotoSrc" alt="object photo">
                </div>
                <div class="col-10 col-md-10">
                    <div class="d-flex justify-content-between mb-4 border-bottom pb-2">
                        <b>{{ object.name }}</b>
                        <span>Rating: {{ object.rating }}</span>
                    </div>

                    <p>
                        <b>Adventages</b>: {{ object.advs }}
                    </p>

                    <p>
                        <b>Disadventages</b>: {{ object.disadvs }}
                    </p>

                    <p v-for="custom_option in object.custom_options">
                        <b>{{ custom_option.name }}</b>: {{ custom_option.value }}
                    </p>
                </div>
            </div>
        </div>
        <div v-if="isHovered" class="row mx-2">
            <button @click="emit('updateButtonClicked', object)" class="btn btn-primary px-1 my-2">&#9998;</button>
            <button @click="emit('deleteButtonClicked', object.id)" class="btn btn-danger px-1 my-2">X</button>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { photoUrl } from '@/api/objects';
const { object } = defineProps(['object'])
const photoSrc = ref(photoUrl(object.id))
const isHovered = ref(false)
const emit = defineEmits(['updateButtonClicked', 'deleteButtonClicked'])

function defaultPhotoSrc() {
    return photoSrc.value = '/images/defaultPhoto.png'
}

function hover() {
    isHovered.value = true
}

function unhover() {
    isHovered.value = false
}

</script>

<style scoped>
p {
    margin: 0.5rem auto;
}
</style>