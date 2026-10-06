<script setup>
</script>

<template>
  <div class="nf-page">
    <div class="nf-code mono">404</div>
    <h1 class="nf-title">这个页面不存在</h1>
    <p class="nf-desc">链接可能已经失效，或者地址输错了。</p>
    <div class="nf-actions">
      <router-link to="/" class="nf-btn">← 回到首页</router-link>
      <router-link to="/works" class="nf-link arrow-link">看看我的作品 <span class="arr">→</span></router-link>
    </div>
  </div>
</template>

<style scoped>
.nf-page {
  max-width: 560px;
  margin: 0 auto;
  padding: 96px 24px 120px;
  text-align: center;
}
.nf-code {
  font-size: 72px;
  font-weight: 700;
  letter-spacing: -0.04em;
  color: #10b981;
  line-height: 1;
}
.nf-title {
  font-size: clamp(24px, 4vw, 32px);
  font-weight: 800;
  letter-spacing: -0.02em;
  margin-top: 16px;
}
.nf-desc {
  font-size: 15px;
  color: rgba(90, 90, 100, 0.95);
  margin-top: 12px;
}
body[data-theme="dark"] .nf-desc { color: rgba(180, 180, 190, 0.95); }
.nf-actions {
  margin-top: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  flex-wrap: wrap;
}
.nf-btn {
  font-size: 14px;
  font-weight: 600;
  color: #10b981;
  text-decoration: none;
  padding: 11px 22px;
  border-radius: 999px;
  border: 1px solid rgba(16, 185, 129, 0.35);
  background: rgba(16, 185, 129, 0.1);
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
.nf-btn:hover { background: rgba(16, 185, 129, 0.18); transform: translateY(-1px); }
.nf-link {
  font-size: 14px;
  font-weight: 500;
  color: inherit;
  text-decoration: none;
  padding-bottom: 3px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.14);
}
body[data-theme="dark"] .nf-link { border-color: rgba(255, 255, 255, 0.16); }
.nf-link:hover { color: #10b981; border-color: #10b981; }
</style>
