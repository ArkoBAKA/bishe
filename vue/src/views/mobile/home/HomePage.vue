<template>
  <div class="home">
    <header class="topbar">
      <div class="brand" @click="router.push('/m/home')">
        <img class="logo" src="@/assets/logo.png" alt="logo" />
        <span class="brand-name">智聚社区</span>
      </div>

      <div class="search">
        <input
          v-model.trim="keyword"
          class="search-input"
          placeholder="搜索贴吧、帖子..."
        />
      </div>

      <nav class="actions">
        <button class="nav-btn" type="button" @click="router.push('/m/home')">
          首页
        </button>
        <button class="nav-btn" type="button" @click="onClickNotifications">
          消息
        </button>
        <button class="nav-btn" type="button" @click="onClickPublish">
          发帖
        </button>
      </nav>

      <div class="user">
        <div v-if="!auth.isAuthed" class="guest-actions">
          <button class="primary" type="button" @click="openLogin">登录</button>
          <button class="ghost" type="button" @click="goRegister">注册</button>
        </div>
        <div v-else class="user-info">
          <img
            v-if="auth.user?.avatarUrl"
            class="avatar"
            :src="auth.user.avatarUrl"
            alt="avatar"
          />
          <div v-else class="avatar placeholder">
            {{ (auth.user?.nickname || "我").slice(0, 1) }}
          </div>
          <div class="user-meta">
            <div class="nickname">
              {{ auth.user?.nickname || `用户${auth.user?.userId}` }}
            </div>
            <div class="user-actions">
              <button
                class="link"
                type="button"
                @click="router.push('/m/profile')"
              >
                个人
              </button>
              <button class="link danger" type="button" @click="onLogout">
                退出
              </button>
            </div>
          </div>
        </div>
      </div>
    </header>

    <main class="layout">
      <aside class="left">
        <div class="panel">
          <div class="panel-title">
            {{ auth.isAuthed ? "我的关注" : "全部贴吧" }}
          </div>
          <div v-if="forumsLoading" class="muted">加载中...</div>
          <div v-else class="forum-list">
            <button
              v-for="f in forums"
              :key="f.forumId"
              class="forum-item"
              type="button"
              @click="toForum(f.forumId)"
            >
              <div class="forum-name">{{ f.name }}</div>
              <div class="forum-sub">
                {{ f.description || `ID: ${f.forumId}` }}
              </div>
              <div class="forum-ops" @click.stop>
                <button
                  class="follow"
                  type="button"
                  @click="toggleFollow(f.forumId)"
                >
                  {{ isFollowed(f.forumId) ? "已关注" : "关注" }}
                </button>
              </div>
            </button>
          </div>
          <div class="panel-footer">
            <button class="ghost" type="button" @click="loadForums()">
              发现更多贴吧
            </button>
          </div>
        </div>
      </aside>

      <section class="center">
        <div class="panel composer">
          <div class="composer-left">
            <img
              v-if="auth.user?.avatarUrl"
              class="avatar"
              :src="auth.user.avatarUrl"
              alt="avatar"
            />
            <div v-else class="avatar placeholder">
              {{ (auth.user?.nickname || "访").slice(0, 1) }}
            </div>
          </div>
          <button class="composer-input" type="button" @click="onClickPublish">
            {{
              auth.isAuthed
                ? "有什么新鲜事想和大家分享？"
                : "登录后即可发布帖子"
            }}
          </button>
          <div class="composer-right">
            <button class="icon-btn" type="button" @click="onClickPublish">
              发布
            </button>
          </div>
        </div>

        <div class="panel">
          <div class="tabs">
            <button
              class="tab"
              :class="{ active: tab === 'latest' }"
              type="button"
              @click="tab = 'latest'"
            >
              最新发布
            </button>
            <button
              class="tab"
              :class="{ active: tab === 'hot' }"
              type="button"
              @click="tab = 'hot'"
            >
              热门推荐
            </button>
            <button
              class="tab"
              :class="{ active: tab === 'follow' }"
              type="button"
              @click="tab = 'follow'"
            >
              精华汇总
            </button>
          </div>

          <div v-if="postsLoading" class="muted pad">加载中...</div>
          <div v-else-if="filteredPosts.length === 0" class="muted pad">
            暂无内容
          </div>
          <div v-else class="post-list">
            <article
              v-for="p in filteredPosts"
              :key="p.postId"
              class="post-card"
            >
              <div class="post-head">
                <div class="post-author">
                  <div class="avatar small placeholder">
                    {{
                      (
                        p.author?.nickname || `U${p.author?.userId || ""}`
                      ).slice(0, 1) || "U"
                    }}
                  </div>
                  <div class="author-meta">
                    <div class="author-name">
                      {{
                        p.author?.nickname ||
                        `用户${p.author?.userId || ""}` ||
                        "匿名"
                      }}
                    </div>
                    <div class="post-sub">
                      <button
                        class="tag"
                        type="button"
                        @click="toForum(p.forum?.forumId || p.forumId)"
                      >
                        {{
                          p.forum?.name ||
                          `贴吧${p.forum?.forumId || p.forumId}`
                        }}
                      </button>
                      <span v-if="p.createdAt"
                        >· {{ formatTime(p.createdAt) }}</span
                      >
                    </div>
                  </div>
                </div>
              </div>

              <h3 class="post-title">{{ p.title }}</h3>
              <div v-if="p.content" class="post-content">{{ p.content }}</div>

              <div class="post-foot">
                <button class="post-btn" type="button" @click="openComments(p)">
                  评论<span v-if="typeof p.commentCount === 'number'">
                    {{ p.commentCount }}</span
                  >
                </button>
                <button
                  class="post-btn"
                  type="button"
                  @click="toForum(p.forum?.forumId || p.forumId)"
                >
                  进入贴吧
                </button>
                <div class="post-metrics">
                  <span v-if="typeof p.viewCount === 'number'"
                    >浏览 {{ p.viewCount }}</span
                  >
                  <span v-if="typeof p.likeCount === 'number'"
                    >点赞 {{ p.likeCount }}</span
                  >
                </div>
              </div>
            </article>
          </div>

          <div class="panel-footer">
            <button class="ghost" type="button" @click="loadPosts()">
              加载更多动态
            </button>
          </div>
        </div>
      </section>

      <aside class="right">
        <div class="panel">
          <div class="panel-head">
            <div class="panel-title">热门贴吧</div>
            <button class="link" type="button" @click="router.push('/m/home')">
              查看全部
            </button>
          </div>
          <div class="hot-list">
            <div
              v-for="(f, idx) in hotForums"
              :key="f.forumId"
              class="hot-item"
            >
              <div class="rank">{{ idx + 1 }}</div>
              <button
                class="hot-main"
                type="button"
                @click="toForum(f.forumId)"
              >
                <div class="hot-name">{{ f.name }}</div>
                <div class="hot-sub">
                  <span v-if="typeof f.followersCount === 'number'"
                    >{{ f.followersCount }} 人关注</span
                  >
                  <span v-else>贴吧 ID {{ f.forumId }}</span>
                </div>
              </button>
              <button
                class="follow"
                type="button"
                @click="toggleFollow(f.forumId)"
              >
                {{ isFollowed(f.forumId) ? "已关注" : "关注" }}
              </button>
            </div>
          </div>
        </div>

        <div class="panel">
          <div class="panel-title">全站活跃趋势</div>
          <div class="trend">
            <svg viewBox="0 0 240 80" class="trend-svg" aria-hidden="true">
              <path
                d="M0 55 C 30 25, 60 70, 90 40 C 120 15, 150 55, 180 30 C 200 18, 220 40, 240 22"
                fill="none"
                stroke="#4f46e5"
                stroke-width="3"
                stroke-linecap="round"
              />
            </svg>
          </div>
        </div>
      </aside>
    </main>

    <button
      class="admin-entry"
      type="button"
      @click="router.push({ name: 'admin-login' })"
    >
      管理端入口
    </button>

    <div v-if="loginOpen" class="modal-mask" @click.self="closeLogin">
      <div class="modal">
        <div class="modal-title">用户登录</div>
        <form class="modal-form" @submit.prevent="onLoginSubmit">
          <label class="field">
            <span>账号</span>
            <input v-model.trim="loginAccount" placeholder="请输入账号" />
          </label>
          <label class="field">
            <span>密码</span>
            <input
              v-model.trim="loginPassword"
              type="password"
              placeholder="请输入密码"
            />
          </label>
          <button class="primary" type="submit" :disabled="loginLoading">
            {{ loginLoading ? "登录中..." : "登录" }}
          </button>
          <button class="ghost" type="button" @click="goRegister">
            注册
          </button>
          <p v-if="loginError" class="error">{{ loginError }}</p>
        </form>
      </div>
    </div>

    <div v-if="commentsOpen" class="drawer-mask" @click.self="closeComments">
      <div class="drawer">
        <div class="drawer-head">
          <div class="drawer-title">评论</div>
          <button class="link" type="button" @click="closeComments">
            关闭
          </button>
        </div>
        <div v-if="commentsLoading" class="muted pad">加载中...</div>
        <div v-else-if="comments.length === 0" class="muted pad">暂无评论</div>
        <div v-else class="comment-list">
          <div v-for="c in comments" :key="c.commentId" class="comment-item">
            <div class="avatar small placeholder">
              {{ (c.author?.nickname || "U").slice(0, 1) }}
            </div>
            <div class="comment-body">
              <div class="comment-head">
                <span class="comment-author">{{
                  c.author?.nickname ||
                  `用户${c.author?.userId || ""}` ||
                  "匿名"
                }}</span>
                <span v-if="c.createdAt" class="comment-time">{{
                  formatTime(c.createdAt)
                }}</span>
              </div>
              <div class="comment-content">{{ c.content }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="notifyOpen" class="drawer-mask" @click.self="notifyOpen = false">
      <div class="drawer">
        <div class="drawer-head">
          <div class="drawer-title">消息</div>
          <button class="link" type="button" @click="notifyOpen = false">
            关闭
          </button>
        </div>
        <div v-if="!auth.isAuthed" class="muted pad">
          <button class="primary" type="button" @click="openLogin">
            登录后查看消息
          </button>
        </div>
        <div v-else-if="notifyLoading" class="muted pad">加载中...</div>
        <div v-else-if="notifications.length === 0" class="muted pad">
          暂无消息
        </div>
        <div v-else class="notify-list">
          <div
            v-for="n in notifications"
            :key="n.notificationId"
            class="notify-item"
          >
            <div class="notify-title">
              <span class="dot" :class="{ read: n.isRead }" />
              <span>{{ n.title || `通知 #${n.notificationId}` }}</span>
            </div>
            <div v-if="n.content" class="notify-content">{{ n.content }}</div>
            <div v-if="n.createdAt" class="notify-time">
              {{ formatTime(n.createdAt) }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRouter } from "vue-router";
import type { CommentItem, Forum, NotificationItem, Post } from "@/types/api";
import { feedApi, followsApi, forumsApi, notificationsApi } from "@/apis";
import { useAuthStore } from "@/stores/auth";

const router = useRouter();
const auth = useAuthStore();

const keyword = ref("");
const tab = ref<"latest" | "hot" | "follow">("latest");

const forums = ref<Forum[]>([]);
const forumsLoading = ref(false);
const myFollowForumIdSet = ref<Set<number>>(new Set());

const posts = ref<Post[]>([]);
const postsLoading = ref(false);

const loginOpen = ref(false);
const loginAccount = ref("");
const loginPassword = ref("");
const loginLoading = ref(false);
const loginError = ref("");

const commentsOpen = ref(false);
const commentsLoading = ref(false);
const comments = ref<CommentItem[]>([]);
const currentPostId = ref<number | null>(null);

const notifyOpen = ref(false);
const notifyLoading = ref(false);
const notifications = ref<NotificationItem[]>([]);

const isFollowed = (forumId: number) => myFollowForumIdSet.value.has(forumId);

const hotForums = computed(() => {
  const list = [...forums.value];
  list.sort((a, b) => {
    const av = typeof a.followersCount === "number" ? a.followersCount : -1;
    const bv = typeof b.followersCount === "number" ? b.followersCount : -1;
    if (bv !== av) return bv - av;
    return b.forumId - a.forumId;
  });
  return list.slice(0, 5);
});

const filteredPosts = computed(() => {
  const key = keyword.value.trim().toLowerCase();
  let list = [...posts.value];

  if (tab.value === "hot") {
    list.sort((a, b) => (b.viewCount || 0) - (a.viewCount || 0));
  } else if (tab.value === "follow") {
    list = list.filter((p) => isFollowed(p.forum?.forumId || p.forumId));
  } else {
    list.sort((a, b) => {
      const at = a.createdAt ? new Date(a.createdAt).getTime() : 0;
      const bt = b.createdAt ? new Date(b.createdAt).getTime() : 0;
      return bt - at;
    });
  }

  if (!key) return list;
  return list.filter((p) => {
    const title = (p.title || "").toLowerCase();
    const content = (p.content || "").toLowerCase();
    const forumName = (p.forum?.name || "").toLowerCase();
    return (
      title.includes(key) || content.includes(key) || forumName.includes(key)
    );
  });
});

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

const openLogin = () => {
  loginOpen.value = true;
  loginError.value = "";
};

const goRegister = () => {
  loginOpen.value = false;
  router.push({ name: "mobile-register" });
};

const closeLogin = () => {
  loginOpen.value = false;
};

const onLoginSubmit = async () => {
  loginError.value = "";
  loginLoading.value = true;
  try {
    await auth.login(loginAccount.value, loginPassword.value);
    loginOpen.value = false;
    loginAccount.value = "";
    loginPassword.value = "";
    await Promise.all([loadFollows(), loadPosts(), loadNotificationsIfOpen()]);
  } catch (e: unknown) {
    loginError.value = e instanceof Error ? e.message : "登录失败";
  } finally {
    loginLoading.value = false;
  }
};

const onLogout = () => {
  auth.logout();
  myFollowForumIdSet.value = new Set();
  posts.value = [];
  notifications.value = [];
  loadPosts();
};

const loadForums = async () => {
  forumsLoading.value = true;
  try {
    const data = await forumsApi.getForums({
      pageNum: 1,
      pageSize: 50,
      keyword: keyword.value.trim() || undefined,
    });
    forums.value = data.list || [];
  } finally {
    forumsLoading.value = false;
  }
};

const loadFollows = async () => {
  if (!auth.isAuthed) {
    myFollowForumIdSet.value = new Set();
    return;
  }
  try {
    const data = await followsApi.getMyFollows({
      pageNum: 1,
      pageSize: 200,
      targetType: "forum",
    });
    const ids = new Set<number>();
    for (const item of data.list || []) ids.add(item.targetId);
    myFollowForumIdSet.value = ids;
  } catch (e) {
    myFollowForumIdSet.value = new Set();
  }
};

const loadGuestLatestPosts = async () => {
  const forumData = await forumsApi.getForums({
    pageNum: 1,
    pageSize: 6,
    keyword: keyword.value.trim() || undefined,
  });
  const topForums = forumData.list || [];
  const pages = await Promise.all(
    topForums.map(async (f) => {
      try {
        const postsData = await forumsApi.getForumPosts(f.forumId, {
          pageNum: 1,
          pageSize: 5,
        });
        const list = postsData.list || [];
        return list.map((p) => ({
          ...p,
          forum: { forumId: f.forumId, name: f.name, coverUrl: f.coverUrl },
        }));
      } catch (e) {
        return [];
      }
    }),
  );
  const merged = pages.flat();
  merged.sort((a, b) => {
    const at = a.createdAt ? new Date(a.createdAt).getTime() : 0;
    const bt = b.createdAt ? new Date(b.createdAt).getTime() : 0;
    return bt - at;
  });
  posts.value = merged;
};

const loadPosts = async () => {
  postsLoading.value = true;
  try {
    if (auth.isAuthed) {
      try {
        const data = await feedApi.getFeed({ pageNum: 1, pageSize: 20 });
        posts.value = data.list || [];
        return;
      } catch (e) {
        await loadGuestLatestPosts();
        return;
      }
    }
    await loadGuestLatestPosts();
  } finally {
    postsLoading.value = false;
  }
};

const toggleFollow = async (forumId: number) => {
  if (!auth.isAuthed) {
    openLogin();
    return;
  }
  if (isFollowed(forumId)) {
    await followsApi.unfollow({ targetType: "forum", targetId: forumId });
  } else {
    await followsApi.follow({ targetType: "forum", targetId: forumId });
  }
  await loadFollows();
};

const toForum = (forumId: number | string) => {
  router.push({ name: "mobile-forum", params: { id: String(forumId) } });
};

const openComments = async (post: Post) => {
  commentsOpen.value = true;
  commentsLoading.value = true;
  comments.value = [];
  currentPostId.value = post.postId;
  try {
    const data = await feedApi.getPostComments(post.postId, {
      pageNum: 1,
      pageSize: 30,
    });
    comments.value = data.list || [];
  } finally {
    commentsLoading.value = false;
  }
};

const closeComments = () => {
  commentsOpen.value = false;
  comments.value = [];
  currentPostId.value = null;
};

const loadNotificationsIfOpen = async () => {
  if (!notifyOpen.value) return;
  if (!auth.isAuthed) return;
  notifyLoading.value = true;
  try {
    const data = await notificationsApi.getNotifications({
      pageNum: 1,
      pageSize: 20,
    });
    notifications.value = data.list || [];
  } finally {
    notifyLoading.value = false;
  }
};

const onClickNotifications = async () => {
  notifyOpen.value = true;
  await loadNotificationsIfOpen();
};

const onClickPublish = () => {
  if (!auth.isAuthed) {
    openLogin();
    return;
  }
  router.push("/m/publish");
};

watch(
  () => auth.isAuthed,
  async (v) => {
    if (v) {
      await loadFollows();
      await loadPosts();
    } else {
      myFollowForumIdSet.value = new Set();
      await loadPosts();
    }
  },
);

let searchTimer: number | null = null;
watch(
  () => keyword.value,
  () => {
    if (searchTimer) window.clearTimeout(searchTimer);
    searchTimer = window.setTimeout(() => {
      loadForums();
      loadPosts();
    }, 250);
  },
);

onMounted(async () => {
  if (auth.token && !auth.user) {
    await auth.fetchMe().catch(() => undefined);
  }
  await loadForums();
  await loadFollows();
  await loadPosts();
});
</script>

<style scoped>
.home {
  min-height: 100vh;
  background: #f5f7fb;
}

.topbar {
  position: sticky;
  top: 0;
  z-index: 10;
  height: 64px;
  display: grid;
  grid-template-columns: 220px 1fr auto 260px;
  gap: 12px;
  align-items: center;
  padding: 0 16px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid #eef0f5;
}

.brand {
  display: inline-flex;
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

.brand-name {
  font-weight: 700;
  letter-spacing: 0.5px;
  color: #0f172a;
}

.search {
  display: flex;
  align-items: center;
}

.search-input {
  width: 100%;
  height: 38px;
  padding: 0 14px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 12px;
  outline: none;
}

.actions {
  display: inline-flex;
  gap: 10px;
  align-items: center;
}

.nav-btn {
  height: 36px;
  padding: 0 12px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 10px;
  cursor: pointer;
}

.user {
  display: flex;
  justify-content: flex-end;
}

.guest-actions {
  display: inline-flex;
  gap: 10px;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-meta {
  display: grid;
  gap: 2px;
}

.nickname {
  font-weight: 600;
  color: #0f172a;
}

.user-actions {
  display: inline-flex;
  gap: 10px;
}

.layout {
  max-width: 1200px;
  margin: 0 auto;
  padding: 16px;
  display: grid;
  grid-template-columns: 260px 1fr 320px;
  gap: 16px;
  align-items: start;
}

.panel {
  background: #fff;
  border: 1px solid #eef0f5;
  border-radius: 16px;
  box-shadow: 0 6px 18px rgba(15, 23, 42, 0.03);
}

.panel-title {
  font-weight: 700;
  padding: 14px 14px 10px;
  color: #0f172a;
}

.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 14px 10px;
}

.panel-footer {
  padding: 12px 14px 14px;
  display: flex;
  justify-content: center;
}

.ghost {
  height: 34px;
  padding: 0 12px;
  border: 1px dashed #d7dbe5;
  background: #fff;
  border-radius: 10px;
  cursor: pointer;
  color: #475569;
}

.forum-list {
  display: grid;
  gap: 10px;
  padding: 0 12px 12px;
}

.forum-item {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 10px;
  align-items: center;
  padding: 10px 12px;
  border-radius: 12px;
  border: 1px solid #eef0f5;
  background: #fff;
  cursor: pointer;
  text-align: left;
}

.forum-name {
  font-weight: 600;
  color: #0f172a;
}

.forum-sub {
  margin-top: 2px;
  font-size: 12px;
  color: #64748b;
}

.forum-ops {
  display: flex;
  justify-content: flex-end;
}

.follow {
  height: 30px;
  padding: 0 12px;
  border: 1px solid #4f46e5;
  color: #4f46e5;
  background: #fff;
  border-radius: 999px;
  cursor: pointer;
}

.center {
  display: grid;
  gap: 16px;
}

.composer {
  display: grid;
  grid-template-columns: 44px 1fr auto;
  gap: 12px;
  align-items: center;
  padding: 12px 14px;
}

.composer-input {
  height: 40px;
  border: 1px solid #e2e8f0;
  background: #f8fafc;
  border-radius: 12px;
  text-align: left;
  padding: 0 14px;
  cursor: pointer;
  color: #64748b;
}

.icon-btn {
  height: 34px;
  padding: 0 12px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 10px;
  cursor: pointer;
}

.tabs {
  display: flex;
  gap: 10px;
  padding: 12px 14px 0;
}

.tab {
  height: 34px;
  padding: 0 12px;
  border: 0;
  background: transparent;
  border-bottom: 2px solid transparent;
  cursor: pointer;
  color: #475569;
}

.tab.active {
  color: #4f46e5;
  border-color: #4f46e5;
  font-weight: 700;
}

.post-list {
  display: grid;
  gap: 14px;
  padding: 14px;
}

.post-card {
  border: 1px solid #eef0f5;
  border-radius: 14px;
  padding: 14px;
  background: #fff;
}

.post-author {
  display: flex;
  gap: 10px;
  align-items: center;
}

.author-meta {
  display: grid;
  gap: 2px;
}

.author-name {
  font-weight: 700;
  color: #0f172a;
}

.post-sub {
  font-size: 12px;
  color: #64748b;
  display: inline-flex;
  gap: 8px;
  align-items: center;
}

.tag {
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 999px;
  padding: 2px 10px;
  height: 24px;
  cursor: pointer;
  color: #475569;
}

.post-title {
  margin: 12px 0 8px;
  color: #0f172a;
  line-height: 1.3;
}

.post-content {
  white-space: pre-wrap;
  color: #334155;
  line-height: 1.65;
}

.post-foot {
  margin-top: 12px;
  display: grid;
  grid-template-columns: auto auto 1fr;
  gap: 10px;
  align-items: center;
}

.post-btn {
  height: 32px;
  padding: 0 12px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 10px;
  cursor: pointer;
}

.post-metrics {
  justify-self: end;
  display: inline-flex;
  gap: 10px;
  font-size: 12px;
  color: #94a3b8;
}

.hot-list {
  padding: 0 12px 12px;
  display: grid;
  gap: 10px;
}

.hot-item {
  display: grid;
  grid-template-columns: 18px 1fr auto;
  gap: 10px;
  align-items: center;
  padding: 10px 12px;
  border: 1px solid #eef0f5;
  border-radius: 12px;
}

.rank {
  font-weight: 800;
  color: #f97316;
  text-align: center;
}

.hot-main {
  border: 0;
  background: transparent;
  text-align: left;
  padding: 0;
  cursor: pointer;
}

.hot-name {
  font-weight: 700;
  color: #0f172a;
}

.hot-sub {
  font-size: 12px;
  color: #64748b;
  margin-top: 2px;
}

.trend {
  padding: 0 14px 14px;
}

.trend-svg {
  width: 100%;
  height: 96px;
}

.admin-entry {
  position: fixed;
  left: 12px;
  bottom: 12px;
  height: 26px;
  padding: 0 10px;
  border-radius: 999px;
  border: 1px solid #e2e8f0;
  background: rgba(255, 255, 255, 0.92);
  cursor: pointer;
  font-size: 12px;
  color: #64748b;
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

.link {
  border: 0;
  background: transparent;
  color: #4f46e5;
  cursor: pointer;
  padding: 0;
}

.link.danger {
  color: #dc2626;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 999px;
  object-fit: cover;
}

.avatar.placeholder {
  display: grid;
  place-items: center;
  background: #eef2ff;
  color: #4f46e5;
  font-weight: 800;
}

.avatar.small {
  width: 34px;
  height: 34px;
  font-size: 12px;
}

.muted {
  color: #94a3b8;
}

.pad {
  padding: 14px;
}

.modal-mask,
.drawer-mask {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.42);
  display: grid;
  place-items: center;
  z-index: 50;
}

.modal {
  width: min(420px, calc(100vw - 32px));
  background: #fff;
  border-radius: 16px;
  border: 1px solid #eef0f5;
  padding: 16px;
}

.modal-title {
  font-weight: 800;
  color: #0f172a;
  margin-bottom: 12px;
}

.modal-form {
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

.field input {
  height: 38px;
  padding: 0 12px;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  outline: none;
}

.error {
  color: #dc2626;
  margin: 0;
}

.drawer-mask {
  place-items: end;
}

.drawer {
  width: min(520px, 100vw);
  height: min(640px, calc(100vh - 64px));
  background: #fff;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  border: 1px solid #eef0f5;
  overflow: auto;
}

.drawer-head {
  position: sticky;
  top: 0;
  background: #fff;
  border-bottom: 1px solid #eef0f5;
  padding: 12px 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.drawer-title {
  font-weight: 800;
  color: #0f172a;
}

.comment-list {
  padding: 12px 14px;
  display: grid;
  gap: 12px;
}

.comment-item {
  display: grid;
  grid-template-columns: 34px 1fr;
  gap: 10px;
  align-items: start;
}

.comment-head {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  align-items: baseline;
}

.comment-author {
  font-weight: 700;
  color: #0f172a;
}

.comment-time {
  font-size: 12px;
  color: #94a3b8;
}

.comment-content {
  margin-top: 4px;
  white-space: pre-wrap;
  color: #334155;
}

.notify-list {
  padding: 12px 14px;
  display: grid;
  gap: 12px;
}

.notify-item {
  border: 1px solid #eef0f5;
  border-radius: 12px;
  padding: 12px;
}

.notify-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 700;
  color: #0f172a;
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 99px;
  background: #4f46e5;
}

.dot.read {
  background: #cbd5e1;
}

.notify-content {
  margin-top: 6px;
  color: #334155;
  white-space: pre-wrap;
}

.notify-time {
  margin-top: 8px;
  font-size: 12px;
  color: #94a3b8;
}

@media (max-width: 1024px) {
  .topbar {
    grid-template-columns: 180px 1fr auto;
    grid-template-rows: auto auto;
    height: auto;
    padding: 12px 12px;
  }
  .user {
    grid-column: 3;
    grid-row: 1;
  }
  .actions {
    grid-column: 1 / -1;
    justify-content: flex-start;
  }
  .layout {
    grid-template-columns: 1fr;
  }
  .left,
  .right {
    display: none;
  }
}
</style>
