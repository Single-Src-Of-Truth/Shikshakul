import { Component, inject, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  CertificateService,
  ClassManagementService,
  StudentService,
  ClassGrade,
  IDCardResponse,
} from '@shikshakul/data-access/academic';
import { IdCardViewComponent } from '../components/id-card-view/id-card-view.component';
import { TcFormComponent } from '../components/tc-form/tc-form.component';
import { BonafideFormComponent } from '../components/bonafide-form/bonafide-form.component';
import { CertificatePreviewComponent } from '../components/certificate-preview/certificate-preview.component';
import { FormsModule } from '@angular/forms';
import { TCResponse, BonafideResponse } from '@shikshakul/data-access/academic';

@Component({
  selector: 'shikshakul-certificate-management-page',
  standalone: true,
  imports: [
    CommonModule,
    IdCardViewComponent,
    TcFormComponent,
    BonafideFormComponent,
    CertificatePreviewComponent,
    FormsModule,
  ],
  templateUrl: './certificate-management-page.component.html',
  styleUrl: './certificate-management-page.component.scss',
})
export class CertificateManagementPageComponent implements OnInit {
  private classService = inject(ClassManagementService);
  private studentService = inject(StudentService);
  private certificateService = inject(CertificateService);

  classes = signal<ClassGrade[]>([]);
  students = signal<any[]>([]);
  idCards = signal<IDCardResponse[]>([]);

  selectedClassId = '';
  selectedClassName = '';
  selectedStudent: any = null;

  viewMode:
    | 'MANAGEMENT'
    | 'ID_CARD_PREVIEW'
    | 'TC_FORM'
    | 'BONAFIDE_FORM'
    | 'CERTIFICATE_PREVIEW' = 'MANAGEMENT';
  generatedCertificate: TCResponse | BonafideResponse | null = null;
  certificateType: 'TC' | 'BONAFIDE' = 'TC';

  ngOnInit() {
    this.loadClasses();
  }

  loadClasses() {
    this.classService.getClasses().subscribe((res: any) => {
      const data = res?.data || res;
      this.classes.set(
        Array.isArray(data)
          ? data.map((c: any) => ({
              ...c,
              id: c.id || c.class_id || c._id,
            }))
          : [],
      );
    });
  }

  onClassChange() {
    if (this.selectedClassId && this.selectedClassId !== 'undefined') {
      const selectedClass = this.classes().find(
        (c) => c.id === this.selectedClassId,
      );
      this.selectedClassName = selectedClass ? selectedClass.name : '';

      this.studentService
        .getStudents({ class_id: this.selectedClassId })
        .subscribe((res: any) => {
          const data = res?.data || res;
          this.students.set(Array.isArray(data) ? data : []);
        });
    } else {
      this.students.set([]);
      this.selectedClassName = '';
    }
  }

  generateIDCards() {
    if (this.selectedClassId && this.students().length > 0) {
      const mappedCards: IDCardResponse[] = this.students().map((st) => {
        const profile = st.profile_data || {};
        const className = st.class?.name || '';
        const sectionName = st.section?.name || '';

        return {
          student_id: st.id,
          admission_no: st.admission_no,
          full_name: `${st.first_name} ${st.last_name}`.trim(),
          class_details: sectionName
            ? `${className} - ${sectionName}`
            : className,
          dob: profile.dob || '',
          blood_group: profile.blood_group || '',
          father_name: profile.father_name || '',
          contact_number: st.mobile || '',
          address: profile.address || '',
          photo_url: profile.photo_url || '',
        };
      });
      this.idCards.set(mappedCards);
      this.viewMode = 'ID_CARD_PREVIEW';
    } else if (this.selectedClassId) {
      alert('No students found in this class to generate ID cards.');
    }
  }

  openTCForm(student: any) {
    this.selectedStudent = student;
    this.viewMode = 'TC_FORM';
  }

  openBonafideForm(student: any) {
    this.selectedStudent = student;
    this.viewMode = 'BONAFIDE_FORM';
  }

  onSuccess(data: TCResponse | BonafideResponse, type: 'TC' | 'BONAFIDE') {
    this.generatedCertificate = data;
    this.certificateType = type;
    this.viewMode = 'CERTIFICATE_PREVIEW';
    this.selectedStudent = null;
  }
}
