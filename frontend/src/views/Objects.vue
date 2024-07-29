<template>
    <div class="d-flex justify-content-end my-3">
        <RouterLink :to="{ name: 'manage_custom_options' }" class="btn btn-primary">
            Manage custom options
        </RouterLink>
    </div>

    <div v-if="!addObjectFormOpened" class="d-flex justify-content-center">
        <button @click="openAddObjectForm" class="btn btn-success px-5">
            New object
        </button>
    </div>

    <div v-if="addObjectFormOpened" class="d-flex justify-content-center">
        <ObjectForm @submitButtonClicked="add" @cancelButtonClicked="hideAddObjectForm" operation="add"
            :comparisonCustomOptions="comparisonCustomOptions" />
    </div>

    <ErrorMessage v-if="error" :error="error" />
    <div v-for="object in objects">
        <ComparedObject v-if="!object.updateFormOpened" :key="object.id" :object="object"
            @updateButtonClicked="openUpdateObjectForm" @deleteButtonClicked="remove" class="my-3 mx-1" />
        <ObjectForm v-else @submitButtonClicked="update" @cancelButtonClicked="hideUpdateObjectForm" :object="object"
            operation="update" />
    </div>
</template>

<script setup>
import ComparedObject from '@/components/ComparedObject.vue'
import ErrorMessage from '@/components/ErrorMessage.vue'
import ObjectForm from '@/components/ObjectForm.vue'
import { ref, onMounted } from 'vue'
import { getAllObjects, updateObject, deleteObject, createObject, uploadPhoto } from '@/api/objects'
import { getComparison } from '@/api/comparisons'
import { getCustomOption } from '@/api/custom_options'
import { useRoute, RouterLink } from 'vue-router'

const objects = ref([])
const error = ref("")
const route = useRoute()
const addObjectFormOpened = ref(false)
const comparisonCustomOptions = ref([])
const comparisonId = route.params.id

onMounted(async () => {
    await fetchInfo()
})

async function fetchInfo() {
    await fetchObjects()
    await fetchCustomOptions()
    await mapOptions()
}

async function fetchObjects() {
    try {
        await getAllObjects(comparisonId).then(response => {
            let result = response.data
            objects.value = result.data
        })
    } catch (e) {
        console.error(e);
        if (e.response?.data?.message) {
            error.value = e.response.data.message
        } else {
            error.value = "Internal error"
        }
    }
}

async function fetchCustomOptions() {
    try {
        let comparisonCustomOptionsIds = []
        await getComparison(comparisonId).then(response => {
            let result = response.data
            comparisonCustomOptionsIds = result.data.custom_option_ids
        })
        for (const customOptionId of comparisonCustomOptionsIds) {
            try {
                await getCustomOption(customOptionId).then(response => {
                    let result_option = response.data
                    comparisonCustomOptions.value.push(result_option.data)
                })
            } catch (e) {
                if (e.response.status === 404) {
                    continue
                }

                throw e
            }
        }
    } catch (e) {
        console.error(e);
        if (e.response?.data?.message) {
            error.value = e.response.data.message
        } else {
            error.value = "Internal error"
        }
    }
}

async function mapOptions() {
    for (const object of objects.value) {
        for (const customOption of object.custom_options) {
            customOption.name = comparisonCustomOptions.value.find(option => option.id === customOption.id).name
        }
    }
}

function openAddObjectForm() {
    addObjectFormOpened.value = true
}

function hideAddObjectForm() {
    addObjectFormOpened.value = false
}

function openUpdateObjectForm(object) {
    object.updateFormOpened = true
}

function hideUpdateObjectForm(object) {
    object.updateFormOpened = false
}

async function update(object) {
    try {
        object.comparison_id = comparisonId
        await updateObject(object)
        if (object.photo) {
            await uploadPhoto(object.id, object.photo)
        }
        await fetchInfo()
    } catch (e) {
        console.error(e);
        if (e.response?.data?.message) {
            error.value = e.response.data.message
        } else {
            error.value = "Internal error"
        }
    }
}

async function remove(object_id) {
    try {
        await deleteObject(object_id)
        await fetchObjects()
    } catch (e) {
        console.error(e);
        if (e.response?.data?.message) {
            error.value = e.response.data.message
        } else {
            error.value = "Internal error"
        }
    }
}

async function add(object) {
    try {
        object.comparison_id = comparisonId
        await createObject(object).then(response => {
            let result = response.data
            object.id = result.data.id
        })
        addObjectFormOpened.value = false
        if (object.photo) {
            await uploadPhoto(object.id, object.photo)
        }
        await fetchInfo()
    } catch (e) {
        console.error(e);
        if (e.response?.data?.message) {
            error.value = e.response.data.message
        } else {
            error.value = "Internal error"
        }
    }
}

</script>

<style scoped></style>