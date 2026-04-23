<template>
  <div class="page">
    <h1>登录</h1>
    <form class="form" @submit.prevent="onSubmit">
      <label class="field">
        <span>账号</span>
        <input v-model.trim="account" placeholder="请输入账号" />
      </label>
      <label class="field">
        <span>密码</span>
        <input v-model.trim="password" type="password" placeholder="请输入密码" />
      </label>
      <button class="btn" type="submit" :disabled="loading">{{ loading ? '登录中...' : '登录' }}</button>
      <p v-if="error" class="error">{{ error }}</p>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const account = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

const onSubmit = async () => {
  error.value = ''
  loading.value = true
  try {
    await auth.login(account.value, password.value)
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/m/home'
    await router.replace(redirect)
  } catch (e: any) {
    error.value = e?.message || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.page {
  padding: 16px;
}

.form {
  display: grid;
  gap: 12px;
  max-width: 360px;
}

.field {
  display: grid;
  gap: 6px;
}

input {
  height: 36px;
  padding: 0 12px;
  border: 1px solid #d9d9d9;
  border-radius: 8px;
}

.btn {
  height: 38px;
  border: 0;
  background: #4f46e5;
  color: #fff;
  border-radius: 10px;
}

.error {
  color: #dc2626;
  margin: 0;
}
</style>

