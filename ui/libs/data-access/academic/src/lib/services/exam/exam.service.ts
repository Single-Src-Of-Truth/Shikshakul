import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { ACAD_API_URL } from '../../academic.config';
import {
  CreateExamScheduleRequest,
  CreateExamTermRequest,
  ExamSchedule,
  ExamTerm,
  MarkSheetResponse,
  SubmitMarksRequest,
  ApiResponse,
  GenerateResultRequest,
  PublishResultRequest,
  ReportCardResponse,
} from '../../models/academic.models';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class ExamService {
  private http = inject(HttpClient);
  private apiUrl = inject(ACAD_API_URL);

  private get baseUrl() {
    return `${this.apiUrl}/academics/exams`;
  }

  getExamTerms(academicYearId?: string): Observable<ApiResponse<ExamTerm[]>> {
    const params: any = {};
    if (academicYearId) params.academic_year_id = academicYearId;
    return this.http.get<ApiResponse<ExamTerm[]>>(`${this.baseUrl}/terms`, {
      params,
    });
  }

  createExamTerm(
    data: CreateExamTermRequest,
  ): Observable<ApiResponse<ExamTerm>> {
    return this.http.post<ApiResponse<ExamTerm>>(`${this.baseUrl}/terms`, data);
  }

  getExamSchedules(filters: any = {}): Observable<ApiResponse<ExamSchedule[]>> {
    return this.http.get<ApiResponse<ExamSchedule[]>>(
      `${this.baseUrl}/schedules`,
      {
        params: filters,
      },
    );
  }

  getExamSchedule(id: string): Observable<ApiResponse<ExamSchedule>> {
    return this.http.get<ApiResponse<ExamSchedule>>(
      `${this.baseUrl}/schedules/${id}`,
    );
  }

  createExamSchedule(
    data: CreateExamScheduleRequest,
  ): Observable<ApiResponse<ExamSchedule>> {
    return this.http.post<ApiResponse<ExamSchedule>>(
      `${this.baseUrl}/schedules`,
      data,
    );
  }

  deleteExamSchedule(id: string): Observable<ApiResponse<void>> {
    return this.http.delete<ApiResponse<void>>(
      `${this.baseUrl}/schedules/${id}`,
    );
  }

  getMarksSheet(
    scheduleId: string,
  ): Observable<ApiResponse<MarkSheetResponse>> {
    return this.http.get<ApiResponse<MarkSheetResponse>>(
      `${this.baseUrl}/schedules/${scheduleId}/marks`,
    );
  }

  submitMarks(data: SubmitMarksRequest): Observable<ApiResponse<void>> {
    return this.http.post<ApiResponse<void>>(`${this.baseUrl}/marks`, data);
  }

  generateResults(data: GenerateResultRequest): Observable<ApiResponse<void>> {
    return this.http.post<ApiResponse<void>>(
      `${this.apiUrl}/results/generate`,
      data,
    );
  }

  publishResults(data: PublishResultRequest): Observable<ApiResponse<void>> {
    return this.http.post<ApiResponse<void>>(
      `${this.apiUrl}/results/publish`,
      data,
    );
  }

  getReportCard(
    studentId: string,
    termId: string,
  ): Observable<ApiResponse<ReportCardResponse>> {
    return this.http.get<ApiResponse<ReportCardResponse>>(
      `${this.apiUrl}/results/students/${studentId}`,
      {
        params: { term_id: termId },
      },
    );
  }
}
