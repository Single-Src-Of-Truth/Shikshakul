import { Component, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '@shikshakul/data-access/iam';

@Component({
  selector: 'skl-reset-password',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './reset-password.component.html',
  styleUrl: './reset-password.component.scss',
})
export class ResetPasswordComponent {
  private authService = inject(AuthService);
  private router = inject(Router);

  identifier = '';
  otp = '';
  newPassword = '';
  isLoading = signal(false);
  error = signal<string | null>(null);

  onSubmit(): void {
    if (!this.identifier || this.otp.length !== 6 || !this.newPassword) return;
    this.isLoading.set(true);
    this.error.set(null);

    this.authService
      .resetPassword({
        identifier: this.identifier,
        otp: this.otp,
        new_password: this.newPassword,
      })
      .subscribe({
        next: () => {
          this.isLoading.set(false);
          this.router.navigate(['/auth/login']);
        },
        error: (err) => {
          this.isLoading.set(false);
          this.error.set(
            err?.error?.error ??
            'Invalid OTP or the code has expired. Please try again.'
          );
        },
      });
  }
}
