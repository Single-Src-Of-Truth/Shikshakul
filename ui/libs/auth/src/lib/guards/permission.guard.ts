import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthStore } from '../auth.store';

/**
 * PermissionGuard factory — protects routes based on IAM permissions.
 *
 * Users without the required permission see a 403 page, not a login redirect.
 * This is the mechanism used to hide Super Admin routes (e.g., /tenants)
 * from regular school staff.
 *
 * Usage in route config:
 * ```ts
 * {
 *   path: 'tenants',
 *   canActivate: [permissionGuard('iam:tenants:read')],
 *   loadComponent: () => import(...)
 * }
 * ```
 */
export const permissionGuard = (requiredPermission: string): CanActivateFn => {
    return () => {
        const store = inject(AuthStore);
        const router = inject(Router);

        if (!store.isAuthenticated()) {
            return router.createUrlTree(['/auth/login']);
        }

        if (store.hasPermission(requiredPermission)) {
            return true;
        }

        // Show a 403 forbidden page rather than redirecting to login
        return router.createUrlTree(['/forbidden']);
    };
};
