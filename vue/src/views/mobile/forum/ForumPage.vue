<template>
  <div class="forum-page">
    <header class="site-top">
      <div class="brand" @click="router.push('/m/home')">
        <img class="logo" src="@/assets/logo.png" alt="logo" />
        <span class="brand-name">智聚社区</span>
      </div>
      <div class="site-actions">
        <button class="ghost" type="button" @click="router.push('/m/home')">
          返回首页
        </button>
        <button
          v-if="!auth.isAuthed"
          class="ghost"
          type="button"
          @click="
            router.push({
              name: 'mobile-login',
              query: { redirect: route.fullPath },
            })
          "
        >
          登录
        </button>
        <button
          v-else
          class="avatar-btn"
          type="button"
          @click="router.push('/m/profile')"
        >
          <img
            v-if="auth.user?.avatarUrl"
            class="avatar-img"
            :src="auth.user.avatarUrl"
            alt="avatar"
          />
          <div v-else class="avatar-fallback">
            {{ (auth.user?.nickname || "我").slice(0, 1) }}
          </div>
        </button>
      </div>
    </header>

    <section class="hero" :style="heroStyle">
      <div class="hero-mask" />
      <div class="hero-inner">
        <div class="hero-left">
          <div class="forum-avatar">
            <img
              v-if="forumCoverUrl"
              class="forum-avatar-img"
              :src="forumCoverUrl"
              alt="cover"
            />
            <div v-else class="forum-avatar-fallback">
              {{ (forum?.name || "吧").slice(0, 1) }}
            </div>
          </div>
          <div class="forum-meta">
            <div class="forum-name-row">
              <div class="forum-name">{{ forum?.name || "贴吧" }}</div>
              <div class="verified" title="认证">✔</div>
            </div>
            <div class="forum-desc">
              {{ forum?.description || `在这里，记录你的兴趣。` }}
            </div>
          </div>
        </div>
        <div class="hero-right">
          <button class="follow" type="button" @click="toggleFollow">
            {{ isFollowed ? "已关注" : "+ 关注贴吧" }}
          </button>
          <button
            v-if="!auth.isAuthed"
            class="ghost"
            type="button"
            @click="
              router.push({
                name: 'mobile-login',
                query: { redirect: route.fullPath },
              })
            "
          >
            签到
          </button>
          <button
            v-else
            class="ghost"
            type="button"
            @click="router.push('/m/publish')"
          >
            发帖
          </button>
        </div>
      </div>
    </section>

    <main class="main">
      <div class="main-inner">
        <section class="left">
          <div class="panel stats">
            <div class="stat">
              <div class="stat-label">今日新增</div>
              <div class="stat-value">{{ todayNewCount }}</div>
            </div>
            <div class="stat">
              <div class="stat-label">累计帖子</div>
              <div class="stat-value">{{ totalPostsText }}</div>
            </div>
            <div class="stat">
              <div class="stat-label">吧友总数</div>
              <div class="stat-value">{{ followersText }}</div>
            </div>
            <div class="stat">
              <div class="stat-label">吧等级</div>
              <div class="stat-value level">{{ levelText }}</div>
            </div>
          </div>

          <div class="panel feed">
            <div class="feed-head">
              <div class="tabs">
                <button
                  class="tab"
                  :class="{ active: activeTab === 'latest' }"
                  type="button"
                  @click="activeTab = 'latest'"
                >
                  新帖
                </button>
                <button
                  class="tab"
                  :class="{ active: activeTab === 'hot' }"
                  type="button"
                  @click="activeTab = 'hot'"
                >
                  精品
                </button>
                <button
                  class="tab"
                  :class="{ active: activeTab === 'media' }"
                  type="button"
                  @click="activeTab = 'media'"
                >
                  视频
                </button>
                <button
                  class="tab"
                  :class="{ active: activeTab === 'help' }"
                  type="button"
                  @click="activeTab = 'help'"
                >
                  吧友互助
                </button>
              </div>
              <div class="feed-tools">
                <button
                  class="tool"
                  type="button"
                  @click="router.push('/m/publish')"
                >
                  发布
                </button>
                <button class="tool" type="button" @click="reload">刷新</button>
              </div>
            </div>

            <div v-if="loading" class="muted pad">加载中...</div>
            <div v-else-if="filteredPosts.length === 0" class="muted pad">
              暂无帖子
            </div>
            <div v-else class="post-list">
              <article
                v-for="p in filteredPosts"
                :key="p.postId"
                class="post-card"
                :class="{ pin: pinnedPostIdSet.has(p.postId) }"
              >
                <div class="post-row">
                  <div class="post-left">
                    <div class="post-tag" v-if="pinnedPostIdSet.has(p.postId)">
                      置顶
                    </div>
                    <div class="post-title">{{ p.title }}</div>
                    <div class="post-sub">
                      <span class="post-author">{{
                        p.author?.nickname ||
                        `用户${p.author?.userId || ""}` ||
                        "匿名"
                      }}</span>
                      <span v-if="p.createdAt" class="dot">·</span>
                      <span v-if="p.createdAt">{{
                        formatTimeShort(p.createdAt)
                      }}</span>
                    </div>
                  </div>
                  <div class="post-right">
                    <div class="metric">
                      <span class="metric-value">{{
                        typeof p.commentCount === "number"
                          ? p.commentCount
                          : "-"
                      }}</span>
                      <span class="metric-label">评论</span>
                    </div>
                  </div>
                </div>
                <div
                  v-if="p.content && pinnedPostIdSet.has(p.postId)"
                  class="post-preview"
                >
                  {{ p.content }}
                </div>
                <div class="post-actions">
                  <button
                    class="btn"
                    type="button"
                    @click="toggleComments(p.postId)"
                  >
                    {{ isCommentsExpanded(p.postId) ? "收起评论" : "评论" }}
                  </button>
                  <button
                    class="btn"
                    type="button"
                    @click="copyShare(p.postId)"
                  >
                    分享
                  </button>
                </div>

                <div v-if="isCommentsExpanded(p.postId)" class="comments-box">
                  <div v-if="isCommentsLoading(p.postId)" class="muted pad-sm">
                    加载中...
                  </div>
                  <div
                    v-else-if="(commentsByPostId[p.postId] || []).length === 0"
                    class="muted pad-sm"
                  >
                    暂无评论
                  </div>
                  <div v-else class="comment-thread">
                    <div
                      v-for="it in commentThread(p.postId)"
                      :key="it.comment.commentId"
                      class="comment-row"
                      :style="{ marginLeft: `${it.depth * 16}px` }"
                    >
                      <div class="avatar small placeholder">
                        {{ (it.comment.author?.nickname || "U").slice(0, 1) }}
                      </div>
                      <div class="comment-body">
                        <div class="comment-head">
                          <span class="comment-author">{{
                            it.comment.author?.nickname ||
                            `用户${it.comment.author?.userId || ""}` ||
                            "匿名"
                          }}</span>
                          <span v-if="it.replyTo" class="reply-to"
                            >回复
                            {{
                              it.replyTo.author?.nickname ||
                              `用户${it.replyTo.author?.userId || ""}` ||
                              "匿名"
                            }}</span
                          >
                          <span
                            v-if="it.comment.createdAt"
                            class="comment-time"
                            >{{ formatTime(it.comment.createdAt) }}</span
                          >
                        </div>
                        <div class="comment-content">
                          {{ it.comment.content }}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </article>
            </div>
          </div>
        </section>

        <aside class="right">
          <div class="panel intro">
            <div class="panel-title">
              <span class="i">i</span>
              <span>贴吧简介</span>
            </div>
            <div class="intro-text">
              {{ forum?.description || "暂无简介" }}
            </div>
            <div class="intro-meta">
              <div class="intro-item">
                <span class="intro-k">现任吧主</span>
                <span class="intro-v">{{ ownerText }}</span>
              </div>
              <div class="intro-item">
                <span class="intro-k">创建时间</span>
                <span class="intro-v">{{ createdAtText }}</span>
              </div>
            </div>
            <button
              class="ghost full"
              type="button"
              @click="
                router.push({
                  name: 'mobile-login',
                  query: { redirect: route.fullPath },
                })
              "
            >
              申请加入管理团队
            </button>
          </div>

          <div class="panel rank">
            <div class="panel-head">
              <div class="panel-title">活跃榜单</div>
              <div class="muted small">{{ rankSubtitle }}</div>
            </div>
            <div class="rank-list">
              <div v-for="(u, idx) in rankList" :key="u.key" class="rank-item">
                <div class="rank-no">
                  {{ String(idx + 1).padStart(2, "0") }}
                </div>
                <div class="rank-user">
                  <div class="rank-avatar">{{ u.avatar }}</div>
                  <div class="rank-name">{{ u.name }}</div>
                </div>
                <div class="rank-score">{{ u.score }}</div>
              </div>
            </div>
          </div>
        </aside>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import type { CommentItem, Forum, Post } from "@/types/api";
import { feedApi, followsApi, forumsApi } from "@/apis";
import { useAuthStore } from "@/stores/auth";

const route = useRoute();
const router = useRouter();
const auth = useAuthStore();

const forumId = Number(route.params.id);
const forum = ref<Forum | null>(null);
const posts = ref<Post[]>([]);
const postsTotal = ref<number | null>(null);
const loading = ref(false);
const isFollowed = ref(false);
const activeTab = ref<"latest" | "hot" | "media" | "help">("latest");

type CommentThreadItem = {
  comment: CommentItem;
  depth: number;
  replyTo?: CommentItem;
};
const commentsByPostId = ref<Record<number, CommentItem[]>>({});
const commentsLoadingByPostId = ref<Record<number, boolean>>({});
const commentsExpandedByPostId = ref<Record<number, boolean>>({});

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

const isCommentsExpanded = (postId: number) =>
  !!commentsExpandedByPostId.value[postId];
const isCommentsLoading = (postId: number) =>
  !!commentsLoadingByPostId.value[postId];

const buildCommentThread = (items: CommentItem[]) => {
  const byId = new Map<number, CommentItem>();
  const children = new Map<number, CommentItem[]>();
  for (const c of items) {
    byId.set(c.commentId, c);
    children.set(c.commentId, []);
  }
  const roots: CommentItem[] = [];
  for (const c of items) {
    const pid = c.parentCommentId || 0;
    if (pid && byId.has(pid)) {
      children.get(pid)!.push(c);
    } else {
      roots.push(c);
    }
  }
  const sortByTimeAsc = (a: CommentItem, b: CommentItem) => {
    const at = a.createdAt ? new Date(a.createdAt).getTime() : 0;
    const bt = b.createdAt ? new Date(b.createdAt).getTime() : 0;
    return at - bt;
  };
  roots.sort(sortByTimeAsc);
  for (const list of children.values()) list.sort(sortByTimeAsc);

  const out: CommentThreadItem[] = [];
  const dfs = (node: CommentItem, depth: number, replyTo?: CommentItem) => {
    out.push({ comment: node, depth, replyTo });
    const kids = children.get(node.commentId) || [];
    for (const k of kids) dfs(k, depth + 1, node);
  };
  for (const r of roots) dfs(r, 0);
  return out;
};

const commentThread = (postId: number) => {
  const list = commentsByPostId.value[postId] || [];
  return buildCommentThread(list);
};

const ensureCommentsLoaded = async (postId: number) => {
  if (commentsByPostId.value[postId]) return;
  commentsLoadingByPostId.value[postId] = true;
  try {
    const data = await feedApi.getPostComments(postId, {
      pageNum: 1,
      pageSize: 50,
    });
    commentsByPostId.value[postId] = data.list || [];
  } finally {
    commentsLoadingByPostId.value[postId] = false;
  }
};

const toggleComments = async (postId: number) => {
  const next = !isCommentsExpanded(postId);
  commentsExpandedByPostId.value[postId] = next;
  if (next) await ensureCommentsLoaded(postId);
};

const formatTimeShort = (value: string) => {
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return value;
  const mm = String(d.getMonth() + 1).padStart(2, "0");
  const dd = String(d.getDate()).padStart(2, "0");
  const hh = String(d.getHours()).padStart(2, "0");
  const mi = String(d.getMinutes()).padStart(2, "0");
  return `${mm}-${dd} ${hh}:${mi}`;
};

const forumCoverUrl = computed(() => {
  const url = forum.value?.coverUrl || "";
  return url || "";
});

const heroStyle = computed(() => {
  if (forumCoverUrl.value)
    return { backgroundImage: `url(${forumCoverUrl.value})` };
  return {
    backgroundImage:
      "linear-gradient(135deg, #c7d2fe 0%, #f1f5f9 60%, #dbeafe 100%)",
  };
});

const totalPostsText = computed(() => {
  const t = postsTotal.value;
  if (typeof t === "number") return formatCount(t);
  return formatCount(posts.value.length);
});

const followersText = computed(() => {
  const n = forum.value?.followersCount;
  if (typeof n !== "number") return "-";
  return formatCount(n);
});

const todayNewCount = computed(() => {
  const now = new Date();
  const y = now.getFullYear();
  const m = now.getMonth();
  const d = now.getDate();
  const start = new Date(y, m, d, 0, 0, 0, 0).getTime();
  return posts.value.filter((p) => {
    if (!p.createdAt) return false;
    const t = new Date(p.createdAt).getTime();
    return Number.isFinite(t) && t >= start;
  }).length;
});

const levelText = computed(() => {
  const n = forum.value?.followersCount ?? 0;
  if (n >= 50000) return "钻石级";
  if (n >= 10000) return "黄金级";
  if (n >= 2000) return "白银级";
  return "青铜级";
});

const pinnedPostIdSet = computed(() => {
  const list = [...posts.value];
  list.sort((a, b) => {
    const at = a.createdAt ? new Date(a.createdAt).getTime() : 0;
    const bt = b.createdAt ? new Date(b.createdAt).getTime() : 0;
    return bt - at;
  });
  const ids = new Set<number>();
  for (const p of list.slice(0, 2)) ids.add(p.postId);
  return ids;
});

const filteredPosts = computed(() => {
  const list = [...posts.value];
  if (activeTab.value === "hot") {
    list.sort((a, b) => (b.viewCount || 0) - (a.viewCount || 0));
    return list;
  }
  if (activeTab.value === "media") {
    return list.filter((p) => (p.content || "").includes("http"));
  }
  if (activeTab.value === "help") {
    return list.filter(
      (p) =>
        (p.title || "").includes("求助") || (p.title || "").includes("帮助"),
    );
  }
  list.sort((a, b) => {
    const at = a.createdAt ? new Date(a.createdAt).getTime() : 0;
    const bt = b.createdAt ? new Date(b.createdAt).getTime() : 0;
    return bt - at;
  });
  return list;
});

const ownerText = computed(() => {
  const raw = forum.value as unknown as { ownerId?: number };
  if (typeof raw?.ownerId === "number") return `用户${raw.ownerId}`;
  return "-";
});

const createdAtText = computed(() => {
  const raw = forum.value as unknown as { createdAt?: string };
  if (!raw?.createdAt) return "-";
  const d = new Date(raw.createdAt);
  if (Number.isNaN(d.getTime())) return raw.createdAt;
  const yy = d.getFullYear();
  const mm = String(d.getMonth() + 1).padStart(2, "0");
  const dd = String(d.getDate()).padStart(2, "0");
  return `${yy}-${mm}-${dd}`;
});

const formatCount = (n: number) => {
  if (n >= 10000) {
    const w = Math.round((n / 10000) * 10) / 10;
    return `${w}w`;
  }
  return String(n);
};

const rankList = computed(() => {
  const map = new Map<
    string,
    { key: string; name: string; avatar: string; score: number }
  >();
  for (const p of posts.value) {
    const id = p.author?.userId ? String(p.author.userId) : "0";
    const name =
      p.author?.nickname ||
      (p.author?.userId ? `用户${p.author.userId}` : "匿名");
    const key = `${id}-${name}`;
    const item = map.get(key) || {
      key,
      name,
      avatar: name.slice(0, 1),
      score: 0,
    };
    item.score += 1;
    map.set(key, item);
  }
  const list = [...map.values()];
  list.sort((a, b) => b.score - a.score);
  return list.slice(0, 5).map((x) => ({ ...x, score: `${x.score} 帖` }));
});

const rankSubtitle = computed(() =>
  rankList.value.length ? "近 30 帖统计" : "暂无数据",
);

const reload = async () => {
  loading.value = true;
  try {
    const [forumData, postsData] = await Promise.all([
      forumsApi.getForumDetail(forumId),
      forumsApi.getForumPosts(forumId, { pageNum: 1, pageSize: 20 }),
    ]);
    forum.value = forumData;
    postsTotal.value = postsData.total ?? null;
    posts.value = (postsData.list || []).map((p) => ({
      ...p,
      forum: {
        forumId: forumData.forumId,
        name: forumData.name,
        coverUrl: forumData.coverUrl,
      },
    }));
    await loadFollowState();
  } finally {
    loading.value = false;
  }
};

const copyShare = async (postId: number) => {
  const text = `${window.location.origin}${window.location.pathname}#/m/forum/${forumId}?postId=${postId}`;
  try {
    await navigator.clipboard.writeText(text);
  } catch (e) {
    return;
  }
};

const loadFollowState = async () => {
  if (!auth.isAuthed) {
    isFollowed.value = false;
    return;
  }
  try {
    const data = await followsApi.getMyFollows({
      pageNum: 1,
      pageSize: 200,
      targetType: "forum",
    });
    isFollowed.value = (data.list || []).some(
      (f) => f.targetId === forumId && f.active !== false,
    );
  } catch (e) {
    isFollowed.value = false;
  }
};

const toggleFollow = async () => {
  if (!auth.isAuthed) {
    router.push({ name: "mobile-login", query: { redirect: route.fullPath } });
    return;
  }
  if (isFollowed.value) {
    await followsApi.unfollow({ targetType: "forum", targetId: forumId });
  } else {
    await followsApi.follow({ targetType: "forum", targetId: forumId });
  }
  await loadFollowState();
};

onMounted(async () => {
  await reload();
});
</script>

<style scoped>
.forum-page {
  min-height: 100vh;
  background: #f5f7fb;
}

.site-top {
  position: sticky;
  top: 0;
  z-index: 30;
  height: 64px;
  padding: 0 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #eef0f5;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(12px);
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
  font-weight: 800;
  color: #0f172a;
}

.site-actions {
  display: inline-flex;
  gap: 10px;
  align-items: center;
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

.avatar-btn {
  width: 36px;
  height: 36px;
  border-radius: 999px;
  border: 1px solid #e2e8f0;
  background: #fff;
  overflow: hidden;
  padding: 0;
  cursor: pointer;
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

.hero {
  position: relative;
  min-height: 150px;
  background-size: cover;
  background-position: center;
}

.hero-mask {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    90deg,
    rgba(15, 23, 42, 0.55),
    rgba(15, 23, 42, 0.15)
  );
}

.hero-inner {
  position: relative;
  max-width: 1200px;
  margin: 0 auto;
  padding: 18px 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
}

.hero-left {
  display: flex;
  gap: 14px;
  align-items: center;
}

.forum-avatar {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.85);
  border: 1px solid rgba(255, 255, 255, 0.7);
  overflow: hidden;
  display: grid;
  place-items: center;
}

.forum-avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.forum-avatar-fallback {
  font-weight: 900;
  color: #0f172a;
}

.forum-meta {
  display: grid;
  gap: 6px;
  color: #fff;
}

.forum-name-row {
  display: inline-flex;
  gap: 8px;
  align-items: center;
}

.forum-name {
  font-weight: 900;
  font-size: 20px;
}

.verified {
  width: 18px;
  height: 18px;
  border-radius: 999px;
  background: #60a5fa;
  color: #0b1020;
  display: grid;
  place-items: center;
  font-size: 12px;
  font-weight: 900;
}

.forum-desc {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.88);
  max-width: 520px;
}

.hero-right {
  display: inline-flex;
  gap: 10px;
  align-items: center;
}

.follow {
  height: 34px;
  padding: 0 12px;
  border: 0;
  background: #4f46e5;
  color: #fff;
  border-radius: 999px;
  cursor: pointer;
}

.main {
  padding: 16px;
}

.main-inner {
  max-width: 1200px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 1fr 340px;
  gap: 16px;
  align-items: start;
}

.panel {
  background: #fff;
  border: 1px solid #eef0f5;
  border-radius: 16px;
  box-shadow: 0 6px 18px rgba(15, 23, 42, 0.03);
}

.stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  padding: 14px;
}

.stat {
  border-right: 1px solid #eef0f5;
  padding-right: 12px;
}

.stat:last-child {
  border-right: 0;
  padding-right: 0;
}

.stat-label {
  font-size: 12px;
  color: #94a3b8;
}

.stat-value {
  margin-top: 6px;
  font-weight: 900;
  color: #0f172a;
}

.stat-value.level {
  color: #f59e0b;
}

.feed {
  margin-top: 16px;
}

.feed-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 14px;
  border-bottom: 1px solid #eef0f5;
}

.tabs {
  display: flex;
  gap: 12px;
  align-items: center;
}

.tab {
  height: 32px;
  padding: 0 12px;
  border: 0;
  background: transparent;
  border-bottom: 2px solid transparent;
  color: #475569;
  cursor: pointer;
}

.tab.active {
  color: #4f46e5;
  border-color: #4f46e5;
  font-weight: 900;
}

.feed-tools {
  display: inline-flex;
  gap: 10px;
  align-items: center;
}

.tool {
  height: 32px;
  padding: 0 12px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 10px;
  cursor: pointer;
  color: #334155;
}

.pad {
  padding: 14px;
}

.post-list {
  display: grid;
  gap: 14px;
  padding: 14px;
}

.post-card {
  border: 1px solid #eef0f5;
  border-radius: 16px;
  padding: 12px 12px 10px;
}

.post-card.pin {
  border-color: rgba(79, 70, 229, 0.55);
  box-shadow: 0 10px 20px rgba(79, 70, 229, 0.08);
}

.post-row {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 12px;
  align-items: start;
}

.post-left {
  min-width: 0;
}

.post-tag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 22px;
  padding: 0 10px;
  border-radius: 999px;
  background: rgba(79, 70, 229, 0.12);
  color: #4f46e5;
  font-weight: 900;
  font-size: 12px;
  margin-bottom: 8px;
}

.post-title {
  font-weight: 900;
  color: #0f172a;
  line-height: 1.35;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.post-sub {
  margin-top: 6px;
  font-size: 12px;
  color: #94a3b8;
  display: inline-flex;
  gap: 6px;
  align-items: center;
}

.post-author {
  color: #475569;
  font-weight: 700;
}

.dot {
  color: #cbd5e1;
}

.muted {
  color: #94a3b8;
}

.small {
  font-size: 12px;
}

.post-right {
  display: grid;
  place-items: center;
}

.metric {
  display: grid;
  place-items: center;
  min-width: 54px;
}

.metric-value {
  font-weight: 900;
  color: #0f172a;
}

.metric-label {
  margin-top: 2px;
  font-size: 12px;
  color: #94a3b8;
}

.post-preview {
  margin-top: 10px;
  color: #334155;
  line-height: 1.7;
  white-space: pre-wrap;
  background: #f8fafc;
  border: 1px solid #eef0f5;
  border-radius: 12px;
  padding: 10px 12px;
  max-height: 160px;
  overflow: auto;
}

.post-actions {
  margin-top: 10px;
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}

.comments-box {
  margin-top: 12px;
  border-top: 1px solid #eef0f5;
  padding-top: 12px;
}

.pad-sm {
  padding: 10px 0 2px;
}

.comment-thread {
  display: grid;
  gap: 10px;
}

.comment-row {
  display: grid;
  grid-template-columns: 34px 1fr;
  gap: 10px;
  align-items: start;
}

.comment-body {
  border: 1px solid #eef0f5;
  background: #f8fafc;
  border-radius: 12px;
  padding: 10px 12px;
}

.comment-head {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  align-items: center;
}

.comment-author {
  font-weight: 800;
  color: #0f172a;
}

.reply-to {
  font-size: 12px;
  color: #64748b;
}

.comment-time {
  font-size: 12px;
  color: #94a3b8;
}

.comment-content {
  margin-top: 6px;
  color: #334155;
  line-height: 1.7;
  white-space: pre-wrap;
}

.btn {
  height: 32px;
  padding: 0 12px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 10px;
  cursor: pointer;
}

.intro {
  padding: 14px;
}

.panel-title {
  font-weight: 900;
  color: #0f172a;
  display: inline-flex;
  gap: 8px;
  align-items: center;
}

.i {
  width: 18px;
  height: 18px;
  border-radius: 999px;
  background: rgba(79, 70, 229, 0.12);
  color: #4f46e5;
  display: grid;
  place-items: center;
  font-size: 12px;
  font-weight: 900;
}

.intro-text {
  margin-top: 10px;
  color: #334155;
  line-height: 1.7;
  white-space: pre-wrap;
}

.intro-meta {
  margin-top: 12px;
  display: grid;
  gap: 8px;
}

.intro-item {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  font-size: 12px;
}

.intro-k {
  color: #94a3b8;
}

.intro-v {
  color: #334155;
  font-weight: 700;
}

.full {
  width: 100%;
  margin-top: 12px;
}

.rank {
  margin-top: 16px;
  padding: 14px;
}

.panel-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
}

.rank-list {
  margin-top: 12px;
  display: grid;
  gap: 10px;
}

.rank-item {
  display: grid;
  grid-template-columns: 30px 1fr auto;
  gap: 10px;
  align-items: center;
  padding: 10px 10px;
  border: 1px solid #eef0f5;
  border-radius: 14px;
}

.rank-no {
  font-weight: 900;
  color: #94a3b8;
  text-align: center;
}

.rank-user {
  display: inline-flex;
  gap: 10px;
  align-items: center;
}

.rank-avatar {
  width: 30px;
  height: 30px;
  border-radius: 999px;
  background: #eef2ff;
  color: #4f46e5;
  display: grid;
  place-items: center;
  font-weight: 900;
}

.rank-name {
  font-weight: 800;
  color: #0f172a;
}

.rank-score {
  font-size: 12px;
  color: #64748b;
  font-weight: 800;
}

.drawer-mask {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.42);
  display: grid;
  place-items: end;
  z-index: 50;
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
  font-weight: 800;
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

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 999px;
}

.avatar.small {
  width: 34px;
  height: 34px;
  font-size: 12px;
}

.avatar.placeholder {
  display: grid;
  place-items: center;
  background: #eef2ff;
  color: #4f46e5;
  font-weight: 900;
}

@media (max-width: 1024px) {
  .main-inner {
    grid-template-columns: 1fr;
  }
  .right {
    display: none;
  }
  .stats {
    grid-template-columns: repeat(2, 1fr);
  }
  .stat {
    border-right: 0;
    padding-right: 0;
  }
  .hero-inner {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
