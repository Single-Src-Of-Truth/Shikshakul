import { CommonModule } from '@angular/common';
import {
  Component,
  EventEmitter,
  inject,
  Input,
  OnInit,
  OnChanges,
  SimpleChanges,
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
import { SelectComponent } from '@shikshakul/shared/ui/select';
import { TimePickerComponent } from '@shikshakul/shared/ui/time-picker';

@Component({
  selector: 'shikshakul-add-routine-drawer',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    SelectComponent,
    TimePickerComponent,
  ],
  templateUrl: './add-routine-drawer.component.html',
  styleUrl: './add-routine-drawer.component.scss',
})
export class AddRoutineDrawerComponent implements OnInit, OnChanges {
  private fb = inject(FormBuilder);
  private classService = inject(ClassManagementService);
  private teacherService = inject(TeacherService);
  private timetableService = inject(TimetableService);
  private snackbar = inject(SnackbarService);
  private cdr = inject(ChangeDetectorRef);

  @Input() currentYear?: AcademicYear;
  @Input() isOpen = false;
  @Input() editData: any = null;
  @Output() close = new EventEmitter<void>();
  @Output() saved = new EventEmitter<void>();

  routineForm: FormGroup;
  classes = signal<ClassGrade[]>([]);
  sections = signal<Section[]>([]);
  subjects = signal<Subject[]>([]);
  teachers = signal<any[]>([]);
  submitting = false;
  isFormSubmitted = false;

  classOptions = signal<{ label: string; value: string }[]>([]);
  sectionOptions = signal<{ label: string; value: string }[]>([]);
  subjectOptions = signal<{ label: string; value: string }[]>([]);
  teacherOptions = signal<{ label: string; value: string }[]>([]);

  dayOptions = signal<{ label: string; value: string }[]>([]);

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
    this.dayOptions.set(
      this.days.map((d) => ({
        label: d.charAt(0) + d.slice(1).toLowerCase(),
        value: d,
      })),
    );
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes['editData'] && this.editData) {
      this.populateForm(this.editData);
    }
  }

  populateForm(data: any) {
    this.routineForm.patchValue({
      class_id: data.class_id,
      subject_id: data.subject_id,
      teacher_id: data.teacher_id,
      day_of_week: data.day_of_week,
      start_time: data.start_time,
      end_time: data.end_time,
      room_number: data.room_number || '',
    });

    if (data.class_id) {
      this.classService
        .getSectionsByClass(data.class_id)
        .subscribe((res: any) => {
          const sectionsData = Array.isArray(res) ? res : res?.data || [];
          this.sections.set(sectionsData);
          this.sectionOptions.set(
            sectionsData.map((s: any) => ({
              label: s.name,
              value: s.id || s._id || s.section_id,
            })),
          );
          this.routineForm.patchValue({ section_id: data.section_id });
          this.cdr.detectChanges();
        });
    }
  }

  loadClasses() {
    this.classService.getClasses().subscribe((res: any) => {
      const data = Array.isArray(res) ? res : res?.data || [];
      this.classes.set(data);
      this.classOptions.set(
        data.map((c: any) => ({
          label: c.name,
          value: c.id || c._id || c.class_id,
        })),
      );
      this.cdr.detectChanges();
    });
  }

  loadTeachers() {
    this.teacherService.getTeachers().subscribe((res: any) => {
      const data = Array.isArray(res) ? res : res?.data || [];
      this.teachers.set(data);
      this.teacherOptions.set(
        data.map((t: any) => ({
          label: `${t.first_name} ${t.last_name}`,
          value: t.id || t._id || t.teacher_id,
        })),
      );
      this.cdr.detectChanges();
    });
  }

  loadSubjects() {
    this.classService.getAllSubjects().subscribe((res: any) => {
      const data = Array.isArray(res) ? res : res?.data || [];
      this.subjects.set(data);
      this.subjectOptions.set(
        data.map((s: any) => ({
          label: s.name,
          value: s.id || s._id || s.subject_id,
        })),
      );
      this.cdr.detectChanges();
    });
  }

  onClassChange() {
    const classId = this.routineForm.get('class_id')?.value;
    this.sections.set([]);
    this.sectionOptions.set([]);
    this.routineForm.patchValue({ section_id: '' });
    if (classId) {
      this.classService.getSectionsByClass(classId).subscribe((res: any) => {
        const data = Array.isArray(res) ? res : res?.data || [];
        this.sections.set(data);
        this.sectionOptions.set(
          data.map((s: any) => ({
            label: s.name,
            value: s.id || s._id || s.section_id,
          })),
        );
        this.cdr.detectChanges();
      });
    }
  }

  onSubmit() {
    this.isFormSubmitted = true;
    if (this.routineForm.invalid || !this.currentYear) {
      Object.keys(this.routineForm.controls).forEach((key) => {
        this.routineForm.get(key)?.markAsTouched();
      });
      this.cdr.detectChanges();
      return;
    }

    this.submitting = true;
    this.cdr.detectChanges();
    const payload = {
      ...this.routineForm.value,
      academic_year_id: this.currentYear.id,
    };

    const request = this.editData
      ? this.timetableService.updateRoutine(this.editData.id, payload)
      : this.timetableService.createRoutine(payload);

    request.subscribe({
      next: () => {
        this.snackbar.success(
          'Success',
          `Routine ${this.editData ? 'updated' : 'created'} successfully.`,
        );
        this.submitting = false;
        this.saved.emit();
        this.close.emit();
        this.cdr.detectChanges();
      },
      error: (err: any) => {
        this.submitting = false;
        this.snackbar.error(
          'Error',
          err.error?.message ||
            `Failed to ${this.editData ? 'update' : 'create'} routine.`,
        );
        this.cdr.detectChanges();
      },
    });
  }

  onReset() {
    if (confirm('Are you sure you want to clear the form?')) {
      this.routineForm.reset({
        class_id: '',
        section_id: '',
        subject_id: '',
        teacher_id: '',
        day_of_week: '',
        start_time: '',
        end_time: '',
        room_number: '',
      });
      this.isFormSubmitted = false;
      this.sections.set([]);
      this.sectionOptions.set([]);
    }
  }
}
