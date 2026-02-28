import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { ClassGrade, Section, Subject } from '../../models/academic.models';
import { ACAD_API_URL } from '../../academic.config';

@Injectable({
  providedIn: 'root',
})
export class ClassManagementService {
  private http = inject(HttpClient);
  private baseUrl = inject(ACAD_API_URL);

  private get classEndpoint(): string {
    return `${this.baseUrl}/academics/setup/classes`;
  }

  private get subjectEndpoint(): string {
    return `${this.baseUrl}/academics/setup/subjects`;
  }

  private getSectionEndpoint(classId: string): string {
    return `${this.classEndpoint}/${classId}/sections`;
  }

  private getSubjectAssignmentEndpoint(classId: string): string {
    return `${this.classEndpoint}/${classId}/subjects`;
  }

  private getClassTeacherEndpoint(sectionId: string): string {
    return `${this.baseUrl}/academics/sections/${sectionId}/class-teacher`;
  }

  private getSubjectTeacherEndpoint(sectionId: string): string {
    return `${this.baseUrl}/academics/sections/${sectionId}/subject-teachers`;
  }

  private getAllocationsEndpoint(sectionId: string): string {
    return `${this.baseUrl}/academics/sections/${sectionId}/allocations`;
  }

  getClasses(): Observable<ClassGrade[]> {
    return this.http.get<ClassGrade[]>(this.classEndpoint);
  }

  createClass(data: {
    name: string;
    sort_order: number;
  }): Observable<ClassGrade> {
    return this.http.post<ClassGrade>(this.classEndpoint, data);
  }

  getSectionsByClass(classId: string): Observable<Section[]> {
    if (!classId || classId === 'undefined') {
      return new Observable((subscriber) => {
        subscriber.next([]);
        subscriber.complete();
      });
    }
    return this.http.get<Section[]>(this.getSectionEndpoint(classId));
  }

  createSection(
    classId: string,
    data: { name: string; capacity: number },
  ): Observable<Section> {
    return this.http.post<Section>(this.getSectionEndpoint(classId), data);
  }

  getAllSubjects(): Observable<Subject[]> {
    return this.http.get<Subject[]>(this.subjectEndpoint);
  }

  createSubject(data: {
    name: string;
    code: string;
    type: string;
  }): Observable<Subject> {
    return this.http.post<Subject>(this.subjectEndpoint, data);
  }

  updateSubject(id: string, data: Partial<Subject>): Observable<Subject> {
    return this.http.put<Subject>(`${this.subjectEndpoint}/${id}`, data);
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
    return this.http.post(this.getSubjectTeacherEndpoint(sectionId), payload);
  }

  getAllocations(sectionId: string, academicYearId: string): Observable<any> {
    return this.http.get(this.getAllocationsEndpoint(sectionId), {
      params: { academic_year_id: academicYearId },
    });
  }

  updateClass(classId: string, data: { name: string; sort_order: number }): Observable<ClassGrade> {
    return this.http.put<ClassGrade>(`${this.classEndpoint}/${classId}`, data);
  }

  deleteClass(classId: string): Observable<void> {
    return this.http.delete<void>(`${this.classEndpoint}/${classId}`);
  }

  updateSection(sectionId: string, data: { name: string; capacity: number }): Observable<Section> {
    return this.http.put<Section>(`${this.baseUrl}/academics/setup/sections/${sectionId}`, data);
  }

  deleteSection(sectionId: string): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/academics/setup/sections/${sectionId}`);
  }
}
