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

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: HomePage, meta: { title: '首页' } },
    { path: '/works', component: WorksPage, meta: { title: '作品' } },
    { path: '/work/:id', component: WorkDetailPage, meta: { title: '作品详情' } },
    { path: '/about', component: AboutPage, meta: { title: '关于我' } },
    { path: '/experience', component: ExperiencePage, meta: { title: '经历' } },
    { path: '/contact', component: ContactPage, meta: { title: '联系' } },
    { path: '/ask', component: AskPage, meta: { title: '问 AI' } },
    { path: '/:pathMatch(.*)*', redirect: '/' }
  ],
  scrollBehavior: () => ({ top: 0 })
})

createApp(App).use(router).mount('#app')
