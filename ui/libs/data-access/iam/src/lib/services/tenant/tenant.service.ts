import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import {
    IamApiResponse,
    Tenant,
    CreateTenantRequest,
    UpdateTenantRequest,
} from '../../models/iam.models';
import { IAM_API_URL } from '../../iam.config';

/**
 * TenantService manages tenants (schools).
 * SUPER ADMIN only — all endpoints require 'iam:tenants:*' permissions.
 * UI should gate access to these using PermissionGuard and AuthStore.hasPermission().
 */
@Injectable({ providedIn: 'root' })
export class TenantService {
    private http = inject(HttpClient);
    private baseUrl = inject(IAM_API_URL);

    listTenants(): Observable<IamApiResponse<Tenant[]>> {
        return this.http.get<IamApiResponse<Tenant[]>>(
            `${this.baseUrl}/tenants`,
            { withCredentials: true },
        );
    }

    getTenant(id: string): Observable<IamApiResponse<Tenant>> {
        return this.http.get<IamApiResponse<Tenant>>(
            `${this.baseUrl}/tenants/${id}`,
            { withCredentials: true },
        );
    }

    createTenant(payload: CreateTenantRequest): Observable<IamApiResponse<Tenant>> {
        return this.http.post<IamApiResponse<Tenant>>(
            `${this.baseUrl}/tenants`,
            payload,
            { withCredentials: true },
        );
    }

    updateTenant(id: string, payload: UpdateTenantRequest): Observable<IamApiResponse<null>> {
        return this.http.patch<IamApiResponse<null>>(
            `${this.baseUrl}/tenants/${id}`,
            payload,
            { withCredentials: true },
        );
    }
}
