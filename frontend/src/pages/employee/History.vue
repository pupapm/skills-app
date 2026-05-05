<template>
  <AppShell
    title="My History"
    subtitle="See yearly trends for score and contribution credit."
  >
    <div class="card toolbar">
      <label>
        <span>Limit</span>
        <input v-model.number="limit" type="number" min="1" max="24" />
      </label>
      <button @click="load" :disabled="loading">
        {{ loading ? "Loading..." : "Load History" }}
      </button>
    </div>

    <div v-if="rows.length" class="stats">
      <div class="stat-card">
        <div class="stat-card__label">Average Skill</div>
        <div class="stat-card__value">{{ averageSkill.toFixed(2) }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-card__label">Average Credit</div>
        <div class="stat-card__value">{{ averageCredit.toFixed(2) }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-card__label">Best Skill</div>
        <div class="stat-card__value">{{ bestSkill.toFixed(2) }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-card__label">Lowest Skill</div>
        <div class="stat-card__value">{{ worstSkill.toFixed(2) }}</div>
      </div>
    </div>

    <div v-if="rows.length" class="grid-two">
      <div class="card">
        <h3>Skill Trend</h3>
        <SparkTrend :values="rows.map(x => x.skill_total)" />
      </div>
      <div class="card">
        <h3>Credit Trend</h3>
        <SparkTrend :values="rows.map(x => x.credit)" />
      </div>
    </div>

    <div class="card" v-if="rows.length">
      <table class="table">
        <thead>
          <tr>
            <th>Year</th>
            <th>Skill Score</th>
            <th>Credit</th>
            <th>Sanity Flag</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="r in rows" :key="r.period_y">
            <td>{{ r.period_y }}</td>
            <td>{{ Number(r.skill_total).toFixed(2) }}</td>
            <td>{{ Number(r.credit).toFixed(2) }}</td>
            <td>{{ r.sanity_flag || "-" }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-else-if="loaded" class="card empty">
      No history found.
    </div>

    <p v-if="error" class="error">{{ error }}</p>
  </AppShell>
</template>

<script setup>
import { computed, ref } from "vue";
import { api } from "../../api/client";
import AppShell from "../../layouts/AppShell.vue";
import SparkTrend from "../../components/SparkTrend.vue";
import { avg, maxValue, minValue } from "../../utils/score";

const limit = ref(12);
const rows = ref([]);
const error = ref("");
const loading = ref(false);
const loaded = ref(false);

const averageSkill = computed(() => avg(rows.value, "skill_total"));
const averageCredit = computed(() => avg(rows.value, "credit"));
const bestSkill = computed(() => maxValue(rows.value, "skill_total"));
const worstSkill = computed(() => minValue(rows.value, "skill_total"));

async function load() {
  error.value = "";
  rows.value = [];
  loading.value = true;
  loaded.value = false;

  try {
    rows.value = await api.history(limit.value);
    loaded.value = true;
  } catch (e) {
    error.value = e.message || String(e);
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.card {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  padding: 18px;
  margin-bottom: 18px;
}
.toolbar {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: end;
  flex-wrap: wrap;
}
label {
  display: grid;
  gap: 6px;
}
button {
  border: none;
  background: #2563eb;
  color: white;
  padding: 10px 14px;
  border-radius: 10px;
  cursor: pointer;
}
.stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(210px, 1fr));
  gap: 14px;
  margin-bottom: 18px;
}
.stat-card {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 18px;
  padding: 18px;
}
.stat-card__label {
  color: #6b7280;
  font-size: 13px;
  margin-bottom: 8px;
}
.stat-card__value {
  font-size: 30px;
  font-weight: 800;
}
.grid-two {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 18px;
}
.table {
  width: 100%;
  border-collapse: collapse;
}
.table th,
.table td {
  padding: 12px 10px;
  border-bottom: 1px solid #f0f0f0;
  text-align: left;
}
.empty {
  color: #6b7280;
}
.error {
  color: #dc2626;
}
@media (max-width: 900px) {
  .grid-two {
    grid-template-columns: 1fr;
  }
}
</style>