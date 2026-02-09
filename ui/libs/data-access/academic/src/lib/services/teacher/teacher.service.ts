import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { ACAD_API_URL } from '../../academic.config';
import { Observable } from 'rxjs';
import { TeacherOnboardData } from '../../models/academic.models';

@Injectable({
  providedIn: 'root',
})
export class TeacherService {
  private http = inject(HttpClient);
  private baseUrl = inject(ACAD_API_URL);

  private get apiUrl(): string {
    return `${this.baseUrl}/academics/teachers`;
  }

  getTeachers(): Observable<any[]> {
    return this.http.get<any[]>(this.apiUrl);
  }

  onboardTeacher(data: TeacherOnboardData): Observable<any> {
    return this.http.post(`${this.apiUrl}/onboard`, data);
  }

  getTeacherById(id: string): Observable<any> {
    return this.http.get<any>(`${this.apiUrl}/${id}`);
  }

  updateTeacher(
    id: string,
    data: Partial<TeacherOnboardData>,
  ): Observable<any> {
    return this.http.put<any>(`${this.apiUrl}/${id}`, data);
  }

  deleteTeacher(id: string): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`);
  }

  approveTeacher(id: string, action: 'APPROVE' | 'REJECT'): Observable<any> {
    return this.http.post(`${this.apiUrl}/${id}/approve`, { action });
  }
}
