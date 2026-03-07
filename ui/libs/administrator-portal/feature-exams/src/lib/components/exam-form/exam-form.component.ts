import { CommonModule } from '@angular/common';
import { Component, inject, OnInit, signal, ChangeDetectorRef, Output, EventEmitter, Input } from '@angular/core';
import {
  FormArray,
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { 
  ExamService, 
  ClassManagementService, 
  TeacherService, 
  ExamTerm, 
  ClassGrade, 
  Subject,
  CreateExamScheduleRequest 
} from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';

@Component({
  selector: 'shikshakul-exam-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './exam-form.component.html',
  styleUrl: './exam-form.component.scss',
})
export class ExamFormComponent implements OnInit {
  private fb = inject(FormBuilder);
  private examService = inject(ExamService);
  private classService = inject(ClassManagementService);
  private snackbar = inject(SnackbarService);
  private cdr = inject(ChangeDetectorRef);

  private _terms: ExamTerm[] = [];
  @Input() set terms(value: any[]) {
    this._terms = (value || []).map(term => ({
      ...term,
      id: term.id || term.exam_term_id || term._id,
      name: term.name || term.term_name
    }));
  }
  get terms(): ExamTerm[] {
    return this._terms;
  }

  @Output() saved = new EventEmitter<void>();

  examForm: FormGroup;
  classes = signal<ClassGrade[]>([]);
  subjects = signal<Subject[]>([]);
  submitting = signal(false);

  constructor() {
    this.examForm = this.fb.group({
      exam_term_id: ['', Validators.required],
      class_id: ['', Validators.required],
      schedules: this.fb.array([]),
    });
  }

  ngOnInit(): void {
    this.loadInitialData();
  }

  loadInitialData() {
    this.classService.getClasses().subscribe((res: any) => {
      const data = res?.data || res;
      const parsedClasses = Array.isArray(data) ? data.map((c: any) => ({
        ...c,
        id: c.id || c.class_id || c._id,
        name: c.name || c.class_name
      })) : [];
      this.classes.set(parsedClasses);
      this.cdr.detectChanges();
    });

    this.classService.getAllSubjects().subscribe((res: any) => {
      const data = res?.data || res;
      const parsedSubjects = Array.isArray(data) ? data.map((s: any) => ({
        ...s,
        subject_id: s.subject_id || s.id || s._id,
        name: s.name || s.subject_name
      })) : [];
      this.subjects.set(parsedSubjects);
      this.cdr.detectChanges();
    });
  }

  // Helper to create a new row
  createScheduleRow(): FormGroup {
    return this.fb.group({
      subject_id: ['', Validators.required],
      exam_date: ['', Validators.required],
      start_time: ['', Validators.required],
      end_time: ['', Validators.required],
      max_marks: [100, [Validators.required, Validators.min(1)]],
      pass_marks: [33, [Validators.required, Validators.min(0)]],
      room_number: [''],
    });
  }

  // Getter for easy access in HTML
  get schedulesArray(): FormArray {
    return this.examForm.get('schedules') as FormArray;
  }

  addSubject() {
    this.schedulesArray.push(this.createScheduleRow());
    this.cdr.detectChanges();
  }

  removeSubject(index: number) {
    this.schedulesArray.removeAt(index);
    this.cdr.detectChanges();
  }

  onSubmit() {
    if (this.examForm.invalid) {
      this.examForm.markAllAsTouched();
      this.snackbar.info('Validation', 'Please fill all required fields highlight in red.');
      return;
    }

    if (this.schedulesArray.length === 0) {
      this.snackbar.error('Error', 'Please add at least one subject to schedule.');
      return;
    }

    this.submitting.set(true);
    const { exam_term_id, class_id, schedules } = this.examForm.value;

    const requests = schedules.map((s: any) => {
      // Calculate duration in minutes from start and end time (Format HH:mm)
      let duration_min = 0;
      if (s.start_time && s.end_time) {
        const [startH, startM] = s.start_time.split(':').map(Number);
        const [endH, endM] = s.end_time.split(':').map(Number);
        duration_min = (endH * 60 + endM) - (startH * 60 + startM);
        if (duration_min < 0) duration_min += 24 * 60; // handle crossing midnight if ever needed
      }

      let formattedDate = s.exam_date;
      if (s.exam_date && !s.exam_date.includes('T')) {
        // Use native Date object to reliably format to ISO8601 UTC
        formattedDate = new Date(s.exam_date).toISOString();
      }

      return {
        ...s,
        exam_date: formattedDate,
        duration_min,
        exam_term_id,
        class_id
      };
    });

    // For simplicity, we'll create them sequentially or use forkJoin if available
    // Here we'll just demonstrate the first one or a loop for MVP
    // Better yet, if the API supports bulk creation, that would be ideal.
    // The Postman snippet shows "Create Exam Schedule" (singular looking)
    
    let completedCount = 0;
    requests.forEach((req: CreateExamScheduleRequest) => {
      this.examService.createExamSchedule(req).subscribe({
        next: () => {
          completedCount++;
          if (completedCount === requests.length) {
            this.handleSuccess();
          }
        },
        error: (err) => {
          this.submitting.set(false);
          this.snackbar.error('Error', err.error?.message || 'Failed to schedule exam.');
          this.cdr.detectChanges();
        }
      });
    });
  }

  private handleSuccess() {
    this.submitting.set(false);
    this.snackbar.success('Success', 'Examination schedule published successfully.');
    this.examForm.reset({
      exam_term_id: '',
      class_id: '',
      schedules: []
    });
    this.schedulesArray.clear();
    this.saved.emit();
    this.cdr.detectChanges();
  }
}
