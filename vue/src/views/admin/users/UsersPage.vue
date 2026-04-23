<template>
  <div class="users">
    <div class="head">
      <div class="title">用户管理</div>
      <button class="ghost" type="button" @click="reset">清空</button>
    </div>

    <div class="panel">
      <div class="panel-head">
        <div class="panel-title">按用户 ID 查询</div>
      </div>
      <div class="form">
        <label class="field">
          <span>用户 ID</span>
          <input v-model.trim="userIdInput" placeholder="请输入 userId" />
        </label>
        <div class="actions">
          <button class="primary" type="button" :disabled="loading || !userIdNumber" @click="fetchUser">
            {{ loading ? "查询中..." : "查询" }}
          </button>
        </div>
        <p v-if="error" class="error">{{ error }}</p>
      </div>
    </div>

    <div v-if="user" class="grid">
      <div class="panel">
        <div class="panel-head">
          <div class="panel-title">用户信息</div>
          <button class="link" type="button" @click="fetchUser">刷新</button>
        </div>
        <div class="user-card">
          <div class="avatar">
            <img v-if="user.avatarUrl" class="avatar-img" :src="user.avatarUrl" alt="avatar" />
            <div v-else class="avatar-fallback">{{ (user.nickname || `U${user.userId}`).slice(0, 1) }}</div>
          </div>
          <div class="meta">
            <div class="name">{{ user.nickname || `用户${user.userId}` }}</div>
            <div class="sub">
              <span>ID {{ user.userId }}</span>
              <span v-if="user.role" class="dot">·</span>
              <span v-if="user.role">{{ user.role }}</span>
              <span v-if="user.status" class="dot">·</span>
              <span v-if="user.status">{{ user.status }}</span>
            </div>
            <div v-if="user.bio" class="bio">{{ user.bio }}</div>
            <div v-if="user.banUntil" class="ban">banUntil: {{ user.banUntil }}</div>
          </div>
        </div>
      </div>

      <div class="panel">
        <div class="panel-head">
          <div class="panel-title">封禁操作</div>
        </div>

        <div class="form">
          <label class="field">
            <span>封禁时长（秒）</span>
            <input v-model.trim="durationInput" placeholder="例如：86400（1 天）" />
          </label>
          <label class="field">
            <span>备注（可选）</span>
            <input v-model.trim="remark" placeholder="处理备注" />
          </label>
          <div class="actions">
            <button class="danger" type="button" :disabled="acting" @click="ban">
              {{ acting ? "处理中..." : "封禁" }}
            </button>
            <button class="ghost" type="button" :disabled="acting" @click="unban">
              解除封禁
            </button>
          </div>
          <p v-if="actionMsg" class="success">{{ actionMsg }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { adminApi, usersApi } from '@/apis'
import type { UserPublicInfo } from '@/apis/modules/users'

const userIdInput = ref('')
const loading = ref(false)
const error = ref('')
const user = ref<UserPublicInfo | null>(null)

const durationInput = ref('86400')
const remark = ref('')
const acting = ref(false)
const actionMsg = ref('')

const userIdNumber = computed(() => {
  const n = Number(userIdInput.value)
  return Number.isFinite(n) && n > 0 ? n : 0
})

const fetchUser = async () => {
  error.value = ''
  actionMsg.value = ''
  if (!userIdNumber.value) {
    error.value = '请输入正确的 userId'
    return
  }
  loading.value = true
  try {
    const data = await usersApi.getUserPublic(userIdNumber.value)
    user.value = data
  } catch (e: unknown) {
    user.value = null
    error.value = e instanceof Error ? e.message : '查询失败'
  } finally {
    loading.value = false
  }
}

const parseDurationSeconds = () => {
  const n = Number(durationInput.value)
  return Number.isFinite(n) && n > 0 ? Math.floor(n) : 0
}

const ban = async () => {
  if (!user.value) return
  const seconds = parseDurationSeconds()
  if (!seconds) {
    error.value = '请输入正确的封禁时长（秒）'
    return
  }
  const ok = window.confirm(`确认封禁用户 #${user.value.userId}，时长 ${seconds} 秒？`)
  if (!ok) return

  acting.value = true
  error.value = ''
  actionMsg.value = ''
  try {
    const data = await adminApi.banUser(user.value.userId, { durationSeconds: seconds, remark: remark.value || undefined })
    actionMsg.value = `已封禁：${data.status}${data.banUntil ? `，banUntil=${data.banUntil}` : ''}`
    await fetchUser()
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : '封禁失败'
  } finally {
    acting.value = false
  }
}

const unban = async () => {
  if (!user.value) return
  const ok = window.confirm(`确认解除封禁用户 #${user.value.userId}？`)
  if (!ok) return

  acting.value = true
  error.value = ''
  actionMsg.value = ''
  try {
    const banUntil = new Date().toISOString()
    const data = await adminApi.banUser(user.value.userId, { banUntil, remark: remark.value || 'unban' })
    actionMsg.value = `已提交解除封禁：${data.status}`
    await fetchUser()
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : '解除失败'
  } finally {
    acting.value = false
  }
}

const reset = () => {
  userIdInput.value = ''
  user.value = null
  error.value = ''
  actionMsg.value = ''
  durationInput.value = '86400'
  remark.value = ''
}
</script>

<style scoped>
.users {
  display: grid;
  gap: 12px;
}

.head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.title {
  font-weight: 900;
  color: #0f172a;
}

.panel {
  background: #fff;
  border: 1px solid #eef0f5;
  border-radius: 16px;
  box-shadow: 0 6px 18px rgba(15, 23, 42, 0.03);
}

.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 14px;
  border-bottom: 1px solid #eef0f5;
}

.panel-title {
  font-weight: 900;
  color: #0f172a;
}

.form {
  padding: 14px;
  display: grid;
  gap: 12px;
}

.field {
  display: grid;
  gap: 6px;
}

.field span {
  font-size: 12px;
  color: #64748b;
}

input {
  height: 36px;
  padding: 0 12px;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  outline: none;
}

.actions {
  display: inline-flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}

.primary {
  height: 34px;
  padding: 0 12px;
  border: 0;
  background: #4f46e5;
  color: #fff;
  border-radius: 12px;
  cursor: pointer;
}

.ghost {
  height: 34px;
  padding: 0 12px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 12px;
  cursor: pointer;
  color: #334155;
}

.danger {
  height: 34px;
  padding: 0 12px;
  border: 1px solid #fecaca;
  background: #fff;
  color: #dc2626;
  border-radius: 12px;
  cursor: pointer;
}

.grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.user-card {
  padding: 14px;
  display: grid;
  grid-template-columns: 54px 1fr;
  gap: 12px;
  align-items: center;
}

.avatar {
  width: 54px;
  height: 54px;
  border-radius: 999px;
  overflow: hidden;
  border: 1px solid #e2e8f0;
  background: #fff;
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.avatar-fallback {
  width: 100%;
  height: 100%;
  display: grid;
  place-items: center;
  background: #eef2ff;
  color: #4f46e5;
  font-weight: 900;
}

.meta {
  display: grid;
  gap: 4px;
}

.name {
  font-weight: 900;
  color: #0f172a;
}

.sub {
  font-size: 12px;
  color: #94a3b8;
  display: inline-flex;
  gap: 6px;
  align-items: center;
  flex-wrap: wrap;
}

.dot {
  color: #cbd5e1;
}

.bio {
  color: #334155;
  line-height: 1.7;
  white-space: pre-wrap;
}

.ban {
  font-size: 12px;
  color: #dc2626;
  font-weight: 800;
}

.error {
  color: #dc2626;
  margin: 0;
  font-weight: 800;
}

.success {
  color: #16a34a;
  margin: 0;
  font-weight: 800;
}

.link {
  border: 0;
  background: transparent;
  color: #4f46e5;
  cursor: pointer;
  padding: 0;
  font-weight: 800;
}

@media (max-width: 960px) {
  .grid {
    grid-template-columns: 1fr;
  }
}
</style>

