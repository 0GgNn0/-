# LinkTree 前端重设计实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将 LinkTree 个人页面升级为治愈蓝风格，添加动态渐变背景、3D 交互按钮和精致微动画

**Architecture:** 更新 CSS 样式实现视觉效果，添加 JavaScript 实现 3D 交互和主题切换

**Tech Stack:** HTML5, CSS3 (动画/渐变/毛玻璃), JavaScript (DOM 事件处理)

---

## 文件结构

| 文件 | 操作 | 职责 |
|------|------|------|
| `static/style.css` | 修改 | 更新色彩系统、添加动画、3D 效果、响应式适配 |
| `templates/index.html` | 修改 | 添加 3D 交互 JavaScript、优化主题切换脚本 |

---

## Task 1: 更新 CSS 变量和色彩系统

**Files:**
- Modify: `static/style.css:1-58`

- [ ] **Step 1: 更新亮色模式 CSS 变量**

将 `static/style.css` 中的 `:root` 部分替换为治愈蓝色彩系统：

```css
:root {
  /* 主色渐变 */
  --primary-start: #4facfe;
  --primary-end: #00f2fe;
  --primary-gradient: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  
  /* 背景渐变 */
  --bg-gradient: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  
  /* 卡片 */
  --bg-start: #e8ecf1;
  --bg-end: #d5dce6;
  --card-bg: rgba(255, 255, 255, 0.85);
  --card-border: rgba(79, 172, 254, 0.2);
  --card-shadow: 0 8px 32px rgba(79, 172, 254, 0.1), 0 2px 8px rgba(0, 0, 0, 0.04);
  
  /* 按钮 */
  --btn-bg: rgba(255, 255, 255, 0.78);
  --btn-hover-bg: rgba(255, 255, 255, 0.92);
  --btn-border: rgba(79, 172, 254, 0.15);
  --btn-hover-border: rgba(79, 172, 254, 0.3);
  --btn-shadow: 0 2px 8px rgba(79, 172, 254, 0.08);
  --btn-hover-shadow: 0 6px 20px rgba(79, 172, 254, 0.2);
  
  /* 文字 */
  --text-primary: #2d3748;
  --text-secondary: #718096;
  --text-btn: #2c3e50;
  
  /* 强调色 */
  --accent: #4facfe;
  --accent-light: rgba(79, 172, 254, 0.10);
  
  /* 其他 */
  --divider: rgba(79, 172, 254, 0.10);
  --input-bg: rgba(79, 172, 254, 0.03);
  --input-border: rgba(79, 172, 254, 0.15);
  --danger: #e74c3c;
  --success: #27ae60;
  --toggle-bg: rgba(79, 172, 254, 0.08);
  --footer-color: #9aa0b0;
}
```

- [ ] **Step 2: 更新暗色模式 CSS 变量**

将 `[data-theme="dark"]` 部分替换为：

```css
[data-theme="dark"] {
  --bg-start: #0f0f1a;
  --bg-end: #1a1a2e;
  --card-bg: rgba(30, 30, 55, 0.75);
  --card-border: rgba(79, 172, 254, 0.25);
  --card-shadow: 0 8px 32px rgba(0, 0, 0, 0.3), 0 2px 8px rgba(79, 172, 254, 0.1);
  --btn-bg: rgba(255, 255, 255, 0.06);
  --btn-hover-bg: rgba(255, 255, 255, 0.10);
  --btn-border: rgba(79, 172, 254, 0.2);
  --btn-hover-border: rgba(79, 172, 254, 0.35);
  --btn-shadow: 0 2px 8px rgba(79, 172, 254, 0.15);
  --btn-hover-shadow: 0 6px 20px rgba(79, 172, 254, 0.3);
  --text-primary: #e2e8f0;
  --text-secondary: #a0aec0;
  --text-btn: #d0d0e0;
  --accent: #4facfe;
  --accent-light: rgba(79, 172, 254, 0.20);
  --divider: rgba(79, 172, 254, 0.12);
  --input-bg: rgba(79, 172, 254, 0.05);
  --input-border: rgba(79, 172, 254, 0.2);
  --danger: #e74c3c;
  --success: #2ecc71;
  --toggle-bg: rgba(79, 172, 254, 0.12);
  --footer-color: #5a5a78;
}
```

- [ ] **Step 3: 验证变量更新**

打开浏览器开发者工具，检查 CSS 变量是否正确加载。

- [ ] **Step 4: Commit**

```bash
git add static/style.css
git commit -m "feat: update CSS variables to healing blue theme"
```

---

## Task 2: 实现动态渐变背景

**Files:**
- Modify: `static/style.css:60-101`

- [ ] **Step 1: 添加背景动画关键帧**

在 `static/style.css` 中添加动画关键帧（在全局重置之后）：

```css
/* ---- 背景动画 ---- */
@keyframes gradientFlow1 {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

@keyframes gradientFlow2 {
  0% { background-position: 100% 0%; }
  50% { background-position: 0% 100%; }
  100% { background-position: 100% 0%; }
}

@keyframes gradientFlow3 {
  0% { background-position: 50% 0%; }
  50% { background-position: 50% 100%; }
  100% { background-position: 50% 0%; }
}
```

- [ ] **Step 2: 更新 body 背景样式**

将 `body` 的背景样式替换为：

```css
body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC",
    "Hiragino Sans GB", "Microsoft YaHei", "Helvetica Neue", Arial, sans-serif;
  min-height: 100vh;
  color: var(--text-primary);
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding: 40px 20px;
  transition: color 0.3s ease;
  position: relative;
  overflow-x: hidden;
  background: 
    linear-gradient(135deg, rgba(79, 172, 254, 0.25) 0%, rgba(0, 242, 254, 0.25) 100%),
    linear-gradient(225deg, rgba(102, 126, 234, 0.25) 0%, rgba(118, 75, 162, 0.25) 100%),
    linear-gradient(315deg, rgba(79, 172, 254, 0.2) 0%, rgba(0, 242, 254, 0.2) 100%),
    linear-gradient(45deg, rgba(102, 126, 234, 0.2) 0%, rgba(118, 75, 162, 0.2) 100%);
  background-size: 200% 200%, 200% 200%, 200% 200%, 200% 200%;
  animation: gradientFlow1 20s ease infinite, gradientFlow2 25s ease infinite, gradientFlow3 30s ease infinite;
}
```

- [ ] **Step 3: 更新背景遮罩样式**

将 `body::before` 替换为：

```css
body::before {
  content: '';
  position: fixed;
  inset: 0;
  z-index: 0;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(2px);
  -webkit-backdrop-filter: blur(2px);
  pointer-events: none;
}

[data-theme="dark"] body::before {
  background: rgba(0, 0, 0, 0.25);
}
```

- [ ] **Step 4: 验证背景动画**

刷新页面，确认渐变背景动画正常运行，无卡顿。

- [ ] **Step 5: Commit**

```bash
git add static/style.css
git commit -m "feat: add animated gradient background"
```

---

## Task 3: 更新头像样式和悬停效果

**Files:**
- Modify: `static/style.css:151-183`

- [ ] **Step 1: 更新头像容器样式**

将 `.avatar-wrap` 样式替换为：

```css
.avatar-wrap {
  margin: 0 auto 20px;
  width: 96px;
  height: 96px;
  border-radius: 50%;
  overflow: visible;
  border: 3px solid var(--card-border);
  box-shadow: 0 4px 16px rgba(79, 172, 254, 0.2);
  transition: transform 0.4s cubic-bezier(0.34, 1.56, 0.64, 1),
              box-shadow 0.4s ease;
  position: relative;
}

.avatar-wrap:hover {
  transform: scale(1.1);
  box-shadow: 0 10px 40px rgba(79, 172, 254, 0.4);
}
```

- [ ] **Step 2: 添加头像发光效果**

在 `.avatar-wrap:hover` 之后添加：

```css
.avatar-wrap::after {
  content: '';
  position: absolute;
  inset: -6px;
  background: var(--primary-gradient);
  border-radius: 50%;
  opacity: 0;
  transition: opacity 0.4s ease;
  z-index: -1;
  filter: blur(8px);
}

.avatar-wrap:hover::after {
  opacity: 0.6;
}
```

- [ ] **Step 3: 更新头像图片样式**

将 `.avatar-wrap img` 样式更新为：

```css
.avatar-wrap img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
  position: relative;
  z-index: 1;
}
```

- [ ] **Step 4: 验证头像效果**

刷新页面，悬停头像，确认放大和蓝色发光效果正常。

- [ ] **Step 5: Commit**

```bash
git add static/style.css
git commit -m "feat: add avatar hover glow effect"
```

---

## Task 4: 更新链接按钮样式和 3D 效果

**Files:**
- Modify: `static/style.css:202-264`

- [ ] **Step 1: 更新链接按钮基础样式**

将 `.link-btn` 样式替换为：

```css
.link-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  height: 54px;
  padding: 0 24px;
  border-radius: 16px;
  background: var(--btn-bg);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  border: 1px solid var(--btn-border);
  box-shadow: var(--btn-shadow);
  color: var(--text-btn);
  font-size: 16px;
  font-weight: 500;
  text-decoration: none;
  letter-spacing: 0.3px;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transform-style: preserve-3d;
}
```

- [ ] **Step 2: 添加按钮发光效果**

将 `.link-btn::after` 替换为：

```css
.link-btn::before {
  content: '';
  position: absolute;
  inset: -2px;
  background: var(--primary-gradient);
  border-radius: inherit;
  opacity: 0;
  transition: opacity 0.3s ease;
  z-index: -1;
  filter: blur(8px);
}

.link-btn::after {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.2), transparent);
  transition: left 0.5s ease;
}
```

- [ ] **Step 3: 更新按钮悬停效果**

将 `.link-btn:hover` 替换为：

```css
.link-btn:hover {
  transform: translateY(-2px) scale(1.02);
  background: var(--btn-hover-bg);
  border-color: var(--btn-hover-border);
  box-shadow: var(--btn-hover-shadow);
}

.link-btn:hover::before {
  opacity: 0.5;
}

.link-btn:hover::after {
  left: 100%;
}
```

- [ ] **Step 4: 验证按钮效果**

刷新页面，悬停按钮，确认发光和缩放效果正常。

- [ ] **Step 5: Commit**

```bash
git add static/style.css
git commit -m "feat: add button glow and hover effects"
```

---

## Task 5: 添加 3D 倾斜交互 JavaScript

**Files:**
- Modify: `templates/index.html:53-95`

- [ ] **Step 1: 在 index.html 中添加 3D 交互脚本**

在 `</body>` 标签之前，主题切换脚本之后添加：

```html
<script>
  // 3D 倾斜效果
  (function() {
    const buttons = document.querySelectorAll('.link-btn');
    
    buttons.forEach(btn => {
      btn.addEventListener('mousemove', (e) => {
        const rect = btn.getBoundingClientRect();
        const x = e.clientX - rect.left;
        const y = e.clientY - rect.top;
        const centerX = rect.width / 2;
        const centerY = rect.height / 2;
        const rotateX = (y - centerY) / 10;
        const rotateY = (centerX - x) / 10;
        
        btn.style.transform = `perspective(1000px) rotateX(${rotateX}deg) rotateY(${rotateY}deg) scale(1.02)`;
      });
      
      btn.addEventListener('mouseleave', () => {
        btn.style.transform = 'perspective(1000px) rotateX(0) rotateY(0) scale(1)';
      });
      
      btn.addEventListener('mousedown', () => {
        btn.style.transform = 'perspective(1000px) rotateX(0) rotateY(0) scale(0.98)';
      });
      
      btn.addEventListener('mouseup', () => {
        btn.style.transform = 'perspective(1000px) rotateX(0) rotateY(0) scale(1)';
      });
    });
  })();
</script>
```

- [ ] **Step 2: 验证 3D 效果**

刷新页面，鼠标在按钮上移动，确认 3D 倾斜效果流畅，鼠标离开时平滑回正。

- [ ] **Step 3: Commit**

```bash
git add templates/index.html
git commit -m "feat: add 3D tilt interaction for buttons"
```

---

## Task 6: 优化主题切换脚本

**Files:**
- Modify: `templates/index.html:53-95`

- [ ] **Step 1: 更新主题切换脚本**

将现有的主题切换脚本替换为优化版本：

```html
<script>
  // 主题切换
  (function() {
    const html = document.documentElement;
    const toggle = document.getElementById('themeToggle');
    const icon = document.getElementById('themeIcon');
    const storedTheme = localStorage.getItem('linktree-theme') || 'auto';

    function getSystemTheme() {
      return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
    }

    function getEffectiveTheme() {
      if (storedTheme === 'auto') {
        return getSystemTheme();
      }
      return storedTheme;
    }

    function applyTheme(theme) {
      html.setAttribute('data-theme', theme);
      icon.textContent = theme === 'dark' ? '☀️' : '🌙';
      localStorage.setItem('linktree-theme', theme);
      
      // 更新 body 背景动画
      if (theme === 'dark') {
        document.body.style.animation = 'none';
        document.body.offsetHeight; // 触发重绘
        document.body.style.animation = '';
      }
    }

    // 初始主题
    applyTheme(getEffectiveTheme());

    // 切换按钮
    toggle.addEventListener('click', function() {
      const current = html.getAttribute('data-theme');
      const next = current === 'dark' ? 'light' : 'dark';
      applyTheme(next);
    });

    // 系统主题变化时自动跟随
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', function() {
      if (localStorage.getItem('linktree-theme') === 'auto') {
        applyTheme(getSystemTheme());
      }
    });
  })();
</script>
```

- [ ] **Step 2: 验证主题切换**

刷新页面，点击主题切换按钮，确认暗色/亮色模式切换正常，背景动画在切换时平滑过渡。

- [ ] **Step 3: Commit**

```bash
git add templates/index.html
git commit -m "feat: optimize theme toggle with smooth transitions"
```

---

## Task 7: 更新卡片和容器样式

**Files:**
- Modify: `static/style.css:133-149`

- [ ] **Step 1: 更新主容器样式**

将 `.container` 样式替换为：

```css
.container {
  position: relative;
  z-index: 1;
  max-width: 500px;
  width: 100%;
  margin: 60px auto 40px;
  background: var(--card-bg);
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  border: 1px solid var(--card-border);
  border-radius: 28px;
  box-shadow: var(--card-shadow);
  padding: 48px 40px 40px;
  text-align: center;
  transition: all 0.3s ease;
}

.container:hover {
  box-shadow: 0 12px 48px rgba(79, 172, 254, 0.15), 0 4px 12px rgba(0, 0, 0, 0.05);
  border-color: rgba(79, 172, 254, 0.3);
}
```

- [ ] **Step 2: 验证卡片效果**

刷新页面，确认卡片毛玻璃效果正常，悬停时有微妙的发光效果。

- [ ] **Step 3: Commit**

```bash
git add static/style.css
git commit -m "feat: enhance card hover glow effect"
```

---

## Task 8: 更新响应式样式

**Files:**
- Modify: `static/style.css:622-662`

- [ ] **Step 1: 更新移动端响应式样式**

将 `@media (max-width: 600px)` 中的样式更新为：

```css
@media (max-width: 600px) {
  body {
    padding: 20px 12px;
  }

  .container {
    padding: 36px 24px 32px;
    border-radius: 22px;
    margin-top: 30px;
  }

  .name {
    font-size: 20px;
  }

  .link-btn {
    height: 48px;
    font-size: 15px;
    border-radius: 14px;
  }

  .theme-toggle {
    top: 16px;
    right: 16px;
    width: 36px;
    height: 36px;
    font-size: 18px;
  }

  .admin-container {
    padding: 24px 20px;
    margin: 20px 0;
    border-radius: 18px;
  }

  .link-table th:nth-child(2),
  .link-table td:nth-child(2) {
    display: none;
  }
}
```

- [ ] **Step 2: 添加平板响应式样式**

在移动端样式之后添加：

```css
@media (min-width: 600px) and (max-width: 900px) {
  .container {
    max-width: 480px;
    padding: 44px 36px 36px;
  }
}
```

- [ ] **Step 3: 验证响应式布局**

使用浏览器开发者工具，切换不同设备尺寸，确认布局正确。

- [ ] **Step 4: Commit**

```bash
git add static/style.css
git commit -m "feat: add tablet responsive breakpoint"
```

---

## Task 9: 最终测试和验证

**Files:**
- None (testing only)

- [ ] **Step 1: 启动本地服务器**

```bash
cd C:\Users\20506\Desktop\中转站\linktree
go run main.go
```

- [ ] **Step 2: 测试动态背景**

访问 http://localhost:8080，确认渐变背景动画正常运行，无卡顿。

- [ ] **Step 3: 测试 3D 按钮效果**

鼠标在按钮上移动，确认 3D 倾斜效果流畅，鼠标离开时平滑回正。

- [ ] **Step 4: 测试头像效果**

悬停头像，确认放大和蓝色发光效果正常。

- [ ] **Step 5: 测试主题切换**

点击主题切换按钮，确认暗色/亮色模式切换正常。

- [ ] **Step 6: 测试响应式布局**

使用浏览器开发者工具，切换移动/平板/桌面尺寸，确认布局正确。

- [ ] **Step 7: 最终 Commit**

```bash
git add .
git commit -m "feat: complete LinkTree frontend redesign with healing blue theme"
```

---

## 完成

所有任务完成后，LinkTree 个人页面将具有：

1. ✅ 动态渐变背景（20-25秒循环动画）
2. ✅ 按钮 3D 倾斜效果（鼠标跟随）
3. ✅ 头像悬停放大 + 蓝色发光
4. ✅ 治愈蓝色系配色
5. ✅ 完整的暗色/亮色模式适配
6. ✅ 响应式设计（移动/平板/桌面）
