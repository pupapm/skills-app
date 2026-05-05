import { createRouter, createWebHistory } from "vue-router";

import Login from "../pages/Login.vue";
import Register from "../pages/Register.vue";

import EmployeeDashboard from "../pages/employee/Dashboard.vue";
import EmployeeSubmit from "../pages/employee/Submit.vue";
import EmployeeHistory from "../pages/employee/History.vue";

import ManagerMembers from "../pages/manager/Members.vue";
import ManagerEmployeeDetail from "../pages/manager/EmployeeDetail.vue";

import { getSession, getToken } from "../api/client";

const routes = [
  { path: "/", redirect: "/login" },
  { path: "/login", component: Login },
  { path: "/register", component: Register },

  { path: "/employee/dashboard", component: EmployeeDashboard },
  { path: "/employee/submit", component: EmployeeSubmit },
  { path: "/employee/history", component: EmployeeHistory },

  { path: "/manager/members", component: ManagerMembers },
  {
    path: "/manager/members/:employee_id",
    component: ManagerEmployeeDetail,
    props: true,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to) => {
  const token = getToken();
  const session = getSession();

  if (to.path === "/login" || to.path === "/register") {
    return true;
  }

  if (!token || !session.employee_id) {
    return "/login";
  }

  if (to.path.startsWith("/manager") && !session.is_manager) {
    return "/employee/dashboard";
  }

  return true;
});

export default router;