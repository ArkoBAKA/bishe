<template>
  <div class="publish-page">
    <header class="top">
      <button class="back" type="button" @click="router.back()">返回</button>
      <div class="title">发布帖子</div>
      <button class="ghost" type="button" @click="router.replace('/m/home')">
        首页
      </button>
    </header>

    <main class="main">
      <div class="card">
        <label class="field">
          <span>选择贴吧</span>
          <select v-model.number="forumId">
            <option :value="0" disabled>请选择</option>
            <option v-for="f in forums" :key="f.forumId" :value="f.forumId">
              {{ f.name }}
            </option>
          </select>
        </label>

        <label class="field">
          <span>标题</span>
          <input v-model.trim="title" placeholder="请输入标题（1~200）" />
        </label>

        <label class="field">
          <span>正文</span>
          <textarea v-model.trim="content" placeholder="请输入正文（1~5000）" />
        </label>

        <div class="actions">
          <button
            class="primary"
            type="button"
            :disabled="submitting"
            @click="onSubmit"
          >
            {{ submitting ? "发布中..." : "发布" }}
          </button>
          <button
            class="danger"
            type="button"
            :disabled="submitting"
            @click="onReset"
          >
            清空
          </button>
        </div>

        <p v-if="error" class="error">{{ error }}</p>
        <p v-if="success" class="success">{{ success }}</p>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import type { Forum } from "@/types/api";
import { feedApi, forumsApi } from "@/apis";

const router = useRouter();

const forums = ref<Forum[]>([]);
const forumId = ref(0);
const title = ref("");
const content = ref("");
const submitting = ref(false);
const error = ref("");
const success = ref("");

const loadForums = async () => {
  const data = await forumsApi.getForums({ pageNum: 1, pageSize: 100 });
  forums.value = data.list || [];
  if (forums.value.length > 0 && forumId.value === 0) {
    forumId.value = forums.value[0].forumId;
  }
};

const validate = () => {
  const t = title.value.trim();
  const c = content.value.trim();
  if (!forumId.value) return "请选择贴吧";
  if (t.length < 1 || t.length > 200) return "标题长度需在 1~200";
  if (c.length < 1 || c.length > 5000) return "正文长度需在 1~5000";
  return "";
};

const onSubmit = async () => {
  error.value = "";
  success.value = "";
  const msg = validate();
  if (msg) {
    error.value = msg;
    return;
  }
  submitting.value = true;
  try {
    const data = await feedApi.createPost({
      forumId: forumId.value,
      title: title.value.trim(),
      content: content.value.trim(),
    });
    success.value = `发布成功：#${data.postId}（状态：${data.status}）`;
    router.replace("/m/home");
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : "发布失败";
  } finally {
    submitting.value = false;
  }
};

const onReset = () => {
  title.value = "";
  content.value = "";
  error.value = "";
  success.value = "";
};

onMounted(async () => {
  await loadForums();
});
</script>

<style scoped>
.publish-page {
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
  display: grid;
  place-items: start center;
}

.card {
  width: min(680px, 100%);
  background: #fff;
  border: 1px solid #eef0f5;
  border-radius: 16px;
  padding: 16px;
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

input,
select {
  height: 38px;
  padding: 0 12px;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  outline: none;
  background: #fff;
}

textarea {
  min-height: 180px;
  padding: 10px 12px;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  outline: none;
  resize: vertical;
  font-family: inherit;
  line-height: 1.6;
}

.actions {
  display: flex;
  gap: 12px;
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

.error {
  margin: 0;
  color: #dc2626;
}

.success {
  margin: 0;
  color: #16a34a;
}
</style>
