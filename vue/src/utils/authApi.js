import { simpleRequest} from './baseApi'
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

    let response = await simpleRequest(endpoint, { ...options, headers })

    if (response.status === 401 && auth.refreshToken) {
        await auth.refresh()

        headers['Authorization'] = `Bearer ${auth.accessToken}`
        response = await simpleRequest(endpoint, { ...options, headers })
    }

    if (!response.ok) {
        throw new Error(`HTTP ${response.status}`)
    }

    return response.status !== 204 ? await response.json() : null
}

export const authApi = {
    get: (url) => authRequest(url),
    post: (url, data) => authRequest(url, { method: 'POST', body: JSON.stringify(data) }),
    put: (url, data) => authRequest(url, { method: 'PUT', body: JSON.stringify(data) }),
    delete: (url) => authRequest(url, { method: 'DELETE' })
}