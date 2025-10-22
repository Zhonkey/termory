<template>
  <div>
    <h3>{{ t('app.register') }}</h3>

    <div v-if="error" class="alert alert-danger">{{ error }}</div>

    <form @submit.prevent="register">
      <div class="mb-3">
        <label>{{ t('app.email') }}:</label>
        <input v-model="email" class="form-control" type="email" required />
      </div>
      <div class="mb-3">
        <label>{{ t('app.password') }}:</label>
        <input v-model="password" class="form-control" type="password" required />
      </div>
      <div class="mb-3">
        <label>{{ t('app.repeat_password') }}:</label>
        <input v-model="repeat_password" class="form-control" type="password" required />
      </div>
      <button class="btn btn-primary">{{ t('app.register') }}</button> <router-link to="/login" class="me-2">{{ t('app.login') }}</router-link>

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
const repeat_password = ref('')
const error = ref(null)

async function handleRegister() {
  try {
    if(password.value !== repeat_password.value) {
      error.value = t('error.passwords_not_same')
    }
    await auth.register(email.value, password.value)
    router.push('/users')
  } catch (e) {
    error.value = e.message
  }
}
</script>