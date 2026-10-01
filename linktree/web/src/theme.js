import { ref, computed } from 'vue'
import { darkTheme } from 'naive-ui'

const stored = ref(localStorage.getItem('linktree-theme') || 'auto')
const media = window.matchMedia('(prefers-color-scheme: dark)')

function effective() {
  if (stored.value === 'auto') return media.matches ? 'dark' : 'light'
  return stored.value
}

export const themeMode = stored
export const isDark = computed(() => effective() === 'dark')
export const naiveTheme = computed(() => (isDark.value ? darkTheme : null))

export function toggleTheme() {
  setTheme(effective() === 'dark' ? 'light' : 'dark')
}

export function setTheme(t) {
  stored.value = t
  localStorage.setItem('linktree-theme', t)
  applyDom()
}

export function applyDom() {
  const t = effective()
  document.documentElement.setAttribute('data-theme', t)
  document.body.setAttribute('data-theme', t)
}

media.addEventListener('change', () => {
  if (stored.value === 'auto') applyDom()
})

export const themeVars = {
  common: {
    primaryColor: '#10b981',
    primaryColorHover: '#0fa97a',
    primaryColorPressed: '#059669',
    primaryColorSuppl: '#10b981',
    borderRadius: '10px',
    fontFamily: "'Geist', -apple-system, BlinkMacSystemFont, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif"
  },
  Card: {
    borderRadius: '16px'
  },
  Button: {
    fontWeight: '600'
  }
}
