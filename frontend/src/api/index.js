import axios from 'axios'
import { APP_URL } from './url'

const instance = axios.create({
    baseURL: APP_URL
})

export default instance