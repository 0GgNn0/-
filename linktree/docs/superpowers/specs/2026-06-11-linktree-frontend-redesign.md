# LinkTree 前端重设计规范

**日期：** 2026-06-11
**项目：** LinkTree 个人链接聚合页
**设计风格：** 活力创意 - 治愈蓝

---

## 概述

将 LinkTree 个人页面从当前的基础毛玻璃风格升级为更具活力和治愈感的设计，采用动态渐变背景、3D 交互按钮和精致的微动画。

---

## 设计目标

1. **治愈感** - 蓝色系配色，柔和的视觉效果
2. **活力创意** - 动态背景、3D 交互、流畅动画
3. **现代感** - 精致的毛玻璃、发光效果、响应式设计
4. **用户体验** - 清晰的视觉层次、直观的交互反馈

---

## 色彩系统

### 主色调 - 治愈蓝

```css
:root {
  /* 主色渐变 */
  --primary-start: #4facfe;
  --primary-end: #00f2fe;
  --primary-gradient: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  
  /* 背景渐变 */
  --bg-gradient: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  
  /* 文字 */
  --text-primary: #2d3748;
  --text-secondary: #718096;
  
  /* 强调色 */
  --accent: #4facfe;
  --accent-glow: rgba(79, 172, 254, 0.4);
  
  /* 卡片 */
  --card-bg: rgba(255, 255, 255, 0.85);
  --card-border: rgba(79, 172, 254, 0.2);
  --card-shadow: 0 8px 32px rgba(79, 172, 254, 0.1);
}
```

### 暗色模式

```css
[data-theme="dark"] {
  --text-primary: #e2e8f0;
  --text-secondary: #a0aec0;
  --card-bg: rgba(30, 30, 55, 0.75);
  --card-border: rgba(79, 172, 254, 0.3);
  --card-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}
```

---

## 动态背景

### 实现方式

使用 CSS 动画实现多层渐变流动效果，3-4 个渐变层以不同速度和方向移动。

### 关键帧动画

```css
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

body {
  background: 
    linear-gradient(135deg, rgba(79, 172, 254, 0.3) 0%, rgba(0, 242, 254, 0.3) 100%),
    linear-gradient(225deg, rgba(102, 126, 234, 0.3) 0%, rgba(118, 75, 162, 0.3) 100%),
    linear-gradient(315deg, rgba(79, 172, 254, 0.2) 0%, rgba(0, 242, 254, 0.2) 100%);
  background-size: 200% 200%, 200% 200%, 200% 200%;
  animation: gradientFlow1 20s ease infinite, gradientFlow2 25s ease infinite;
}
```

### 背景遮罩

```css
body::before {
  content: '';
  position: fixed;
  inset: 0;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(2px);
  pointer-events: none;
}

[data-theme="dark"] body::before {
  background: rgba(0, 0, 0, 0.2);
}
```

---

## 按钮 3D 倾斜效果

### CSS 样式

```css
.link-btn {
  position: relative;
  transform-style: preserve-3d;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.link-btn:hover {
  box-shadow: 0 10px 40px rgba(79, 172, 254, 0.3);
}
```

### JavaScript 交互

```javascript
document.querySelectorAll('.link-btn').forEach(btn => {
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
});
```

### 发光效果

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

.link-btn:hover::before {
  opacity: 0.6;
}
```

---

## 头像效果

### CSS 样式

```css
.avatar-wrap {
  width: 96px;
  height: 96px;
  border-radius: 50%;
  overflow: hidden;
  border: 3px solid var(--card-border);
  box-shadow: 0 4px 16px rgba(79, 172, 254, 0.2);
  transition: transform 0.4s cubic-bezier(0.34, 1.56, 0.64, 1),
              box-shadow 0.4s ease;
}

.avatar-wrap:hover {
  transform: scale(1.1);
  box-shadow: 0 10px 40px rgba(79, 172, 254, 0.4);
}
```

### 边框动画

```css
.avatar-wrap::after {
  content: '';
  position: absolute;
  inset: -4px;
  background: var(--primary-gradient);
  border-radius: 50%;
  opacity: 0;
  transition: opacity 0.3s ease;
  z-index: -1;
}

.avatar-wrap:hover::after {
  opacity: 0.5;
}
```

---

## 卡片与容器

### 主容器

```css
.container {
  position: relative;
  z-index: 1;
  max-width: 500px;
  width: 100%;
  margin: 60px auto 40px;
  background: var(--card-bg);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--card-border);
  border-radius: 28px;
  box-shadow: var(--card-shadow);
  padding: 48px 40px 40px;
  text-align: center;
  transition: all 0.3s ease;
}
```

### 悬停效果

```css
.container:hover {
  box-shadow: 0 12px 48px rgba(79, 172, 254, 0.15);
  border-color: rgba(79, 172, 254, 0.3);
}
```

---

## 响应式设计

### 移动端 (< 600px)

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
}
```

### 平板 (600-900px)

```css
@media (min-width: 600px) and (max-width: 900px) {
  .container {
    max-width: 480px;
    padding: 44px 36px 36px;
  }
}
```

---

## 暗色模式适配

### 背景调整

```css
[data-theme="dark"] body {
  background: 
    linear-gradient(135deg, rgba(79, 172, 254, 0.2) 0%, rgba(0, 242, 254, 0.2) 100%),
    linear-gradient(225deg, rgba(102, 126, 234, 0.2) 0%, rgba(118, 75, 162, 0.2) 100%);
}
```

### 发光效果增强

```css
[data-theme="dark"] .link-btn:hover {
  box-shadow: 0 10px 40px rgba(79, 172, 254, 0.4);
}

[data-theme="dark"] .avatar-wrap:hover {
  box-shadow: 0 10px 40px rgba(79, 172, 254, 0.5);
}
```

---

## 动画时长规范

| 元素 | 状态变化 | 时长 | 缓动函数 |
|------|----------|------|----------|
| 按钮 | 悬停 | 0.3s | cubic-bezier(0.34, 1.56, 0.64, 1) |
| 头像 | 悬停 | 0.4s | cubic-bezier(0.34, 1.56, 0.64, 1) |
| 卡片 | 悬停 | 0.3s | ease |
| 背景 | 循环 | 20-25s | ease |
| 主题切换 | 切换 | 0.3s | ease |

---

## 文件变更

### 修改文件

1. `static/style.css` - 更新样式，添加动画和 3D 效果
2. `templates/index.html` - 添加必要的 HTML 结构和 JavaScript

### 新增文件

无

---

## 验收标准

1. ✅ 动态渐变背景正常运行，无卡顿
2. ✅ 按钮 3D 倾斜效果流畅，鼠标离开时平滑回正
3. ✅ 头像悬停放大效果正常，带蓝色发光
4. ✅ 暗色/亮色模式切换正常
5. ✅ 移动端响应式布局正确
6. ✅ 所有动画性能良好，不影响页面加载速度

---

## 下一步

批准此规范后，将进入实现阶段：
1. 更新 `static/style.css` 添加新样式
2. 更新 `templates/index.html` 添加 JavaScript 交互
3. 测试所有效果在不同设备和浏览器上的表现
