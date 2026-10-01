<script setup>
import { onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { NConfigProvider, NMessageProvider, NDialogProvider, NButton } from 'naive-ui'
import { naiveTheme, themeVars, isDark, toggleTheme, applyDom } from './theme'
import { profile, loadProfile } from './profile'

const route = useRoute()

const nav = [
  { to: '/works', label: '作品' },
  { to: '/about', label: '关于' },
  { to: '/experience', label: '经历' },
  { to: '/ask', label: '问 AI' },
  { to: '/contact', label: '联系' }
]

const pageTitle = computed(() =>
  route.meta.title && route.path !== '/'
    ? route.meta.title + ' - ' + profile.value.display_name
    : profile.value.display_name
)

onMounted(() => {
  applyDom()
  loadProfile()
})
</script>

<template>
  <n-config-provider :theme="naiveTheme" :theme-vars="themeVars">
    <n-dialog-provider>
      <n-message-provider placement="top">
        <div class="app-wrap">
          <header class="topbar">
            <router-link to="/" class="topbar-name mono">{{ profile.display_name }}</router-link>
            <nav class="topnav">
              <router-link v-for="item in nav" :key="item.to" :to="item.to">{{ item.label }}</router-link>
            </nav>
            <n-button
              size="small"
              quaternary
              circle
              class="theme-btn"
              :aria-label="isDark ? '切换到亮色' : '切换到暗色'"
              @click="toggleTheme"
            >
              <template #icon>
                <svg v-if="isDark" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg>
                <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>
              </template>
            </n-button>
          </header>

          <main class="app-main">
            <router-view v-slot="{ Component }">
              <transition name="page" mode="out-in">
                <component :is="Component" :key="route.path" />
              </transition>
            </router-view>
          </main>

          <footer class="app-footer mono">
            <router-link to="/ask">问 AI 关于我</router-link>
            <a href="/admin/login">管理后台</a>
            <span class="copy">© 2026 {{ profile.display_name }}</span>
          </footer>
        </div>
      </n-message-provider>
    </n-dialog-provider>
  </n-config-provider>
</template>

<style scoped>
.app-wrap {
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  position: relative;
  z-index: 1;
}

.topbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  padding: 24px 40px;
}
.topbar-name {
  font-size: 13px;
  font-weight: 600;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  text-decoration: none;
  color: inherit;
}
.topnav { display: flex; gap: 4px; align-items: center; flex-wrap: wrap; }
.topnav a {
  padding: 7px 14px;
  border-radius: 999px;
  font-size: 13px;
  color: rgba(80, 80, 90, 0.9);
  text-decoration: none;
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
body[data-theme="dark"] .topnav a { color: rgba(200, 200, 210, 0.85); }
.topnav a:hover { color: #10b981; background: rgba(16, 185, 129, 0.08); }
.topnav a.router-link-active {
  color: #10b981;
  background: rgba(16, 185, 129, 0.1);
}
.theme-btn { flex-shrink: 0; }

.app-main { flex: 1; }

.app-footer {
  padding: 24px 40px 32px;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 24px;
  flex-wrap: wrap;
}
.app-footer a, .app-footer .copy {
  font-size: 12px;
  color: rgba(120, 120, 130, 0.85);
  text-decoration: none;
  letter-spacing: 0.04em;
}
.app-footer a:hover { color: #10b981; }

@media (max-width: 768px) {
  .topbar { padding: 16px 20px; flex-wrap: wrap; }
  .topnav a { padding: 6px 10px; font-size: 12px; }
  .app-footer { padding: 20px; gap: 16px; }
}
</style>
