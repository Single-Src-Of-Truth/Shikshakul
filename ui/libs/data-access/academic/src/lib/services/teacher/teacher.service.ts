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

  private get endpoint(): string {
    return `${this.baseUrl}/teachers`;
  }

  getTeachers(): Observable<any[]> {
    return this.http.get<any[]>(this.endpoint);
  }

  onboardTeacher(data: TeacherOnboardData): Observable<any> {
    return this.http.post(`${this.endpoint}/onboard`, data);
  }

  getTeacherById(id: string): Observable<any> {
    return this.http.get<any>(`${this.endpoint}/${id}`);
  }

  updateTeacher(
    id: string,
    data: Partial<TeacherOnboardData>,
  ): Observable<any> {
    return this.http.put<any>(`${this.endpoint}/${id}`, data);
  }

  deleteTeacher(id: string): Observable<void> {
    return this.http.delete<void>(`${this.endpoint}/${id}`);
  }

  approveTeacher(id: string, action: 'APPROVE' | 'REJECT'): Observable<any> {
    return this.http.post<any>(`${this.endpoint}/${id}/approve`, { action });
  }

  createSchema(data: any): Observable<any> {
    return this.http.post<any>(`${this.endpoint}/schema`, data);
  }

  getActiveSchema(): Observable<any> {
    return this.http.get<any>(`${this.endpoint}/schema`);
  }

  deleteSchema(): Observable<void> {
    return this.http.delete<void>(`${this.endpoint}/schema`);
  }

  assignClassTeacher(sectionId: string, teacherId: string): Observable<any> {
    return this.http.post<any>(`${this.endpoint}/sections/${sectionId}/class-teacher`, {
      teacher_id: teacherId,
    });
  }

  assignSubjectTeachers(sectionId: string, allocations: any[]): Observable<any> {
    return this.http.post<any>(
      `${this.endpoint}/sections/${sectionId}/subject-teachers`,
      { allocations },
    );
  }

  getSectionAllocations(sectionId: string): Observable<any> {
    return this.http.get<any>(
      `${this.endpoint}/sections/${sectionId}/allocations`,
    );
  }
}
