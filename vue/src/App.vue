<template>
  <div class="container mt-4">
    <h2 class="mb-3">🧠 {{ t('app.title') }}</h2>
    <nav class="mb-4">
      <router-link v-if="auth.user" to="/users">{{ t('app.users') }}</router-link>
    </nav>

    <router-view />

    <div class="mt-3">
      <button @click="setLang('ru')">🇷🇺 RU</button>
      <button @click="setLang('en')">🇬🇧 EN</button>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from './stores/auth'

import { useI18n } from 'vue-i18n'
const { t, locale } = useI18n()

const auth = useAuthStore()

function setLang(lang) {
  locale.value = lang
}
</script>

<style>
body {
  font-family: sans-serif;
  background: #f7f7f7;
}
.container {
  max-width: 800px;
  margin: 0 auto;
  background: white;
  padding: 20px;
  border-radius: 10px;
}
</style>