<template>
  <AppShell
    title="Employee Detail"
    subtitle="Inspect one employee’s current yearly result, component breakdown, history, and profile portfolio."
  >
    <div class="card filters">
      <label>
        <span>Year</span>
        <input v-model="periodY" />
      </label>
      <button @click="load" :disabled="loadingDetail">
        {{ loadingDetail ? "Loading..." : "Load Detail" }}
      </button>
      <button class="secondary" @click="loadHistory" :disabled="loadingHistory">
        {{ loadingHistory ? "Loading..." : "Load History" }}
      </button>
    </div>

    <div v-if="detail" class="detail-grid">
      <div class="card">
        <h3>{{ detail.full_name }}</h3>
        <p><strong>Employee ID:</strong> {{ detail.employee_id }}</p>
        <p><strong>Role:</strong> {{ detail.role }}</p>
        <p><strong>Team:</strong> {{ detail.team }}</p>
        <p><strong>Year:</strong> {{ detail.period_y }}</p>
        <p><strong>Sanity Flag:</strong> {{ detail.sanity_flag || "-" }}</p>
      </div>

      <div class="card score-card">
        <div class="score-label">Final Score</div>
        <div class="score-value">{{ Number(detail.skill_total).toFixed(2) }}</div>
        <div class="score-sub">
          base {{ detail?.breakdown?.base?.toFixed(1)}}
         + credit {{ Number(detail.credit).toFixed(1) }}</div>
      </div>
    </div>

    <div v-if="breakdownItems.length" class="grid-two">
      <div class="card">
        <h3>Component Breakdown</h3>
        <RadarChart :items="breakdownItems" :max="maxComponent" />
      </div>

      <div class="card">
        <h3>Component Scores</h3>
        <ScoreBarList :items="breakdownItems" :max="maxComponent" />
      </div>
    </div>

    <div v-if="history.length" class="grid-two">
      <div class="card">
        <h3>Skill Trend</h3>
        <SparkTrend :values="history.map(x => x.skill_total)" />
      </div>
      <div class="card">
        <h3>Credit Trend</h3>
        <SparkTrend :values="history.map(x => x.credit)" />
      </div>
    </div>

    <div v-if="history.length" class="card">
      <h3>History</h3>
      <ul class="history-list">
        <li v-for="h in history" :key="h.period_y">
          <span>{{ h.period_y }}</span>
          <strong>{{ Number(h.skill_total).toFixed(2) }}</strong>
          <span>credit {{ Number(h.credit).toFixed(2) }}</span>
        </li>
      </ul>
    </div>

    <section class="profile-section">
      <div class="section-head">
        <div>
          <div class="section-kicker">PROFILE</div>
          <h3 class="section-title">Project Portfolio</h3>
          <p class="section-subtitle">
            Supporting project history and practical experience for assignment decisions.
          </p>
        </div>
      </div>

      <div class="table-card">
        <table v-if="projects.length" class="portfolio-table">
          <thead>
            <tr>
              <th>Project Name</th>
              <th>Year</th>
              <th>Responsibility</th>
              <th>Tools Used</th>
              <th class="actions-col"></th>
            </tr>
          </thead>
          <tbody>
            <template v-for="project in projects" :key="project.id">
              <tr>
                <td>{{ project.project_name }}</td>
                <td>{{ project.project_year }}</td>
                <td>{{ project.responsibility || "-" }}</td>
                <td>{{ project.tools_used || "-" }}</td>
                <td class="actions-col">
                  <button type="button" class="link-btn" @click="toggleExpand(project.id)">
                    {{ expandedId === project.id ? "Hide" : "View" }}
                  </button>
                </td>
              </tr>

              <tr v-if="expandedId === project.id" class="detail-row">
                <td colspan="5">
                  <div class="detail-grid-2">
                    <div>
                      <div class="detail-label">What they did</div>
                      <div class="detail-value">{{ project.what_you_did || "-" }}</div>
                    </div>
                    <div>
                      <div class="detail-label">Outcome</div>
                      <div class="detail-value">{{ project.outcome || "-" }}</div>
                    </div>
                  </div>
                </td>
              </tr>
            </template>
          </tbody>
        </table>

        <div v-else-if="loadingProjects" class="empty-state">
          Loading projects...
        </div>

        <div v-else class="empty-state">
          No projects found for this employee.
        </div>
      </div>
    </section>

    <p v-if="error" class="error">{{ error }}</p>
  </AppShell>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { api } from "../../api/client";
import { useRoute } from "vue-router";
import AppShell from "../../layouts/AppShell.vue";
import ScoreBarList from "../../components/ScoreBarList.vue";
import SparkTrend from "../../components/SparkTrend.vue";
import RadarChart from "../../components/RadarChart.vue";
import { formatRoleBreakdown, parseBreakdown } from "../../utils/score";

const route = useRoute();
const employee_id = route.params.employee_id;

const periodY = ref(route.query.period_y || String(new Date().getFullYear()));
const detail = ref(null);
const history = ref([]);
const projects = ref([]);
const error = ref("");
const loadingDetail = ref(false);
const loadingHistory = ref(false);
const loadingProjects = ref(false);
const expandedId = ref(null);

const parsedBreakdown = computed(() => parseBreakdown(detail.value?.breakdown));
const breakdownItems = computed(() =>
  formatRoleBreakdown(detail.value?.role, detail.value?.breakdown).filter(
    (x) => x.key !== "credit"
  )
);
const maxComponent = computed(() =>
  Math.max(...breakdownItems.value.map((x) => x.value), 1)
);

async function load() {
  error.value = "";
  detail.value = null;
  loadingDetail.value = true;
  try {
    detail.value = await api.managerEmployee(employee_id, periodY.value);
  } catch (e) {
    error.value = e.message || String(e);
  } finally {
    loadingDetail.value = false;
  }
}

async function loadHistory() {
  error.value = "";
  history.value = [];
  loadingHistory.value = true;
  try {
    history.value = await api.managerEmployeeHistory(employee_id, 12);
  } catch (e) {
    error.value = e.message || String(e);
  } finally {
    loadingHistory.value = false;
  }
}

async function loadProjects() {
  error.value = "";
  projects.value = [];
  loadingProjects.value = true;
  try {
    projects.value = await api.managerEmployeeProjects(employee_id);
  } catch (e) {
    error.value = e.message || String(e);
  } finally {
    loadingProjects.value = false;
  }
}

function toggleExpand(id) {
  expandedId.value = expandedId.value === id ? null : id;
}

onMounted(() => {
  loadProjects();
});
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
  gap: 12px;
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

button.secondary {
  background: #111827;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1.2fr 0.8fr;
  gap: 16px;
}

.grid-two {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 18px;
}

.score-card {
  display: grid;
  align-content: center;
}

.score-label {
  color: #6b7280;
  font-size: 13px;
  margin-bottom: 8px;
}

.score-value {
  font-size: 38px;
  font-weight: 800;
}

.score-sub {
  margin-top: 8px;
  color: #6b7280;
}

.history-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: grid;
  gap: 10px;
}

.history-list li {
  display: flex;
  justify-content: space-between;
  gap: 14px;
  border-bottom: 1px solid #f0f0f0;
  padding-bottom: 10px;
}

.profile-section {
  margin-top: 28px;
}

.section-head {
  display: flex;
  justify-content: space-between;
  align-items: end;
  gap: 16px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.section-kicker {
  font-size: 12px;
  letter-spacing: 0.12em;
  color: #6b7280;
  margin-bottom: 8px;
}

.section-title {
  margin: 0;
  font-size: 18px;
}

.section-subtitle {
  margin: 4px 0 0;
  color: #6b7280;
}

.table-card {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  overflow: hidden;
}

.portfolio-table {
  width: 100%;
  border-collapse: collapse;
}

.portfolio-table thead th {
  font-size: 13px;
  font-weight: 600;
  color: #374151;
  text-align: left;
  padding: 16px 20px;
  border-bottom: 1px solid #e5e7eb;
  background: white;
}

.portfolio-table tbody td {
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
  vertical-align: top;
}

.actions-col {
  width: 120px;
}

.link-btn {
  border: none;
  background: none;
  padding: 0;
  color: #2563eb;
  cursor: pointer;
  font: inherit;
}

.detail-row td {
  background: #fafafa;
  padding-top: 0;
}

.detail-grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 18px;
  padding-top: 10px;
}

.detail-label {
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #6b7280;
  margin-bottom: 6px;
}

.detail-value {
  color: #111827;
  line-height: 1.5;
}

.empty-state {
  padding: 24px 20px;
  color: #6b7280;
}

.error {
  color: #dc2626;
}

@media (max-width: 900px) {
  .detail-grid,
  .grid-two,
  .detail-grid-2 {
    grid-template-columns: 1fr;
  }

  .portfolio-table {
    display: block;
    overflow-x: auto;
    white-space: nowrap;
  }
}
</style>