import http from '@/api'
import { APP_URL } from './url'

export const photoUrl = (id) => {
    return `${APP_URL}/objects/${id}/photo`
}

export function getAllObjects(comparison_id) {
    return http.get(`/objects?comparison_id=${comparison_id}`)
}

export function updateObject(object) {
    return http.put(`/objects/${object.id}`, object)
}

export function deleteObject(id) {
    return http.delete(`/objects/${id}`)
}

export function createObject(object) {
    return http.post('/objects', object)
}

export function uploadPhoto(id, photo) {
    let data = new FormData();
    data.append('photo', photo, photo.name);
    return http.post(`/objects/${id}/photo`, data, {
        headers: {
            'Content-Type': `multipart/form-data;boundary=${data._boundary}`
        }
    })
}