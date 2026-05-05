<template>
  <div class="auth-page">
    <div class="auth-card">
      <h2>Login</h2>
      <p class="muted">Sign in with your employee account.</p>

      <form @submit.prevent="login" class="form">
        <label>
          Employee ID
          <input v-model.trim="employee_id" placeholder="EMP001" required />
        </label>

        <label>
          Password
          <input v-model="password" type="password" placeholder="Password" required />
        </label>

        <button :disabled="loading">
          {{ loading ? "Logging in..." : "Login" }}
        </button>

        <router-link to="/register" class="link">Create account</router-link>

        <p v-if="error" class="error">{{ error }}</p>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { api, setSession, setToken } from "../api/client";
import { useRouter } from "vue-router";

const router = useRouter();

const employee_id = ref("");
const password = ref("");
const loading = ref(false);
const error = ref("");

async function login() {
  error.value = "";
  loading.value = true;

  try {
    const res = await api.login({
      employee_id: employee_id.value,
      password: password.value,
    });

    setToken(res.token);
    setSession(res.user);

    router.push(res.user.is_manager ? "/manager/members" : "/employee/dashboard");
  } catch (e) {
    error.value = e.message || "Login failed";
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  background: #f6f8fb;
  font-family: system-ui;
}
.auth-card {
  width: 100%;
  max-width: 420px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 18px;
  padding: 24px;
}
.muted {
  color: #6b7280;
}
.form {
  display: grid;
  gap: 14px;
}
label {
  display: grid;
  gap: 6px;
}
input {
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 10px;
}
button {
  background: #2563eb;
  color: white;
  border: none;
  padding: 11px 14px;
  border-radius: 10px;
  cursor: pointer;
}
.link {
  color: #2563eb;
  text-decoration: none;
  text-align: center;
}
.error {
  color: #dc2626;
}
</style>