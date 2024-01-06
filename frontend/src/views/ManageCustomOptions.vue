<template>
    <div>
        <ErrorMessage v-if="error" :error="error" />

        <h2 class="text-center my-5">Custom options for {{ comparison.name }}</h2>
        <ul>
            <li @mouseover="hover(customOption)" @mouseleave="unhover(customOption)"
                v-for="customOption in comparisonCustomOptions" :key="customOption.id">
                {{ customOption.name }}
            </li>
        </ul>

        <h2 class="text-center my-5">All custom options</h2>
        <div v-if="!customOptionFormOpened" class="d-flex justify-content-center mt-5">
            <button @click="openCustomOptionForm" class="btn btn-success px-5">
                New custom option
            </button>
        </div>

        <div v-if="customOptionFormOpened" class="d-flex justify-content-center">
            <CustomOptionForm @addButtonClicked="addCustomOption" @cancelButtonClicked="hideCustomOptionForm" />
        </div>

        <ul>
            <li @mouseover="hover(customOption)" @mouseleave="unhover(customOption)"
                v-for="customOption in allCustomOptions" :key="customOption.id">
                {{ customOption.name }}
                <div class="d-inline" v-if="customOption.isHovered">
                    <button v-if="!comparisonCustomOptions.includes(customOption)"
                        @click="addCustomOptionToComparison(customOption)" class="btn btn-primary px-1 py-0 mx-1">
                        +
                    </button>
                    <button @click="removeCustomOption(customOption)" class="btn btn-danger px-1 py-0">
                        X
                    </button>
                </div>

            </li>
        </ul>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getAllCustomOptions, createCustomOption, deleteCustomOption } from '@/api/custom_options'
import { getComparison, updateComparison } from '@/api/comparisons'
import ErrorMessage from '@/components/ErrorMessage.vue'
import CustomOptionForm from '@/components/CustomOptionForm.vue'

const allCustomOptions = ref([])
const comparisonCustomOptions = ref([])
const comparison = ref({})
const route = useRoute()
const comparisonId = route.params.id
const error = ref("")
const customOptionFormOpened = ref(false)

onMounted(async () => {
    await fetchComparison()
    await fetchCustomOptions()
})

async function fetchComparison() {
    try {
        await getComparison(comparisonId).then(response => {
            let result = response.data
            comparison.value = result.data
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
        await getAllCustomOptions().then(response => {
            let result = response.data
            allCustomOptions.value = result.data
        })

        let tempComparisonCustomOptions = []
        for (const customOption of allCustomOptions.value) {
            if (comparison.value.custom_option_ids.includes(customOption.id)) {
                tempComparisonCustomOptions.push(customOption);
            }
            comparisonCustomOptions.value = tempComparisonCustomOptions
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

function hover(item) {
    item.isHovered = true
}

function unhover(item) {
    item.isHovered = false
}

async function removeCustomOption(customOption) {
    try {
        await deleteCustomOption(customOption.id)
        await fetchCustomOptions()
    } catch (e) {
        console.error(e);
        if (e.response?.data?.message) {
            error.value = e.response.data.message
        } else {
            error.value = "Internal error"
        }
    }
}

async function addCustomOptionToComparison(customOption) {
    try {
        comparison.value.custom_option_ids = [...comparison.value.custom_option_ids, customOption.id]
        await updateComparison(comparison.value)
        await fetchComparison()
        await fetchCustomOptions()
    } catch (e) {
        console.error(e);
        if (e.response?.data?.message) {
            error.value = e.response.data.message
        } else {
            error.value = "Internal error"
        }
    }
}

async function addCustomOption(customOption) {
    try {
        await createCustomOption(customOption)
        customOptionFormOpened.value = false
        await fetchCustomOptions()
    } catch (e) {
        console.error(e);
        if (e.response?.data?.message) {
            error.value = e.response.data.message
        } else {
            error.value = "Internal error"
        }
    }
}

function hideCustomOptionForm() {
    customOptionFormOpened.value = false
}

function openCustomOptionForm() {
    customOptionFormOpened.value = true
}

</script>

<style scoped></style>