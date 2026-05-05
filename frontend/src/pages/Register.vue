<template>
  <div class="auth-page">
    <div class="auth-card">
      <h2>Register</h2>
      <p class="muted">Create an employee account. Manager access is controlled separately by the managers table.</p>

      <form @submit.prevent="register" class="form">
        <label>
          Employee ID
          <input v-model.trim="form.employee_id" placeholder="EMP001" required />
        </label>

        <label>
          Full Name
          <input v-model.trim="form.full_name" placeholder="Full Name" required />
        </label>

        <label>
          Role
          <select v-model="form.role">
            <option value="ux">UX/UI</option>
            <option value="qa">QA / Tester</option>
            <option value="ba">Business Analyst</option>
          </select>
        </label>

        <label>
          Team
          <input v-model.trim="form.team" placeholder="Team" required />
        </label>

        <label>
          Password
          <input v-model="form.password" type="password" placeholder="At least 6 characters" required />
        </label>

        <button :disabled="loading">
          {{ loading ? "Creating..." : "Create Account" }}
        </button>

        <router-link to="/login" class="link">Back to login</router-link>

        <p v-if="error" class="error">{{ error }}</p>
      </form>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from "vue";
import { api, setSession, setToken } from "../api/client";
import { useRouter } from "vue-router";

const router = useRouter();

const form = reactive({
  employee_id: "",
  full_name: "",
  role: "ux",
  team: "",
  password: "",
});

const loading = ref(false);
const error = ref("");

async function register() {
  error.value = "";
  loading.value = true;

  try {
    const res = await api.register(form);
    setToken(res.token);
    setSession(res.user);
    router.push(res.user.is_manager ? "/manager/members" : "/employee/dashboard");
  } catch (e) {
    error.value = e.message || "Register failed";
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
  max-width: 460px;
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
input,
select {
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