import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import {
    IamApiResponse,
    UserProfile,
    Session,
    UpdateProfileRequest,
} from '../../models/iam.models';
import { IAM_API_URL } from '../../iam.config';

@Injectable({ providedIn: 'root' })
export class ProfileService {
    private http = inject(HttpClient);
    private baseUrl = inject(IAM_API_URL);

    getMyProfile(): Observable<IamApiResponse<UserProfile>> {
        return this.http.get<IamApiResponse<UserProfile>>(
            `${this.baseUrl}/profile/me`,
            { withCredentials: true },
        );
    }

    updateMyProfile(payload: UpdateProfileRequest): Observable<IamApiResponse<null>> {
        return this.http.patch<IamApiResponse<null>>(
            `${this.baseUrl}/profile/me`,
            payload,
            { withCredentials: true },
        );
    }

    getMySessions(): Observable<IamApiResponse<Session[]>> {
        return this.http.get<IamApiResponse<Session[]>>(
            `${this.baseUrl}/profile/sessions`,
            { withCredentials: true },
        );
    }
}
