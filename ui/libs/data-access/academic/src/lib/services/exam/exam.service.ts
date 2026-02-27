import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { ACAD_API_URL } from '../../academic.config';
import {
  CreateExamScheduleRequest,
  CreateExamTermRequest,
  ExamSchedule,
  ExamTerm,
} from '../../models/academic.models';

@Injectable({
  providedIn: 'root',
})
export class ExamService {
  private http = inject(HttpClient);
  private apiUrl = inject(ACAD_API_URL);

  private get baseUrl() {
    return `${this.apiUrl}/academics/exams`;
  }

  getExamTerms(academicYearId?: string) {
    const params: any = {};
    if (academicYearId) params.academic_year_id = academicYearId;
    return this.http.get<ExamTerm[]>(`${this.baseUrl}/terms`, { params });
  }

  createExamTerm(data: CreateExamTermRequest) {
    return this.http.post<ExamTerm>(`${this.baseUrl}/terms`, data);
  }

  getExamSchedules(filters: any = {}) {
    return this.http.get<ExamSchedule[]>(`${this.baseUrl}/schedules`, {
      params: filters,
    });
  }

  createExamSchedule(data: CreateExamScheduleRequest) {
    return this.http.post<ExamSchedule>(`${this.baseUrl}/schedules`, data);
  }

  deleteExamSchedule(id: string) {
    return this.http.delete(`${this.baseUrl}/schedules/${id}`);
  }
}
