import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { CreateRoutineRequest, Routine } from '../../models/academic.models';
import { ACAD_API_URL } from '../../academic.config';

@Injectable({
  providedIn: 'root',
})
export class TimetableService {
  private http = inject(HttpClient);
  private baseUrl = inject(ACAD_API_URL);

  private get apiUrl(): string {
    return `${this.baseUrl}/academics/timetables`;
  }

  getTimetable(filters: {
    section_id?: string;
    teacher_id?: string;
    day?: string;
  }): Observable<Routine[]> {
    return this.http.get<Routine[]>(this.apiUrl, { params: filters });
  }

  createRoutine(data: CreateRoutineRequest): Observable<any> {
    return this.http.post<any>(this.apiUrl, data);
  }

  updateRoutine(id: string, data: CreateRoutineRequest): Observable<any> {
    return this.http.put<any>(`${this.apiUrl}/${id}`, data);
  }

  deleteRoutine(id: string): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`);
  }
}
