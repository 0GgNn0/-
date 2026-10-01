<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NButton, useMessage } from 'naive-ui'
import { api } from '../api'

const route = useRoute()
const router = useRouter()
const message = useMessage()

const work = ref(null)
const loading = ref(true)
const notFound = ref(false)

const images = computed(() => {
  if (!work.value) return []
  const list = []
  if (work.value.cover_url) list.push(work.value.cover_url)
  for (const m of work.value.media || []) {
    if (m.media_type === 'image' && !list.includes(m.media_url)) list.push(m.media_url)
  }
  return list
})

const videos = computed(() =>
  ((work.value && work.value.media) || []).filter((m) => m.media_type === 'video')
)
const audios = computed(() =>
  ((work.value && work.value.media) || []).filter((m) => m.media_type === 'audio')
)

onMounted(async () => {
  try {
    const works = await api.works()
    work.value = works.find((w) => String(w.id) === String(route.params.id)) || null
    if (!work.value) notFound.value = true
  } catch (_) {
    notFound.value = true
  } finally {
    loading.value = false
  }
})

function openExternal() {
  if (work.value && work.value.external_url) {
    window.open(work.value.external_url, '_blank', 'noopener')
  } else {
    message.warning('暂无外部链接')
  }
}
</script>

<template>
  <div class="detail-page">
    <div class="reveal" style="--d:0s">
      <router-link to="/works" class="page-back">← 返回作品</router-link>
    </div>

    <div v-if="loading" class="state-box">加载中...</div>

    <template v-else-if="work">
      <article class="detail glass reveal" style="--d:.1s">
        <span v-if="work.category_name" class="cat-tag">{{ work.category_name }}</span>
        <h1 class="detail-title">{{ work.title }}</h1>
        <p v-if="work.description" class="detail-desc">{{ work.description }}</p>

        <div class="gallery" v-if="images.length">
          <img
            v-for="(src, i) in images"
            :key="i"
            :src="src"
            :alt="work.title"
            class="gallery-img"
          >
        </div>

        <div class="media-block" v-if="videos.length">
          <div class="block-label mono">视频</div>
          <video v-for="v in videos" :key="v.id" :src="v.media_url" controls class="media-video"></video>
        </div>

        <div class="media-block" v-if="audios.length">
          <div class="block-label mono">音频</div>
          <audio v-for="a in audios" :key="a.id" :src="a.media_url" controls class="media-audio"></audio>
        </div>

        <div class="detail-actions">
          <n-button v-if="work.external_url" type="primary" round @click="openExternal">
            打开项目 ↗
          </n-button>
          <router-link to="/works" class="back-link arrow-link">浏览更多作品 <span class="arr">→</span></router-link>
        </div>
      </article>
    </template>

    <div v-else class="state-box">
      作品不存在
      <div style="margin-top: 16px">
        <router-link to="/works" class="back-link">← 返回作品列表</router-link>
      </div>
    </div>
  </div>
</template>

<style scoped>
.detail-page { max-width: 760px; margin: 0 auto; padding: 28px 24px 72px; }

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
.page-back:hover { color: #10b981; border-color: #10b981; }

.detail { margin-top: 20px; border-radius: 16px; padding: 36px; }
.cat-tag {
  display: inline-block;
  padding: 3px 10px;
  border-radius: 999px;
  font-size: 12px;
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
  margin-bottom: 14px;
}
.detail-title { font-size: clamp(28px, 5vw, 40px); font-weight: 800; letter-spacing: -0.03em; line-height: 1.1; }
.detail-desc { font-size: 15px; line-height: 1.8; color: rgba(90, 90, 100, 0.95); margin-top: 14px; text-wrap: pretty; }
body[data-theme="dark"] .detail-desc { color: rgba(180, 180, 190, 0.95); }

.gallery { margin-top: 24px; display: flex; flex-direction: column; gap: 12px; }
.gallery-img {
  width: 100%;
  border-radius: 12px;
  background: rgba(0, 0, 0, 0.04);
  border: 1px solid rgba(0, 0, 0, 0.06);
}
body[data-theme="dark"] .gallery-img { background: rgba(255, 255, 255, 0.04); border-color: rgba(255, 255, 255, 0.08); }

.media-block { margin-top: 24px; display: flex; flex-direction: column; gap: 10px; }
.block-label { font-size: 11px; letter-spacing: 0.12em; text-transform: uppercase; color: rgba(120, 120, 130, 0.95); }
.media-video { width: 100%; border-radius: 12px; }
.media-audio { width: 100%; }

.detail-actions { margin-top: 32px; display: flex; align-items: center; gap: 20px; flex-wrap: wrap; }
.back-link { font-size: 14px; color: rgba(90, 90, 100, 0.95); text-decoration: none; }
body[data-theme="dark"] .back-link { color: rgba(190, 190, 200, 0.95); }
.back-link:hover { color: #10b981; }

.state-box { text-align: center; padding: 80px 20px; color: rgba(120, 120, 130, 0.95); font-size: 14px; }
</style>
