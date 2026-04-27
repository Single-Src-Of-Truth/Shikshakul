import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import {
    IamApiResponse,
    GenerateInviteRequest,
    GenerateInviteResponse,
    AcceptInviteRequest,
    AcceptInviteResponse,
} from '../../models/iam.models';
import { IAM_API_URL } from '../../iam.config';

@Injectable({ providedIn: 'root' })
export class InviteService {
    private http = inject(HttpClient);
    private baseUrl = inject(IAM_API_URL);

    /**
     * Generates a one-time invite link/token for a user to join a school.
     * Requires 'iam:invites:execute' permission.
     */
    generateInvite(payload: GenerateInviteRequest): Observable<IamApiResponse<GenerateInviteResponse>> {
        return this.http.post<IamApiResponse<GenerateInviteResponse>>(
            `${this.baseUrl}/invites/generate`,
            payload,
            { withCredentials: true },
        );
    }

    /**
     * PUBLIC endpoint. Accepts an invite token (from email link) and creates the user account.
     * This is the entry point for the client-side verification flow.
     */
    acceptInvite(payload: AcceptInviteRequest): Observable<IamApiResponse<AcceptInviteResponse>> {
        return this.http.post<IamApiResponse<AcceptInviteResponse>>(
            `${this.baseUrl}/invites/accept`,
            payload,
        );
    }

    /**
     * Extends the expiry of a pending invite.
     * Requires 'iam:invites:execute' permission.
     */
    extendInvite(inviteId: string): Observable<IamApiResponse<null>> {
        return this.http.patch<IamApiResponse<null>>(
            `${this.baseUrl}/invites/${inviteId}/extend`,
            {},
            { withCredentials: true },
        );
    }

    /**
     * Cancels a pending invite.
     * Requires 'iam:invites:delete' permission.
     */
    cancelInvite(inviteId: string): Observable<IamApiResponse<null>> {
        return this.http.delete<IamApiResponse<null>>(
            `${this.baseUrl}/invites/${inviteId}`,
            { withCredentials: true },
        );
    }
}
