import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { AcademicYear } from '../../models/academic.models';
import { ACAD_API_URL } from '../../academic.config';

@Injectable({
  providedIn: 'root',
})
export class AcademicYearService {
  private http = inject(HttpClient);
  private baseUrl = inject(ACAD_API_URL);

  private get endpoint() {
    return `${this.baseUrl}/academics/setup/academic-years`;
  }

  getAcademicYears(): Observable<AcademicYear[]> {
    return this.http.get<AcademicYear[]>(this.endpoint);
  }

  getCurrentAcademicYear(): Observable<AcademicYear> {
    return this.http.get<AcademicYear>(`${this.endpoint}/current`);
  }

  createAcademicYear(data: Partial<AcademicYear>): Observable<AcademicYear> {
    return this.http.post<AcademicYear>(this.endpoint, data);
  }

  updateAcademicYear(
    id: string,
    data: Partial<AcademicYear>,
  ): Observable<AcademicYear> {
    return this.http.put<AcademicYear>(`${this.endpoint}/${id}`, data);
  }

  deleteAcademicYear(id: string): Observable<void> {
    return this.http.delete<void>(`${this.endpoint}/${id}`);
  }
}
