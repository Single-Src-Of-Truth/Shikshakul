import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ACAD_API_URL } from '../../academic.config';
import {
    ApiResponse,
    IDCardResponse,
    IssueTCRequest,
    TCResponse,
    IssueBonafideRequest,
    BonafideResponse,
} from '../../models/academic.models';

@Injectable({
    providedIn: 'root',
})
export class CertificateService {
    private http = inject(HttpClient);
    private apiUrl = inject(ACAD_API_URL);

    private get baseUrl() {
        return `${this.apiUrl}/academics/certificates`;
    }

    getIDCards(classId: string): Observable<ApiResponse<IDCardResponse[]>> {
        return this.http.get<ApiResponse<IDCardResponse[]>>(`${this.baseUrl}/id-cards`, {
            params: { class_id: classId },
        });
    }

    generateTC(data: IssueTCRequest): Observable<ApiResponse<TCResponse>> {
        return this.http.post<ApiResponse<TCResponse>>(`${this.baseUrl}/tc`, data);
    }

    generateBonafide(data: IssueBonafideRequest): Observable<ApiResponse<BonafideResponse>> {
        return this.http.post<ApiResponse<BonafideResponse>>(`${this.baseUrl}/bonafide`, data);
    }
}
