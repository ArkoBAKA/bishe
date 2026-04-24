<template>
  <div class="messages-page">
    <header class="top">
      <button class="back" type="button" @click="router.back()">返回</button>
      <div class="title">消息</div>
      <button class="ghost" type="button" @click="router.replace('/m/profile')">
        我的
      </button>
    </header>

    <main class="main">
      <div class="card">
        <div class="head">
          <div class="tabs">
            <button
              class="tab"
              :class="{ active: tab === 'unread' }"
              type="button"
              @click="tab = 'unread'"
            >
              未读
            </button>
            <button
              class="tab"
              :class="{ active: tab === 'all' }"
              type="button"
              @click="tab = 'all'"
            >
              全部
            </button>
          </div>
          <div class="tools">
            <button
              class="link"
              type="button"
              :disabled="loading"
              @click="reload"
            >
              刷新
            </button>
            <button
              class="link"
              type="button"
              :disabled="loading || list.length === 0"
              @click="markAllRead"
            >
              全部已读
            </button>
          </div>
        </div>

        <div v-if="loading" class="muted pad">加载中...</div>
        <div v-else-if="list.length === 0" class="muted pad">暂无消息</div>
        <div v-else class="list">
          <button
            v-for="n in list"
            :key="n.notificationId"
            class="item"
            type="button"
            @click="onClickItem(n)"
          >
            <div class="row">
              <div class="left">
                <span class="dot" :class="{ read: n.isRead }" />
                <div class="t">
                  {{ n.title || `通知 #${n.notificationId}` }}
                </div>
              </div>
              <div v-if="n.createdAt" class="time">
                {{ formatTime(n.createdAt) }}
              </div>
            </div>
            <div v-if="expandedId === n.notificationId" class="content">
              {{ n.content || "（无内容）" }}
            </div>
          </button>
        </div>

        <div class="foot">
          <button
            class="ghost full"
            type="button"
            :disabled="loadingMore || !hasMore"
            @click="loadMore"
          >
            {{
              hasMore ? (loadingMore ? "加载中..." : "加载更多") : "没有更多了"
            }}
          </button>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import { useRouter } from "vue-router";
import type { NotificationItem } from "@/types/api";
import { notificationsApi } from "@/apis";

const router = useRouter();

const tab = ref<"unread" | "all">("unread");
const list = ref<NotificationItem[]>([]);
const expandedId = ref<number | null>(null);

const pageNum = ref(1);
const pageSize = 20;
const hasMore = ref(true);
const loading = ref(false);
const loadingMore = ref(false);

const formatTime = (value: string) => {
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return value;
  const yy = d.getFullYear();
  const mm = String(d.getMonth() + 1).padStart(2, "0");
  const dd = String(d.getDate()).padStart(2, "0");
  const hh = String(d.getHours()).padStart(2, "0");
  const mi = String(d.getMinutes()).padStart(2, "0");
  return `${yy}-${mm}-${dd} ${hh}:${mi}`;
};

const fetchPage = async (nextPageNum: number) => {
  const data = await notificationsApi.getNotifications({
    pageNum: nextPageNum,
    pageSize,
    isRead: tab.value === "unread" ? false : undefined,
  });
  return data;
};

const reload = async () => {
  loading.value = true;
  expandedId.value = null;
  try {
    pageNum.value = 1;
    const data = await fetchPage(1);
    list.value = data.list || [];
    hasMore.value = (data.list || []).length >= pageSize;
  } finally {
    loading.value = false;
  }
};

const loadMore = async () => {
  if (!hasMore.value || loadingMore.value) return;
  loadingMore.value = true;
  try {
    const next = pageNum.value + 1;
    const data = await fetchPage(next);
    const nextList = data.list || [];
    list.value = [...list.value, ...nextList];
    pageNum.value = next;
    hasMore.value = nextList.length >= pageSize;
  } finally {
    loadingMore.value = false;
  }
};

const onClickItem = async (n: NotificationItem) => {
  expandedId.value =
    expandedId.value === n.notificationId ? null : n.notificationId;
  if (n.isRead) return;
  try {
    await notificationsApi.readNotification(n.notificationId);
    n.isRead = true;
    if (tab.value === "unread") {
      list.value = list.value.filter(
        (x) => x.notificationId !== n.notificationId,
      );
      if (expandedId.value === n.notificationId) expandedId.value = null;
    }
  } catch {
    return;
  }
};

const markAllRead = async () => {
  if (loading.value) return;
  loading.value = true;
  try {
    await notificationsApi.readAllNotifications();
    await reload();
  } finally {
    loading.value = false;
  }
};

watch(
  () => tab.value,
  async () => {
    await reload();
  },
);

onMounted(async () => {
  await reload();
});
</script>

<style scoped>
.messages-page {
  min-height: 100vh;
  background: transparent;
}

.top {
  position: sticky;
  top: 0;
  z-index: 10;
  display: grid;
  grid-template-columns: auto 1fr auto;
  gap: 12px;
  align-items: center;
  padding: 12px 18px;
  background: rgba(255, 255, 255, 0.82);
  backdrop-filter: blur(18px);
  border-bottom: 1px solid rgba(226, 232, 240, 0.9);
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.04);
}

.back {
  height: 34px;
  padding: 0 12px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 10px;
  cursor: pointer;
}

.title {
  font-weight: 900;
  color: #0f172a;
}

.ghost {
  height: 34px;
  padding: 0 12px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 10px;
  cursor: pointer;
  color: #334155;
}

.main {
  padding: 24px 18px 34px;
  max-width: 760px;
  margin: 0 auto;
}

.card {
  background: var(--card-bg);
  border: 1px solid var(--card-border);
  border-radius: 22px;
  box-shadow: var(--card-shadow);
  overflow: hidden;
}

.head {
  padding: 14px 16px 10px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  border-bottom: 1px solid rgba(226, 232, 240, 0.8);
}

.tabs {
  display: flex;
  gap: 10px;
  align-items: center;
}

.tab {
  height: 32px;
  padding: 0 12px;
  border: 1px solid rgba(226, 232, 240, 0.95);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.88);
  cursor: pointer;
  font-weight: 800;
  color: #334155;
}

.tab.active {
  border-color: rgba(79, 70, 229, 0.32);
  color: #4f46e5;
  background: rgba(79, 70, 229, 0.08);
}

.tools {
  display: flex;
  gap: 10px;
  align-items: center;
}

.link {
  border: 0;
  background: transparent;
  color: #4f46e5;
  cursor: pointer;
  padding: 0;
  font-weight: 800;
}

.muted {
  color: #94a3b8;
}

.pad {
  padding: 16px;
}

.list {
  display: grid;
}

.item {
  border: 0;
  background: transparent;
  text-align: left;
  padding: 12px 16px;
  cursor: pointer;
  border-bottom: 1px solid rgba(226, 232, 240, 0.65);
}

.row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 10px;
}

.left {
  display: inline-flex;
  gap: 10px;
  align-items: center;
  min-width: 0;
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 999px;
  background: #4f46e5;
  flex: 0 0 auto;
}

.dot.read {
  background: rgba(148, 163, 184, 0.6);
}

.t {
  font-weight: 900;
  color: #0f172a;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.time {
  font-size: 12px;
  color: #94a3b8;
  flex: 0 0 auto;
}

.content {
  margin-top: 8px;
  color: #475569;
  font-size: 13px;
  line-height: 1.5;
  white-space: pre-wrap;
}

.foot {
  padding: 12px 16px;
  display: grid;
}

.full {
  width: 100%;
}
</style>
