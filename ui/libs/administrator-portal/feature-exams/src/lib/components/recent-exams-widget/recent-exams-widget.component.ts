import { CommonModule } from '@angular/common';
import { Component, Input, Output, EventEmitter, inject, signal } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { ExamTerm, ExamService, AcademicYearService } from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';

@Component({
  selector: 'shikshakul-recent-exams-widget',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './recent-exams-widget.component.html',
  styleUrl: './recent-exams-widget.component.scss',
})
export class RecentExamsWidgetComponent {
  private fb = inject(FormBuilder);
  private examService = inject(ExamService);
  private snackbar = inject(SnackbarService);
  private acadYearService = inject(AcademicYearService);

  @Input() terms: ExamTerm[] = [];
  @Input() loading = false;
  @Output() termSaved = new EventEmitter<void>();

  showForm = signal(false);
  submitting = signal(false);
  termForm: FormGroup;

  constructor() {
    this.termForm = this.fb.group({
      name: ['', Validators.required],
      start_date: ['', Validators.required],
      end_date: ['', Validators.required],
    });
  }

  toggleForm() {
    this.showForm.set(!this.showForm());
  }

  onSubmit() {
    if (this.termForm.invalid) return;

    this.submitting.set(true);
    
    this.acadYearService.getCurrentAcademicYear().subscribe((year: any) => {
      const yearData = year?.data || year;
      const academic_year_id = yearData?.id || yearData?.academic_year_id;
      if (!academic_year_id) {
        this.snackbar.error('Error', 'Current academic year not found.');
        this.submitting.set(false);
        return;
      }

      const request = {
        ...this.termForm.value,
        academic_year_id
      };

      this.examService.createExamTerm(request).subscribe({
        next: () => {
          this.snackbar.success('Success', 'Exam term created successfully.');
          this.submitting.set(false);
          this.showForm.set(false);
          this.termForm.reset();
          this.termSaved.emit();
        },
        error: (err) => {
          this.submitting.set(false);
          this.snackbar.error('Error', err.error?.message || 'Failed to create term.');
        }
      });
    });
  }
}
