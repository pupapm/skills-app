<template>
  <div class="shell">
    <aside class="sidebar">
      <div class="brand">
        <div class="brand__title">Skills Intelligence</div>
        <div class="brand__sub">Internal Scoring System</div>
      </div>

      <div v-if="session.employee_id" class="profile">
        <div class="profile__id">{{ session.employee_id }}</div>
        <div class="profile__meta">
          <span class="pill">{{ session.role || "unknown" }}</span>
          <span v-if="session.is_manager" class="pill pill--manager">manager</span>
        </div>
      </div>

      <nav class="nav">
        <template v-if="session.is_manager">
          <router-link to="/manager/members" class="nav__link">Members</router-link>
        </template>

        <template v-else>
          <router-link to="/employee/dashboard" class="nav__link">Dashboard</router-link>
          <router-link to="/employee/submit" class="nav__link">Submit</router-link>
          <router-link to="/employee/history" class="nav__link">History</router-link>
        </template>
      </nav>

      <button type="button" class="logout" @click="logout">Logout</button>
    </aside>

    <div class="main-wrap">
      <header class="topbar">
        <div class="topbar__inner">
          <h1 class="page-title">{{ title }}</h1>
          <p v-if="subtitle" class="page-subtitle">{{ subtitle }}</p>
        </div>
      </header>

      <main class="content">
        <section class="page-body">
          <slot />
        </section>
      </main>
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue";
import { useRouter } from "vue-router";
import { clearSession } from "../api/client";

defineProps({
  title: { type: String, default: "" },
  subtitle: { type: String, default: "" },
});

const router = useRouter();

const session = computed(() => {
  try {
    return JSON.parse(localStorage.getItem("session") || "{}");
  } catch {
    return {};
  }
});

function logout() {
  clearSession();
  router.push("/login");
}
</script>

<style scoped>
.shell {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 260px minmax(0, 1fr);
  background: #f6f8fb;
  position: relative;
  isolation: isolate;
}

.sidebar {
  background: #111827;
  color: white;
  padding: 24px 18px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  position: relative;
  z-index: 2;
}

.brand__title {
  font-size: 20px;
  font-weight: 700;
}

.brand__sub {
  font-size: 13px;
  color: #9ca3af;
  margin-top: 4px;
}

.profile {
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.04);
  border-radius: 14px;
  padding: 14px;
}

.profile__id {
  font-weight: 700;
  margin-bottom: 8px;
}

.profile__meta {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.pill {
  display: inline-block;
  padding: 4px 10px;
  border-radius: 999px;
  background: #1f2937;
  font-size: 12px;
  text-transform: uppercase;
}

.pill--manager {
  background: #0f766e;
}

.nav {
  display: grid;
  gap: 8px;
}

.nav__link {
  color: #d1d5db;
  text-decoration: none;
  padding: 10px 12px;
  border-radius: 10px;
  transition: background 0.15s ease, color 0.15s ease;
}

.nav__link:hover {
  background: rgba(255, 255, 255, 0.06);
  color: white;
}

.nav__link.router-link-active {
  background: #2563eb;
  color: white;
}

.logout {
  margin-top: auto;
  border: none;
  background: #dc2626;
  color: white;
  padding: 10px 12px;
  border-radius: 10px;
  cursor: pointer;
  font: inherit;
}

.main-wrap {
  min-width: 0;
  display: flex;
  flex-direction: column;
  position: relative;
  z-index: 1;
}

.topbar {
  background: white;
  border-bottom: 1px solid #e5e7eb;
  position: relative;
  z-index: 1;
}

.topbar__inner {
  padding: 24px 28px;
}

.page-title {
  margin: 0;
  font-size: 28px;
  line-height: 1.2;
}

.page-subtitle {
  margin: 6px 0 0;
  color: #6b7280;
}

.content {
  min-width: 0;
  position: relative;
  z-index: 1;
}

.page-body {
  padding: 24px 28px;
  position: relative;
  z-index: 1;
}

@media (max-width: 900px) {
  .shell {
    grid-template-columns: 1fr;
  }

  .sidebar {
    gap: 14px;
  }

  .topbar__inner,
  .page-body {
    padding: 18px;
  }
}
</style>