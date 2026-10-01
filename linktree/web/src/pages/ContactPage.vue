<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../api'

const links = ref([])

onMounted(async () => {
  try {
    links.value = await api.links()
  } catch (_) { /* empty state */ }
})
</script>

<template>
  <div class="contact-page">
    <div class="reveal" style="--d:0s">
      <router-link to="/" class="page-back">← 返回首页</router-link>
    </div>

    <div class="card glass reveal" style="--d:.1s">
      <div class="kicker mono">Contact</div>
      <h1 class="h1">联系我</h1>
      <p class="sub">对我感兴趣？通过以下方式找到我</p>

      <div class="contact-list">
        <a
          v-for="(l, i) in links"
          :key="l.id"
          :href="l.url"
          target="_blank"
          rel="noopener noreferrer"
          class="contact-item reveal"
          :style="{ '--d': i * 0.08 + 's' }"
        >
          <div class="ico">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/></svg>
          </div>
          <div class="meta">
            <div class="t">{{ l.title }}</div>
            <div class="u mono">{{ l.url }}</div>
          </div>
          <svg class="arrow" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="7" y1="17" x2="17" y2="7"/><polyline points="7 7 17 7 17 17"/></svg>
        </a>
      </div>
      <div v-if="!links.length" class="empty">暂无联系方式</div>
    </div>
  </div>
</template>

<style scoped>
.contact-page { max-width: 560px; margin: 0 auto; padding: 28px 24px 72px; }

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
.sub { font-size: 14px; color: rgba(90, 90, 100, 0.95); margin-top: 10px; }
body[data-theme="dark"] .sub { color: rgba(180, 180, 190, 0.95); }

.contact-list { margin-top: 24px; display: flex; flex-direction: column; gap: 10px; }
.contact-item {
  display: flex; align-items: center; gap: 12px;
  padding: 14px 16px;
  border: 1px solid rgba(0, 0, 0, 0.07);
  border-radius: 12px;
  background: rgba(0, 0, 0, 0.035);
  text-decoration: none;
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
body[data-theme="dark"] .contact-item { background: rgba(255, 255, 255, 0.045); border-color: rgba(255, 255, 255, 0.08); }
.contact-item:hover { border-color: #10b981; transform: translateY(-2px); box-shadow: 0 12px 28px -10px rgba(0, 0, 0, 0.12); }
.contact-item:active { transform: translateY(0) scale(0.99); }
.ico {
  width: 36px; height: 36px;
  border-radius: 8px;
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
}
.meta { min-width: 0; }
.t { font-size: 14px; font-weight: 600; }
.u {
  font-size: 12px; color: rgba(120, 120, 130, 0.95);
  margin-top: 1px;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}
.arrow {
  margin-left: auto; color: rgba(120, 120, 130, 0.8);
  transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  flex-shrink: 0;
}
.contact-item:hover .arrow { transform: translate(3px, -3px); color: #10b981; }
.empty { margin-top: 20px; color: rgba(120, 120, 130, 0.95); font-size: 14px; }
</style>
