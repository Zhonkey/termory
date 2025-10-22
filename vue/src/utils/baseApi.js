const API_BASE = 'http://localhost:8080/api'

export async function baseRequest(endpoint, options = {}) {
    const headers = {
        'Content-Type': 'application/json',
        ...(options.headers || {})
    }

    return fetch(`${API_BASE}${endpoint}`, {
        ...options,
        headers
    })
}

export const baseApi = {
    get: (url) => baseRequest(url),
    post: (url, data) => baseRequest(url, { method: 'POST', body: JSON.stringify(data) }),
    put: (url, data) => baseRequest(url, { method: 'PUT', body: JSON.stringify(data) }),
    delete: (url) => baseRequest(url, { method: 'DELETE' })
}