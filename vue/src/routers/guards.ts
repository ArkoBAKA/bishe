/*
 * @Date: 2026-04-23 14:47:28
 * @LastEditors: ArongWang 3312428832@qq.com
 * @LastEditTime: 2026-04-23 14:47:29
 * @FilePath: /vue/src/routers/guards.ts
 * @Description: 
 */
import type { NavigationGuardNext, RouteLocationNormalized } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const getIsAuthed = () => {
    const auth = useAuthStore()
    return auth.isAuthed
}

const getIsAdmin = () => {
    const auth = useAuthStore()
    return auth.isAdmin
}

export const authGuard = (to: RouteLocationNormalized, from: RouteLocationNormalized, next: NavigationGuardNext) => {
    if (!to.meta.requiresAuth) return next()
    if (getIsAuthed()) return next()
    next({ name: 'mobile-login', query: { redirect: to.fullPath } })
}

export const adminGuard = (to: RouteLocationNormalized, from: RouteLocationNormalized, next: NavigationGuardNext) => {
    if (!to.meta.requiresAdmin) return next()
    if (getIsAuthed() && getIsAdmin()) return next()
    next({ name: 'mobile-home' })
}

