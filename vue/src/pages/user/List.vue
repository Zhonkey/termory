<template>
  <div>
    <h3>Список пользователей</h3>
    <ul>
      <li v-for="user in users" :key="user.id">
        <router-link :to="'/users/' + user.id">{{ user.email }}</router-link>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { UserService } from '../../services/users'

const users = ref([])

onMounted(async () => {
  const res = await UserService.list()
  users.value = res.users
})
</script>