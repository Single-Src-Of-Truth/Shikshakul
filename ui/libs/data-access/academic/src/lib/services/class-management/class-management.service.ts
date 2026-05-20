import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import {
  ClassGrade,
  Section,
  Subject,
  ApiResponse,
} from '../../models/academic.models';
import { ACAD_API_URL } from '../../academic.config';

@Injectable({
  providedIn: 'root',
})
export class ClassManagementService {
  private http = inject(HttpClient);
  private baseUrl = inject(ACAD_API_URL);

  private get classEndpoint(): string {
    return `${this.baseUrl}/setup/classes`;
  }

  private get subjectEndpoint(): string {
    return `${this.baseUrl}/setup/subjects`;
  }

  private getSectionEndpoint(classId: string): string {
    return `${this.classEndpoint}/${classId}/sections`;
  }

  private getSubjectAssignmentEndpoint(classId: string): string {
    return `${this.classEndpoint}/${classId}/subjects`;
  }

  private getClassTeacherEndpoint(sectionId: string): string {
    return `${this.baseUrl}/sections/${sectionId}/class-teacher`;
  }

  private get(sectionId: string): string {
    return `${this.baseUrl}/sections/${sectionId}/subject-teachers`;
  }

  private getAllocationsEndpoint(sectionId: string): string {
    return `${this.baseUrl}/sections/${sectionId}/allocations`;
  }

  getClasses(): Observable<ApiResponse<ClassGrade[]>> {
    return this.http.get<ApiResponse<ClassGrade[]>>(this.classEndpoint);
  }

  createClass(data: {
    name: string;
    sort_order: number;
  }): Observable<ApiResponse<ClassGrade>> {
    return this.http.post<ApiResponse<ClassGrade>>(this.classEndpoint, data);
  }

  getSectionsByClass(classId: string): Observable<ApiResponse<Section[]>> {
    if (!classId || classId === 'undefined') {
      return new Observable((subscriber) => {
        subscriber.next({ data: [], message: '', success: true });
        subscriber.complete();
      });
    }
    return this.http.get<ApiResponse<Section[]>>(
      this.getSectionEndpoint(classId),
    );
  }

  createSection(
    classId: string,
    data: { name: string; capacity: number; state?: string },
  ): Observable<ApiResponse<Section>> {
    return this.http.post<ApiResponse<Section>>(
      this.getSectionEndpoint(classId),
      data,
    );
  }

  getAllSubjects(): Observable<ApiResponse<Subject[]>> {
    return this.http.get<ApiResponse<Subject[]>>(this.subjectEndpoint);
  }

  createSubject(data: {
    name: string;
    code: string;
    type: string;
  }): Observable<ApiResponse<Subject>> {
    return this.http.post<ApiResponse<Subject>>(this.subjectEndpoint, data);
  }

  updateSubject(
    id: string,
    data: Partial<Subject>,
  ): Observable<ApiResponse<Subject>> {
    return this.http.put<ApiResponse<Subject>>(
      `${this.subjectEndpoint}/${id}`,
      data,
    );
  }

  assignSubjectToClass(
    classId: string,
    data: {
      subjects: {
        subject_id: string;
        is_optional: boolean;
        weekly_lectures: number;
      }[];
    },
  ): Observable<void> {
    return this.http.post<void>(
      this.getSubjectAssignmentEndpoint(classId),
      data,
    );
  }

  unassignSubjectFromClass(
    classId: string,
    subjectId: string,
  ): Observable<void> {
    return this.http.delete<void>(
      `${this.getSubjectAssignmentEndpoint(classId)}/${subjectId}`,
    );
  }

  assignClassTeacher(sectionId: string, teacherId: string): Observable<any> {
    return this.http.post(this.getClassTeacherEndpoint(sectionId), {
      teacher_id: teacherId,
    });
  }

  assignSubjectTeacher(
    sectionId: string,
    payload: {
      academic_year_id: string;
      subject_id: string;
      teacher_id: string;
    },
  ): Observable<any> {
    return this.http.post(`${this.baseUrl}/sections/${sectionId}/subject-teachers`, payload);
  }

  getAllocations(sectionId: string, academicYearId: string): Observable<any> {
    return this.http.get(this.getAllocationsEndpoint(sectionId), {
      params: { academic_year_id: academicYearId },
    });
  }

  updateClass(
    classId: string,
    data: { name: string; sort_order: number },
  ): Observable<ApiResponse<ClassGrade>> {
    return this.http.put<ApiResponse<ClassGrade>>(
      `${this.classEndpoint}/${classId}`,
      data,
    );
  }

  deleteClass(classId: string): Observable<void> {
    return this.http.delete<void>(`${this.classEndpoint}/${classId}`);
  }

  updateSection(
    sectionId: string,
    data: { name: string; capacity: number },
  ): Observable<ApiResponse<Section>> {
    return this.http.put<ApiResponse<Section>>(
      `${this.baseUrl}/setup/sections/${sectionId}`,
      data,
    );
  }

  deleteSection(sectionId: string): Observable<void> {
    return this.http.delete<void>(
      `${this.baseUrl}/setup/sections/${sectionId}`,
    );
  }
}
