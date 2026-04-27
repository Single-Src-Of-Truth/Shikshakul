import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable, shareReplay } from 'rxjs';
import { AcademicYear, ApiResponse } from '../../models/academic.models';
import { ACAD_API_URL } from '../../academic.config';

@Injectable({
  providedIn: 'root',
})
export class AcademicYearService {
  private http = inject(HttpClient);
  private baseUrl = inject(ACAD_API_URL);

  private get endpoint() {
    return `${this.baseUrl}/setup/academic-years`;
  }

  getAcademicYears(): Observable<ApiResponse<AcademicYear[]>> {
    return this.http.get<ApiResponse<AcademicYear[]>>(this.endpoint);
  }

  private currentYear$?: Observable<ApiResponse<AcademicYear>>;

  getCurrentAcademicYear(): Observable<ApiResponse<AcademicYear>> {
    if (!this.currentYear$) {
      this.currentYear$ = this.http
        .get<ApiResponse<AcademicYear>>(`${this.endpoint}/current`)
        .pipe(shareReplay(1));
    }
    return this.currentYear$;
  }

  clearCurrentYearCache(): void {
    this.currentYear$ = undefined;
  }

  createAcademicYear(
    data: Partial<AcademicYear>,
  ): Observable<ApiResponse<AcademicYear>> {
    return this.http.post<ApiResponse<AcademicYear>>(this.endpoint, data);
  }

  updateAcademicYear(
    id: string,
    data: Partial<AcademicYear>,
  ): Observable<ApiResponse<AcademicYear>> {
    return this.http.put<ApiResponse<AcademicYear>>(
      `${this.endpoint}/${id}`,
      data,
    );
  }

  deleteAcademicYear(id: string): Observable<ApiResponse<void>> {
    return this.http.delete<ApiResponse<void>>(`${this.endpoint}/${id}`);
  }
}
