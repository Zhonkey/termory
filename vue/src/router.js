import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from './stores/auth'
import Login from './pages/login/Login.vue'
import Register from './pages/login/Register.vue'
import UserList from './pages/user/List.vue'
import UserView from './pages/user/View.vue'
import UserEdit from './pages/user/Edit.vue'

const routes = [
    { path: '/', name: 'Home', redirect: '/login' },
    { path: '/login', name: 'Login', component: Login },
    { path: '/register', name: 'Register', component: Register },
    { path: '/users', name: 'UserList', component: UserList, meta: { requiresAuth: true } },
    { path: '/users/:id', name: 'UserView', component: UserView , meta: { requiresAuth: true } },
    { path: '/users/:id/edit', name: 'UserEdit', component: UserEdit , meta: { requiresAuth: true } }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach((to, from, next) => {
    const auth = useAuthStore()

    if (to.meta.requiresAuth && !auth.accessToken) {
        next({ name: 'Login' })
    } else if (to.name === 'Login' && auth.accessToken) {
        next({ name: 'UserList' })
    } else {
        next()
    }
})

export default router
