<template>
  <div class="notify">
    <div class="head">
      <div class="title">通知公告</div>
      <div class="tools">
        <button class="ghost" type="button" @click="refresh">刷新</button>
      </div>
    </div>

    <div class="grid">
      <div class="panel">
        <div class="panel-head">
          <div class="panel-title">通知中心</div>
          <div class="actions">
            <select v-model="filter" class="select">
              <option value="all">全部</option>
              <option value="unread">未读</option>
              <option value="read">已读</option>
            </select>
            <button class="ghost" type="button" :disabled="loading" @click="readAll">全部已读</button>
          </div>
        </div>

        <div v-if="loading" class="muted pad">加载中...</div>
        <div v-else-if="list.length === 0" class="muted pad">暂无通知</div>
        <div v-else class="rows">
          <div v-for="n in list" :key="n.notificationId" class="row">
            <div class="row-main">
              <div class="row-title">
                <span class="dot" :class="{ read: n.isRead }" />
                <span>{{ n.title || `通知 #${n.notificationId}` }}</span>
              </div>
              <div v-if="n.content" class="row-content">{{ n.content }}</div>
              <div v-if="n.createdAt" class="row-sub">{{ formatTime(n.createdAt) }}</div>
            </div>
            <div class="row-actions">
              <button class="btn" type="button" :disabled="n.isRead || actingId === n.notificationId" @click="markRead(n.notificationId)">
                标记已读
              </button>
            </div>
          </div>
        </div>

        <div v-if="error" class="error">{{ error }}</div>
      </div>

      <div class="panel">
        <div class="panel-head">
          <div class="panel-title">发布公告（发帖闭环）</div>
        </div>
        <div class="form">
          <label class="field">
            <span>选择贴吧</span>
            <select v-model.number="forumId" class="select">
              <option :value="0" disabled>请选择</option>
              <option v-for="f in forums" :key="f.forumId" :value="f.forumId">
                {{ f.name }}
              </option>
            </select>
          </label>
          <label class="field">
            <span>标题</span>
            <input v-model.trim="titleInput" placeholder="公告标题" />
          </label>
          <label class="field">
            <span>正文</span>
            <textarea v-model.trim="contentInput" placeholder="公告正文" />
          </label>
          <div class="actions">
            <button class="primary" type="button" :disabled="publishing" @click="publish">
              {{ publishing ? "发布中..." : "发布（进入待审）" }}
            </button>
            <button class="ghost" type="button" :disabled="publishing" @click="clearPublish">清空</button>
          </div>
          <p v-if="publishMsg" class="success">{{ publishMsg }}</p>
          <p v-if="publishErr" class="error">{{ publishErr }}</p>
          <button class="link" type="button" @click="router.push({ name: 'admin-posts-pending' })">去“帖子审核”里通过公告</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import type { Forum, NotificationItem } from '@/types/api'
import { feedApi, forumsApi, notificationsApi } from '@/apis'

const router = useRouter()

const filter = ref<'all' | 'unread' | 'read'>('all')
const notifications = ref<NotificationItem[]>([])
const loading = ref(false)
const actingId = ref<number | null>(null)
const error = ref('')

const forums = ref<Forum[]>([])
const forumId = ref(0)
const titleInput = ref('')
const contentInput = ref('')
const publishing = ref(false)
const publishMsg = ref('')
const publishErr = ref('')

const list = computed(() => {
  if (filter.value === 'unread') return notifications.value.filter((n) => !n.isRead)
  if (filter.value === 'read') return notifications.value.filter((n) => n.isRead)
  return notifications.value
})

const formatTime = (value: string) => {
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return value
  const yy = d.getFullYear()
  const mm = String(d.getMonth() + 1).padStart(2, '0')
  const dd = String(d.getDate()).padStart(2, '0')
  const hh = String(d.getHours()).padStart(2, '0')
  const mi = String(d.getMinutes()).padStart(2, '0')
  return `${yy}-${mm}-${dd} ${hh}:${mi}`
}

const loadForums = async () => {
  const data = await forumsApi.getForums({ pageNum: 1, pageSize: 200 })
  forums.value = data.list || []
  if (forums.value.length > 0 && forumId.value === 0) forumId.value = forums.value[0].forumId
}

const loadNotifications = async () => {
  error.value = ''
  loading.value = true
  try {
    const data = await notificationsApi.getNotifications({ pageNum: 1, pageSize: 50 })
    notifications.value = data.list || []
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : '加载失败'
  } finally {
    loading.value = false
  }
}

const refresh = async () => {
  await Promise.all([loadForums(), loadNotifications()])
}

const markRead = async (notificationId: number) => {
  actingId.value = notificationId
  error.value = ''
  try {
    await notificationsApi.readNotification(notificationId)
    await loadNotifications()
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : '操作失败'
  } finally {
    actingId.value = null
  }
}

const readAll = async () => {
  const ok = window.confirm('确认将全部通知标记为已读？')
  if (!ok) return
  error.value = ''
  loading.value = true
  try {
    await notificationsApi.readAllNotifications()
    await loadNotifications()
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : '操作失败'
  } finally {
    loading.value = false
  }
}

const validatePublish = () => {
  const t = titleInput.value.trim()
  const c = contentInput.value.trim()
  if (!forumId.value) return '请选择贴吧'
  if (!t) return '请输入标题'
  if (!c) return '请输入正文'
  return ''
}

const publish = async () => {
  publishMsg.value = ''
  publishErr.value = ''
  const msg = validatePublish()
  if (msg) {
    publishErr.value = msg
    return
  }
  publishing.value = true
  try {
    const data = await feedApi.createPost({
      forumId: forumId.value,
      title: `【公告】${titleInput.value.trim()}`,
      content: contentInput.value.trim()
    })
    publishMsg.value = `已发布：#${data.postId}（状态：${data.status}），请到“帖子审核”里通过`
    clearPublish()
  } catch (e: unknown) {
    publishErr.value = e instanceof Error ? e.message : '发布失败'
  } finally {
    publishing.value = false
  }
}

const clearPublish = () => {
  titleInput.value = ''
  contentInput.value = ''
}

onMounted(refresh)
</script>

<style scoped>
.notify {
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

.tools {
  display: inline-flex;
  gap: 10px;
  align-items: center;
}

.grid {
  display: grid;
  grid-template-columns: 1fr 420px;
  gap: 16px;
  align-items: start;
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
  gap: 12px;
}

.panel-title {
  font-weight: 900;
  color: #0f172a;
}

.actions {
  display: inline-flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.select {
  height: 34px;
  padding: 0 12px;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  background: #fff;
  color: #334155;
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

.rows {
  padding: 10px 14px 14px;
  display: grid;
  gap: 10px;
}

.row {
  border: 1px solid #eef0f5;
  border-radius: 14px;
  padding: 12px;
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 12px;
}

.row-title {
  font-weight: 900;
  color: #0f172a;
  display: inline-flex;
  gap: 10px;
  align-items: center;
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 99px;
  background: #4f46e5;
}

.dot.read {
  background: #cbd5e1;
}

.row-content {
  margin-top: 8px;
  color: #334155;
  white-space: pre-wrap;
  line-height: 1.7;
}

.row-sub {
  margin-top: 8px;
  font-size: 12px;
  color: #94a3b8;
}

.row-actions {
  display: grid;
  justify-items: end;
  align-content: start;
}

.btn {
  height: 30px;
  padding: 0 10px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 12px;
  cursor: pointer;
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

textarea {
  min-height: 140px;
  padding: 10px 12px;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  outline: none;
  resize: vertical;
  font-family: inherit;
  line-height: 1.6;
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

.muted {
  color: #94a3b8;
}

.pad {
  padding: 14px;
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
  text-align: left;
}

@media (max-width: 1100px) {
  .grid {
    grid-template-columns: 1fr;
  }
}
</style>

