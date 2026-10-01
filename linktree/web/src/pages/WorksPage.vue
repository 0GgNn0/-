<script setup>
import { ref, onMounted, computed } from 'vue'
import { api } from '../api'

const works = ref([])
const activeFilter = ref('all')

const filters = [
  { key: 'all', label: '全部' },
  { key: 'web', label: 'Web 应用' },
  { key: 'ai', label: 'AI 应用' },
  { key: 'tool', label: '工具类' }
]

const filtered = computed(() =>
  activeFilter.value === 'all'
    ? works.value
    : works.value.filter((w) => w.category_slug === activeFilter.value)
)

function cover(w) {
  if (w.cover_url) return w.cover_url
  const m = (w.media || []).find((x) => x.media_type === 'image')
  return m ? m.media_url : ''
}

function trackSpot(el, e) {
  const r = el.getBoundingClientRect()
  el.style.setProperty('--mx', e.clientX - r.left + 'px')
  el.style.setProperty('--my', e.clientY - r.top + 'px')
}

onMounted(async () => {
  try {
    works.value = await api.works()
  } catch (_) { /* empty state */ }
})
</script>

<template>
  <div class="works-page">
    <header class="page-header reveal" style="--d:0s">
      <div>
        <div class="page-kicker mono">Works</div>
        <h1 class="page-title">作品</h1>
      </div>
      <router-link to="/" class="page-back">← 返回首页</router-link>
    </header>

    <div class="works-filter reveal" style="--d:.1s">
      <button
        v-for="f in filters"
        :key="f.key"
        class="filter-chip"
        :class="{ active: activeFilter === f.key }"
        @click="activeFilter = f.key"
      >
        {{ f.label }}
      </button>
    </div>

    <div v-if="filtered.length" class="works-grid">
      <router-link
        v-for="(w, i) in filtered"
        :key="w.id"
        :to="'/work/' + w.id"
        class="work-card spot reveal"
        :style="{ '--d': i * 0.07 + 's' }"
        @mousemove="trackSpot($event.currentTarget, $event)"
      >
        <img v-if="cover(w)" :src="cover(w)" :alt="w.title" class="work-card-cover">
        <div v-else class="work-card-cover ph">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="32" height="32"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
        </div>
        <div class="work-card-body">
          <span v-if="w.category_name" class="work-card-cat">{{ w.category_name }}</span>
          <div class="work-card-title">{{ w.title }}</div>
          <div v-if="w.description" class="work-card-desc">{{ w.description }}</div>
        </div>
      </router-link>
    </div>
    <div v-else class="empty-hint">暂无作品，等待添加</div>
  </div>
</template>

<style scoped>
.works-page { max-width: 960px; margin: 0 auto; padding: 28px 24px 72px; }

.page-header { display: flex; justify-content: space-between; align-items: flex-end; margin-bottom: 28px; gap: 16px; }
.page-kicker { font-size: 11px; letter-spacing: 0.14em; text-transform: uppercase; color: #10b981; margin-bottom: 8px; }
.page-title { font-size: clamp(32px, 5vw, 44px); font-weight: 800; letter-spacing: -0.03em; line-height: 1.05; }
.page-back {
  font-size: 13px; color: rgba(90, 90, 100, 0.95);
  text-decoration: none; padding: 6px 14px;
  border-radius: 999px; border: 1px solid rgba(0, 0, 0, 0.1);
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(10px);
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  white-space: nowrap;
}
body[data-theme="dark"] .page-back { background: rgba(24, 24, 27, 0.6); color: rgba(190, 190, 200, 0.95); border-color: rgba(255, 255, 255, 0.12); }
.page-back:hover { color: #10b981; border-color: #10b981; transform: translateY(-1px); }

.works-filter { display: flex; gap: 8px; margin-bottom: 24px; flex-wrap: wrap; }
.filter-chip {
  padding: 6px 16px;
  border-radius: 999px;
  border: 1px solid rgba(0, 0, 0, 0.1);
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(10px);
  color: rgba(90, 90, 100, 0.95);
  font-size: 13px;
  cursor: pointer;
  font-family: inherit;
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
body[data-theme="dark"] .filter-chip { background: rgba(24, 24, 27, 0.6); color: rgba(190, 190, 200, 0.95); border-color: rgba(255, 255, 255, 0.12); }
.filter-chip:hover { border-color: #10b981; color: #10b981; transform: translateY(-1px); }
.filter-chip:active { transform: translateY(0) scale(0.97); }
.filter-chip.active { background: #10b981; color: #fff; border-color: #10b981; }

.works-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap: 16px; }
.work-card {
  display: block;
  border-radius: 16px;
  overflow: hidden;
  text-decoration: none;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(0, 0, 0, 0.06);
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
body[data-theme="dark"] .work-card { background: rgba(24, 24, 27, 0.85); border-color: rgba(255, 255, 255, 0.06); }
.work-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.15);
  border-color: rgba(16, 185, 129, 0.4);
}
.work-card:active { transform: translateY(-2px) scale(0.99); }
.work-card-cover {
  width: 100%; height: 165px;
  object-fit: contain; display: block;
  background: rgba(0, 0, 0, 0.04);
  transition: transform 0.5s cubic-bezier(0.16, 1, 0.3, 1);
}
body[data-theme="dark"] .work-card-cover { background: rgba(255, 255, 255, 0.05); }
.work-card:hover .work-card-cover { transform: scale(1.03); }
.work-card-cover.ph {
  display: flex; align-items: center; justify-content: center;
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.1), rgba(0, 0, 0, 0.04));
  color: #10b981;
}
.work-card-body { padding: 16px; }
.work-card-cat {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 999px;
  font-size: 11px;
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
  margin-bottom: 8px;
}
.work-card-title { font-size: 15px; font-weight: 700; letter-spacing: -0.01em; }
.work-card-desc {
  font-size: 13px; color: rgba(90, 90, 100, 0.95);
  margin-top: 4px; line-height: 1.6;
  display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden;
}
body[data-theme="dark"] .work-card-desc { color: rgba(170, 170, 180, 0.95); }
.empty-hint { text-align: center; padding: 60px 20px; color: rgba(120, 120, 130, 0.9); font-size: 14px; }
</style>
