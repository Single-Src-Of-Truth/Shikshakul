import { Component, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { InviteService } from '@shikshakul/data-access/iam';

@Component({
    selector: 'skl-accept-invite',
    standalone: true,
    imports: [CommonModule, FormsModule, RouterLink],
    templateUrl: './accept-invite.component.html',
    styleUrl: './accept-invite.component.scss',
})
export class AcceptInviteComponent {
    private inviteService = inject(InviteService);
    private route = inject(ActivatedRoute);
    private router = inject(Router);

    inviteToken = this.route.snapshot.queryParamMap.get('token') ?? '';
    firstName = '';
    lastName = '';
    password = '';
    isLoading = signal(false);
    error = signal<string | null>(null);
    success = signal(false);

    onSubmit(): void {
        if (!this.inviteToken || !this.firstName || !this.password) return;
        this.isLoading.set(true);
        this.error.set(null);

        this.inviteService
            .acceptInvite({
                invite_token: this.inviteToken,
                first_name: this.firstName,
                last_name: this.lastName,
                password: this.password,
            })
            .subscribe({
                next: () => {
                    this.isLoading.set(false);
                    this.success.set(true);
                },
                error: (err) => {
                    this.isLoading.set(false);
                    this.error.set(
                        err?.error?.error ??
                        'Invalid or expired invite. Please contact your administrator.'
                    );
                },
            });
    }
}
