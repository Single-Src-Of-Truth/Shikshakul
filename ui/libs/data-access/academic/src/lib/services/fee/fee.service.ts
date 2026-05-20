import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { ACAD_API_URL } from '../../academic.config';
import { Observable } from 'rxjs';
import { AcademicYear, FeeHead, FeeStructure, ApiResponse } from '../../models/academic.models';

@Injectable({
  providedIn: 'root',
})
export class FeeService {
  private http = inject(HttpClient);
  private baseUrl = inject(ACAD_API_URL);

  private get endpoint(): string {
    return `${this.baseUrl}/fees`;
  }

  getFeeHeads(): Observable<FeeHead[]> {
    return this.http.get<FeeHead[]>(`${this.endpoint}/heads`);
  }

  createFeeHead(data: FeeHead): Observable<FeeHead> {
    return this.http.post<FeeHead>(`${this.endpoint}/heads`, data);
  }

  deleteFeeHead(id: string): Observable<void> {
    return this.http.delete<void>(`${this.endpoint}/heads/${id}`);
  }

  getFeeStructure(classId: string, academicYearId: string): Observable<any[]> {
    return this.http.get<any[]>(`${this.endpoint}/structure`, {
      params: { class_id: classId, academic_year_id: academicYearId },
    });
  }

  createFeeStructure(data: FeeStructure): Observable<any> {
    return this.http.post(`${this.endpoint}/structure`, data);
  }

  deleteFeeStructure(id: string): Observable<void> {
    return this.http.delete<void>(`${this.endpoint}/structure/${id}`);
  }

  generateDemands(data: {
    class_id: string;
    academic_year_id: string;
    month: number;
    year: number;
  }): Observable<any> {
    return this.http.post(`${this.endpoint}/generate-demands`, data);
  }

  getStudentDues(studentId: string): Observable<any[]> {
    return this.http.get<any[]>(`${this.endpoint}/students/${studentId}/dues`);
  }

  getTransactionHistory(studentId: string): Observable<any[]> {
    return this.http.get<any[]>(
      `${this.endpoint}/students/${studentId}/history`,
    );
  }

  getStudentHistory(studentId: string): Observable<ApiResponse<any>> {
    return this.http.get<ApiResponse<any>>(`${this.endpoint}/students/${studentId}/history`);
  }

  collectFee(data: {
    student_fee_id: string;
    amount: number;
    payment_mode: string;
    reference_no?: string;
  }): Observable<any> {
    return this.http.post(`${this.endpoint}/collect`, data);
  }

  updateFeeHead(id: string, data: Partial<FeeHead>): Observable<FeeHead> {
    return this.http.put<FeeHead>(`${this.endpoint}/heads/${id}`, data);
  }

  updateFeeStructure(id: string, data: Partial<FeeStructure>): Observable<any> {
    return this.http.put<any>(`${this.endpoint}/structure/${id}`, data);
  }
}
