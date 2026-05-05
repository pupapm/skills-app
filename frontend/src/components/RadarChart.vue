<template>
  <div class="radar-card">
    <svg :viewBox="`0 0 ${size} ${size}`" class="radar-svg">
      <polygon
        v-for="r in ringSteps"
        :key="`ring-${r}`"
        :points="polygonPoints((radius * r) / ringLevels)"
        class="grid-ring"
      />

      <line
        v-for="(item, i) in items"
        :key="`axis-${item.key}`"
        :x1="center"
        :y1="center"
        :x2="axisPoint(i, radius).x"
        :y2="axisPoint(i, radius).y"
        class="axis-line"
      />

      <polygon
        v-if="items.length >= 3"
        :points="dataPolygon"
        class="data-fill"
      />

      <polygon
        v-if="items.length >= 3"
        :points="dataPolygon"
        class="data-line"
      />

      <circle
        v-for="(item, i) in items"
        :key="`point-${item.key}`"
        :cx="valuePoint(i, normalizedValue(item)).x"
        :cy="valuePoint(i, normalizedValue(item)).y"
        r="4"
        class="data-point"
      />

      <g v-for="(item, i) in items" :key="`value-${item.key}`">
        <text
          :x="valuePoint(i, normalizedValue(item)).x"
          :y="valuePoint(i, normalizedValue(item)).y - 10"
          text-anchor="middle"
          class="value-label"
        >
          {{ displayValue(item) }}
        </text>
      </g>

      <g v-for="(item, i) in items" :key="`label-${item.key}`">
        <text
          :x="labelPoint(i).x"
          :y="labelPoint(i).y"
          :text-anchor="labelAnchor(labelPoint(i).x)"
          dominant-baseline="middle"
          class="outer-label"
        >
          {{ item.label }}
        </text>
      </g>
    </svg>
  </div>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  items: { type: Array, default: () => [] },
  max: { type: Number, default: 100 },
  size: { type: Number, default: 420 },
  ringLevels: { type: Number, default: 5 },
});

const size = computed(() => props.size);
const center = computed(() => props.size / 2);
const radius = computed(() => props.size * 0.28);

const ringSteps = computed(() =>
  Array.from({ length: props.ringLevels }, (_, i) => i + 1)
);

function angleFor(i) {
  const n = props.items.length || 1;
  return -Math.PI / 2 + (i * 2 * Math.PI) / n;
}

function axisPoint(i, r) {
  const a = angleFor(i);
  return {
    x: center.value + Math.cos(a) * r,
    y: center.value + Math.sin(a) * r,
  };
}

function valuePoint(i, value01) {
  return axisPoint(i, radius.value * value01);
}

function labelPoint(i) {
  return axisPoint(i, radius.value + 42);
}

function labelAnchor(x) {
  if (Math.abs(x - center.value) < 10) return "middle";
  return x < center.value ? "end" : "start";
}

function itemMax(item) {
  return Number(item.max || props.max || 1);
}

function normalizedValue(item) {
  const max = itemMax(item);
  const value = Number(item.value || 0);
  return Math.max(0, Math.min(1, value / max));
}

function displayValue(item) {
  return `${Number(item.value || 0).toFixed(1)}/${Number(itemMax(item)).toFixed(0)}`;
}

function polygonPoints(r) {
  return props.items
    .map((_, i) => {
      const p = axisPoint(i, r);
      return `${p.x},${p.y}`;
    })
    .join(" ");
}

const dataPolygon = computed(() =>
  props.items
    .map((item, i) => {
      const p = valuePoint(i, normalizedValue(item));
      return `${p.x},${p.y}`;
    })
    .join(" ")
);
</script>

<style scoped>
.radar-card {
  width: 100%;
  overflow: hidden;
}

.radar-svg {
  width: 100%;
  height: auto;
  display: block;
}

.grid-ring {
  fill: none;
  stroke: #e5e7eb;
  stroke-width: 1;
}

.axis-line {
  stroke: #e5e7eb;
  stroke-width: 1;
}

.data-fill {
  fill: rgba(37, 99, 235, 0.18);
  stroke: none;
}

.data-line {
  fill: none;
  stroke: #2563eb;
  stroke-width: 2.5;
}

.data-point {
  fill: #2563eb;
}

.outer-label {
  font-size: 11px;
  fill: #374151;
}

.value-label {
  font-size: 10px;
  fill: #111827;
  font-weight: 600;
}
</style>