import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { ACAD_API_URL } from '../../academic.config';
import { Observable } from 'rxjs';
import { FeeHead, FeeStructure } from '../../models/academic.models';

@Injectable({
  providedIn: 'root',
})
export class FeeService {
  private http = inject(HttpClient);
  private baseUrl = inject(ACAD_API_URL);

  private get feeBase(): string {
    return `${this.baseUrl}/academics/fees`;
  }

  getFeeHeads(): Observable<FeeHead[]> {
    return this.http.get<FeeHead[]>(`${this.feeBase}/heads`);
  }

  createFeeHead(data: FeeHead): Observable<FeeHead> {
    return this.http.post<FeeHead>(`${this.feeBase}/heads`, data);
  }

  deleteFeeHead(id: string): Observable<void> {
    return this.http.delete<void>(`${this.feeBase}/heads/${id}`);
  }

  getFeeStructure(classId: string, academicYearId: string): Observable<any[]> {
    return this.http.get<any[]>(`${this.feeBase}/structure`, {
      params: { class_id: classId, academic_year_id: academicYearId },
    });
  }

  createFeeStructure(data: FeeStructure): Observable<any> {
    return this.http.post(`${this.feeBase}/structure`, data);
  }

  deleteFeeStructure(id: string): Observable<void> {
    return this.http.delete<void>(`${this.feeBase}/structure/${id}`);
  }

  generateDemands(data: {
    class_id: string;
    academic_year_id: string;
    month: number;
    year: number;
  }): Observable<any> {
    return this.http.post(`${this.feeBase}/generate-demands`, data);
  }

  getStudentDues(studentId: string): Observable<any[]> {
    return this.http.get<any[]>(`${this.feeBase}/students/${studentId}/dues`);
  }

  getTransactionHistory(studentId: string): Observable<any[]> {
    return this.http.get<any[]>(
      `${this.feeBase}/students/${studentId}/history`,
    );
  }

  collectFee(data: {
    student_fee_id: string;
    amount: number;
    payment_mode: import('../../models/academic.models').PaymentMode;
    remarks?: string;
  }): Observable<any> {
    return this.http.post(`${this.feeBase}/collect`, data);
  }

  updateFeeHead(id: string, data: Partial<FeeHead>): Observable<FeeHead> {
    return this.http.put<FeeHead>(`${this.feeBase}/heads/${id}`, data);
  }

  updateFeeStructure(id: string, data: Partial<FeeStructure>): Observable<any> {
    return this.http.put<any>(`${this.feeBase}/structure/${id}`, data);
  }
}
