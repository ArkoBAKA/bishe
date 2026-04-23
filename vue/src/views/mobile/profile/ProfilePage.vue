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
        <button class="ghost" type="button" @click="router.replace('/m/home')">
          去首页查看
        </button>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import type { Forum } from "@/types/api";
import { followsApi, forumsApi, notificationsApi } from "@/apis";
import { useAuthStore } from "@/stores/auth";

const router = useRouter();
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
  const data = await forumsApi.getForums({ pageNum: 1, pageSize: 200 });
  allForums.value = data.list || [];
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
  } finally {
    loadingUnread.value = false;
  }
};

const onLogout = () => {
  auth.logout();
  router.replace("/m/home");
};

onMounted(async () => {
  await loadAllForums();
  await loadFollows();
  await loadUnread();
});
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  background: #f5f7fb;
}

.top {
  position: sticky;
  top: 0;
  z-index: 10;
  display: grid;
  grid-template-columns: auto 1fr auto;
  gap: 12px;
  align-items: center;
  padding: 12px 14px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid #eef0f5;
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
  padding: 14px;
  max-width: 760px;
  margin: 0 auto;
  display: grid;
  gap: 14px;
}

.card {
  background: #fff;
  border: 1px solid #eef0f5;
  border-radius: 16px;
  padding: 14px;
  display: grid;
  gap: 12px;
}

.me {
  display: grid;
  grid-template-columns: 54px 1fr;
  gap: 12px;
  align-items: center;
}

.avatar {
  width: 54px;
  height: 54px;
  border-radius: 999px;
  object-fit: cover;
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
  gap: 10px;
  flex-wrap: wrap;
}

.primary {
  height: 36px;
  padding: 0 14px;
  border: 0;
  background: #4f46e5;
  color: #fff;
  border-radius: 10px;
  cursor: pointer;
}

.danger {
  height: 36px;
  padding: 0 14px;
  border: 1px solid #fecaca;
  background: #fff;
  color: #dc2626;
  border-radius: 10px;
  cursor: pointer;
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
  border-radius: 14px;
  padding: 12px;
  background: #fff;
  cursor: pointer;
  text-align: left;
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
