<template>
  <div>
    <h3>Редактирование пользователя</h3>
    <p>ID: {{ id }}</p>
    <input v-model="name" class="form-control" />
    <button class="btn btn-success mt-2" @click="save">Сохранить</button>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const id = computed(() => route.params.id)
const name = ref('')

function save() {
  alert(`Пользователь ${id.value} сохранён как ${name.value}`)
  router.push(`/users/${id.value}`)
}
</script>