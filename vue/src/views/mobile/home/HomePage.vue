<template>
  <div class="home">
    <header class="topbar">
      <div class="brand" @click="goHome">
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
        <button class="nav-btn" type="button" @click="goHome">首页</button>
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
          <div class="panel-title">贴吧列表</div>
          <div v-if="forumsLoading" class="muted">加载中...</div>
          <div v-else class="forum-list">
            <button
              v-for="f in forums"
              :key="f.forumId"
              class="forum-item"
              type="button"
              @click="toForum(f.forumId)"
            >
              <div class="forum-left">
                <div class="forum-icon">
                  <div class="forum-icon-fallback">
                    {{ (f.name || "吧").slice(0, 1) }}
                  </div>
                </div>
                <div class="forum-text">
                  <div class="forum-name">{{ f.name }}</div>
                  <div class="forum-sub">
                    {{ f.description || `ID: ${f.forumId}` }}
                  </div>
                </div>
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
            <button class="ghost" type="button" @click="openCreateForum">
              创建贴吧
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
                <button
                  class="post-btn"
                  type="button"
                  @click="toggleComments(p.postId)"
                >
                  {{ isCommentsExpanded(p.postId) ? "收起评论" : "评论"
                  }}<span v-if="typeof p.commentCount === 'number'">
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
                <button
                  class="post-btn"
                  type="button"
                  @click="openReport('post', p.postId)"
                >
                  举报
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
                      <div class="comment-actions">
                        <button
                          class="comment-btn"
                          type="button"
                          @click="openReply(p.postId, it.comment.commentId)"
                        >
                          回复
                        </button>
                        <button
                          class="comment-btn"
                          type="button"
                          @click="openReport('comment', it.comment.commentId)"
                        >
                          举报
                        </button>
                      </div>

                      <div
                        v-if="
                          replyParentByPostId[p.postId] === it.comment.commentId
                        "
                        class="reply-box"
                      >
                        <textarea
                          v-model.trim="replyContentByPostId[p.postId]"
                          class="reply-input"
                          placeholder="写下你的回复..."
                        />
                        <div class="reply-actions">
                          <button
                            class="comment-send"
                            type="button"
                            :disabled="
                              sendingCommentId === it.comment.commentId
                            "
                            @click="submitReply(p.postId)"
                          >
                            {{
                              sendingCommentId === it.comment.commentId
                                ? "发送中..."
                                : "发送"
                            }}
                          </button>
                          <button
                            class="comment-cancel"
                            type="button"
                            @click="cancelReply(p.postId)"
                          >
                            取消
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="add-comment">
                  <input
                    v-model.trim="newCommentByPostId[p.postId]"
                    class="add-input"
                    placeholder="新增评论..."
                    @focus="prepareComment(p.postId)"
                  />
                  <button
                    class="add-send"
                    type="button"
                    :disabled="sendingPostId === p.postId"
                    @click="submitNewComment(p.postId)"
                  >
                    {{ sendingPostId === p.postId ? "发送中..." : "发送" }}
                  </button>
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
            <button class="link" type="button" @click="goHome">查看全部</button>
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
          <button class="ghost" type="button" @click="goRegister">注册</button>
          <p v-if="loginError" class="error">{{ loginError }}</p>
        </form>
      </div>
    </div>

    <div
      v-if="createForumOpen"
      class="modal-mask"
      @click.self="closeCreateForum"
    >
      <div class="modal">
        <div class="modal-title">创建贴吧</div>
        <form class="modal-form" @submit.prevent="onCreateForumSubmit">
          <label class="field">
            <span>名称</span>
            <input
              v-model.trim="createForumName"
              placeholder="请输入贴吧名称"
            />
          </label>
          <label class="field">
            <span>简介</span>
            <textarea
              v-model.trim="createForumDescription"
              placeholder="一句话介绍这个贴吧（可选）"
            />
          </label>
          <label class="field">
            <span>封面</span>
            <div class="cover-box">
              <div class="cover-preview">
                <img
                  v-if="createForumCoverUrl"
                  class="cover-img"
                  :src="createForumCoverUrl"
                  alt="cover"
                />
                <div v-else class="cover-placeholder">未选择图片</div>
              </div>
              <div class="cover-actions">
                <input
                  ref="coverFileInput"
                  class="cover-file"
                  type="file"
                  accept="image/*"
                  @change="onCoverFileChange"
                />
                <button
                  class="ghost"
                  type="button"
                  :disabled="createForumCoverUploading"
                  @click="pickCoverFile"
                >
                  {{ createForumCoverUploading ? "上传中..." : "上传图片" }}
                </button>
                <button
                  class="ghost"
                  type="button"
                  :disabled="createForumCoverUploading || !createForumCoverUrl"
                  @click="clearCover"
                >
                  清除
                </button>
              </div>
              <p v-if="createForumCoverError" class="error">
                {{ createForumCoverError }}
              </p>
            </div>
          </label>

          <button
            class="primary"
            type="submit"
            :disabled="createForumLoading || !createForumName.trim()"
          >
            {{ createForumLoading ? "创建中..." : "创建" }}
          </button>
          <button class="ghost" type="button" @click="closeCreateForum">
            取消
          </button>
          <p v-if="createForumError" class="error">{{ createForumError }}</p>
        </form>
      </div>
    </div>

    <div v-if="reportOpen" class="modal-mask" @click.self="closeReport">
      <div class="modal">
        <div class="modal-title">举报</div>
        <form class="modal-form" @submit.prevent="submitReport">
          <label class="field">
            <span>原因</span>
            <input v-model.trim="reportReason" placeholder="请输入举报原因" />
          </label>
          <label class="field">
            <span>补充描述（可选）</span>
            <textarea
              v-model.trim="reportDetail"
              placeholder="补充说明（最多 1000 字）"
            />
          </label>
          <button
            class="primary"
            type="submit"
            :disabled="reportLoading || !reportReason.trim()"
          >
            {{ reportLoading ? "提交中..." : "确认举报" }}
          </button>
          <button class="ghost" type="button" @click="closeReport">取消</button>
          <p v-if="reportError" class="error">{{ reportError }}</p>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import type { CommentItem, Forum, Post, ReportTargetType } from "@/types/api";
import { feedApi, followsApi, forumsApi, reportsApi, usersApi } from "@/apis";
import { useAuthStore } from "@/stores/auth";

const router = useRouter();
const route = useRoute();
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

const createForumOpen = ref(false);
const createForumName = ref("");
const createForumDescription = ref("");
const createForumCoverUrl = ref("");
const createForumLoading = ref(false);
const createForumError = ref("");
const coverFileInput = ref<HTMLInputElement | null>(null);
const createForumCoverUploading = ref(false);
const createForumCoverError = ref("");

const reportOpen = ref(false);
const reportTargetType = ref<ReportTargetType>("post");
const reportTargetId = ref(0);
const reportReason = ref("");
const reportDetail = ref("");
const reportLoading = ref(false);
const reportError = ref("");

type CommentThreadItem = {
  comment: CommentItem;
  depth: number;
  replyTo?: CommentItem;
};
const commentsByPostId = ref<Record<number, CommentItem[]>>({});
const commentsLoadingByPostId = ref<Record<number, boolean>>({});
const commentsExpandedByPostId = ref<Record<number, boolean>>({});
const newCommentByPostId = ref<Record<number, string>>({});
const replyParentByPostId = ref<Record<number, number>>({});
const replyContentByPostId = ref<Record<number, string>>({});
const sendingPostId = ref<number | null>(null);
const sendingCommentId = ref<number | null>(null);

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

const goHome = async () => {
  if (route.name === "mobile-home") {
    await Promise.all([loadForums(), loadFollows(), loadPosts()]);
    return;
  }
  router.push("/m/home");
};

const goRegister = () => {
  loginOpen.value = false;
  router.push({ name: "mobile-register" });
};

const closeLogin = () => {
  loginOpen.value = false;
};

const openCreateForum = () => {
  if (!auth.isAuthed) {
    openLogin();
    return;
  }
  createForumOpen.value = true;
  createForumError.value = "";
  createForumCoverError.value = "";
};

const closeCreateForum = () => {
  createForumOpen.value = false;
  createForumLoading.value = false;
  createForumError.value = "";
  createForumCoverUploading.value = false;
  createForumCoverError.value = "";
};

const openReport = (targetType: ReportTargetType, targetId: number) => {
  if (!auth.isAuthed) {
    openLogin();
    return;
  }
  reportTargetType.value = targetType;
  reportTargetId.value = targetId;
  reportReason.value = "";
  reportDetail.value = "";
  reportError.value = "";
  reportOpen.value = true;
};

const closeReport = () => {
  reportOpen.value = false;
};

const submitReport = async () => {
  const reason = reportReason.value.trim();
  if (!reason) {
    reportError.value = "请输入举报原因";
    return;
  }
  if (reason.length > 64) {
    reportError.value = "原因最多 64 字";
    return;
  }
  const detail = reportDetail.value.trim();
  if (detail.length > 1000) {
    reportError.value = "补充描述最多 1000 字";
    return;
  }

  const ok = window.confirm("确认提交举报？");
  if (!ok) return;

  reportLoading.value = true;
  reportError.value = "";
  try {
    await reportsApi.createReport({
      targetType: reportTargetType.value,
      targetId: reportTargetId.value,
      reason,
      detail: detail ? detail : undefined,
    });
    reportOpen.value = false;
    window.alert("已提交举报");
  } catch (e: unknown) {
    reportError.value = e instanceof Error ? e.message : "提交失败";
  } finally {
    reportLoading.value = false;
  }
};

const pickCoverFile = () => {
  if (!auth.isAuthed) {
    openLogin();
    return;
  }
  createForumCoverError.value = "";
  coverFileInput.value?.click();
};

const clearCover = () => {
  createForumCoverUrl.value = "";
  createForumCoverError.value = "";
  if (coverFileInput.value) coverFileInput.value.value = "";
};

const onCoverFileChange = async (e: Event) => {
  if (!auth.isAuthed) {
    openLogin();
    return;
  }
  const input = e.target as HTMLInputElement | null;
  const file = input?.files?.[0];
  if (!file) return;

  createForumCoverUploading.value = true;
  createForumCoverError.value = "";
  try {
    const data = await usersApi.upload({ file, scene: "forum_cover" });
    createForumCoverUrl.value = data.url;
  } catch (err: unknown) {
    createForumCoverError.value =
      err instanceof Error ? err.message : "上传失败";
    if (coverFileInput.value) coverFileInput.value.value = "";
  } finally {
    createForumCoverUploading.value = false;
  }
};

const onCreateForumSubmit = async () => {
  if (!auth.isAuthed) {
    openLogin();
    return;
  }
  const name = createForumName.value.trim();
  if (!name) {
    createForumError.value = "请输入贴吧名称";
    return;
  }

  createForumLoading.value = true;
  createForumError.value = "";
  try {
    await forumsApi.createForum({
      name,
      description: createForumDescription.value.trim() || undefined,
      coverUrl: createForumCoverUrl.value.trim() || undefined,
    });
    createForumOpen.value = false;
    createForumName.value = "";
    createForumDescription.value = "";
    createForumCoverUrl.value = "";
    await Promise.all([loadForums(), loadFollows(), loadPosts()]);
  } catch (e: unknown) {
    createForumError.value = e instanceof Error ? e.message : "创建失败";
  } finally {
    createForumLoading.value = false;
  }
};

const onLoginSubmit = async () => {
  loginError.value = "";
  loginLoading.value = true;
  try {
    await auth.login(loginAccount.value, loginPassword.value);
    loginOpen.value = false;
    loginAccount.value = "";
    loginPassword.value = "";
    await Promise.all([loadFollows(), loadPosts()]);
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
    for (const item of data.list || []) {
      if (item.active !== false) ids.add(item.targetId);
    }
    myFollowForumIdSet.value = ids;
  } catch (e) {
    myFollowForumIdSet.value = new Set();
  }
};

const loadPosts = async () => {
  postsLoading.value = true;
  try {
    const data = await feedApi.getPublicPosts({ pageNum: 1, pageSize: 30 });
    const forumMap = new Map<
      number,
      { forumId: number; name: string; coverUrl?: string }
    >();
    for (const f of forums.value)
      forumMap.set(f.forumId, {
        forumId: f.forumId,
        name: f.name,
        coverUrl: f.coverUrl,
      });
    posts.value = (data.list || []).map((p) => {
      const maybeForum = forumMap.get(p.forumId);
      return {
        ...p,
        forum: p.forum || maybeForum,
      };
    });
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

const ensureCommentsLoaded = async (postId: number, force = false) => {
  if (commentsByPostId.value[postId] && !force) return;
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
  if (!auth.isAuthed) {
    openLogin();
    return;
  }
  const next = !isCommentsExpanded(postId);
  commentsExpandedByPostId.value[postId] = next;
  if (next) await ensureCommentsLoaded(postId);
};

const prepareComment = async (postId: number) => {
  if (!auth.isAuthed) {
    openLogin();
    return;
  }
  commentsExpandedByPostId.value[postId] = true;
  await ensureCommentsLoaded(postId);
};

const submitNewComment = async (postId: number) => {
  if (!auth.isAuthed) {
    openLogin();
    return;
  }
  const text = (newCommentByPostId.value[postId] || "").trim();
  if (!text) return;
  sendingPostId.value = postId;
  try {
    await feedApi.createComment(postId, { content: text });
    newCommentByPostId.value[postId] = "";
    await prepareComment(postId);
    await ensureCommentsLoaded(postId, true);
  } finally {
    sendingPostId.value = null;
  }
};

const openReply = async (postId: number, commentId: number) => {
  if (!auth.isAuthed) {
    openLogin();
    return;
  }
  await prepareComment(postId);
  replyParentByPostId.value[postId] = commentId;
  if (!replyContentByPostId.value[postId])
    replyContentByPostId.value[postId] = "";
};

const cancelReply = (postId: number) => {
  replyParentByPostId.value[postId] = 0 as unknown as number;
  replyContentByPostId.value[postId] = "";
};

const submitReply = async (postId: number) => {
  if (!auth.isAuthed) {
    openLogin();
    return;
  }
  const parentId = replyParentByPostId.value[postId];
  if (!parentId) return;
  const text = (replyContentByPostId.value[postId] || "").trim();
  if (!text) return;
  sendingCommentId.value = parentId;
  try {
    await feedApi.createComment(postId, {
      content: text,
      parentCommentId: parentId,
    });
    cancelReply(postId);
    await ensureCommentsLoaded(postId, true);
  } finally {
    sendingCommentId.value = null;
  }
};

const onClickNotifications = () => {
  if (!auth.isAuthed) {
    openLogin();
    return;
  }
  router.push({ name: "mobile-messages" });
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
  background: transparent;
}

.topbar {
  position: sticky;
  top: 0;
  z-index: 10;
  min-height: 68px;
  display: grid;
  grid-template-columns: 220px 1fr auto 260px;
  gap: 14px;
  align-items: center;
  padding: 10px 18px;
  background: rgba(255, 255, 255, 0.82);
  backdrop-filter: blur(18px);
  border-bottom: 1px solid rgba(226, 232, 240, 0.9);
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.04);
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
  font-weight: 900;
  letter-spacing: 0.5px;
  color: #111827;
}

.search {
  display: flex;
  align-items: center;
}

.search-input {
  width: 100%;
  height: 42px;
  padding: 0 16px;
  border: 1px solid rgba(226, 232, 240, 0.95);
  background: rgba(255, 255, 255, 0.92);
  border-radius: 14px;
  outline: none;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.8);
}

.actions {
  display: inline-flex;
  gap: 10px;
  align-items: center;
}

.nav-btn {
  height: 38px;
  padding: 0 14px;
  border: 1px solid rgba(226, 232, 240, 0.95);
  background: rgba(255, 255, 255, 0.88);
  border-radius: 12px;
  cursor: pointer;
  color: #334155;
  font-weight: 700;
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
  padding: 22px 18px 32px;
  display: grid;
  grid-template-columns: 260px 1fr 320px;
  gap: 18px;
  align-items: start;
}

.panel {
  background: var(--card-bg);
  border: 1px solid var(--card-border);
  border-radius: 20px;
  box-shadow: var(--card-shadow);
  overflow: hidden;
  backdrop-filter: blur(10px);
}

.panel-title {
  font-weight: 900;
  padding: 16px 16px 13px;
  color: #111827;
  border-bottom: 1px solid rgba(241, 245, 249, 0.95);
}

.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 16px 12px;
  border-bottom: 1px solid rgba(241, 245, 249, 0.95);
}

.panel-footer {
  padding: 14px 16px 16px;
  display: flex;
  justify-content: center;
}

.ghost {
  height: 36px;
  padding: 0 14px;
  border: 1px dashed rgba(203, 213, 225, 0.95);
  background: rgba(255, 255, 255, 0.88);
  border-radius: 12px;
  cursor: pointer;
  color: #475569;
  font-weight: 700;
}

.forum-list {
  display: grid;
  gap: 10px;
  padding: 2px 14px 14px;
}

.forum-item {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 10px;
  align-items: center;
  padding: 12px 12px;
  border-radius: 14px;
  border: 1px solid rgba(238, 240, 245, 0.92);
  background: linear-gradient(180deg, #ffffff 0%, #fbfdff 100%);
  cursor: pointer;
  text-align: left;
  transition:
    box-shadow 0.2s ease,
    border-color 0.2s ease,
    transform 0.2s ease;
}

.forum-item:hover {
  border-color: rgba(79, 70, 229, 0.24);
  box-shadow: 0 14px 32px rgba(15, 23, 42, 0.08);
  transform: translateY(-2px);
}

.forum-left {
  display: flex;
  gap: 10px;
  align-items: center;
  min-width: 0;
}

.forum-icon {
  width: 38px;
  height: 38px;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid #e2e8f0;
  background: #fff;
  flex: 0 0 auto;
}

.forum-icon-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.forum-icon-fallback {
  width: 100%;
  height: 100%;
  display: grid;
  place-items: center;
  background: #eef2ff;
  color: #4f46e5;
  font-weight: 900;
}

.forum-text {
  min-width: 0;
}

.forum-name {
  font-weight: 600;
  color: #0f172a;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.forum-sub {
  margin-top: 2px;
  font-size: 12px;
  color: #64748b;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.forum-ops {
  display: flex;
  justify-content: flex-end;
}

.follow {
  height: 32px;
  padding: 0 12px;
  border: 1px solid #4f46e5;
  color: #4f46e5;
  background: rgba(255, 255, 255, 0.92);
  border-radius: 999px;
  cursor: pointer;
  font-weight: 800;
  box-shadow: 0 4px 12px rgba(79, 70, 229, 0.08);
}

.center {
  display: grid;
  gap: 16px;
}

.composer {
  display: grid;
  grid-template-columns: 44px 1fr auto;
  gap: 14px;
  align-items: center;
  padding: 16px;
  background:
    radial-gradient(
      260px 120px at 0% 0%,
      rgba(79, 70, 229, 0.06),
      transparent 60%
    ),
    rgba(255, 255, 255, 0.9);
}

.composer-input {
  height: 44px;
  border: 1px solid rgba(226, 232, 240, 0.96);
  background: linear-gradient(180deg, #fbfdff 0%, #f8fafc 100%);
  border-radius: 14px;
  text-align: left;
  padding: 0 16px;
  cursor: pointer;
  color: #64748b;
  font-weight: 600;
}

.icon-btn {
  height: 38px;
  padding: 0 14px;
  border: 0;
  background: linear-gradient(135deg, #4f46e5 0%, #6366f1 100%);
  color: #fff;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 800;
  box-shadow: 0 10px 22px rgba(79, 70, 229, 0.22);
}

.tabs {
  display: flex;
  gap: 10px;
  padding: 12px 16px 0;
}

.tab {
  height: 36px;
  padding: 0 14px;
  border: 0;
  background: transparent;
  border-bottom: 2px solid transparent;
  cursor: pointer;
  color: #475569;
  font-weight: 700;
}

.tab.active {
  color: #4f46e5;
  border-color: #4f46e5;
  font-weight: 700;
}

.post-list {
  display: grid;
  gap: 16px;
  padding: 16px;
}

.post-card {
  border: 1px solid rgba(238, 240, 245, 0.92);
  border-radius: 18px;
  padding: 16px;
  background: linear-gradient(180deg, #ffffff 0%, #fcfdff 100%);
  transition:
    box-shadow 0.2s ease,
    border-color 0.2s ease;
}

.post-card:hover {
  border-color: rgba(79, 70, 229, 0.18);
  box-shadow: 0 18px 36px rgba(15, 23, 42, 0.08);
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
  background: linear-gradient(180deg, #f8fafc 0%, #ffffff 100%);
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

.comment-actions {
  margin-top: 8px;
  display: flex;
  justify-content: flex-end;
}

.comment-btn {
  height: 28px;
  padding: 0 10px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 999px;
  cursor: pointer;
  font-size: 12px;
  color: #475569;
  font-weight: 800;
}

.reply-box {
  margin-top: 10px;
  border-top: 1px dashed #e2e8f0;
  padding-top: 10px;
  display: grid;
  gap: 10px;
}

.reply-input {
  min-height: 70px;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 10px 12px;
  outline: none;
  resize: vertical;
  font-family: inherit;
  line-height: 1.6;
  background: #fff;
}

.reply-actions {
  display: inline-flex;
  gap: 10px;
  justify-content: flex-end;
  flex-wrap: wrap;
}

.comment-send {
  height: 32px;
  padding: 0 12px;
  border: 0;
  background: #4f46e5;
  color: #fff;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 900;
}

.comment-cancel {
  height: 32px;
  padding: 0 12px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 800;
  color: #334155;
}

.add-comment {
  margin-top: 12px;
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 10px;
  align-items: center;
  border-top: 1px solid #eef0f5;
  padding-top: 12px;
}

.add-input {
  height: 38px;
  padding: 0 12px;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  outline: none;
  background: #fff;
}

.add-send {
  height: 38px;
  padding: 0 12px;
  border: 0;
  background: #4f46e5;
  color: #fff;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 900;
}

@media (max-width: 1100px) {
  .layout {
    grid-template-columns: 240px 1fr;
    grid-template-areas:
      "left center"
      "right right";
  }
  .left {
    grid-area: left;
  }
  .center {
    grid-area: center;
  }
  .right {
    grid-area: right;
  }
}

@media (max-width: 820px) {
  .layout {
    grid-template-columns: 1fr;
    grid-template-areas:
      "center"
      "left"
      "right";
  }
  .topbar {
    grid-template-columns: 1fr;
    height: auto;
    padding: 10px 12px;
    gap: 10px;
  }
  .actions {
    justify-content: space-between;
  }
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
  font-weight: 800;
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

.modal-mask {
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

.field textarea {
  min-height: 84px;
  padding: 10px 12px;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  outline: none;
  resize: vertical;
  font-family: inherit;
  line-height: 1.4;
}

.cover-box {
  display: grid;
  gap: 10px;
}

.cover-preview {
  height: 130px;
  border: 1px dashed #e2e8f0;
  border-radius: 12px;
  background: #f8fafc;
  overflow: hidden;
  display: grid;
  place-items: center;
}

.cover-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.cover-placeholder {
  color: #94a3b8;
  font-size: 12px;
  font-weight: 700;
}

.cover-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.cover-file {
  display: none;
}

.error {
  color: #dc2626;
  margin: 0;
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
