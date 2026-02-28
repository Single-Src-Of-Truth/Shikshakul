import { CommonModule } from '@angular/common';
import { Component, inject, OnInit, signal, ChangeDetectorRef } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ExamFormComponent } from '../components/exam-form/exam-form.component';
import { RecentExamsWidgetComponent } from '../components/recent-exams-widget/recent-exams-widget.component';
import { ExamGuidelinesWidgetComponent } from '../components/exam-guidelines-widget/exam-guidelines-widget.component';
import { 
  ExamService, 
  ExamTerm, 
  ExamSchedule, 
  AcademicYearService, 
  AcademicYear,
  ClassManagementService,
  ClassGrade
} from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';

@Component({
  selector: 'shikshakul-exam-setup-page',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ExamFormComponent,
    RecentExamsWidgetComponent,
    ExamGuidelinesWidgetComponent,
  ],
  templateUrl: './exam-setup-page.component.html',
  styleUrl: './exam-setup-page.component.scss',
})
export class ExamSetupPageComponent implements OnInit {
  private examService = inject(ExamService);
  private acadYearService = inject(AcademicYearService);
  private classService = inject(ClassManagementService);
  private snackbar = inject(SnackbarService);
  private cdr = inject(ChangeDetectorRef);

  examTerms = signal<ExamTerm[]>([]);
  examSchedules = signal<ExamSchedule[]>([]);
  classes = signal<ClassGrade[]>([]);
  currentYear = signal<AcademicYear | undefined>(undefined);
  loading = signal(false);

  selectedTermId = signal('');
  selectedClassId = signal('');

  ngOnInit(): void {
    this.loadCurrentYear();
    this.loadClasses();
  }

  loadClasses() {
    this.classService.getClasses().subscribe((res: any) => {
      this.classes.set(Array.isArray(res) ? res : res?.data || []);
      this.cdr.detectChanges();
    });
  }

  loadCurrentYear() {
    this.acadYearService.getCurrentAcademicYear().subscribe((res: any) => {
      // Ensure we handle both {data: year} and year formats
      const year = res?.data || res;
      this.currentYear.set(year);
      
      const yearId = year?.id || year?.academic_year_id;
      if (yearId) {
        this.loadExamTerms(yearId);
      } else {
        console.error('No academic year ID found', res);
      }
      this.cdr.detectChanges();
    });
  }

  loadExamTerms(yearId: string) {
    this.loading.set(true);
    this.examService.getExamTerms(yearId).subscribe({
      next: (res: any) => {
        this.examTerms.set(Array.isArray(res) ? res : res?.data || []);
        this.loading.set(false);
        this.cdr.detectChanges();
      },
      error: () => {
        this.loading.set(false);
        this.snackbar.error('Error', 'Failed to load exam terms.');
        this.cdr.detectChanges();
      }
    });
  }

  onFilterChange() {
    this.loadExamSchedules();
  }

  loadExamSchedules() {
    const termId = this.selectedTermId();
    const classId = this.selectedClassId();

    if (!termId || !classId || termId === 'undefined' || classId === 'undefined') {
      this.examSchedules.set([]);
      this.cdr.detectChanges();
      return;
    }

    this.loading.set(true);
    this.examService.getExamSchedules({ term_id: termId, class_id: classId }).subscribe({
      next: (res: any) => {
        this.examSchedules.set(Array.isArray(res) ? res : res?.data || []);
        this.loading.set(false);
        this.cdr.detectChanges();
      },
      error: () => {
        this.loading.set(false);
        this.snackbar.error('Error', 'Failed to load exam schedules.');
        this.cdr.detectChanges();
      }
    });
  }

  onTermSaved() {
    const year = this.currentYear();
    const yearId = year?.id || year?.academic_year_id;
    if (yearId) this.loadExamTerms(yearId);
  }

  onScheduleSaved() {
    this.loadExamSchedules();
  }
}
