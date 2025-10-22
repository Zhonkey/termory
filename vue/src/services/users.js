import { authApi } from '../utils/authApi'

export const UserService = {
    list: () => authApi.get('/users'),
    view: (id) => authApi.get(`/users/${id}`),
    update: (id, data) => authApi.put(`/users/${id}`, data),
    delete: (id) => authApi.delete(`/users/${id}`)
}