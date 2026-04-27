import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ACAD_API_URL } from '../../academic.config';
import {
    ApiResponse,
    CreateEventRequest,
    UpdateEventRequest,
    EventResponse,
} from '../../models/academic.models';

@Injectable({
    providedIn: 'root',
})
export class CalendarService {
    private http = inject(HttpClient);
    private apiUrl = inject(ACAD_API_URL);

    private get endpoint() {
        return `${this.apiUrl}/calendar/events`;
    }

    getEvents(month?: number, year?: number): Observable<ApiResponse<EventResponse[]>> {
        const params: any = {};
        if (month) params.month = month;
        if (year) params.year = year;
        return this.http.get<ApiResponse<EventResponse[]>>(this.endpoint, { params });
    }

    createEvent(event: Partial<EventResponse>): Observable<ApiResponse<EventResponse>> {
        return this.http.post<ApiResponse<EventResponse>>(this.endpoint, event);
    }

    updateEvent(id: string, event: Partial<EventResponse>): Observable<ApiResponse<EventResponse>> {
        return this.http.put<ApiResponse<EventResponse>>(`${this.endpoint}/${id}`, event);
    }

    deleteEvent(id: string): Observable<ApiResponse<void>> {
        return this.http.delete<ApiResponse<void>>(`${this.endpoint}/${id}`);
    }
}
