<template>
  <div class="reports">
    <div class="head">
      <div class="title">举报审计</div>
      <div class="tools">
        <select v-model="status" class="select">
          <option value="pending">待处理</option>
          <option value="processed">已处理</option>
        </select>
        <button class="ghost" type="button" @click="refresh">刷新</button>
      </div>
    </div>

    <div class="panel">
      <div class="panel-head">
        <div class="muted">共 {{ total }} 条</div>
        <div class="pager">
          <button
            class="pager-btn"
            type="button"
            :disabled="pageNum <= 1 || loading"
            @click="prevPage"
          >
            上一页
          </button>
          <div class="pager-text">第 {{ pageNum }} 页</div>
          <button
            class="pager-btn"
            type="button"
            :disabled="isLastPage || loading"
            @click="nextPage"
          >
            下一页
          </button>
        </div>
      </div>

      <div v-if="loading" class="muted pad">加载中...</div>
      <div v-else-if="list.length === 0" class="muted pad">暂无数据</div>
      <div v-else class="rows">
        <div v-for="r in list" :key="r.reportId" class="row">
          <div class="row-main">
            <div class="row-title">
              #{{ r.reportId }} 举报 {{ r.targetType }} {{ r.targetId }}
              <span class="pill" :class="{ done: r.status === 'processed' }">{{
                r.status === "processed" ? "已处理" : "待处理"
              }}</span>
            </div>
            <div class="row-sub">
              <span>{{ r.reason || "无原因" }}</span>
              <span v-if="r.createdAt" class="dot">·</span>
              <span v-if="r.createdAt">{{ formatTime(r.createdAt) }}</span>
            </div>
            <div v-if="r.detail" class="row-detail">{{ r.detail }}</div>
            <div v-if="r.status === 'processed'" class="row-processed">
              <span>动作：{{ r.action || "-" }}</span>
              <span class="dot">·</span>
              <span>备注：{{ r.processRemark || "-" }}</span>
            </div>
          </div>

          <div class="row-actions">
            <template v-if="r.status === 'pending'">
              <button
                class="btn ok"
                type="button"
                :disabled="actingId === r.reportId"
                @click="process(r, 'close')"
              >
                关闭
              </button>
              <button
                class="btn warn"
                type="button"
                :disabled="actingId === r.reportId || r.targetType !== 'post'"
                @click="process(r, 'hidePost')"
              >
                隐藏帖
              </button>
              <button
                class="btn danger"
                type="button"
                :disabled="actingId === r.reportId || r.targetType !== 'post'"
                @click="process(r, 'deletePost')"
              >
                删帖
              </button>
              <button
                class="btn danger"
                type="button"
                :disabled="
                  actingId === r.reportId || r.targetType !== 'comment'
                "
                @click="process(r, 'deleteComment')"
              >
                删评
              </button>
              <button
                class="btn"
                type="button"
                :disabled="actingId === r.reportId"
                @click="process(r, 'banUser')"
              >
                封禁
              </button>
            </template>
            <template v-else>
              <button class="btn" type="button" @click="status = 'pending'">
                查看待处理
              </button>
            </template>
          </div>
        </div>
      </div>
    </div>

    <div v-if="error" class="error">{{ error }}</div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { adminApi } from "@/apis";
import type { AdminReportItem } from "@/apis/modules/admin";

const status = ref<"pending" | "processed">("pending");
const pageNum = ref(1);
const pageSize = ref(10);
const total = ref(0);
const list = ref<AdminReportItem[]>([]);
const loading = ref(false);
const actingId = ref<number | null>(null);
const error = ref("");

const isLastPage = computed(
  () => pageNum.value * pageSize.value >= total.value,
);

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

const fetchList = async () => {
  error.value = "";
  loading.value = true;
  try {
    const data = await adminApi.getReports({
      pageNum: pageNum.value,
      pageSize: pageSize.value,
      status: status.value,
    });
    list.value = data.list || [];
    total.value = data.total || 0;
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : "加载失败";
  } finally {
    loading.value = false;
  }
};

const refresh = async () => {
  await fetchList();
};

const prevPage = async () => {
  if (pageNum.value <= 1) return;
  pageNum.value -= 1;
  await fetchList();
};

const nextPage = async () => {
  if (isLastPage.value) return;
  pageNum.value += 1;
  await fetchList();
};

const process = async (
  report: AdminReportItem,
  action: "close" | "deletePost" | "deleteComment" | "hidePost" | "banUser",
) => {
  const remark = window.prompt("处理备注（可选）") || undefined;
  const payload: {
    action: "close" | "deletePost" | "deleteComment" | "hidePost" | "banUser";
    processRemark?: string;
    durationSeconds?: number;
    banUntil?: string;
  } = { action, processRemark: remark };

  if (action === "banUser") {
    const s = window.prompt("封禁时长（秒，可选；不填则 1 天）");
    const durationSeconds = s ? Number(s) : 86400;
    payload.durationSeconds = Number.isFinite(durationSeconds)
      ? durationSeconds
      : 86400;
  }

  const ok = window.confirm(`确认执行动作：${action} ?`);
  if (!ok) return;

  actingId.value = report.reportId;
  error.value = "";
  try {
    await adminApi.processReport(report.reportId, payload);
    await fetchList();
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : "操作失败";
  } finally {
    actingId.value = null;
  }
};

watch(
  () => status.value,
  async () => {
    pageNum.value = 1;
    await fetchList();
  },
);

onMounted(fetchList);
</script>

<style scoped>
.reports {
  display: grid;
  gap: 12px;
}

.head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.title {
  font-weight: 900;
  color: #0f172a;
}

.tools {
  display: inline-flex;
  gap: 10px;
  align-items: center;
}

.select {
  height: 34px;
  padding: 0 12px;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  background: #fff;
  color: #334155;
}

.ghost {
  height: 34px;
  padding: 0 12px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 12px;
  cursor: pointer;
  color: #334155;
}

.panel {
  background: #fff;
  border: 1px solid #eef0f5;
  border-radius: 16px;
  box-shadow: 0 6px 18px rgba(15, 23, 42, 0.03);
}

.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 14px;
  border-bottom: 1px solid #eef0f5;
}

.pager {
  display: inline-flex;
  gap: 10px;
  align-items: center;
}

.pager-btn {
  height: 30px;
  padding: 0 10px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 12px;
  cursor: pointer;
}

.pager-text {
  font-size: 12px;
  color: #64748b;
  font-weight: 800;
}

.rows {
  padding: 10px 14px 14px;
  display: grid;
  gap: 10px;
}

.row {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 12px;
  align-items: start;
  border: 1px solid #eef0f5;
  border-radius: 14px;
  padding: 12px;
}

.row-title {
  font-weight: 900;
  color: #0f172a;
  display: inline-flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}

.pill {
  height: 22px;
  padding: 0 10px;
  border-radius: 999px;
  background: rgba(234, 88, 12, 0.12);
  color: #ea580c;
  font-size: 12px;
  font-weight: 900;
  display: inline-flex;
  align-items: center;
}

.pill.done {
  background: rgba(22, 163, 74, 0.12);
  color: #16a34a;
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

.row-detail {
  margin-top: 10px;
  color: #334155;
  white-space: pre-wrap;
  line-height: 1.7;
  background: #f8fafc;
  border: 1px solid #eef0f5;
  border-radius: 12px;
  padding: 10px 12px;
}

.row-processed {
  margin-top: 10px;
  font-size: 12px;
  color: #64748b;
  display: inline-flex;
  gap: 6px;
  align-items: center;
  flex-wrap: wrap;
}

.row-actions {
  display: grid;
  gap: 8px;
  justify-items: end;
}

.btn {
  height: 30px;
  padding: 0 10px;
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 12px;
  cursor: pointer;
  min-width: 72px;
}

.btn.ok {
  border-color: rgba(22, 163, 74, 0.25);
  color: #16a34a;
}

.btn.warn {
  border-color: rgba(234, 88, 12, 0.25);
  color: #ea580c;
}

.btn.danger {
  border-color: rgba(220, 38, 38, 0.25);
  color: #dc2626;
}

.muted {
  color: #94a3b8;
}

.pad {
  padding: 14px;
}

.error {
  color: #dc2626;
  font-weight: 800;
}

@media (max-width: 920px) {
  .row {
    grid-template-columns: 1fr;
  }
  .row-actions {
    grid-auto-flow: column;
    justify-content: flex-start;
    justify-items: start;
    flex-wrap: wrap;
  }
}
</style>
