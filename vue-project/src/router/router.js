import { createRouter,createWebHistory} from "vue-router";

// 路由信息
const routes = [
    {
        path: "/",
        name: "Login",
        component:  () => import('../views/Login.vue'),
    },
    {
        path: "/login",
        name: "Login",
        component:  () => import('../views/Login.vue'),
    },
    {
        path: "/register",
        name: "Register",
        component:  () => import('../views/Register.vue'),
    },
    {
        path: "/home",
        name: "Home",
        component:  () => import('../views/Home.vue'),
    },
];

// 导出路由
const router = createRouter({
    history: createWebHistory(),
    routes
});

export default router;