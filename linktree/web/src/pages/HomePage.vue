<script setup>
import { ref, onMounted, computed } from 'vue'
import { NButton } from 'naive-ui'
import { api } from '../api'
import { profile } from '../profile'

const works = ref([])
const featured = computed(() => works.value.slice(0, 2))

const skills = ['Go', 'Python', 'TypeScript', 'React', 'Vue', 'LLM · AI', 'Docker', 'SQL']
const avatarUrl = '/static/avatar.jpg'

const bioText = computed(() => {
  const b = (profile.value.bio || '').replace(/。$/, '')
  const tail = '用代码把想法变成产品 · 全栈开发 · AI 应用实践者'
  return b ? b + '。' + tail : tail
})

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
  <div class="home">
    <section class="hero">
      <div class="reveal" style="--d:0s">
        <span class="status-badge"><span class="dot"></span>Open to Work · 可实习 / 全职</span>
      </div>
      <img :src="avatarUrl" :alt="profile.display_name" class="hero-avatar reveal" style="--d:.08s">
      <h1 class="hero-name reveal" style="--d:.16s">{{ profile.display_name }}<span class="accent">.</span></h1>
      <div class="hero-title mono reveal" style="--d:.24s">Full-Stack / AI Developer</div>
      <p class="hero-bio reveal" style="--d:.32s">{{ bioText }}</p>
      <div class="hero-skills mono reveal" style="--d:.40s">
        <template v-for="(s, i) in skills" :key="s">
          <span>{{ s }}</span><span v-if="i < skills.length - 1" class="sep">/</span>
        </template>
      </div>
      <div class="hero-actions reveal" style="--d:.48s">
        <n-button type="primary" round tag="router-link" to="/ask" size="large">💬 问 AI 关于我</n-button>
        <router-link to="/works" class="btn-text arrow-link">查看作品 <span class="arr">→</span></router-link>
      </div>
    </section>

    <section v-if="featured.length" class="featured reveal" style="--d:.56s">
      <div class="featured-label mono">
        Selected Works · 精选作品
        <router-link to="/works" class="arrow-link">全部作品 <span class="arr">→</span></router-link>
      </div>
      <div class="featured-grid">
        <router-link
          v-for="(w, i) in featured"
          :key="w.id"
          :to="'/work/' + w.id"
          class="featured-card spot"
          @mousemove="trackSpot($event.currentTarget, $event)"
        >
          <img v-if="cover(w)" :src="cover(w)" :alt="w.title" class="featured-card-cover">
          <div v-else class="featured-card-cover ph">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="32" height="32"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
          </div>
          <div class="featured-card-body">
            <div class="featured-card-idx mono">0{{ i + 1 }}<template v-if="w.category_name"> · {{ w.category_name }}</template></div>
            <div class="featured-card-title">{{ w.title }}</div>
            <div v-if="w.description" class="featured-card-desc">{{ w.description }}</div>
          </div>
        </router-link>
      </div>
    </section>
  </div>
</template>

<style scoped>
.home { max-width: 860px; margin: 0 auto; padding: 0 40px; }

.hero { padding: 40px 0 72px; display: flex; flex-direction: column; align-items: flex-start; }
.hero-avatar {
  width: 76px; height: 76px;
  border-radius: 22px;
  object-fit: cover;
  border: 1px solid rgba(0, 0, 0, 0.08);
  margin: 26px 0 28px;
  box-shadow: 0 8px 24px -8px rgba(0, 0, 0, 0.18);
}
body[data-theme="dark"] .hero-avatar { border-color: rgba(255, 255, 255, 0.1); }
.hero-name {
  font-size: clamp(48px, 9vw, 92px);
  font-weight: 800;
  letter-spacing: -0.045em;
  line-height: 0.98;
  text-wrap: balance;
}
.hero-title {
  font-size: 13px;
  font-weight: 500;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: #10b981;
  margin-top: 18px;
}
.hero-bio {
  font-size: 17px;
  color: rgba(80, 80, 90, 0.95);
  line-height: 1.75;
  margin-top: 18px;
  max-width: 540px;
  text-wrap: pretty;
}
body[data-theme="dark"] .hero-bio { color: rgba(190, 190, 200, 0.95); }
.hero-skills { font-size: 13px; color: rgba(120, 120, 130, 0.95); margin-top: 26px; letter-spacing: 0.02em; }
.hero-skills .sep { color: #10b981; margin: 0 10px; opacity: 0.7; }
.hero-actions { display: flex; align-items: center; gap: 20px; margin-top: 36px; flex-wrap: wrap; }
.btn-text {
  display: inline-flex; align-items: center; gap: 6px;
  font-size: 14px; font-weight: 500;
  text-decoration: none;
  padding-bottom: 3px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.14);
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
body[data-theme="dark"] .btn-text { border-color: rgba(255, 255, 255, 0.16); }
.btn-text:hover { color: #10b981; border-color: #10b981; }

.featured { padding-bottom: 64px; }
.featured-label {
  font-size: 11px; font-weight: 500;
  color: rgba(120, 120, 130, 0.95);
  text-transform: uppercase; letter-spacing: 0.12em;
  margin-bottom: 16px;
  display: flex; justify-content: space-between; align-items: center;
}
.featured-label a { color: #10b981; text-decoration: none; font-size: 13px; text-transform: none; letter-spacing: 0; display: inline-flex; align-items: center; gap: 4px; }
.featured-grid { display: grid; grid-template-columns: 1.12fr 1fr; gap: 16px; align-items: start; }
.featured-card {
  display: block;
  border-radius: 16px;
  overflow: hidden;
  text-decoration: none;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(0, 0, 0, 0.06);
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
body[data-theme="dark"] .featured-card { background: rgba(24, 24, 27, 0.85); border-color: rgba(255, 255, 255, 0.06); }
.featured-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.15);
  border-color: rgba(16, 185, 129, 0.4);
}
.featured-card-cover {
  width: 100%; height: 170px;
  object-fit: contain;
  display: block;
  background: rgba(0, 0, 0, 0.04);
  transition: transform 0.5s cubic-bezier(0.16, 1, 0.3, 1);
}
body[data-theme="dark"] .featured-card-cover { background: rgba(255, 255, 255, 0.05); }
.featured-card:hover .featured-card-cover { transform: scale(1.03); }
.featured-card-cover.ph {
  display: flex; align-items: center; justify-content: center;
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.1), rgba(0, 0, 0, 0.04));
  color: #10b981;
}
.featured-card-body { padding: 18px; }
.featured-card-idx { font-size: 11px; color: #10b981; letter-spacing: 0.1em; margin-bottom: 6px; }
.featured-card-title { font-size: 16px; font-weight: 700; letter-spacing: -0.01em; }
.featured-card-desc {
  font-size: 13px; color: rgba(90, 90, 100, 0.95);
  margin-top: 6px; line-height: 1.6;
  display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden;
}
body[data-theme="dark"] .featured-card-desc { color: rgba(170, 170, 180, 0.95); }

@media (max-width: 768px) {
  .home { padding: 0 24px; }
  .hero { padding: 24px 0 56px; }
  .hero-name { font-size: clamp(40px, 12vw, 56px); }
  .featured-grid { grid-template-columns: 1fr; }
  .hero-skills { font-size: 12px; }
  .hero-skills .sep { margin: 0 7px; }
}
</style>
