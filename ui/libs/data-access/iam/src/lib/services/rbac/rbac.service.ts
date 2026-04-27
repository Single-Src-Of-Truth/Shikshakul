import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import {
    IamApiResponse,
    Role,
    Permission,
    CreateRoleRequest,
    UpdateRoleRequest,
} from '../../models/iam.models';
import { IAM_API_URL } from '../../iam.config';

@Injectable({ providedIn: 'root' })
export class RbacService {
    private http = inject(HttpClient);
    private baseUrl = inject(IAM_API_URL);

    listRoles(): Observable<IamApiResponse<Role[]>> {
        return this.http.get<IamApiResponse<Role[]>>(
            `${this.baseUrl}/roles`,
            { withCredentials: true },
        );
    }

    createRole(payload: CreateRoleRequest): Observable<IamApiResponse<Role>> {
        return this.http.post<IamApiResponse<Role>>(
            `${this.baseUrl}/roles`,
            payload,
            { withCredentials: true },
        );
    }

    updateRole(id: string, payload: UpdateRoleRequest): Observable<IamApiResponse<null>> {
        return this.http.patch<IamApiResponse<null>>(
            `${this.baseUrl}/roles/${id}`,
            payload,
            { withCredentials: true },
        );
    }

    deleteRole(id: string): Observable<IamApiResponse<null>> {
        return this.http.delete<IamApiResponse<null>>(
            `${this.baseUrl}/roles/${id}`,
            { withCredentials: true },
        );
    }

    listPermissions(): Observable<IamApiResponse<Permission[]>> {
        return this.http.get<IamApiResponse<Permission[]>>(
            `${this.baseUrl}/permissions`,
            { withCredentials: true },
        );
    }
}
