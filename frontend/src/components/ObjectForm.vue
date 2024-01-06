<template>
    <form class="d-flex row rounded border border-3 p-2 mx-2 my-1" @submit.prevent>
        <b class="text-center my-2">{{ formTitle }}</b>

        <input type="text" class="form-control my-1" placeholder="Name" v-model="object.name">
        <input type="number" class="form-control my-1" placeholder="Rating" v-model="object.rating">
        <input type="text" class="form-control my-1" placeholder="Adventages" v-model="object.advs">
        <input type="text" class="form-control my-1" placeholder="Disadventages" v-model="object.disadvs">
        <input v-for="customOption in customOptions" type="text" class="form-control my-1" :placeholder="customOption.name"
            v-model="customOption.value">
        <div>
            <label for="objectPhoto my-2">Object photo:</label>
            <input type="file" class="form-control-file mx-2" id="objectPhoto"
                @change="object.photo = $event.target.files[0]">
        </div>

        <div class="d-flex justify-content-center my-2">
            <button @click="emit('submitButtonClicked', object)" class="btn btn-success col-4">
                Submit
            </button>
            <button @click="emit('cancelButtonClicked', object)" class="btn btn-danger col-4 mx-2">
                Cancel
            </button>
        </div>
    </form>
</template>

<script setup>
import { ref, onMounted } from 'vue'
const emit = defineEmits(['submitButtonClicked', 'cancelButtonClicked'])
const object = ref({
    name: "",
    rating: "",
    advs: "",
    disadvs: "",
    custom_options: [],
    photo: null,
})
const props = defineProps(['operation', 'object', 'comparisonCustomOptions'])
const formTitle = ref('')
const customOptions = ref([])

onMounted(() => {
    if (props.operation === 'update') {
        formTitle.value = 'Update object'
        object.value = props.object
        customOptions.value = props.object.custom_options
    } else {
        formTitle.value = 'New object'
        customOptions.value = props.comparisonCustomOptions
    }
})

</script>

<style scoped></style>
