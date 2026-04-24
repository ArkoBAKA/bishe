<template>
  <div class="profile-page">
    <header class="top">
      <button class="back" type="button" @click="router.back()">返回</button>
      <div class="title">我的</div>
      <button class="ghost" type="button" @click="router.replace('/m/home')">
        首页
      </button>
    </header>

    <main class="main">
      <div class="card">
        <div class="me">
          <img
            v-if="auth.user?.avatarUrl"
            class="avatar"
            :src="auth.user.avatarUrl"
            alt="avatar"
          />
          <div v-else class="avatar placeholder">
            {{ (auth.user?.nickname || "我").slice(0, 1) }}
          </div>
          <div class="meta">
            <div class="name">
              {{ auth.user?.nickname || `用户${auth.user?.userId}` }}
            </div>
            <div class="sub">
              <span>ID {{ auth.user?.userId }}</span>
              <span v-if="auth.user?.role">· {{ auth.user.role }}</span>
            </div>
          </div>
        </div>

        <div class="row">
          <button
            class="primary"
            type="button"
            @click="router.push('/m/publish')"
          >
            发帖
          </button>
          <button class="danger" type="button" @click="onLogout">
            退出登录
          </button>
          <button
            v-if="auth.isAdmin"
            class="ghost"
            type="button"
            @click="router.push('/admin/dashboard')"
          >
            管理端
          </button>
        </div>
      </div>

      <div class="card">
        <div class="section-head">
          <div class="section-title">我的关注（贴吧）</div>
          <button class="link" type="button" @click="loadFollows">刷新</button>
        </div>
        <div v-if="loadingFollows" class="muted">加载中...</div>
        <div v-else-if="followForums.length === 0" class="muted">暂无关注</div>
        <div v-else class="list">
          <button
            v-for="f in followForums"
            :key="f.forumId"
            class="item"
            type="button"
            @click="
              router.push({
                name: 'mobile-forum',
                params: { id: String(f.forumId) },
              })
            "
          >
            <div class="item-title">{{ f.name }}</div>
            <div class="item-sub">
              {{ f.description || `贴吧 ID ${f.forumId}` }}
            </div>
          </button>
        </div>
      </div>

      <div class="card">
        <div class="section-head">
          <div class="section-title">未读消息</div>
          <button class="link" type="button" @click="loadUnread">刷新</button>
        </div>
        <div v-if="loadingUnread" class="muted">加载中...</div>
        <div v-else class="muted">未读 {{ unreadCount }} 条</div>
        <button class="ghost" type="button" @click="router.push('/m/messages')">
          查看消息
        </button>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onActivated, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import type { Forum } from "@/types/api";
import { followsApi, forumsApi, notificationsApi } from "@/apis";
import { useAuthStore } from "@/stores/auth";

const router = useRouter();
const route = useRoute();
const auth = useAuthStore();

const loadingFollows = ref(false);
const followForumIds = ref<number[]>([]);
const allForums = ref<Forum[]>([]);

const loadingUnread = ref(false);
const unreadCount = ref(0);

const followForums = computed(() => {
  const map = new Map<number, Forum>();
  for (const f of allForums.value) map.set(f.forumId, f);
  return followForumIds.value.map(
    (id) => map.get(id) || { forumId: id, name: `贴吧${id}` },
  );
});

const loadAllForums = async () => {
  try {
    const data = await forumsApi.getForums({ pageNum: 1, pageSize: 200 });
    allForums.value = data.list || [];
  } catch {
    allForums.value = [];
  }
};

const loadFollows = async () => {
  loadingFollows.value = true;
  try {
    const data = await followsApi.getMyFollows({
      pageNum: 1,
      pageSize: 200,
      targetType: "forum",
    });
    followForumIds.value = (data.list || [])
      .filter((x) => x.active !== false)
      .map((x) => x.targetId);
  } catch {
    followForumIds.value = [];
  } finally {
    loadingFollows.value = false;
  }
};

const loadUnread = async () => {
  loadingUnread.value = true;
  try {
    const data = await notificationsApi.getNotifications({
      pageNum: 1,
      pageSize: 1,
      isRead: false,
    });
    unreadCount.value = data.total || 0;
  } catch {
    unreadCount.value = 0;
  } finally {
    loadingUnread.value = false;
  }
};

const onLogout = () => {
  auth.logout();
  router.replace("/m/home");
};

const refresh = async () => {
  if (!auth.isAuthed) {
    followForumIds.value = [];
    allForums.value = [];
    unreadCount.value = 0;
    return;
  }
  await Promise.all([loadAllForums(), loadFollows(), loadUnread()]);
};

onMounted(async () => {
  await refresh();
});

onActivated(async () => {
  if (route.name !== "mobile-profile") return;
  await refresh();
});

watch(
  () => auth.isAuthed,
  async () => {
    if (route.name !== "mobile-profile") return;
    await refresh();
  },
);
</script>

<style scoped>
.profile-page {
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
  display: grid;
  gap: 16px;
}

.card {
  background: var(--card-bg);
  border: 1px solid var(--card-border);
  border-radius: 22px;
  padding: 18px;
  display: grid;
  gap: 14px;
  box-shadow: var(--card-shadow);
  backdrop-filter: blur(10px);
}

.me {
  display: grid;
  grid-template-columns: 54px 1fr;
  gap: 14px;
  align-items: center;
}

.avatar {
  width: 54px;
  height: 54px;
  border-radius: 999px;
  object-fit: cover;
  border: 2px solid rgba(255, 255, 255, 0.85);
  box-shadow: 0 10px 24px rgba(79, 70, 229, 0.12);
}

.avatar.placeholder {
  display: grid;
  place-items: center;
  background: #eef2ff;
  color: #4f46e5;
  font-weight: 900;
}

.meta {
  display: grid;
  gap: 2px;
}

.name {
  font-weight: 900;
  color: #0f172a;
}

.sub {
  color: #64748b;
  font-size: 12px;
}

.row {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.primary {
  height: 40px;
  padding: 0 16px;
  border: 0;
  background: linear-gradient(135deg, #4f46e5 0%, #6366f1 100%);
  color: #fff;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 900;
  box-shadow: 0 10px 24px rgba(79, 70, 229, 0.2);
}

.danger {
  height: 40px;
  padding: 0 16px;
  border: 1px solid #fecaca;
  background: rgba(255, 255, 255, 0.9);
  color: #dc2626;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 800;
}

.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.section-title {
  font-weight: 900;
  color: #0f172a;
}

.link {
  border: 0;
  background: transparent;
  color: #4f46e5;
  cursor: pointer;
  padding: 0;
}

.muted {
  color: #94a3b8;
}

.list {
  display: grid;
  gap: 10px;
}

.item {
  border: 1px solid #eef0f5;
  border-radius: 16px;
  padding: 14px;
  background: linear-gradient(180deg, #ffffff 0%, #fbfdff 100%);
  cursor: pointer;
  text-align: left;
  transition:
    border-color 0.2s ease,
    box-shadow 0.2s ease,
    transform 0.2s ease;
}

.item:hover {
  border-color: rgba(79, 70, 229, 0.2);
  box-shadow: 0 14px 32px rgba(15, 23, 42, 0.08);
  transform: translateY(-2px);
}

.item-title {
  font-weight: 800;
  color: #0f172a;
}

.item-sub {
  margin-top: 2px;
  font-size: 12px;
  color: #64748b;
}
</style>
