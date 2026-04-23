<template>
  <div class="auth-page" :class="{ admin: mode === 'admin' }">
    <div class="card">
      <div class="card-head">
        <div class="brand" @click="router.replace('/m/home')">
          <img class="logo" src="@/assets/logo.png" alt="logo" />
          <div class="brand-meta">
            <div class="brand-name">智聚社区</div>
            <div class="brand-sub">{{ mode === "admin" ? "管理端登录" : "用户中心" }}</div>
          </div>
        </div>
        <div class="head-actions">
          <div v-if="mode === 'admin'" class="chip danger">Admin</div>
          <div v-else-if="mode === 'register'" class="chip">Register</div>
          <button class="text-btn" type="button" @click="router.replace('/m/home')">返回首页</button>
        </div>
      </div>

      <div class="card-title">{{ title }}</div>
      <div class="card-sub">{{ subtitle }}</div>

      <div v-if="mode === 'admin'" class="notice">
        <div class="notice-title">提示</div>
        <div class="notice-text">仅管理员账号可登录后台管理系统。</div>
      </div>

      <form v-if="mode !== 'register'" class="form" @submit.prevent="onSubmitLogin">
      <label class="field">
        <span>账号</span>
          <input v-model.trim="account" autocomplete="username" placeholder="请输入账号" />
      </label>
      <label class="field">
        <span>密码</span>
          <input v-model.trim="password" type="password" autocomplete="current-password" placeholder="请输入密码" />
      </label>
        <button class="btn" type="submit" :disabled="loading">
          {{ loading ? "登录中..." : "登录" }}
        </button>
      <p v-if="error" class="error">{{ error }}</p>
    </form>

    <form v-else class="form" @submit.prevent="onSubmitRegister">
      <label class="field">
        <span>账号</span>
          <input v-model.trim="regAccount" autocomplete="username" placeholder="请输入账号" />
      </label>
      <label class="field">
        <span>密码</span>
          <input v-model.trim="regPassword" type="password" autocomplete="new-password" placeholder="请输入密码" />
      </label>
      <label class="field">
        <span>确认密码</span>
          <input v-model.trim="regPassword2" type="password" autocomplete="new-password" placeholder="请再次输入密码" />
      </label>
      <label class="field">
        <span>昵称（可选）</span>
        <input v-model.trim="nickname" placeholder="请输入昵称" />
      </label>

      <div class="field">
        <span>头像（可选）</span>
        <div class="uploader">
            <button class="preview-btn" type="button" :disabled="uploading" @click="triggerPick">
              <img v-if="avatarUrl" class="preview-img" :src="avatarUrl" alt="avatar" />
              <div v-else class="preview-empty">{{ uploading ? "上传中..." : "上传头像" }}</div>
              <div class="preview-tip">{{ avatarUrl ? "点击更换（已回显）" : "支持 png/jpg/webp" }}</div>
            </button>
          <div class="upload-actions">
            <input ref="fileInputRef" class="file" type="file" accept="image/*" @change="onPickFile" />
              <button class="ghost" type="button" :disabled="uploading" @click="triggerPick">
                {{ uploading ? "上传中..." : "从本地选择" }}
              </button>
              <button class="ghost" type="button" :disabled="uploading || !avatarUrl" @click="clearAvatar">
              清除
            </button>
          </div>
        </div>
      </div>

      <button class="btn" type="submit" :disabled="loading">
        {{ loading ? "注册中..." : "注册" }}
      </button>
      <p v-if="error" class="error">{{ error }}</p>
      <p v-if="success" class="success">{{ success }}</p>
    </form>

    <div class="switch">
      <button v-if="mode === 'login'" class="link" type="button" @click="goRegister">
        没有账号？去注册
      </button>
      <button v-else-if="mode === 'register'" class="link" type="button" @click="goLogin">
        已有账号？去登录
      </button>
    </div>
      <div class="bottom-actions">
        <button v-if="mode === 'admin'" class="ghost full" type="button" @click="router.replace('/admin/login')">
          管理员登录入口
        </button>
        <button v-else class="ghost full" type="button" @click="router.replace('/m/home')">返回首页</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import { usersApi } from "@/apis";

const router = useRouter();
const route = useRoute();
const auth = useAuthStore();

const mode = computed<"login" | "admin" | "register">(() => {
  if (route.meta.adminLogin) return "admin";
  if (route.meta.register) return "register";
  return "login";
});

const title = computed(() => {
  if (mode.value === "admin") return "管理员登录";
  if (mode.value === "register") return "用户注册";
  return "登录";
});

const subtitle = computed(() => {
  if (mode.value === "admin") return "请输入管理员账号密码以进入后台。";
  if (mode.value === "register") return "创建账号后即可发帖、评论、关注与接收通知。";
  return "欢迎回来，登录后继续浏览与互动。";
});

const account = ref("");
const password = ref("");
const loading = ref(false);
const error = ref("");
const success = ref("");

const regAccount = ref("");
const regPassword = ref("");
const regPassword2 = ref("");
const nickname = ref("");
const avatarUrl = ref("");
const uploading = ref(false);
const fileInputRef = ref<HTMLInputElement | null>(null);

const onSubmitLogin = async () => {
  error.value = "";
  loading.value = true;
  try {
    await auth.login(account.value, password.value);
    if (mode.value === "admin") {
      if (!auth.isAdmin) {
        auth.logout();
        error.value = "非管理员账号";
        return;
      }
      await router.replace("/admin/dashboard");
      return;
    }
    const redirect =
      typeof route.query.redirect === "string"
        ? route.query.redirect
        : "/m/home";
    await router.replace(redirect);
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : "登录失败";
  } finally {
    loading.value = false;
  }
};

const validateRegister = () => {
  if (!regAccount.value.trim()) return "请输入账号";
  if (!regPassword.value.trim()) return "请输入密码";
  if (regPassword.value !== regPassword2.value) return "两次输入的密码不一致";
  return "";
};

const onSubmitRegister = async () => {
  error.value = "";
  success.value = "";
  const msg = validateRegister();
  if (msg) {
    error.value = msg;
    return;
  }

  loading.value = true;
  try {
    await usersApi.register({
      account: regAccount.value.trim(),
      password: regPassword.value.trim(),
      nickname: nickname.value.trim() || undefined,
      avatarUrl: avatarUrl.value || undefined,
    });
    success.value = "注册成功，请登录";
    await router.replace({ name: "mobile-login", query: { account: regAccount.value.trim() } });
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : "注册失败";
  } finally {
    loading.value = false;
  }
};

const triggerPick = () => {
  fileInputRef.value?.click();
};

const clearAvatar = () => {
  avatarUrl.value = "";
  if (fileInputRef.value) fileInputRef.value.value = "";
};

const onPickFile = async (evt: Event) => {
  const input = evt.target as HTMLInputElement | null;
  const file = input?.files?.[0];
  if (!file) return;
  uploading.value = true;
  error.value = "";
  try {
    const data = await usersApi.upload({ file, bucket: "public", scene: "avatar" });
    avatarUrl.value = data.url;
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : "上传失败";
    if (fileInputRef.value) fileInputRef.value.value = "";
  } finally {
    uploading.value = false;
  }
};

const goRegister = () => {
  router.push({ name: "mobile-register" });
};

const goLogin = () => {
  router.push({ name: "mobile-login" });
};

onMounted(() => {
  const qAccount = typeof route.query.account === "string" ? route.query.account : "";
  if (qAccount && mode.value !== "register") account.value = qAccount;
});
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  padding: 22px 16px;
  display: grid;
  place-items: center;
  background:
    radial-gradient(900px 300px at 15% 15%, rgba(79, 70, 229, 0.20), transparent 60%),
    radial-gradient(900px 300px at 85% 25%, rgba(56, 189, 248, 0.16), transparent 55%),
    linear-gradient(180deg, #f8fafc 0%, #f5f7fb 100%);
}

.auth-page.admin {
  background:
    radial-gradient(900px 300px at 15% 15%, rgba(220, 38, 38, 0.12), transparent 60%),
    radial-gradient(900px 300px at 85% 25%, rgba(79, 70, 229, 0.18), transparent 55%),
    linear-gradient(180deg, #f8fafc 0%, #f5f7fb 100%);
}

.card {
  width: min(420px, 100%);
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(14px);
  border: 1px solid #eef0f5;
  border-radius: 18px;
  box-shadow: 0 20px 60px rgba(15, 23, 42, 0.08);
  padding: 16px;
}

.card-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.brand {
  display: inline-flex;
  gap: 10px;
  align-items: center;
  cursor: pointer;
  user-select: none;
}

.logo {
  width: 34px;
  height: 34px;
  border-radius: 10px;
}

.brand-meta {
  display: grid;
  gap: 2px;
}

.brand-name {
  font-weight: 900;
  color: #0f172a;
  line-height: 1;
}

.brand-sub {
  font-size: 12px;
  color: #94a3b8;
}

.head-actions {
  display: inline-flex;
  gap: 10px;
  align-items: center;
}

.chip {
  height: 24px;
  padding: 0 10px;
  border-radius: 999px;
  background: rgba(79, 70, 229, 0.12);
  color: #4f46e5;
  font-weight: 900;
  font-size: 12px;
  display: inline-flex;
  align-items: center;
}

.chip.danger {
  background: rgba(220, 38, 38, 0.12);
  color: #dc2626;
}

.text-btn {
  border: 0;
  background: transparent;
  color: #64748b;
  cursor: pointer;
  padding: 0;
  font-weight: 800;
}

.card-title {
  margin-top: 12px;
  font-size: 22px;
  font-weight: 1000;
  color: #0f172a;
}

.card-sub {
  margin-top: 6px;
  color: #64748b;
  font-size: 12px;
  line-height: 1.5;
}

.notice {
  margin-top: 12px;
  border: 1px solid rgba(220, 38, 38, 0.18);
  background: rgba(220, 38, 38, 0.06);
  border-radius: 14px;
  padding: 10px 12px;
}

.notice-title {
  font-weight: 900;
  color: #dc2626;
}

.notice-text {
  margin-top: 4px;
  font-size: 12px;
  color: #7f1d1d;
  line-height: 1.5;
}

.form {
  display: grid;
  gap: 12px;
  margin-top: 14px;
}

.field {
  display: grid;
  gap: 6px;
}

input {
  height: 38px;
  padding: 0 12px;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  background: #fff;
  outline: none;
}

.btn {
  height: 40px;
  border: 0;
  background: #4f46e5;
  color: #fff;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 900;
}

.ghost {
  height: 38px;
  padding: 0 14px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 12px;
  cursor: pointer;
  color: #334155;
}

.error {
  color: #dc2626;
  margin: 0;
}

.success {
  color: #16a34a;
  margin: 0;
}

.switch {
  margin-top: 10px;
}

.link {
  border: 0;
  background: transparent;
  color: #4f46e5;
  cursor: pointer;
  padding: 0;
  font-weight: 900;
}

.uploader {
  display: grid;
  grid-template-columns: 1fr;
  gap: 12px;
  align-items: center;
}

.preview-btn {
  width: 100%;
  min-height: 104px;
  border-radius: 16px;
  border: 1px dashed #cbd5e1;
  background: #f8fafc;
  display: grid;
  place-items: center;
  gap: 6px;
  cursor: pointer;
  padding: 12px;
}

.preview-img {
  width: 76px;
  height: 76px;
  border-radius: 999px;
  object-fit: cover;
  border: 1px solid #e2e8f0;
  background: #fff;
}

.preview-empty {
  font-weight: 900;
  color: #334155;
}

.preview-tip {
  font-size: 12px;
  color: #94a3b8;
}

.upload-actions {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}

.file {
  display: none;
}

.bottom-actions {
  margin-top: 14px;
  display: grid;
  gap: 10px;
}

.full {
  width: 100%;
}
</style>
