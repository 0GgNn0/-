<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../api'

const experiences = ref([])

onMounted(async () => {
  try {
    experiences.value = await api.experiences()
  } catch (_) { /* empty state */ }
})
</script>

<template>
  <div class="exp-page">
    <div class="reveal" style="--d:0s">
      <router-link to="/" class="page-back">← 返回首页</router-link>
    </div>

    <div class="card glass reveal" style="--d:.1s">
      <div class="kicker mono">Experience</div>
      <h1 class="h1">经历</h1>

      <div v-if="experiences.length" class="timeline">
        <div
          v-for="(e, i) in experiences"
          :key="e.id"
          class="tl-item reveal"
          :style="{ '--d': i * 0.1 + 's' }"
        >
          <div class="tl-dot"></div>
          <div class="tl-period mono">{{ e.period }}</div>
          <div class="tl-title">{{ e.title }}</div>
          <div v-if="e.description" class="tl-desc">{{ e.description }}</div>
        </div>
      </div>
      <div v-else class="empty">暂无经历</div>
    </div>
  </div>
</template>

<style scoped>
.exp-page { max-width: 680px; margin: 0 auto; padding: 28px 24px 72px; }

.page-back {
  display: inline-block;
  font-size: 13px; color: rgba(90, 90, 100, 0.95);
  text-decoration: none; padding: 6px 14px;
  border-radius: 999px; border: 1px solid rgba(0, 0, 0, 0.1);
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(10px);
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
body[data-theme="dark"] .page-back { background: rgba(24, 24, 27, 0.6); color: rgba(190, 190, 200, 0.95); border-color: rgba(255, 255, 255, 0.12); }
.page-back:hover { color: #10b981; border-color: #10b981; transform: translateY(-1px); }

.card { margin-top: 20px; border-radius: 16px; padding: 36px; }
.kicker { font-size: 11px; letter-spacing: 0.14em; text-transform: uppercase; color: #10b981; margin-bottom: 10px; }
.h1 { font-size: clamp(32px, 5vw, 44px); font-weight: 800; letter-spacing: -0.03em; line-height: 1.05; }

.timeline { margin-top: 28px; position: relative; padding-left: 22px; }
.timeline::before {
  content: '';
  position: absolute; left: 5px; top: 6px; bottom: 6px;
  width: 1px;
  background: linear-gradient(#10b981, rgba(0, 0, 0, 0.12));
}
.tl-item { position: relative; padding-bottom: 30px; }
.tl-item:last-child { padding-bottom: 0; }
.tl-dot {
  position: absolute; left: -22px; top: 5px;
  width: 11px; height: 11px;
  border-radius: 50%;
  background: #10b981;
  border: 2px solid #f9fafb;
  box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.15);
}
body[data-theme="dark"] .tl-dot { border-color: #09090b; }
.tl-period { font-size: 12px; color: #10b981; letter-spacing: 0.06em; font-variant-numeric: tabular-nums; }
.tl-title { font-size: 17px; font-weight: 700; margin-top: 5px; letter-spacing: -0.01em; }
.tl-desc { font-size: 14px; line-height: 1.7; margin-top: 5px; color: rgba(90, 90, 100, 0.95); text-wrap: pretty; }
body[data-theme="dark"] .tl-desc { color: rgba(180, 180, 190, 0.95); }
.empty { margin-top: 24px; color: rgba(120, 120, 130, 0.95); font-size: 14px; }
</style>
