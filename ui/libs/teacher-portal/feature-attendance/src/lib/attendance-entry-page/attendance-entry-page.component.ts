import { CommonModule, DatePipe } from '@angular/common';
import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { AttendanceFiltersComponent } from '../components/attendance-filters/attendance-filters.component';
import { StudentListComponent } from '../components/student-list/student-list.component';
import { AttendanceFooterComponent } from '../components/attendance-footer/attendance-footer.component';
import {
  AttendanceService,
  ClassManagementService,
  StudentService,
  AttendanceStatus,
  Section,
  ClassGrade
} from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';
import { HttpErrorResponse } from '@angular/common/http';

export interface AttendanceStudentModel {
  student_id: string;
  rollNo: string;
  name: string;
  initials: string;
  status: AttendanceStatus | string;
  remark: string;
}

@Component({
  selector: 'shikshakul-attendance-entry-page',
  standalone: true,
  imports: [
    CommonModule,
    AttendanceFiltersComponent,
    StudentListComponent,
    AttendanceFooterComponent,
  ],
  providers: [DatePipe],
  templateUrl: './attendance-entry-page.component.html',
  styleUrl: './attendance-entry-page.component.scss',
})
export class AttendanceEntryPageComponent implements OnInit {
  private attendanceService = inject(AttendanceService);
  private classService = inject(ClassManagementService);
  private studentService = inject(StudentService);
  private snackbar = inject(SnackbarService);
  private datePipe = inject(DatePipe);

  // Filter state
  classes = signal<ClassGrade[]>([]);
  availableSections = signal<Section[]>([]);

  selectedClassId = signal('');
  selectedSectionId = signal('');
  selectedDate = signal(this.datePipe.transform(new Date(), 'yyyy-MM-dd') as string);

  // Data state
  students = signal<AttendanceStudentModel[]>([]);
  loading = signal(false);
  processing = signal(false);

  // Computed properties
  total = computed(() => this.students().length);
  present = computed(() => this.students().filter(s => s.status === 'PRESENT' || s.status === 'P').length);
  absent = computed(() => this.students().filter(s => s.status === 'ABSENT' || s.status === 'A').length);
  late = computed(() => this.students().filter(s => s.status === 'LATE' || s.status === 'L').length);

  ngOnInit(): void {
    this.loadClasses();
  }

  loadClasses() {
    this.classService.getClasses().subscribe({
      next: (res: any) => {
        const data = res?.data || res;
        this.classes.set(Array.isArray(data) ? data.map((c: any) => ({
          ...c,
          id: c.id || c.class_id || c._id
        })) : []);
      },
      error: () => this.snackbar.error('Error', 'Failed to load classes.')
    });
  }

  onFilterChange(event: { classId: string, sectionId: string, date: string }) {
    const prevClassId = this.selectedClassId();
    this.selectedClassId.set(event.classId);
    this.selectedSectionId.set(event.sectionId);
    this.selectedDate.set(event.date);

    if (event.classId && event.classId !== prevClassId) {
      this.classService.getSectionsByClass(event.classId).subscribe({
        next: (res: any) => {
          const data = res?.data || res;
          this.availableSections.set(Array.isArray(data) ? data : []);
        },
        error: () => this.snackbar.error('Error', 'Failed to load sections.')
      });
    }

    if (event.classId && event.sectionId && event.date) {
      this.loadAttendanceData();
    } else {
      this.students.set([]);
    }
  }

  loadAttendanceData() {
    this.loading.set(true);
    // 1. Try to load existing register for the date
    this.attendanceService.getClassRegister(this.selectedSectionId(), this.selectedDate()).subscribe({
      next: (res: any) => {
        const data = res?.data || res;
        if (data && data.records && data.records.length > 0) {
          // Existing attendance found
          this.students.set(data.records.map((r: any) => ({
            student_id: r.student_id,
            rollNo: r.roll_no,
            name: r.student_name,
            initials: this.getInitials(r.student_name),
            status: r.status,
            remark: r.remarks || ''
          })));
          this.loading.set(false);
        } else {
          // No attendance found, load raw students list
          this.loadStudentList();
        }
      },
      error: () => {
        // Fallback or API not returning 404 cleanly, load raw students
        this.loadStudentList();
      }
    });
  }

  private loadStudentList() {
    this.studentService.getStudents({ class_id: this.selectedClassId() }).subscribe({
      next: (res: any) => {
        const data = res?.data || res;
        let stds = Array.isArray(data) ? data : (data?.students || []);

        // Map to attendance model, default all to Present
        this.students.set(stds.map((s: any) => ({
          student_id: s.id || s.student_id || s._id,
          rollNo: s.roll_number || s.admission_no || 'N/A',
          name: `${s.first_name || s.name || ''} ${s.last_name || ''}`.trim(),
          initials: this.getInitials(`${s.first_name || s.name || ''} ${s.last_name || ''}`),
          status: 'PRESENT',
          remark: ''
        })));
        this.loading.set(false);
      },
      error: () => {
        this.snackbar.error('Error', 'Failed to load student list.');
        this.loading.set(false);
      }
    });
  }

  private getInitials(name: string): string {
    if (!name) return 'S';
    const parts = name.trim().split(' ');
    if (parts.length >= 2) return `${parts[0][0]}${parts[parts.length - 1][0]}`.toUpperCase();
    return parts[0].substring(0, 2).toUpperCase();
  }

  updateStats() {
    // Computed properties handle this automatically in Angular 17. 
    // We just need to trigger a shallow copy of the signal if mutating internally.
    this.students.set([...this.students()]);
  }

  onSave() {
    if (!this.selectedClassId() || !this.selectedSectionId() || !this.selectedDate()) {
      this.snackbar.error('Validation Error', 'Please ensure class, section, and date are selected.');
      return;
    }

    if (this.students().length === 0) {
      this.snackbar.warning('No Data', 'No students to mark attendance for.');
      return;
    }

    this.processing.set(true);

    const payload = {
      class_id: this.selectedClassId(),
      section_id: this.selectedSectionId(),
      date: new Date(this.selectedDate()).toISOString(),
      students: this.students().map(s => ({
        student_id: s.student_id,
        status: s.status,
        remarks: s.remark
      }))
    };

    this.attendanceService.markAttendance(payload).subscribe({
      next: () => {
        this.snackbar.success('Success', 'Attendance saved successfully.');
        this.processing.set(false);
      },
      error: (err: HttpErrorResponse) => {
        this.snackbar.error('Error', err.error?.message || 'Failed to save attendance.');
        this.processing.set(false);
      }
    });
  }
}
