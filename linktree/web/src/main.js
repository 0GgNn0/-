import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import './style.css'

import HomePage from './pages/HomePage.vue'
import WorksPage from './pages/WorksPage.vue'
import WorkDetailPage from './pages/WorkDetailPage.vue'
import AboutPage from './pages/AboutPage.vue'
import ExperiencePage from './pages/ExperiencePage.vue'
import ContactPage from './pages/ContactPage.vue'
import AskPage from './pages/AskPage.vue'
import NotFoundPage from './pages/NotFoundPage.vue'

const SITE_NAME = '啊芃'
const DEFAULT_DESC = '全栈 / AI 开发者个人主页：作品集、经历与 AI 问答'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: HomePage, meta: { title: '', desc: DEFAULT_DESC } },
    { path: '/works', component: WorksPage, meta: { title: '作品', desc: '我做过的东西：Web 应用、AI 应用与工具类项目' } },
    { path: '/work/:id', component: WorkDetailPage, meta: { title: '作品详情', desc: '项目详情：技术选型、实现方式与线上地址' } },
    { path: '/about', component: AboutPage, meta: { title: '关于我', desc: '技术栈、个人简介与职业方向' } },
    { path: '/experience', component: ExperiencePage, meta: { title: '经历', desc: '工作与项目经历时间线' } },
    { path: '/contact', component: ContactPage, meta: { title: '联系', desc: '通过邮箱或社交账号联系我' } },
    { path: '/ask', component: AskPage, meta: { title: '问 AI', desc: '用 AI 问答了解我的技能与项目经历' } },
    { path: '/:pathMatch(.*)*', component: NotFoundPage, meta: { title: '页面不存在', noindex: true } }
  ],
  scrollBehavior: () => ({ top: 0 })
})

// 路由级标题与描述：爬虫和分享卡片之外，也让浏览器标签页/历史记录有意义
router.afterEach((to) => {
  const title = to.meta.title ? `${to.meta.title} - ${SITE_NAME}` : `${SITE_NAME} · 全栈 / AI 开发者`
  document.title = title

  const setMeta = (selector, attr, value) => {
    let el = document.head.querySelector(selector)
    if (!el) {
      el = document.createElement('meta')
      if (selector.includes('property=')) el.setAttribute('property', selector.match(/"([^"]+)"/)[1])
      else el.setAttribute('name', selector.match(/"([^"]+)"/)[1])
      document.head.appendChild(el)
    }
    el.setAttribute(attr, value)
  }
  const desc = to.meta.desc || DEFAULT_DESC
  setMeta('meta[name="description"]', 'content', desc)
  setMeta('meta[property="og:title"]', 'content', title)
  setMeta('meta[property="og:description"]', 'content', desc)
  setMeta('meta[property="og:url"]', 'content', window.location.href)
  setMeta('meta[name="robots"]', 'content', to.meta.noindex ? 'noindex,follow' : 'index,follow')
})

createApp(App).use(router).mount('#app')
