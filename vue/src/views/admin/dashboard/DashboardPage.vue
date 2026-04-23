<template>
  <div class="dashboard">
    <div class="cards">
      <div class="card">
        <div class="card-head">
          <div class="card-title">待审帖子</div>
          <div class="badge" :class="{ warn: pendingPostsCount > 0 }">
            {{ pendingPostsCount }}
          </div>
        </div>
        <div class="card-sub">需要审核处理的帖子数量</div>
      </div>
      <div class="card">
        <div class="card-head">
          <div class="card-title">待处理举报</div>
          <div class="badge" :class="{ warn: pendingReportsCount > 0 }">
            {{ pendingReportsCount }}
          </div>
        </div>
        <div class="card-sub">需要管理员处理的举报单</div>
      </div>
      <div class="card">
        <div class="card-head">
          <div class="card-title">系统健康</div>
          <div class="badge" :class="{ ok: health?.ok }">
            {{ health?.ok ? "正常" : "异常" }}
          </div>
        </div>
        <div class="card-sub">MySQL / Redis 连通性</div>
      </div>
      <div class="card">
        <div class="card-head">
          <div class="card-title">快捷入口</div>
          <button class="link" type="button" @click="refresh">刷新</button>
        </div>
        <div class="quick">
          <button
            class="quick-btn"
            type="button"
            @click="router.push({ name: 'admin-posts-pending' })"
          >
            去审核
          </button>
          <button
            class="quick-btn"
            type="button"
            @click="router.push({ name: 'admin-reports' })"
          >
            看举报
          </button>
        </div>
      </div>
    </div>

    <div class="grid">
      <div class="panel chart">
        <div class="panel-head">
          <div class="panel-title">活跃趋势分析（占位）</div>
          <div class="muted">最近 7 天</div>
        </div>
        <div class="chart-body">
          <svg viewBox="0 0 640 180" class="chart-svg" aria-hidden="true">
            <path
              d="M10 150 C 80 120, 120 160, 180 120 C 240 80, 280 120, 340 90 C 400 60, 450 95, 520 70 C 560 55, 600 70, 630 62"
              fill="none"
              stroke="#4f46e5"
              stroke-width="4"
              stroke-linecap="round"
            />
          </svg>
        </div>
      </div>

      <div class="panel chart">
        <div class="panel-head">
          <div class="panel-title">内容分类（占位）</div>
          <div class="muted">按帖子类型</div>
        </div>
        <div class="chart-body center">
          <svg viewBox="0 0 120 120" class="donut" aria-hidden="true">
            <circle
              cx="60"
              cy="60"
              r="46"
              fill="none"
              stroke="#e2e8f0"
              stroke-width="14"
            />
            <circle
              cx="60"
              cy="60"
              r="46"
              fill="none"
              stroke="#4f46e5"
              stroke-width="14"
              stroke-dasharray="120 289"
              stroke-linecap="round"
            />
            <circle
              cx="60"
              cy="60"
              r="46"
              fill="none"
              stroke="#f59e0b"
              stroke-width="14"
              stroke-dasharray="80 329"
              stroke-dashoffset="-120"
              stroke-linecap="round"
            />
            <circle
              cx="60"
              cy="60"
              r="46"
              fill="none"
              stroke="#ec4899"
              stroke-width="14"
              stroke-dasharray="60 349"
              stroke-dashoffset="-200"
              stroke-linecap="round"
            />
          </svg>
        </div>
      </div>
    </div>

    <div class="panel table">
      <div class="panel-head">
        <div class="panel-title">待审核帖子（快捷）</div>
        <button
          class="link"
          type="button"
          @click="router.push({ name: 'admin-posts-pending' })"
        >
          查看全部
        </button>
      </div>
      <div v-if="loading" class="muted pad">加载中...</div>
      <div v-else-if="pendingPosts.length === 0" class="muted pad">
        暂无待审帖子
      </div>
      <div v-else class="rows">
        <div v-for="p in pendingPosts" :key="p.postId" class="row">
          <div class="row-main">
            <div class="row-title">{{ p.title }}</div>
            <div class="row-sub">
              <span>postId: {{ p.postId }}</span>
              <span class="dot">·</span>
              <span>forumId: {{ p.forumId }}</span>
              <span v-if="p.createdAt" class="dot">·</span>
              <span v-if="p.createdAt">{{ formatTime(p.createdAt) }}</span>
            </div>
          </div>
          <div class="row-actions">
            <button
              class="btn ok"
              type="button"
              @click="review(p.postId, 'approve')"
            >
              通过
            </button>
            <button
              class="btn"
              type="button"
              @click="review(p.postId, 'reject')"
            >
              驳回
            </button>
            <button
              class="btn warn"
              type="button"
              @click="review(p.postId, 'hide')"
            >
              隐藏
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="panel table">
      <div class="panel-head">
        <div class="panel-title">待处理举报（快捷）</div>
        <button
          class="link"
          type="button"
          @click="router.push({ name: 'admin-reports' })"
        >
          查看全部
        </button>
      </div>
      <div v-if="loadingReports" class="muted pad">加载中...</div>
      <div v-else-if="pendingReports.length === 0" class="muted pad">
        暂无待处理举报
      </div>
      <div v-else class="rows">
        <div v-for="r in pendingReports" :key="r.reportId" class="row">
          <div class="row-main">
            <div class="row-title">
              #{{ r.reportId }} 举报 {{ r.targetType }} {{ r.targetId }}
            </div>
            <div class="row-sub">
              <span>{{ r.reason || "无原因" }}</span>
              <span v-if="r.createdAt" class="dot">·</span>
              <span v-if="r.createdAt">{{ formatTime(r.createdAt) }}</span>
            </div>
          </div>
          <div class="row-actions">
            <button
              class="btn ok"
              type="button"
              @click="process(r.reportId, 'close')"
            >
              关闭
            </button>
            <button
              class="btn warn"
              type="button"
              @click="process(r.reportId, 'hidePost')"
              :disabled="r.targetType !== 'post'"
            >
              隐藏帖
            </button>
            <button
              class="btn"
              type="button"
              @click="process(r.reportId, 'deletePost')"
              :disabled="r.targetType !== 'post'"
            >
              删帖
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { adminApi } from "@/apis";
import type { Post } from "@/types/api";
import type { AdminReportItem, HealthResponse } from "@/apis/modules/admin";

const router = useRouter();

const loading = ref(false);
const pendingPosts = ref<Post[]>([]);
const postsTotal = ref(0);

const loadingReports = ref(false);
const pendingReports = ref<AdminReportItem[]>([]);
const reportsTotal = ref(0);

const health = ref<HealthResponse | null>(null);

const pendingPostsCount = computed(() => postsTotal.value);
const pendingReportsCount = computed(() => reportsTotal.value);

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

const refresh = async () => {
  loading.value = true;
  loadingReports.value = true;
  try {
    const [pendingPostsData, pendingReportsData, healthData] =
      await Promise.all([
        adminApi.getPendingPosts({ pageNum: 1, pageSize: 5 }),
        adminApi.getReports({ pageNum: 1, pageSize: 5, status: "pending" }),
        adminApi.getHealth(),
      ]);
    pendingPosts.value = pendingPostsData.list || [];
    postsTotal.value = pendingPostsData.total || 0;
    pendingReports.value = pendingReportsData.list || [];
    reportsTotal.value = pendingReportsData.total || 0;
    health.value = healthData;
  } finally {
    loading.value = false;
    loadingReports.value = false;
  }
};

const review = async (
  postId: number,
  action: "approve" | "reject" | "hide",
) => {
  const remark = window.prompt("处理备注（可选）") || undefined;
  await adminApi.reviewPost(postId, { action, reviewRemark: remark });
  await refresh();
};

const process = async (
  reportId: number,
  action: "close" | "deletePost" | "hidePost",
) => {
  const remark = window.prompt("处理备注（可选）") || undefined;
  await adminApi.processReport(reportId, { action, processRemark: remark });
  await refresh();
};

onMounted(refresh);
</script>

<style scoped>
.dashboard {
  display: grid;
  gap: 16px;
}

.cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}

.card {
  background: #fff;
  border: 1px solid #eef0f5;
  border-radius: 16px;
  padding: 14px;
  box-shadow: 0 6px 18px rgba(15, 23, 42, 0.03);
}

.card-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.card-title {
  font-weight: 900;
  color: #0f172a;
}

.card-sub {
  margin-top: 8px;
  font-size: 12px;
  color: #94a3b8;
}

.badge {
  height: 26px;
  padding: 0 10px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 999px;
  background: rgba(148, 163, 184, 0.16);
  color: #475569;
  font-weight: 900;
}

.badge.warn {
  background: rgba(234, 88, 12, 0.12);
  color: #ea580c;
}

.badge.ok {
  background: rgba(22, 163, 74, 0.12);
  color: #16a34a;
}

.quick {
  margin-top: 12px;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.quick-btn {
  height: 32px;
  padding: 0 12px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 12px;
  cursor: pointer;
}

.grid {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 16px;
}

.panel {
  background: #fff;
  border: 1px solid #eef0f5;
  border-radius: 16px;
  box-shadow: 0 6px 18px rgba(15, 23, 42, 0.03);
}

.panel-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  padding: 12px 14px;
  border-bottom: 1px solid #eef0f5;
}

.panel-title {
  font-weight: 900;
  color: #0f172a;
}

.chart-body {
  padding: 14px;
}

.chart-body.center {
  display: grid;
  place-items: center;
}

.chart-svg {
  width: 100%;
  height: 200px;
}

.donut {
  width: 140px;
  height: 140px;
}

.table .panel-head {
  border-bottom: 1px solid #eef0f5;
}

.rows {
  padding: 8px 14px 14px;
  display: grid;
  gap: 10px;
}

.row {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 12px;
  align-items: center;
  border: 1px solid #eef0f5;
  border-radius: 14px;
  padding: 10px 12px;
}

.row-title {
  font-weight: 800;
  color: #0f172a;
}

.row-sub {
  margin-top: 6px;
  font-size: 12px;
  color: #94a3b8;
  display: inline-flex;
  gap: 6px;
  align-items: center;
  flex-wrap: wrap;
}

.dot {
  color: #cbd5e1;
}

.row-actions {
  display: inline-flex;
  gap: 8px;
  align-items: center;
}

.btn {
  height: 30px;
  padding: 0 10px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 12px;
  cursor: pointer;
}

.btn.ok {
  border-color: rgba(22, 163, 74, 0.25);
  color: #16a34a;
}

.btn.warn {
  border-color: rgba(234, 88, 12, 0.25);
  color: #ea580c;
}

.muted {
  color: #94a3b8;
}

.pad {
  padding: 14px;
}

.link {
  border: 0;
  background: transparent;
  color: #4f46e5;
  cursor: pointer;
  padding: 0;
  font-weight: 800;
}

@media (max-width: 1100px) {
  .cards {
    grid-template-columns: repeat(2, 1fr);
  }
  .grid {
    grid-template-columns: 1fr;
  }
}
</style>
