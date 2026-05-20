import { HttpClient, HttpParams } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { ACAD_API_URL } from '../../academic.config';
import { Observable } from 'rxjs';
import {
  StudentOnboardData,
  ApiResponse,
  StudentStatus,
  StudentAction,
} from '../../models/academic.models';

@Injectable({
  providedIn: 'root',
})
export class StudentService {
  private http = inject(HttpClient);
  private baseUrl = inject(ACAD_API_URL);

  private get endpoint(): string {
    return `${this.baseUrl}/students`;
  }

  getStudents(filters?: {
    class_id?: string;
    status?: StudentStatus;
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
    return this.http.get<ApiResponse<any[]>>(this.endpoint, { params });
  }

  getActiveStudents(class_id?: string): Observable<ApiResponse<any[]>> {
    let params = new HttpParams();
    if (class_id && class_id !== 'undefined') {
      params = params.set('class_id', class_id);
    }
    return this.http.get<ApiResponse<any[]>>(`${this.endpoint}`, {
      params,
    });
  }

  onboardStudent(data: StudentOnboardData): Observable<any> {
    return this.http.post(`${this.endpoint}/onboard`, data);
  }

  getStudentById(id: string): Observable<any> {
    return this.http.get<any>(`${this.endpoint}/${id}`);
  }

  updateStudent(
    id: string,
    data: Partial<StudentOnboardData>,
  ): Observable<any> {
    return this.http.put<any>(`${this.endpoint}/${id}`, data);
  }

  approveStudent(
    id: string,
    data: { action: StudentAction; section_id?: string },
  ): Observable<any> {
    return this.http.post<any>(`${this.endpoint}/${id}/approve`, data);
  }

  deleteStudent(id: string): Observable<void> {
    return this.http.delete<void>(`${this.endpoint}/${id}`);
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

  configureSequence(data: any): Observable<any> {
    return this.http.post<any>(`${this.endpoint}/setup/admission-sequence`, data);
  }
}
