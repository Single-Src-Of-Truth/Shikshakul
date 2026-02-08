import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Subject } from '../../models/academic.models';
import { ACAD_API_URL } from '../../academic.config';

@Injectable({
  providedIn: 'root',
})
export class SubjectService {
  private http = inject(HttpClient);
  private baseUrl = inject(ACAD_API_URL);

  private get apiUrl(): string {
    return `${this.baseUrl}/academics/setup/subjects`;
  }

  getSubjects(): Observable<Subject[]> {
    return this.http.get<Subject[]>(this.apiUrl);
  }

  createSubject(subject: Subject): Observable<Subject> {
    return this.http.post<Subject>(this.apiUrl, subject);
  }

  updateSubject(id: string, subject: Partial<Subject>): Observable<Subject> {
    return this.http.put<Subject>(`${this.apiUrl}/${id}`, subject);
  }

  deleteSubject(id: string): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`);
  }
}
