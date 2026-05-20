import { Route } from '@angular/router';

export const authRoutes: Route[] = [
    {
        path: 'login',
        loadComponent: () =>
            import('./pages/login/login.component').then((m) => m.LoginComponent),
    },
    {
        path: 'forgot-password',
        loadComponent: () =>
            import('./pages/forgot-password/forgot-password.component').then(
                (m) => m.ForgotPasswordComponent,
            ),
    },
    {
        path: 'reset-password',
        loadComponent: () =>
            import('./pages/reset-password/reset-password.component').then(
                (m) => m.ResetPasswordComponent,
            ),
    },
    {
        path: 'accept-invite',
        loadComponent: () =>
            import('./pages/accept-invite/accept-invite.component').then(
                (m) => m.AcceptInviteComponent,
            ),
    },
    {
        path: '',
        redirectTo: 'login',
        pathMatch: 'full',
    },
];
