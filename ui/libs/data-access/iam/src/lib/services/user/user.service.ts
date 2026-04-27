import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import {
    IamApiResponse,
    UpdateUserStatusRequest,
    ChangeUserRoleRequest,
} from '../../models/iam.models';
import { IAM_API_URL } from '../../iam.config';

/**
 * UserService manages users within a school (tenant).
 * All actions require 'iam:users:execute' or 'iam:users:delete' permissions.
 */
@Injectable({ providedIn: 'root' })
export class UserService {
    private http = inject(HttpClient);
    private baseUrl = inject(IAM_API_URL);

    updateStatus(userId: string, payload: UpdateUserStatusRequest): Observable<IamApiResponse<null>> {
        return this.http.patch<IamApiResponse<null>>(
            `${this.baseUrl}/users/${userId}/status`,
            payload,
            { withCredentials: true },
        );
    }

    changeRole(userId: string, payload: ChangeUserRoleRequest): Observable<IamApiResponse<null>> {
        return this.http.patch<IamApiResponse<null>>(
            `${this.baseUrl}/users/${userId}/role`,
            payload,
            { withCredentials: true },
        );
    }

    deleteUser(userId: string): Observable<IamApiResponse<null>> {
        return this.http.delete<IamApiResponse<null>>(
            `${this.baseUrl}/users/${userId}`,
            { withCredentials: true },
        );
    }
}
