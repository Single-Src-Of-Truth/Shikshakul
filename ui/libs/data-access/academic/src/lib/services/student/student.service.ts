import { HttpClient, HttpParams } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { ACAD_API_URL } from '../../academic.config';
import { Observable } from 'rxjs';
import { StudentOnboardData, ApiResponse } from '../../models/academic.models';

@Injectable({
  providedIn: 'root',
})
export class StudentService {
  private http = inject(HttpClient);
  private baseUrl = inject(ACAD_API_URL);

  private get apiUrl(): string {
    return `${this.baseUrl}/academics/students`;
  }

  getStudents(filters?: {
    class_id?: string;
    status?: string;
  }): Observable<ApiResponse<any[]>> {
    let params = new HttpParams();
    if (filters) {
      if (filters.class_id && filters.class_id !== 'undefined') {
        params = params.set('class_id', filters.class_id);
      }
      if (filters.status) {
        params = params.set('status', filters.status);
      }
    }
    return this.http.get<ApiResponse<any[]>>(this.apiUrl, { params });
  }

  onboardStudent(data: StudentOnboardData): Observable<any> {
    return this.http.post(`${this.apiUrl}/onboard`, data);
  }

  getStudentById(id: string): Observable<any> {
    return this.http.get<any>(`${this.apiUrl}/${id}`);
  }

  updateStudent(
    id: string,
    data: Partial<StudentOnboardData>,
  ): Observable<any> {
    return this.http.put<any>(`${this.apiUrl}/${id}`, data);
  }

  approveStudent(id: string): Observable<any> {
    return this.http.post<any>(`${this.apiUrl}/${id}/approve`, {});
  }

  deleteStudent(id: string): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`);
  }

  createSchema(data: any): Observable<any> {
    return this.http.post<any>(`${this.apiUrl}/schema`, data);
  }

  getActiveSchema(): Observable<any> {
    return this.http.get<any>(`${this.apiUrl}/schema`);
  }

  deleteSchema(): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/schema`);
  }

  configureSequence(data: any): Observable<any> {
    return this.http.post<any>(`${this.apiUrl}/setup/admission-sequence`, data);
  }
}
