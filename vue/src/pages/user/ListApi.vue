<template>
  <div>
    <h3>Список пользователей</h3>
    <ul>
      <li v-for="user in users" :key="user.id">
        <router-link :to="'/users/' + user.id">{{ user.name }}</router-link>
      </li>
    </ul>
  </div>
</template>

<script setup>
const users = [
  { id: 1, name: 'Иван' },
  { id: 2, name: 'Анна' },
  { id: 3, name: 'Петр' }
]
</script>