import http from '@/api'

export function getCustomOption(id) {
    return http.get(`/custom_options/${id}`)
}

export function getAllCustomOptions() {
    return http.get(`/custom_options`)
}

export function deleteCustomOption(id) {
    return http.delete(`/custom_options/${id}`)
}

export function createCustomOption(customOption) {
    return http.post('/custom_options', customOption)
}
