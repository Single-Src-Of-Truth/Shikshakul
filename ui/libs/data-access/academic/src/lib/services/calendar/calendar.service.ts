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

    private get baseUrl() {
        return `${this.apiUrl}/academics/calendar/events`;
    }

    getEvents(month?: number, year?: number): Observable<ApiResponse<EventResponse[]>> {
        let params = new HttpParams();
        if (month !== undefined) {
            params = params.set('month', month.toString());
        }
        if (year !== undefined) {
            params = params.set('year', year.toString());
        }

        return this.http.get<ApiResponse<EventResponse[]>>(this.baseUrl, { params });
    }

    createEvent(data: CreateEventRequest): Observable<ApiResponse<EventResponse>> {
        return this.http.post<ApiResponse<EventResponse>>(this.baseUrl, data);
    }

    updateEvent(id: string, data: UpdateEventRequest): Observable<ApiResponse<EventResponse>> {
        return this.http.put<ApiResponse<EventResponse>>(`${this.baseUrl}/${id}`, data);
    }

    deleteEvent(id: string): Observable<ApiResponse<void>> {
        return this.http.delete<ApiResponse<void>>(`${this.baseUrl}/${id}`);
    }
}
