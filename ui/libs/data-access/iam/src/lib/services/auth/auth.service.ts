import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import {
    IamApiResponse,
    LoginRequest,
    LoginResponse,
    ForgotPasswordRequest,
    ResetPasswordRequest,
    ChangePasswordRequest,
} from '../../models/iam.models';
import { IAM_API_URL } from '../../iam.config';

@Injectable({ providedIn: 'root' })
export class AuthService {
    private http = inject(HttpClient);
    private baseUrl = inject(IAM_API_URL);

    login(payload: LoginRequest): Observable<IamApiResponse<LoginResponse>> {
        return this.http.post<IamApiResponse<LoginResponse>>(
            `${this.baseUrl}/login`,
            payload,
            { withCredentials: true },
        );
    }

    logout(): Observable<IamApiResponse<null>> {
        return this.http.post<IamApiResponse<null>>(
            `${this.baseUrl}/logout`,
            {},
            { withCredentials: true },
        );
    }

    logoutAll(): Observable<IamApiResponse<null>> {
        return this.http.post<IamApiResponse<null>>(
            `${this.baseUrl}/sessions/logout-all`,
            {},
            { withCredentials: true },
        );
    }

    /**
     * Validates the current session cookie with the IAM service.
     * Returns user context (id, tenant, permissions) via response headers.
     * See: Skl-User-Id, Skl-Tenant-Id, Skl-Permissions
     */
    verify(): Observable<void> {
        return this.http.get<void>(`${this.baseUrl}/verify`, {
            withCredentials: true,
        });
    }

    forgotPassword(payload: ForgotPasswordRequest): Observable<IamApiResponse<null>> {
        return this.http.post<IamApiResponse<null>>(
            `${this.baseUrl}/password/forgot`,
            payload,
        );
    }

    resetPassword(payload: ResetPasswordRequest): Observable<IamApiResponse<null>> {
        return this.http.post<IamApiResponse<null>>(
            `${this.baseUrl}/password/reset`,
            payload,
        );
    }

    changePassword(payload: ChangePasswordRequest): Observable<IamApiResponse<null>> {
        return this.http.post<IamApiResponse<null>>(
            `${this.baseUrl}/profile/password/change`,
            payload,
            { withCredentials: true },
        );
    }
}
