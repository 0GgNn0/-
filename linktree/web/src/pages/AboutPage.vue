<script setup>
const stats = [
  { n: '3+', l: '年开发经验' },
  { n: '20+', l: '项目作品' },
  { n: '5+', l: 'AI 应用实践' }
]
const skills = ['Go', 'Python', 'TypeScript', 'React', 'Vue', 'LLM / AI', 'Docker', 'SQL']
</script>

<template>
  <div class="about-page">
    <div class="reveal" style="--d:0s">
      <router-link to="/" class="page-back">← 返回首页</router-link>
    </div>

    <div class="card glass reveal" style="--d:.1s">
      <div class="kicker mono">About</div>
      <h1 class="h1">关于我</h1>
      <p class="body">
        全栈开发者，热衷于用 AI 技术构建实用产品。擅长从需求分析到部署上线的全流程开发，关注大模型应用落地与自动化工具链建设。
      </p>

      <div class="stats">
        <div v-for="s in stats" :key="s.l" class="stat">
          <div class="n mono">{{ s.n }}</div>
          <div class="l">{{ s.l }}</div>
        </div>
      </div>

      <div class="section-label mono">技术栈</div>
      <div class="skills">
        <span v-for="sk in skills" :key="sk">{{ sk }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.about-page { max-width: 680px; margin: 0 auto; padding: 28px 24px 72px; }

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
.body { font-size: 15px; line-height: 1.8; color: rgba(90, 90, 100, 0.95); margin-top: 18px; text-wrap: pretty; }
body[data-theme="dark"] .body { color: rgba(180, 180, 190, 0.95); }

.stats { display: flex; gap: 32px; margin-top: 28px; flex-wrap: wrap; }
.stat .n {
  font-size: 30px; font-weight: 600; color: #10b981;
  font-variant-numeric: tabular-nums;
  letter-spacing: -0.02em;
}
.stat .l { font-size: 13px; color: rgba(90, 90, 100, 0.95); margin-top: 2px; }
body[data-theme="dark"] .stat .l { color: rgba(180, 180, 190, 0.95); }

.section-label {
  font-size: 11px; font-weight: 500;
  color: rgba(120, 120, 130, 0.95);
  text-transform: uppercase; letter-spacing: 0.12em;
  margin-top: 28px; margin-bottom: 12px;
}
.skills { display: flex; flex-wrap: wrap; gap: 8px; }
.skills span {
  padding: 5px 14px;
  border-radius: 999px;
  font-size: 13px;
  background: rgba(0, 0, 0, 0.045);
  color: rgba(90, 90, 100, 0.95);
  border: 1px solid rgba(0, 0, 0, 0.07);
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
body[data-theme="dark"] .skills span {
  background: rgba(255, 255, 255, 0.05);
  color: rgba(190, 190, 200, 0.95);
  border-color: rgba(255, 255, 255, 0.09);
}
.skills span:hover { border-color: #10b981; color: #10b981; transform: translateY(-1px); }
</style>
