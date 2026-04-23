<template>
  <div class="admin-layout">
    <header class="header">
      <span class="title">管理端</span>
      <nav class="nav">
        <router-link class="link" :to="{ name: 'admin-dashboard' }">概览</router-link>
        <router-link class="link" :to="{ name: 'admin-posts-pending' }">待审帖子</router-link>
        <router-link class="link" :to="{ name: 'admin-reports' }">举报</router-link>
      </nav>
    </header>
    <main class="main">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts"></script>

<style scoped>
.admin-layout {
  min-height: 100vh;
  background: #f5f7fb;
}

.header {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  background: #fff;
  border-bottom: 1px solid #eef0f5;
}

.nav {
  display: flex;
  gap: 12px;
}

.link {
  text-decoration: none;
  color: #334155;
}

.main {
  padding: 16px;
}
</style>

