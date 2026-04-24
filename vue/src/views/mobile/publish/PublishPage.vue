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

        <div class="field">
          <div class="field-row">
            <span>正文</span>
            <span class="hint">支持富文本；建议正文纯文本长度 1~5000</span>
          </div>

          <div class="toolbar">
            <button class="tool" type="button" @mousedown.prevent @click="exec('bold')">加粗</button>
            <button class="tool" type="button" @mousedown.prevent @click="exec('italic')">斜体</button>
            <button class="tool" type="button" @mousedown.prevent @click="exec('underline')">下划线</button>
            <div class="sep" />
            <button class="tool" type="button" @mousedown.prevent @click="exec('insertUnorderedList')">无序</button>
            <button class="tool" type="button" @mousedown.prevent @click="exec('insertOrderedList')">有序</button>
            <div class="sep" />
            <button class="tool" type="button" @mousedown.prevent @click="onLink">链接</button>
            <button class="tool" type="button" @mousedown.prevent @click="triggerImage">图片</button>
            <button class="tool ghosty" type="button" @mousedown.prevent @click="clearFormat">清除格式</button>
          </div>

          <div class="editor-wrap">
            <div
              ref="editorRef"
              class="editor"
              contenteditable="true"
              :data-placeholder="'请输入正文（支持图片、链接、列表）...'"
              @input="onEditorInput"
            />
          </div>
          <input
            ref="imageInputRef"
            class="file"
            type="file"
            accept="image/*"
            @change="onPickImage"
          />
        </div>

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
import { feedApi, forumsApi, usersApi } from "@/apis";

const router = useRouter();

const forums = ref<Forum[]>([]);
const forumId = ref(0);
const title = ref("");
const contentHtml = ref("");
const submitting = ref(false);
const error = ref("");
const success = ref("");
const editorRef = ref<HTMLDivElement | null>(null);
const imageInputRef = ref<HTMLInputElement | null>(null);

const loadForums = async () => {
  const data = await forumsApi.getForums({ pageNum: 1, pageSize: 100 });
  forums.value = data.list || [];
  if (forums.value.length > 0 && forumId.value === 0) {
    forumId.value = forums.value[0].forumId;
  }
};

const sanitizeRichHtml = (html: string) => {
  if (!html) return "";
  const doc = new DOMParser().parseFromString(html, "text/html");
  const blocked = ["script", "style", "iframe", "object", "embed", "link", "meta"];
  for (const sel of blocked) {
    for (const el of Array.from(doc.querySelectorAll(sel))) el.remove();
  }
  for (const el of Array.from(doc.body.querySelectorAll("*"))) {
    for (const attr of Array.from(el.attributes)) {
      const name = attr.name.toLowerCase();
      const value = (attr.value || "").trim();
      if (name.startsWith("on")) {
        el.removeAttribute(attr.name);
        continue;
      }
      if (name === "style") {
        el.removeAttribute(attr.name);
        continue;
      }
      if ((name === "href" || name === "src") && /^javascript:/i.test(value)) {
        el.removeAttribute(attr.name);
        continue;
      }
    }
  }
  return doc.body.innerHTML;
};

const getPlainTextLength = () => {
  const text = editorRef.value?.innerText || "";
  return text.replace(/\s+/g, " ").trim().length;
};

const validate = () => {
  const t = title.value.trim();
  const plainLen = getPlainTextLength();
  if (!forumId.value) return "请选择贴吧";
  if (t.length < 1 || t.length > 200) return "标题长度需在 1~200";
  if (plainLen < 1 || plainLen > 5000) return "正文纯文本长度需在 1~5000";
  return "";
};

const onEditorInput = () => {
  const html = editorRef.value?.innerHTML || "";
  contentHtml.value = html;
};

const focusEditor = () => {
  editorRef.value?.focus();
};

const exec = (command: string) => {
  focusEditor();
  document.execCommand(command);
  onEditorInput();
};

const clearFormat = () => {
  focusEditor();
  document.execCommand("removeFormat");
  document.execCommand("unlink");
  onEditorInput();
};

const onLink = () => {
  focusEditor();
  const url = window.prompt("请输入链接 URL（https://...）");
  if (!url) return;
  document.execCommand("createLink", false, url.trim());
  onEditorInput();
};

const insertHtmlAtCursor = (html: string) => {
  focusEditor();
  const sel = window.getSelection();
  if (!sel || sel.rangeCount === 0) {
    editorRef.value!.insertAdjacentHTML("beforeend", html);
    onEditorInput();
    return;
  }
  const range = sel.getRangeAt(0);
  range.deleteContents();
  const frag = range.createContextualFragment(html);
  const last = frag.lastChild;
  range.insertNode(frag);
  if (last) {
    const r = document.createRange();
    r.setStartAfter(last);
    r.collapse(true);
    sel.removeAllRanges();
    sel.addRange(r);
  }
  onEditorInput();
};

const triggerImage = () => {
  imageInputRef.value?.click();
};

const onPickImage = async (e: Event) => {
  const input = e.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;
  try {
    const data = await usersApi.upload({ file, bucket: "public", scene: "post-content" });
    if (data.url) {
      const safeUrl = data.url.replace(/"/g, "");
      insertHtmlAtCursor(`<img src="${safeUrl}" alt="image" />`);
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : "图片上传失败";
  } finally {
    input.value = "";
  }
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
    const raw = editorRef.value?.innerHTML || contentHtml.value || "";
    const html = sanitizeRichHtml(raw);
    const data = await feedApi.createPost({
      forumId: forumId.value,
      title: title.value.trim(),
      content: html,
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
  contentHtml.value = "";
  if (editorRef.value) editorRef.value.innerHTML = "";
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

.field-row {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  align-items: baseline;
}

.hint {
  font-size: 12px;
  color: #94a3b8;
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

.toolbar {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  align-items: center;
  padding: 10px 10px 0;
}

.tool {
  height: 30px;
  padding: 0 10px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 999px;
  cursor: pointer;
  color: #334155;
  font-size: 12px;
  font-weight: 800;
}

.tool.ghosty {
  border-style: dashed;
  color: #64748b;
}

.sep {
  width: 1px;
  height: 18px;
  background: #e2e8f0;
}

.editor-wrap {
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  overflow: hidden;
  background: #fff;
}

.editor {
  min-height: 220px;
  padding: 12px;
  outline: none;
  line-height: 1.75;
  color: #0f172a;
  white-space: normal;
  word-break: break-word;
}

.editor:empty:before {
  content: attr(data-placeholder);
  color: #94a3b8;
}

:deep(.editor img) {
  max-width: 100%;
  height: auto;
  display: block;
  border-radius: 12px;
  margin: 8px 0;
  border: 1px solid #eef0f5;
}

:deep(.editor a) {
  color: #4f46e5;
}

.file {
  display: none;
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
