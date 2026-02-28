import { CommonModule } from '@angular/common';
import {
  Component,
  EventEmitter,
  inject,
  Input,
  OnInit,
  Output,
  signal,
  ChangeDetectorRef,
} from '@angular/core';
import {
  FormsModule,
  ReactiveFormsModule,
  FormBuilder,
  FormGroup,
  Validators,
} from '@angular/forms';
import {
  AcademicYear,
  ClassGrade,
  ClassManagementService,
  Section,
  Subject,
  TeacherService,
  TimetableService,
} from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';

@Component({
  selector: 'shikshakul-add-routine-drawer',
  standalone: true,
  imports: [CommonModule, FormsModule, ReactiveFormsModule],
  templateUrl: './add-routine-drawer.component.html',
  styleUrl: './add-routine-drawer.component.scss',
})
export class AddRoutineDrawerComponent implements OnInit {
  private fb = inject(FormBuilder);
  private classService = inject(ClassManagementService);
  private teacherService = inject(TeacherService);
  private timetableService = inject(TimetableService);
  private snackbar = inject(SnackbarService);
  private cdr = inject(ChangeDetectorRef);

  @Input() currentYear?: AcademicYear;
  @Output() close = new EventEmitter<void>();
  @Output() saved = new EventEmitter<void>();

  routineForm: FormGroup;
  classes = signal<ClassGrade[]>([]);
  sections = signal<Section[]>([]);
  subjects = signal<Subject[]>([]);
  teachers = signal<any[]>([]);
  submitting = false;

  days = ['MONDAY', 'TUESDAY', 'WEDNESDAY', 'THURSDAY', 'FRIDAY', 'SATURDAY'];

  constructor() {
    this.routineForm = this.fb.group({
      class_id: ['', Validators.required],
      section_id: ['', Validators.required],
      subject_id: ['', Validators.required],
      teacher_id: ['', Validators.required],
      day_of_week: ['', Validators.required],
      start_time: ['', Validators.required],
      end_time: ['', Validators.required],
      room_number: [''],
    });
  }

  ngOnInit(): void {
    this.loadClasses();
    this.loadTeachers();
    this.loadSubjects();
  }

  loadClasses() {
    this.classService.getClasses().subscribe((res: any) => {
      this.classes.set(Array.isArray(res) ? res : res?.data || []);
      this.cdr.detectChanges();
    });
  }

  loadTeachers() {
    this.teacherService.getTeachers().subscribe((res: any) => {
      this.teachers.set(Array.isArray(res) ? res : res?.data || []);
      this.cdr.detectChanges();
    });
  }

  loadSubjects() {
    this.classService.getAllSubjects().subscribe((res: any) => {
      this.subjects.set(Array.isArray(res) ? res : res?.data || []);
      this.cdr.detectChanges();
    });
  }

  onClassChange() {
    const classId = this.routineForm.get('class_id')?.value;
    this.sections.set([]);
    this.routineForm.patchValue({ section_id: '' });
    if (classId) {
      this.classService.getSectionsByClass(classId).subscribe((res: any) => {
        this.sections.set(Array.isArray(res) ? res : res?.data || []);
        this.cdr.detectChanges();
      });
    }
  }

  onSubmit() {
    if (this.routineForm.invalid || !this.currentYear) {
      // Mark all as touched to show validation errors
      this.routineForm.markAllAsTouched();
      this.snackbar.info(
        'Validation Error',
        'Please complete all required fields highlight in red.',
      );
      return;
    }

    this.submitting = true;
    this.cdr.detectChanges();
    const payload = {
      ...this.routineForm.value,
      academic_year_id: this.currentYear.id,
    };

    this.timetableService.createRoutine(payload).subscribe({
      next: () => {
        this.snackbar.success('Success', 'Routine created successfully.');
        this.submitting = false;
        this.saved.emit();
        this.close.emit();
        this.cdr.detectChanges();
      },
      error: (err) => {
        this.submitting = false;
        this.snackbar.error(
          'Error',
          err.error?.message || 'Failed to create routine.',
        );
        this.cdr.detectChanges();
      },
    });
  }
}
