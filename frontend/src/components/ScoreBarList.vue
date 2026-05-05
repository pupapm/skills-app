<template>
  <div class="score-list">
    <div v-for="item in items" :key="item.key" class="score-item">
      <div class="score-head">
        <div>
          <div class="score-label">{{ item.label }}</div>
          <div class="score-sub">
            {{ formatNumber(item.value) }} / {{ formatNumber(item.max || max) }}
          </div>
        </div>

        <div class="score-percent">
          {{ percent(item).toFixed(0) }}%
        </div>
      </div>

      <div class="bar">
        <div
          class="bar-fill"
          :style="{ width: `${percent(item)}%` }"
        ></div>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  items: { type: Array, default: () => [] },
  max: { type: Number, default: 100 },
});

function itemMax(item) {
  return Number(item.max || props.max || 1);
}

function percent(item) {
  const max = itemMax(item);
  const value = Number(item.value || 0);
  return Math.max(0, Math.min(100, (value / max) * 100));
}

function formatNumber(value) {
  return Number(value || 0).toFixed(1);
}
</script>

<style scoped>
.score-list {
  display: grid;
  gap: 14px;
}

.score-item {
  display: grid;
  gap: 8px;
}

.score-head {
  display: flex;
  justify-content: space-between;
  gap: 14px;
  align-items: center;
}

.score-label {
  font-weight: 700;
  color: #111827;
}

.score-sub {
  margin-top: 2px;
  font-size: 12px;
  color: #6b7280;
}

.score-percent {
  font-size: 13px;
  font-weight: 700;
  color: #374151;
}

.bar {
  width: 100%;
  height: 10px;
  background: #e5e7eb;
  border-radius: 999px;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  background: #2563eb;
  border-radius: 999px;
  transition: width 0.2s ease;
}
</style>