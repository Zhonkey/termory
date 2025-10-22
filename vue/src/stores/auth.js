import { defineStore } from 'pinia'
import { AuthService } from '../services/auth'

export const useAuthStore = defineStore('auth', {
    state: () => ({
        accessToken: localStorage.getItem('access_token') || null,
        refreshToken: localStorage.getItem('refresh_token') || null,
        user: JSON.parse(localStorage.getItem('user') || 'null'),
        refreshing: false
    }),

    actions: {
        async login(email, password) {
            const data = AuthService.login(email, password)

            if (!res.ok) throw new Error(data.message || 'Ошибка авторизации')

            this.accessToken = data.access_token
            this.refreshToken = data.refresh_token
            this.user = data.user

            localStorage.setItem('access_token', data.access_token)
            localStorage.setItem('refresh_token', data.refresh_token)
            localStorage.setItem('user', JSON.stringify(data.user))
        },
        async register(email, password) {
            const data = AuthService.register(email, password)

            if (!res.ok) throw new Error(data.message || 'Ошибка авторизации')

            this.accessToken = data.access_token
            this.refreshToken = data.refresh_token
            this.user = data.user

            localStorage.setItem('access_token', data.access_token)
            localStorage.setItem('refresh_token', data.refresh_token)
            localStorage.setItem('user', JSON.stringify(data.user))
        },
        logout() {
            this.accessToken = null
            this.refreshToken = null
            this.user = null
            localStorage.clear()
            window.location.href = '/login'
        },
        async refresh() {
            if (this.refreshing || !this.refreshToken) return
            this.refreshing = true

            try {
                const data = AuthService.refresh(this.refreshToken)

                if (res.ok && data.access_token) {
                    this.accessToken = data.access_token
                    localStorage.setItem('access_token', data.access_token)
                } else {
                    this.logout()
                }
            } catch (e) {
                this.logout()
            } finally {
                this.refreshing = false
            }
        }
    }
})