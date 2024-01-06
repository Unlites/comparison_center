<template>
    <div class="mx-3">
        <div v-if="!comparisonFormOpened" class="d-flex justify-content-center">
            <button @click="openComparisonForm" class="btn btn-success">
                New comparison
            </button>
        </div>

        <div v-if="comparisonFormOpened" class="d-flex justify-content-center">
            <ComparisonForm @addButtonClicked="addComparison" @cancelButtonClicked="hideComparisonForm" />
        </div>

        <ErrorMessage v-if="error" :error="error" />

        <div>
            <RouterLink v-for="comparison in comparisons" class="d-block my-2 text-black text-decoration-none"
                :key="comparison.id" :to="{ name: 'comparison_objects', params: { id: comparison.id } }">
                <Comparison :comparison="comparison" @deleteButtonClicked="removeComparison" />
            </RouterLink>
        </div>
    </div>
</template>

<script setup>
import Comparison from '@/components/Comparison.vue'
import ComparisonForm from '@/components/ComparisonForm.vue'
import ErrorMessage from '@/components/ErrorMessage.vue'
import { RouterLink } from 'vue-router'
import { ref, onMounted } from 'vue'
import { getAllComparisons, createComparison, deleteComparison } from '@/api/comparisons'

const comparisons = ref([])
const comparisonFormOpened = ref(false)
const error = ref("")

onMounted(async () => {
    await fetchComparisons()
})

async function fetchComparisons() {
    try {
        await getAllComparisons().then(response => {
            let result = response.data
            comparisons.value = result.data
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

function openComparisonForm() {
    comparisonFormOpened.value = true
}

function hideComparisonForm() {
    comparisonFormOpened.value = false
}

async function addComparison(comparison) {
    try {
        await createComparison(comparison)
        comparisonFormOpened.value = false
        await fetchComparisons()
    } catch (e) {
        console.error(e);
        if (e.response?.data?.message) {
            error.value = e.response.data.message
        } else {
            error.value = "Internal error"
        }
    }
}

async function removeComparison(id) {
    try {
        await deleteComparison(id)
        await fetchComparisons()
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