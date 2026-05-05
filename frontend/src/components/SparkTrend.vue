<template>
  <div class="spark-wrap">
    <svg viewBox="0 0 100 32" preserveAspectRatio="none" class="spark">
      <polyline
        v-if="points.length > 1"
        :points="polyline"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
    </svg>
  </div>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  values: { type: Array, default: () => [] },
});

const points = computed(() => {
  const vals = props.values.map((v) => Number(v || 0));
  if (!vals.length) return [];

  const max = Math.max(...vals, 1);
  const min = Math.min(...vals, 0);
  const range = max - min || 1;

  return vals.map((v, i) => {
    const x = vals.length === 1 ? 50 : (i / (vals.length - 1)) * 100;
    const y = 30 - ((v - min) / range) * 26;
    return [x, y];
  });
});

const polyline = computed(() =>
  points.value.map(([x, y]) => `${x},${y}`).join(" ")
);
</script>

<style scoped>
.spark-wrap {
  width: 100%;
  height: 48px;
  color: #2563eb;
}
.spark {
  width: 100%;
  height: 100%;
  display: block;
}
</style>