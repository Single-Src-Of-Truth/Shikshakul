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

  private get endpoint(): string {
    return `${this.baseUrl}/timetables`;
  }

  getTimetable(filters: {
    section_id?: string;
    teacher_id?: string;
    day?: string;
  }): Observable<Routine[]> {
    return this.http.get<Routine[]>(this.endpoint, { params: filters });
  }

  createRoutine(data: CreateRoutineRequest): Observable<any> {
    return this.http.post<any>(this.endpoint, data);
  }

  updateRoutine(id: string, data: CreateRoutineRequest): Observable<any> {
    return this.http.put<any>(`${this.endpoint}/${id}`, data);
  }

  deleteRoutine(id: string): Observable<void> {
    return this.http.delete<void>(`${this.endpoint}/${id}`);
  }
}
