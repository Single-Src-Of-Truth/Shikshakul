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

  private get endpoint(): string {
    return `${this.baseUrl}/setup/subjects`;
  }

  getSubjects(): Observable<Subject[]> {
    return this.http.get<Subject[]>(this.endpoint);
  }

  createSubject(subject: Subject): Observable<Subject> {
    return this.http.post<Subject>(this.endpoint, subject);
  }

  updateSubject(id: string, subject: Partial<Subject>): Observable<Subject> {
    return this.http.put<Subject>(`${this.endpoint}/${id}`, subject);
  }

  deleteSubject(id: string): Observable<void> {
    return this.http.delete<void>(`${this.endpoint}/${id}`);
  }
}
