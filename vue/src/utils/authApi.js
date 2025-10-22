import { baseRequest } from './baseApi'
import { useAuthStore } from '../stores/auth'

export async function authRequest(endpoint, options = {}) {
    const auth = useAuthStore()

    const headers = {
        'Content-Type': 'application/json',
        ...(options.headers || {})
    }

    if (auth.accessToken) {
        headers['Authorization'] = `Bearer ${auth.accessToken}`
    }

    let res = await baseRequest(endpoint, { ...options, headers })

    if (res.status === 401 && auth.refreshToken) {
        await auth.refresh()

        headers['Authorization'] = `Bearer ${auth.accessToken}`
        res = await baseRequest(endpoint, { ...options, headers })
    }

    return res
}

export const authApi = {
    get: (url) => authRequest(url),
    post: (url, data) => authRequest(url, { method: 'POST', body: JSON.stringify(data) }),
    put: (url, data) => authRequest(url, { method: 'PUT', body: JSON.stringify(data) }),
    delete: (url) => authRequest(url, { method: 'DELETE' })
}