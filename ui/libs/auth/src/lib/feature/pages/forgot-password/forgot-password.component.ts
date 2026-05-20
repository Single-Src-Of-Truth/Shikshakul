import { Component, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { AuthService } from '@shikshakul/data-access/iam';

@Component({
  selector: 'skl-forgot-password',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './forgot-password.component.html',
  styleUrl: './forgot-password.component.scss',
})
export class ForgotPasswordComponent {
  private authService = inject(AuthService);

  identifier = '';
  isLoading = signal(false);
  submitted = signal(false);

  onSubmit(): void {
    if (!this.identifier) return;
    this.isLoading.set(true);

    this.authService.forgotPassword({ identifier: this.identifier }).subscribe({
      next: () => {
        this.isLoading.set(false);
        this.submitted.set(true);
      },
      error: () => {
        // Backend always returns success for security (doesn't leak if user exists)
        this.isLoading.set(false);
        this.submitted.set(true);
      },
    });
  }
}
