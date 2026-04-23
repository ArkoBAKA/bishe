<template>
  <div class="pending">
    <div class="head">
      <div class="title">待审核帖子</div>
      <div class="tools">
        <button class="ghost" type="button" @click="refresh">刷新</button>
      </div>
    </div>

    <div class="panel">
      <div class="panel-head">
        <div class="muted">共 {{ total }} 条</div>
        <div class="pager">
          <button class="pager-btn" type="button" :disabled="pageNum <= 1 || loading" @click="prevPage">
            上一页
          </button>
          <div class="pager-text">第 {{ pageNum }} 页</div>
          <button class="pager-btn" type="button" :disabled="isLastPage || loading" @click="nextPage">
            下一页
          </button>
        </div>
      </div>

      <div v-if="loading" class="muted pad">加载中...</div>
      <div v-else-if="list.length === 0" class="muted pad">暂无待审帖子</div>
      <div v-else class="rows">
        <div v-for="p in list" :key="p.postId" class="row">
          <div class="row-main">
            <div class="row-title">{{ p.title }}</div>
            <div class="row-sub">
              <span>postId: {{ p.postId }}</span>
              <span class="dot">·</span>
              <span>forumId: {{ p.forumId }}</span>
              <span v-if="p.createdAt" class="dot">·</span>
              <span v-if="p.createdAt">{{ formatTime(p.createdAt) }}</span>
            </div>
            <div v-if="p.content" class="row-content">{{ p.content }}</div>
          </div>
          <div class="row-actions">
            <button class="btn ok" type="button" :disabled="actingId === p.postId" @click="review(p.postId, 'approve')">
              通过
            </button>
            <button class="btn" type="button" :disabled="actingId === p.postId" @click="review(p.postId, 'reject')">
              驳回
            </button>
            <button class="btn warn" type="button" :disabled="actingId === p.postId" @click="review(p.postId, 'hide')">
              隐藏
            </button>
            <button class="btn danger" type="button" :disabled="actingId === p.postId" @click="remove(p.postId)">
              删除
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="error" class="error">{{ error }}</div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { adminApi } from '@/apis'
import type { Post } from '@/types/api'

const pageNum = ref(1)
const pageSize = ref(10)
const total = ref(0)
const list = ref<Post[]>([])
const loading = ref(false)
const actingId = ref<number | null>(null)
const error = ref('')

const isLastPage = computed(() => pageNum.value * pageSize.value >= total.value)

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

const fetchList = async () => {
  error.value = ''
  loading.value = true
  try {
    const data = await adminApi.getPendingPosts({ pageNum: pageNum.value, pageSize: pageSize.value })
    list.value = data.list || []
    total.value = data.total || 0
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : '加载失败'
  } finally {
    loading.value = false
  }
}

const refresh = async () => {
  await fetchList()
}

const prevPage = async () => {
  if (pageNum.value <= 1) return
  pageNum.value -= 1
  await fetchList()
}

const nextPage = async () => {
  if (isLastPage.value) return
  pageNum.value += 1
  await fetchList()
}

const review = async (postId: number, action: 'approve' | 'reject' | 'hide') => {
  const remark = window.prompt('处理备注（可选）') || undefined
  actingId.value = postId
  error.value = ''
  try {
    await adminApi.reviewPost(postId, { action, reviewRemark: remark })
    await fetchList()
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : '操作失败'
  } finally {
    actingId.value = null
  }
}

const remove = async (postId: number) => {
  const ok = window.confirm(`确认删除帖子 #${postId} ?`)
  if (!ok) return
  actingId.value = postId
  error.value = ''
  try {
    await adminApi.deletePost(postId)
    await fetchList()
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : '删除失败'
  } finally {
    actingId.value = null
  }
}

onMounted(fetchList)
</script>

<style scoped>
.pending {
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

.ghost {
  height: 34px;
  padding: 0 12px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 12px;
  cursor: pointer;
  color: #334155;
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

.pager {
  display: inline-flex;
  gap: 10px;
  align-items: center;
}

.pager-btn {
  height: 30px;
  padding: 0 10px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 12px;
  cursor: pointer;
}

.pager-text {
  font-size: 12px;
  color: #64748b;
  font-weight: 800;
}

.rows {
  padding: 10px 14px 14px;
  display: grid;
  gap: 10px;
}

.row {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 12px;
  align-items: start;
  border: 1px solid #eef0f5;
  border-radius: 14px;
  padding: 12px;
}

.row-title {
  font-weight: 900;
  color: #0f172a;
}

.row-sub {
  margin-top: 6px;
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

.row-content {
  margin-top: 10px;
  color: #334155;
  line-height: 1.7;
  white-space: pre-wrap;
  background: #f8fafc;
  border: 1px solid #eef0f5;
  border-radius: 12px;
  padding: 10px 12px;
  max-height: 180px;
  overflow: auto;
}

.row-actions {
  display: grid;
  gap: 8px;
  justify-items: end;
}

.btn {
  height: 30px;
  padding: 0 10px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 12px;
  cursor: pointer;
  min-width: 64px;
}

.btn.ok {
  border-color: rgba(22, 163, 74, 0.25);
  color: #16a34a;
}

.btn.warn {
  border-color: rgba(234, 88, 12, 0.25);
  color: #ea580c;
}

.btn.danger {
  border-color: rgba(220, 38, 38, 0.25);
  color: #dc2626;
}

.muted {
  color: #94a3b8;
}

.pad {
  padding: 14px;
}

.error {
  color: #dc2626;
  font-weight: 800;
}

@media (max-width: 920px) {
  .row {
    grid-template-columns: 1fr;
  }
  .row-actions {
    grid-auto-flow: column;
    justify-content: flex-start;
    justify-items: start;
  }
}
</style>
