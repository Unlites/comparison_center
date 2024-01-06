import http from '@/api'

export function getAllComparisons() {
    return http.get('/comparisons')
}

export function getComparison(id) {
    return http.get(`/comparisons/${id}`)
}

export function createComparison(comparison) {
    return http.post('/comparisons', comparison)
}

export function updateComparison(comparison) {
    return http.put(`/comparisons/${comparison.id}`, comparison)
}

export function deleteComparison(id) {
    return http.delete(`/comparisons/${id}`)
}