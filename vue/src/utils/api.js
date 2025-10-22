import { useAuthStore } from '../stores/auth'

const API_BASE = 'http://localhost:8080/api'

export async function apiRequest(endpoint, options = {}) {
    const auth = useAuthStore()
    const headers = {
        'Content-Type': 'application/json',
        ...(options.headers || {})
    }

    if (auth.accessToken) {
        headers['Authorization'] = `Bearer ${auth.accessToken}`
    }

    const res = await fetch(`${API_BASE}${endpoint}`, {
        ...options,
        headers
    })

    // Если токен просрочен, пробуем обновить
    if (res.status === 401 && auth.refreshToken) {
        await auth.refresh()

        // Повторяем запрос
        headers['Authorization'] = `Bearer ${auth.accessToken}`
        return fetch(`${API_BASE}${endpoint}`, { ...options, headers })
    }

    return res
}

// Хелперы для удобства
export const api = {
    get: (url) => apiRequest(url),
    post: (url, data) => apiRequest(url, { method: 'POST', body: JSON.stringify(data) }),
    put: (url, data) => apiRequest(url, { method: 'PUT', body: JSON.stringify(data) }),
    delete: (url) => apiRequest(url, { method: 'DELETE' })
}