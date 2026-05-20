import { CommonModule } from '@angular/common';
import {
  Component,
  inject,
  OnInit,
  signal,
  ChangeDetectorRef,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import {
  ExamService,
  ExamTerm,
  ClassManagementService,
  ClassGrade,
  AcademicYearService,
} from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';

@Component({
  selector: 'lib-result-management',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './result-management.component.html',
  styleUrl: './result-management.component.scss',
})
export class ResultManagementComponent implements OnInit {
  private examService = inject(ExamService);
  private classService = inject(ClassManagementService);
  private acadYearService = inject(AcademicYearService);
  private snackbar = inject(SnackbarService);
  private cdr = inject(ChangeDetectorRef);

  examTerms = signal<ExamTerm[]>([]);
  classes = signal<ClassGrade[]>([]);

  selectedTermId = signal('');
  selectedClassId = signal('');

  loading = signal(false);
  processing = signal(false);

  ngOnInit(): void {
    this.loadInitialData();
  }

  loadInitialData() {
    this.loading.set(true);

    // Get current year first to filter terms
    this.acadYearService.getCurrentAcademicYear().subscribe((yearRes) => {
      const year = yearRes?.data;
      const yearId = year?.id || year?.academic_year_id;

      if (yearId) {
        this.examService.getExamTerms(yearId).subscribe((termRes) => {
          this.examTerms.set(termRes?.data || []);
          this.checkLoadingState();
        });
      }

      this.classService.getClasses().subscribe((res) => {
        const data = res?.data || res;
        this.classes.set(
          Array.isArray(data)
            ? data.map(
                (c: {
                  id?: string;
                  class_id?: string;
                  _id?: string;
                  name?: string;
                  class_name?: string;
                  sort_order?: number;
                }) => ({
                  id: c.id || c.class_id || c._id || '',
                  name: c.name || c.class_name || '',
                  sort_order: c.sort_order || 0,
                }),
              )
            : [],
        );
        this.checkLoadingState();
      });
    });
  }

  private checkLoadingState() {
    if (this.examTerms().length >= 0 && this.classes().length >= 0) {
      this.loading.set(false);
      this.cdr.detectChanges();
    }
  }

  generateResults() {
    if (!this.selectedTermId() || !this.selectedClassId()) {
      this.snackbar.warning(
        'Selection Required',
        'Please select both Exam Term and Class.',
      );
      return;
    }

    this.processing.set(true);
    this.examService
      .generateResults({
        exam_term_id: this.selectedTermId(),
        class_id: this.selectedClassId(),
      })
      .subscribe({
        next: () => {
          this.snackbar.success('Success', 'Results generated successfully.');
          this.processing.set(false);
          this.cdr.detectChanges();
        },
        error: (err) => {
          this.snackbar.error(
            'Error',
            err.error?.message || 'Failed to generate results.',
          );
          this.processing.set(false);
          this.cdr.detectChanges();
        },
      });
  }

  publishResults(publish: boolean) {
    if (!this.selectedTermId() || !this.selectedClassId()) {
      this.snackbar.warning(
        'Selection Required',
        'Please select both Exam Term and Class.',
      );
      return;
    }

    this.processing.set(true);
    this.examService
      .publishResults({
        exam_term_id: this.selectedTermId(),
        class_id: this.selectedClassId(),
        publish,
      })
      .subscribe({
        next: () => {
          this.snackbar.success(
            'Success',
            publish ? 'Results published.' : 'Results unpublished.',
          );
          this.processing.set(false);
          this.cdr.detectChanges();
        },
        error: (err) => {
          this.snackbar.error(
            'Error',
            err.error?.message || 'Failed to update publication status.',
          );
          this.processing.set(false);
          this.cdr.detectChanges();
        },
      });
  }
}
