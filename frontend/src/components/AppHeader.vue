<template>
  <header class="app-header">
    <div class="app-header__left">
      <strong>Skills Intelligence System</strong>
      <span v-if="session.employee_id" class="muted">
        {{ session.employee_id }} <span v-if="session.role">· {{ session.role }}</span>
      </span>
    </div>

    <nav class="app-header__nav">
      <template v-if="session.is_manager">
        <router-link to="/manager/members">Members</router-link>
      </template>

      <template v-else-if="session.employee_id">
        <router-link to="/employee/dashboard">Dashboard</router-link>
        <router-link to="/employee/submit">Submit</router-link>
        <router-link to="/employee/history">History</router-link>
      </template>

      <button v-if="session.employee_id" @click="logout">Logout</button>
    </nav>
  </header>
</template>

<script setup>
import { computed } from "vue";
import { useRouter } from "vue-router";
import { clearSession } from "../api/client";

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
.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  padding: 14px 20px;
  border-bottom: 1px solid #e5e7eb;
  margin-bottom: 20px;
}
.app-header__left {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}
.app-header__nav {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}
.app-header__nav a {
  text-decoration: none;
}
.muted {
  color: #666;
  font-size: 14px;
}
button {
  cursor: pointer;
}
</style>