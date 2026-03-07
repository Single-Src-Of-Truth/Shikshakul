import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ACAD_API_URL } from '../../academic.config';
import {
    ApiResponse,
    MarkAttendanceRequest,
    AttendanceRegisterResponse,
    StudentAttendanceHistory,
} from '../../models/academic.models';

@Injectable({
    providedIn: 'root',
})
export class AttendanceService {
    private http = inject(HttpClient);
    private apiUrl = inject(ACAD_API_URL);

    private get baseUrl() {
        return `${this.apiUrl}/academics/attendance`;
    }

    markAttendance(data: MarkAttendanceRequest): Observable<ApiResponse<void>> {
        return this.http.post<ApiResponse<void>>(this.baseUrl, data);
    }

    getClassRegister(
        sectionId: string,
        date?: string
    ): Observable<ApiResponse<AttendanceRegisterResponse>> {
        const params: any = {};
        if (date) {
            params.date = date; // Backend likely expects a date parameter query
        }
        return this.http.get<ApiResponse<AttendanceRegisterResponse>>(
            `${this.baseUrl}/section/${sectionId}`,
            { params }
        );
    }

    getStudentHistory(
        studentId: string,
        startDate?: string,
        endDate?: string
    ): Observable<ApiResponse<StudentAttendanceHistory[]>> {
        const params: any = {};
        if (startDate) params.start_date = startDate;
        if (endDate) params.end_date = endDate;

        return this.http.get<ApiResponse<StudentAttendanceHistory[]>>(
            `${this.baseUrl}/student/${studentId}`,
            { params }
        );
    }
}
