import { baseApi } from '../utils/baseApi'

export const AuthService = {
    async login(email, password) {
        const { data } = await baseApi.post('/login', { email: email, password: password })
        return data
    },

    async register(email, password) {
        const { data } = await baseApi.post('/register', { email: email, password: password })
        return data
    },

    async refresh(refreshToken) {
        const { data } = await baseApi.post('/refresh', { refresh_token: refreshToken })
        return data
    }
}