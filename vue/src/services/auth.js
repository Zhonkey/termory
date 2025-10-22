import { baseApi } from '../utils/baseApi'

export const AuthService = {
    async login(email, password) {
        const { data } = await baseApi.post('/auth/access_token', { email: email, password: password })
        return data
    },

    async register(email, password) {
        const { data } = await baseApi.put('/users', { email: email, password: password })
        return data
    },

    async refresh(refreshToken) {
        const { data } = await baseApi.post('/auth/refresh_token', { refresh_token: refreshToken })
        return data
    }
}