import { baseApi } from '../utils/baseApi'

export const AuthService = {
    async login(email, password) {
        return await baseApi.post('/auth/access_token', { email: email, password: password })
    },

    async register(email, password) {
        return await baseApi.put('/users', { email: email, password: password })
    },

    async refresh(refreshToken) {
        return await baseApi.post('/auth/refresh_token', { refresh_token: refreshToken })
    }
}