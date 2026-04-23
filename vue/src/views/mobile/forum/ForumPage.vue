<template>
  <div class="page">
    <h1>贴吧详情</h1>
    <p>forumId: {{ route.params.id }}</p>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'

const route = useRoute()
</script>

<style scoped>
.page {
  padding: 16px;
}
</style>

