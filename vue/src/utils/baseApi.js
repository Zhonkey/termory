const API_BASE = import.meta.env.VITE_API_BASE

export async function simpleRequest(endpoint, options = {}) {
    const headers = {
        'Content-Type': 'application/json',
        ...(options.headers || {})
    }

    return await fetch(`${API_BASE}${endpoint}`, {
        ...options,
        headers
    })
}

async function baseRequest(endpoint, options = {}) {
    const response = await simpleRequest(endpoint, options)

    if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`HTTP ${response.status}: ${errorText}`)
    }

    return response.status !== 204 ? await response.json() : null
}

export const baseApi = {
    get: (url) => baseRequest(url),
    post: (url, data) => baseRequest(url, { method: 'POST', body: JSON.stringify(data) }),
    put: (url, data) => baseRequest(url, { method: 'PUT', body: JSON.stringify(data) }),
    delete: (url) => baseRequest(url, { method: 'DELETE' })
}