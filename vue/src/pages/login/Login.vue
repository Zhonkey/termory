<template>
  <div>
    <h3>{{ t('app.auth') }}</h3>

    <div v-if="error" class="alert alert-danger">{{ error }}</div>

    <form @submit.prevent="handleLogin">
      <div class="mb-3">
        <label>{{ t('app.email') }}:</label>
        <input v-model="email" class="form-control" type="email" required />
      </div>
      <div class="mb-3">
        <label>{{ t('app.password') }}:</label>
        <input v-model="password" class="form-control" type="password" required />
      </div>
      <button class="btn btn-primary">{{ t('app.login') }}</button> <router-link to="/register" class="me-2">{{ t('app.register') }}</router-link>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
const { t, locale } = useI18n()

import { useAuthStore } from '../../stores/auth'

const router = useRouter()
const auth = useAuthStore()
const email = ref('')
const password = ref('')
const error = ref(null)

async function handleLogin() {
  try {
    await auth.login(email.value, password.value)
    router.push('/users')
  } catch (e) {
    error.value = e.message
  }
}
</script>