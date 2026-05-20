import { inject, Injectable, signal, computed } from '@angular/core';
import { Router } from '@angular/router';
import { Observable, of } from 'rxjs';
import { tap, map, catchError } from 'rxjs/operators';
import { AuthService, ProfileService, UserProfile, IamApiResponse } from '@shikshakul/data-access/iam';

/**
 * AuthStore — the single source of truth for authentication state.
 *
 * Uses Angular Signals for reactive, zoneless-compatible state management.
 * Consumed by guards (AuthGuard, PermissionGuard) and the layout/sidenav
 * to drive permission-based UI rendering.
 */
@Injectable({ providedIn: 'root' })
export class AuthStore {
    private authService = inject(AuthService);
    private profileService = inject(ProfileService);
    private router = inject(Router);

    // ── Private state signals ────────────────────────────────────────────────
    private _profile = signal<UserProfile | null>(null);
    private _isLoading = signal(false);
    private _error = signal<string | null>(null);

    // ── Public computed signals ──────────────────────────────────────────────
    readonly profile = this._profile.asReadonly();
    readonly isAuthenticated = computed(() => this._profile() !== null);
    readonly permissions = computed(() => this._profile()?.permissions ?? []);
    readonly roles = computed(() => this._profile()?.roles ?? []);
    readonly isLoading = this._isLoading.asReadonly();

    /**
     * Checks if the current user has a specific IAM permission.
     * Use this to conditionally render UI elements.
     *
     * @example
     * // In a component template:
     * // @if(authStore.hasPermission('iam:tenants:read')) { <a>Tenants</a> }
     */
    hasPermission(permission: string): boolean {
        return this.permissions().includes(permission);
    }

    /**
     * Checks if the current user has at least one of the given roles.
     */
    hasRole(...roles: string[]): boolean {
        return roles.some((r) => this.roles().includes(r));
    }

    /**
     * Loads the current user profile from the IAM service and populates state.
     * Call this once on app initialization (ideally in an APP_INITIALIZER).
     */
    loadProfile(): Observable<UserProfile | null> {
        this._isLoading.set(true);
        this._error.set(null);

        return this.profileService.getMyProfile().pipe(
            tap({
                next: (res: IamApiResponse<UserProfile>) => {
                    this._profile.set(res.data);
                    this._isLoading.set(false);
                },
                error: () => {
                    this._profile.set(null);
                    this._isLoading.set(false);
                }
            }),
            map(res => res.data),
            catchError(() => of(null))
        );
    }

    /**
     * Logs the user out: calls the IAM service, clears state, and redirects to login.
     */
    logout(): void {
        this.authService.logout().pipe(
            tap(() => {
                this._profile.set(null);
                this.router.navigate(['/auth/login']);
            }),
        ).subscribe();
    }
}
