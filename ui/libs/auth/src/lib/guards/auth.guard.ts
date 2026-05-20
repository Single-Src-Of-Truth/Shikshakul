import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthStore } from '../auth.store';

/**
 * AuthGuard — prevents unauthenticated users from accessing protected routes.
 *
 * Usage in route config:
 * ```ts
 * { path: 'dashboard', canActivate: [authGuard], component: DashboardComponent }
 * ```
 */
export const authGuard: CanActivateFn = () => {
    const store = inject(AuthStore);
    const router = inject(Router);

    if (store.isAuthenticated()) {
        return true;
    }

    return router.createUrlTree(['/auth/login']);
};
