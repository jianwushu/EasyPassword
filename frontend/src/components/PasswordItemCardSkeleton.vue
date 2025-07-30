<template>
  <div class="skeleton-card">
    <div class="skeleton-header">
      <div class="skeleton-logo"></div>
      <div class="skeleton-website">
        <div class="skeleton-line title"></div>
        <div class="skeleton-line subtitle"></div>
      </div>
    </div>
    <div class="skeleton-body">
      <div class="skeleton-line"></div>
      <div class="skeleton-line"></div>
    </div>
    <div class="skeleton-footer">
      <div class="skeleton-tag"></div>
      <div class="skeleton-dots"></div>
    </div>
  </div>
</template>

<style scoped>
.skeleton-card {
  background-color: var(--card-bg-color, #fff);
  border: 1px solid var(--card-border-color, #e0e0e0);
  border-radius: 12px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.skeleton-line, .skeleton-logo, .skeleton-tag, .skeleton-dots {
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: pulse 1.5s infinite ease-in-out;
  border-radius: 4px;
}

html[data-theme='dark'] .skeleton-line,
html[data-theme='dark'] .skeleton-logo,
html[data-theme='dark'] .skeleton-tag,
html[data-theme='dark'] .skeleton-dots {
  background: linear-gradient(90deg, #3a3a3a 25%, #424242 50%, #3a3a3a 75%);
  background-size: 200% 100%;
}


@keyframes pulse {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}

.skeleton-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.skeleton-logo {
  width: 40px;
  height: 40px;
  border-radius: 8px;
}

.skeleton-website {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.skeleton-line {
  height: 12px;
}

.skeleton-line.title {
  width: 60%;
  height: 14px;
}

.skeleton-line.subtitle {
  width: 40%;
}

.skeleton-body {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 8px;
}

.skeleton-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 8px;
}

.skeleton-tag {
  width: 80px;
  height: 24px;
  border-radius: 16px;
}

.skeleton-dots {
  width: 24px;
  height: 24px;
  border-radius: 50%;
}
</style>