import { Component, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '@shikshakul/data-access/iam';
import { AuthStore } from '../../../auth.store';

@Component({
  selector: 'skl-login',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss',
})
export class LoginComponent {
  private authService = inject(AuthService);
  private authStore = inject(AuthStore);
  private router = inject(Router);

  tenantId = '';
  identifier = '';
  password = '';
  showPassword = signal(false);
  isLoading = signal(false);
  error = signal<string | null>(null);

  onSubmit(): void {
    if (!this.identifier || !this.password) return;

    this.isLoading.set(true);
    this.error.set(null);

    const tenantId = this.tenantId.trim();
    this.authService
      .login({
        tenant_id: tenantId || undefined,
        identifier: this.identifier,
        password: this.password,
      })
      .subscribe({
        next: (res) => {
          // Wait for the profile (and permissions) to be loaded before navigating
          // This prevents AuthGuard from bricking the navigation while state is null
          this.authStore.loadProfile().subscribe(() => {
            this.isLoading.set(false);

            // Map backend commands to UI routes
            const command = res.data?.redirect_command;
            let targetRoute = '/dashboard';

            if (command === 'NAV_PORTAL_ADMIN') targetRoute = '/dashboard';
            if (command === 'NAV_PORTAL_TEACHER') targetRoute = '/teacher/dashboard';
            if (command === 'NAV_PORTAL_PARENT') targetRoute = '/parent/dashboard';

            this.router.navigateByUrl(targetRoute);
          });
        },
        error: (err) => {
          this.isLoading.set(false);
          this.error.set(
            err?.error?.error ?? 'Invalid credentials. Please try again.'
          );
        },
      });
  }
}
