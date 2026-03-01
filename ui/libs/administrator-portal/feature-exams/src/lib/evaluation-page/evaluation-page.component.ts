import { CommonModule } from '@angular/common';
import {
  Component,
  inject,
  OnInit,
  signal,
  ChangeDetectorRef,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import {
  ExamService,
  ExamTerm,
  ExamSchedule,
  AcademicYearService,
  ClassManagementService,
  ClassGrade,
  Subject,
} from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';

@Component({
  selector: 'shikshakul-evaluation-list',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './evaluation-page.component.html',
  styleUrl: './evaluation-page.component.scss',
})
export class EvaluationListComponent implements OnInit {
  private examService = inject(ExamService);
  private acadYearService = inject(AcademicYearService);
  private classService = inject(ClassManagementService);
  private router = inject(Router);
  private cdr = inject(ChangeDetectorRef);

  examTerms = signal<ExamTerm[]>([]);
  classes = signal<ClassGrade[]>([]);
  schedules = signal<ExamSchedule[]>([]);
  subjects = signal<Subject[]>([]);
  loading = signal(false);

  selectedTermId = signal('');
  selectedClassId = signal('');

  ngOnInit(): void {
    this.loadInitialData();
  }

  loadInitialData() {
    this.classService.getClasses().subscribe((res) => {
      const parsedClasses = (res.data || []).map((c: any) => ({
        ...c,
        id: c.id || c.class_id || c._id,
        name: c.name || c.class_name,
      }));
      this.classes.set(parsedClasses);
      this.cdr.detectChanges();
    });

    this.classService.getAllSubjects().subscribe((res) => {
      this.subjects.set(res.data || []);
      this.cdr.detectChanges();
    });

    this.acadYearService.getCurrentAcademicYear().subscribe({
      next: (res: any) => {
        const year = res?.data || res;
        const yearId = year?.id || year?.academic_year_id;
        if (yearId) {
          this.loadExamTerms(yearId);
        } else {
          console.warn('Current academic year ID was undefined', res);
        }
      },
      error: (err) => {
        console.warn(
          'Could not fetch current academic year, falling back to all years.',
          err,
        );
        this.acadYearService.getAcademicYears().subscribe({
          next: (yearsRes: any) => {
            const years = yearsRes?.data || yearsRes || [];
            if (years.length > 0) {
              const yearId = years[0]?.id || years[0]?.academic_year_id;
              if (yearId) {
                this.loadExamTerms(yearId);
              }
            } else {
              console.warn('No academic years found in the system.');
            }
          },
        });
      },
    });
  }

  loadExamTerms(yearId: string) {
    if (!yearId) return;
    this.examService.getExamTerms(yearId).subscribe({
      next: (termRes) => {
        const parsedTerms = (termRes.data || []).map((t: any) => ({
          ...t,
          id: t.id || t.exam_term_id || t._id,
          name: t.name || t.term_name,
        }));
        this.examTerms.set(parsedTerms);
        this.cdr.detectChanges();
      },
      error: (err) => {
        console.error('Error fetching exam terms:', err);
      },
    });
  }

  onFilterChange() {
    const termId = this.selectedTermId();
    const classId = this.selectedClassId();
    console.log('Evaluations filter changed:', { termId, classId });

    if (
      !termId ||
      !classId ||
      termId === 'undefined' ||
      classId === 'undefined'
    ) {
      console.log('Skipping API call because term or class is missing');
      this.schedules.set([]);
      this.cdr.detectChanges();
      return;
    }

    console.log('Fetching exam schedules for', { termId, classId });
    this.loading.set(true);
    this.examService
      .getExamSchedules({ term_id: termId, class_id: classId })
      .subscribe({
        next: (res) => {
          const rawSchedules = res.data || [];

          const currentTerm = this.examTerms().find((t) => t.id === termId);

          const mapped = rawSchedules.map((s) => {
            return {
              ...s,
              exam_term: s.exam_term || currentTerm,
            };
          });

          this.schedules.set(mapped);
          this.loading.set(false);
          this.cdr.detectChanges();
        },
        error: () => {
          this.loading.set(false);
          this.cdr.detectChanges();
        },
      });
  }

  enterMarks(scheduleId: string) {
    this.router.navigate(['/exams/evaluation', scheduleId]);
  }

  getStatusClass(status: string | undefined): string {
    return (status || 'NOT_STARTED').toLowerCase();
  }

  getButtonClass(status: string | undefined): string {
    if (status === 'FINALIZED') return 'btn-outline';
    if (status === 'DRAFT') return 'btn-outline-blue';
    return 'btn-solid-blue';
  }

  getButtonText(status: string | undefined): string {
    if (status === 'FINALIZED') return 'View';
    if (status === 'DRAFT') return 'Resume';
    return 'Start';
  }

  // getClassInitials(classId: string): string {
  //   if (!classId) return 'N/A';
  //   const c = this.classes().find((cls) => cls.id === classId);
  //   if (!c) return 'N/A';
  //   // Example: "Class 7A" -> "7A", "Nursery" -> "Nu"
  //   const words = c.name.split(' ');
  //   if (words.length > 1) {
  //     const code = words[1];
  //     return code.length <= 2 ? code : code.substring(0, 2).toUpperCase();
  //   }
  //   return c.name.substring(0, 2).toUpperCase();
  // }

  // Stats Helpers
  getPendingCount() {
    return this.schedules().filter(
      (s) => !s.evaluation_status || s.evaluation_status === 'NOT_STARTED',
    ).length;
  }

  getDraftCount() {
    return this.schedules().filter((s) => s.evaluation_status === 'DRAFT')
      .length;
  }

  getCompletedCount() {
    return this.schedules().filter((s) => s.evaluation_status === 'FINALIZED')
      .length;
  }
}
