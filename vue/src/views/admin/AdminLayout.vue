<template>
  <div class="admin-layout">
    <aside class="sider">
      <div class="brand" @click="router.push('/m/home')">
        <img class="logo" src="@/assets/logo.png" alt="logo" />
        <div class="brand-text">
          <div class="brand-name">智聚</div>
          <div class="brand-sub">管理</div>
        </div>
      </div>

      <nav class="menu">
        <router-link class="menu-item" :to="{ name: 'admin-dashboard' }">
          <span class="dot" />
          <span>管理概览</span>
        </router-link>
        <router-link class="menu-item" :to="{ name: 'admin-posts-pending' }">
          <span class="dot" />
          <span>帖子审核</span>
        </router-link>
        <router-link class="menu-item" :to="{ name: 'admin-comments-review' }">
          <span class="dot" />
          <span>评论审核</span>
        </router-link>
        <router-link class="menu-item" :to="{ name: 'admin-users' }">
          <span class="dot" />
          <span>用户管理</span>
        </router-link>
        <router-link class="menu-item" :to="{ name: 'admin-notifications' }">
          <span class="dot" />
          <span>通知公告</span>
        </router-link>
        <router-link class="menu-item" :to="{ name: 'admin-reports' }">
          <span class="dot" />
          <span>举报审计</span>
        </router-link>
      </nav>

      <div class="sider-footer">
        <button class="ghost" type="button" @click="router.push('/m/home')">
          返回前台
        </button>
      </div>
    </aside>

    <div class="content">
      <header class="topbar">
        <div class="top-left">
          <div class="page-title">管理系统</div>
          <div class="page-sub">{{ todayText }}</div>
        </div>
        <div class="top-right">
          <div class="me">
            <div class="avatar">
              <img
                v-if="auth.user?.avatarUrl"
                class="avatar-img"
                :src="auth.user.avatarUrl"
                alt="avatar"
              />
              <div v-else class="avatar-fallback">
                {{ (auth.user?.nickname || "管").slice(0, 1) }}
              </div>
            </div>
            <div class="me-meta">
              <div class="me-name">
                {{ auth.user?.nickname || `管理员${auth.user?.userId || ""}` }}
              </div>
              <div class="me-role">{{ auth.user?.role || "admin" }}</div>
            </div>
          </div>
          <button class="danger" type="button" @click="onLogout">退出</button>
        </div>
      </header>

      <main class="main">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";

const router = useRouter();
const auth = useAuthStore();

const todayText = computed(() => {
  const d = new Date();
  const yy = d.getFullYear();
  const mm = String(d.getMonth() + 1).padStart(2, "0");
  const dd = String(d.getDate()).padStart(2, "0");
  return `${yy}-${mm}-${dd}`;
});

const onLogout = () => {
  auth.logout();
  router.replace("/m/home");
};
</script>

<style scoped>
.admin-layout {
  min-height: 100vh;
  background: #f5f7fb;
  display: grid;
  grid-template-columns: 240px 1fr;
}

.sider {
  background: #fff;
  border-right: 1px solid #eef0f5;
  display: grid;
  grid-template-rows: auto 1fr auto;
  padding: 14px 12px;
  gap: 12px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  user-select: none;
}

.logo {
  width: 34px;
  height: 34px;
  border-radius: 10px;
}

.brand-text {
  display: grid;
  gap: 2px;
}

.brand-name {
  font-weight: 900;
  color: #0f172a;
}

.brand-sub {
  font-size: 12px;
  color: #64748b;
}

.menu {
  display: grid;
  gap: 8px;
  align-content: start;
}

.menu-item {
  display: grid;
  grid-template-columns: 10px 1fr;
  gap: 10px;
  align-items: center;
  height: 42px;
  padding: 0 12px;
  border-radius: 12px;
  text-decoration: none;
  color: #334155;
}

.menu-item.router-link-active {
  background: rgba(79, 70, 229, 0.12);
  color: #4f46e5;
  font-weight: 900;
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 99px;
  background: #cbd5e1;
}

.menu-item.router-link-active .dot {
  background: #4f46e5;
}

.sider-footer {
  padding-top: 6px;
  border-top: 1px solid #eef0f5;
}

.ghost {
  width: 100%;
  height: 36px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 12px;
  cursor: pointer;
  color: #334155;
}

.content {
  display: grid;
  grid-template-rows: 64px 1fr;
  min-width: 0;
}

.topbar {
  position: sticky;
  top: 0;
  z-index: 20;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid #eef0f5;
}

.top-left {
  display: grid;
  gap: 2px;
}

.page-title {
  font-weight: 900;
  color: #0f172a;
}

.page-sub {
  font-size: 12px;
  color: #94a3b8;
}

.top-right {
  display: inline-flex;
  gap: 12px;
  align-items: center;
}

.me {
  display: flex;
  gap: 10px;
  align-items: center;
}

.avatar {
  width: 36px;
  height: 36px;
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

.me-meta {
  display: grid;
  gap: 2px;
}

.me-name {
  font-weight: 800;
  color: #0f172a;
}

.me-role {
  font-size: 12px;
  color: #94a3b8;
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

.main {
  padding: 16px;
}

@media (max-width: 960px) {
  .admin-layout {
    grid-template-columns: 1fr;
  }
  .sider {
    display: none;
  }
}
</style>
