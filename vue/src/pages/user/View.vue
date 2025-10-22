<template>
  <div>
    <h3>Просмотр пользователя</h3>
    <p>ID: {{ id }}</p>
    <p>Имя: {{ user?.name || 'Не найден' }}</p>
    <router-link :to="`/users/${id}/edit`">Редактировать</router-link>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const id = computed(() => route.params.id)
const users = [
  { id: '1', name: 'Иван' },
  { id: '2', name: 'Анна' },
  { id: '3', name: 'Петр' }
]
const user = users.find(u => u.id === id.value)
</script>