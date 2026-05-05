<template>
  <AppShell
    title="Manager Members"
    subtitle="Browse members by role and year, with quick team-level summary."
  >
    <div class="card filters">
      <label>
        <span>Year</span>
        <input v-model="periodY" />
      </label>

      <label>
        <span>Role</span>
        <select v-model="role">
          <option value="ux">ux</option>
          <option value="qa">qa</option>
          <option value="ba">ba</option>
        </select>
      </label>

      <button @click="load" :disabled="loading">
        {{ loading ? "Loading..." : "Load Members" }}
      </button>
    </div>

    <div v-if="rows.length" class="stats">
      <div class="stat-card">
        <div class="stat-card__label">Members</div>
        <div class="stat-card__value">{{ rows.length }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-card__label">Average Skill</div>
        <div class="stat-card__value">{{ averageSkill.toFixed(2) }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-card__label">Average Credit</div>
        <div class="stat-card__value">{{ averageCredit.toFixed(2) }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-card__label">Flagged Rows</div>
        <div class="stat-card__value">{{ flaggedCount }}</div>
      </div>
    </div>

    <div v-if="rows.length" class="grid-two">
      <div class="card">
        <h3>Skill Distribution</h3>
        <SparkTrend :values="rows.map(x => x.skill_total)" />
      </div>
      <div class="card">
        <h3>Credit Distribution</h3>
        <SparkTrend :values="rows.map(x => x.credit)" />
      </div>
    </div>

    <div class="card" v-if="rows.length">
      <table class="table">
        <thead>
          <tr>
            <th>employee_id</th>
            <th>name</th>
            <th>team</th>
            <th>skill</th>
            <th>credit</th>
            <th>flag</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="r in rows" :key="r.employee_id">
            <td>{{ r.employee_id }}</td>
            <td>{{ r.full_name }}</td>
            <td>{{ r.team }}</td>
            <td>{{ Number(r.skill_total).toFixed(2) }}</td>
            <td>{{ Number(r.credit).toFixed(2) }}</td>
            <td>{{ r.sanity_flag || "-" }}</td>
            <td>
              <button class="small" @click="go(r.employee_id)">View</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-else class="card empty">No member rows loaded yet.</div>

    <p v-if="error" class="error">{{ error }}</p>
  </AppShell>
</template>

<script setup>
import { computed, ref } from "vue";
import { api } from "../../api/client";
import { useRouter } from "vue-router";
import AppShell from "../../layouts/AppShell.vue";
import SparkTrend from "../../components/SparkTrend.vue";
import { avg } from "../../utils/score";

const router = useRouter();
const periodY = ref(String(new Date().getFullYear()));
const role = ref("qa");
const rows = ref([]);
const error = ref("");
const loading = ref(false);

const averageSkill = computed(() => avg(rows.value, "skill_total"));
const averageCredit = computed(() => avg(rows.value, "credit"));
const flaggedCount = computed(() => rows.value.filter((x) => x.sanity_flag).length);

async function load() {
  error.value = "";
  rows.value = [];
  loading.value = true;
  try {
    rows.value = await api.managerMembers(periodY.value, role.value);
  } catch (e) {
    error.value = e.message || String(e);
  } finally {
    loading.value = false;
  }
}

function go(id) {
  router.push(`/manager/members/${id}?period_y=${encodeURIComponent(periodY.value)}`);
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
.filters {
  display: flex;
  gap: 14px;
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
button.small {
  padding: 8px 12px;
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