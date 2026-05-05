<template>
  <AppShell
    title="My Dashboard"
    subtitle="View current yearly score, component breakdown, and maintain a project portfolio as part of your profile."
  >
    <div class="toolbar card">
      <label>
        <span>Year</span>
        <input v-model="periodY" placeholder="2026" />
      </label>

      <div class="toolbar__actions">
        <button type="button" @click="loadLatest" :disabled="loadingLatest">
          {{ loadingLatest ? "Loading..." : "Load Score" }}
        </button>

        <router-link to="/employee/submit" class="btn-link secondary-link">
          Go to Submit
        </router-link>
      </div>
    </div>

    <div v-if="latest" class="stats">
      <div class="stat-card stat-card--main">
        <div class="stat-card__label">Final Score</div>
        <div class="stat-card__value">{{ Number(latest.skill_total).toFixed(2) }}</div>
        <div class="stat-card__sub">
          base {{ Number(latest?.breakdown?.base || 0).toFixed(1) }}
          + credit {{ Number(latest.credit).toFixed(1) }}
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-card__label">Performance Badge</div>
        <div class="rank-badge" :class="rankBadge.className">
          {{ rankBadge.label }}
        </div>
        <div class="stat-card__sub">{{ rankBadge.description }}</div>
      </div>

      <div class="stat-card">
        <div class="stat-card__label">Contribution Credit</div>
        <div class="stat-card__value">{{ Number(latest.credit).toFixed(2) }}</div>
      </div>

      <div class="stat-card">
        <div class="stat-card__label">Year</div>
        <div class="stat-card__text">{{ latest.period_y }}</div>
      </div>
    </div>

    <div v-if="latest?.sanity_flag" class="card warning-card">
      ⚠️ Sanity Flag: {{ latest.sanity_flag }}
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

    <section class="profile-section">
      <div class="section-head">
        <div>
          <div class="section-kicker">PROFILE</div>
          <h3 class="section-title">Project Portfolio</h3>
          <p class="section-subtitle">
            This appears as part of your profile and can be viewed by managers for assignment decisions.
          </p>
        </div>

        <button type="button" @click="openCreateForm">
          + Add Project
        </button>
      </div>

      <div v-if="showForm" class="editor card">
        <div class="editor__head">
          <h4>{{ editingId ? "Edit Project" : "Add Project" }}</h4>

          <div class="editor__actions">
            <button type="button" class="ghost" @click="cancelEdit">
              Cancel
            </button>

            <button type="button" @click="saveForm" :disabled="savingForm">
              {{ savingForm ? "Saving..." : editingId ? "Update Project" : "Save Project" }}
            </button>
          </div>
        </div>

        <div class="editor__grid">
          <label>
            <span>Project Name</span>
            <input v-model="form.project_name" placeholder="Customer Portal Revamp" />
          </label>

          <label>
            <span>Year</span>
            <input v-model="form.project_year" placeholder="2025" />
          </label>

          <label>
            <span>Responsibility</span>
            <input v-model="form.responsibility" placeholder="UX Designer / QA Engineer / BA" />
          </label>

          <label>
            <span>Tools Used</span>
            <input v-model="form.tools_used" placeholder="Figma, Jira, Postman" />
          </label>

          <label class="full">
            <span>What did you do?</span>
            <textarea
              v-model="form.what_you_did"
              rows="4"
              placeholder="Describe responsibilities and contribution"
            />
          </label>

          <label class="full">
            <span>Outcome</span>
            <textarea
              v-model="form.outcome"
              rows="3"
              placeholder="What changed, improved, or was delivered"
            />
          </label>
        </div>
      </div>

      <div v-if="loadingProjects" class="card empty">
        Loading projects...
      </div>

      <div v-else-if="projects.length" class="project-list">
        <div v-for="project in projects" :key="project.id" class="project-row card">
          <div class="project-main">
            <div>
              <div class="project-name">{{ project.project_name }}</div>
              <div class="project-meta">
                {{ project.project_year }} · {{ project.responsibility || "No responsibility added" }}
              </div>
            </div>

            <div class="project-tools">
              {{ project.tools_used || "-" }}
            </div>

            <div class="project-actions">
              <button type="button" class="link-btn" @click="toggleExpand(project.id)">
                {{ expandedId === project.id ? "Hide" : "View" }}
              </button>
              <button type="button" class="link-btn" @click="startEdit(project)">
                Edit
              </button>
              <button type="button" class="link-btn danger-text" @click="removeProject(project.id)">
                Delete
              </button>
            </div>
          </div>

          <div v-if="expandedId === project.id" class="project-detail">
            <div>
              <div class="detail-label">What did you do?</div>
              <div class="detail-value">{{ project.what_you_did || "-" }}</div>
            </div>

            <div>
              <div class="detail-label">Outcome</div>
              <div class="detail-value">{{ project.outcome || "-" }}</div>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="card empty">
        No projects added yet.
      </div>
    </section>

    <div v-if="emptyLatest" class="card empty">
      No score found for this year yet.
    </div>

    <p v-if="error" class="error">{{ error }}</p>
  </AppShell>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { api, getSession } from "../../api/client";
import AppShell from "../../layouts/AppShell.vue";
import ScoreBarList from "../../components/ScoreBarList.vue";
import RadarChart from "../../components/RadarChart.vue";
import { formatRoleBreakdown } from "../../utils/score";

const session = getSession();
const periodY = ref(String(new Date().getFullYear()));

const latest = ref(null);
const error = ref("");
const loadingLatest = ref(false);
const emptyLatest = ref(false);

const projects = ref([]);
const loadingProjects = ref(false);

const showForm = ref(false);
const savingForm = ref(false);
const editingId = ref(null);
const expandedId = ref(null);

const form = ref(blankForm());
onMounted(() => {
  loadProjects();
});

const breakdownItems = computed(() =>
  formatRoleBreakdown(session.role, latest.value?.breakdown).filter(
    (x) => !["credit", "base", "final"].includes(x.key)
  )
);

const maxComponent = computed(() =>
  Math.max(...breakdownItems.value.map((x) => x.value), 1)
);

const rankBadge = computed(() => {
  const score = Number(latest.value?.skill_total || 0);

  if (score >= 95) {
    return {
      label: "Top Performer",
      className: "rank-top",
      description: "Excellent yearly performance and contribution.",
    };
  }

  if (score >= 85) {
    return {
      label: "Strong Performer",
      className: "rank-strong",
      description: "High capability with strong yearly contribution.",
    };
  }

  if (score >= 70) {
    return {
      label: "Solid Performer",
      className: "rank-solid",
      description: "Reliable performance with good baseline capability.",
    };
  }

  if (score >= 50) {
    return {
      label: "Developing",
      className: "rank-developing",
      description: "Shows progress but still has improvement areas.",
    };
  }

  return {
    label: "Needs Support",
    className: "rank-support",
    description: "Needs support, training, or closer assignment fit.",
  };
});

function blankForm() {
  return {
    project_name: "",
    project_year: "",
    responsibility: "",
    what_you_did: "",
    tools_used: "",
    outcome: "",
  };
}

async function loadLatest() {
  error.value = "";
  latest.value = null;
  emptyLatest.value = false;
  loadingLatest.value = true;

  try {
    latest.value = await api.latestScore(periodY.value);
  } catch (e) {
    if ((e.message || "").toLowerCase().includes("not found")) {
      emptyLatest.value = true;
    } else {
      error.value = e.message || String(e);
    }
  } finally {
    loadingLatest.value = false;
  }
}

async function loadProjects() {
  error.value = "";
  loadingProjects.value = true;

  try {
    const data = await api.myProjects();
    projects.value = Array.isArray(data) ? data : [];
  } catch (e) {
    projects.value = [];
    error.value = e.message || String(e);
  } finally {
    loadingProjects.value = false;
  }
}

function openCreateForm() {
  editingId.value = null;
  form.value = blankForm();
  showForm.value = true;
}

function startEdit(project) {
  editingId.value = project.id;
  form.value = {
    project_name: project.project_name || "",
    project_year: project.project_year || "",
    responsibility: project.responsibility || "",
    what_you_did: project.what_you_did || "",
    tools_used: project.tools_used || "",
    outcome: project.outcome || "",
  };
  showForm.value = true;
}

function cancelEdit() {
  showForm.value = false;
  editingId.value = null;
  form.value = blankForm();
}

function toggleExpand(id) {
  expandedId.value = expandedId.value === id ? null : id;
}

async function saveForm() {
  error.value = "";
  savingForm.value = true;

  try {
    if (editingId.value) {
      await api.updateProject(editingId.value, form.value);
    } else {
      await api.createProject(form.value);
    }

    await loadProjects();
    cancelEdit();
  } catch (e) {
    error.value = e.message || String(e);
  } finally {
    savingForm.value = false;
  }
}

async function removeProject(id) {
  error.value = "";

  try {
    await api.deleteProject(id);
    if (expandedId.value === id) expandedId.value = null;
    await loadProjects();
  } catch (e) {
    error.value = e.message || String(e);
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
  margin-bottom: 18px;
}

.toolbar label {
  display: grid;
  gap: 6px;
}

.toolbar__actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  align-items: center;
}

button,
.btn-link {
  border: none;
  background: #2563eb;
  color: white;
  padding: 10px 14px;
  border-radius: 10px;
  cursor: pointer;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font: inherit;
  user-select: none;
}

button.secondary,
.secondary-link {
  background: #111827;
}

button.ghost {
  background: #6b7280;
}

button:disabled {
  opacity: 0.65;
  cursor: not-allowed;
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

.stat-card--main {
  border-color: #bfdbfe;
}

.stat-card__label {
  color: #6b7280;
  font-size: 13px;
  margin-bottom: 8px;
}

.stat-card__value {
  font-size: 34px;
  font-weight: 800;
}

.stat-card__text {
  font-size: 18px;
  font-weight: 600;
}

.stat-card__sub {
  margin-top: 6px;
  font-size: 12px;
  color: #6b7280;
}

.rank-badge {
  display: inline-flex;
  padding: 8px 12px;
  border-radius: 999px;
  font-size: 14px;
  font-weight: 700;
}

.rank-top {
  background: #dcfce7;
  color: #166534;
}

.rank-strong {
  background: #dbeafe;
  color: #1d4ed8;
}

.rank-solid {
  background: #fef3c7;
  color: #92400e;
}

.rank-developing {
  background: #ffedd5;
  color: #c2410c;
}

.rank-support {
  background: #fee2e2;
  color: #b91c1c;
}

.warning-card {
  background: #fffbeb;
  border-color: #fde68a;
  color: #92400e;
}

.grid-two {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 18px;
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

.editor {
  margin-bottom: 18px;
}

.editor__head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 14px;
  flex-wrap: wrap;
}

.editor__head h4 {
  margin: 0;
}

.editor__actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.editor__grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.editor__grid label {
  display: grid;
  gap: 6px;
}

.editor__grid label.full {
  grid-column: 1 / -1;
}

.editor__grid input,
.editor__grid textarea,
.toolbar input {
  width: 100%;
  border: 1px solid #d1d5db;
  border-radius: 10px;
  padding: 10px 12px;
  font: inherit;
  background: white;
  box-sizing: border-box;
}

.project-list {
  display: grid;
  gap: 12px;
}

.project-row {
  margin-bottom: 0;
}

.project-main {
  display: grid;
  grid-template-columns: 1.4fr 1fr auto;
  gap: 16px;
  align-items: center;
}

.project-name {
  font-weight: 700;
}

.project-meta {
  margin-top: 4px;
  color: #6b7280;
  font-size: 13px;
}

.project-tools {
  color: #374151;
  font-size: 14px;
}

.project-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.link-btn {
  border: none;
  background: none;
  padding: 0;
  color: #2563eb;
  cursor: pointer;
  font: inherit;
}

.danger-text {
  color: #dc2626;
}

.project-detail {
  margin-top: 14px;
  padding-top: 14px;
  border-top: 1px solid #e5e7eb;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 18px;
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

.empty {
  color: #6b7280;
}

.error {
  color: #dc2626;
  margin-top: 12px;
}

@media (max-width: 900px) {
  .grid-two,
  .editor__grid,
  .project-main,
  .project-detail {
    grid-template-columns: 1fr;
  }

  .editor__grid label.full {
    grid-column: auto;
  }

  .project-actions {
    justify-content: flex-start;
  }
}
</style>