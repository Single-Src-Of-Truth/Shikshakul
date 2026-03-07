import { CommonModule, DatePipe } from '@angular/common';
import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { FormsModule } from '@angular/forms';
import {
  AttendanceService,
  ClassManagementService,
  AttendanceStatus,
  Section,
  ClassGrade,
} from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';

export interface AttendanceRecordModel {
  student_id: string;
  rollNo: string;
  name: string;
  initials: string;
  status: AttendanceStatus | string;
  remark: string;
}

@Component({
  selector: 'lib-attendance-page',
  standalone: true,
  imports: [CommonModule, FormsModule],
  providers: [DatePipe],
  templateUrl: './attendance-page.component.html',
  styleUrl: './attendance-page.component.scss',
})
export class AttendancePageComponent implements OnInit {
  private attendanceService = inject(AttendanceService);
  private classService = inject(ClassManagementService);
  private snackbar = inject(SnackbarService);
  private datePipe = inject(DatePipe);

  // Filter state
  classes = signal<ClassGrade[]>([]);
  availableSections = signal<Section[]>([]);

  selectedClassId = signal('');
  selectedSectionId = signal('');
  selectedDate = signal(
    this.datePipe.transform(new Date(), 'yyyy-MM-dd') as string,
  );

  // Data state
  records = signal<AttendanceRecordModel[]>([]);
  loading = signal(false);

  // Computed properties
  total = computed(() => this.records().length);
  present = computed(
    () =>
      this.records().filter((s) => s.status === 'PRESENT' || s.status === 'P')
        .length,
  );
  absent = computed(
    () =>
      this.records().filter((s) => s.status === 'ABSENT' || s.status === 'A')
        .length,
  );
  late = computed(
    () =>
      this.records().filter((s) => s.status === 'LATE' || s.status === 'L')
        .length,
  );

  ngOnInit(): void {
    this.loadClasses();
  }

  loadClasses() {
    this.classService.getClasses().subscribe({
      next: (res: any) => {
        const data = res?.data || res;
        this.classes.set(
          Array.isArray(data)
            ? data.map((c: any) => ({
                ...c,
                id: c.id || c.class_id || c._id,
              }))
            : [],
        );
      },
      error: () => this.snackbar.error('Error', 'Failed to load classes.'),
    });
  }

  onClassChange(classId: string) {
    this.selectedClassId.set(classId);
    this.selectedSectionId.set('');
    this.availableSections.set([]);
    this.records.set([]);

    if (classId) {
      this.classService.getSectionsByClass(classId).subscribe({
        next: (res: any) => {
          const data = res?.data || res;
          this.availableSections.set(Array.isArray(data) ? data : []);
        },
        error: () => this.snackbar.error('Error', 'Failed to load sections.'),
      });
    }
  }

  onFilterChange() {
    if (
      this.selectedClassId() &&
      this.selectedSectionId() &&
      this.selectedDate()
    ) {
      this.loadAttendanceData();
    } else {
      this.records.set([]);
    }
  }

  loadAttendanceData() {
    this.loading.set(true);
    this.attendanceService
      .getClassRegister(this.selectedSectionId(), this.selectedDate())
      .subscribe({
        next: (res: any) => {
          const data = res?.data || res;
          if (data && data.records && data.records.length > 0) {
            this.records.set(
              data.records.map((r: any) => ({
                student_id: r.student_id,
                rollNo: r.roll_no,
                name: r.student_name,
                initials: this.getInitials(r.student_name),
                status: r.status,
                remark: r.remarks || '',
              })),
            );
          } else {
            this.records.set([]);
          }
          this.loading.set(false);
        },
        error: () => {
          this.records.set([]);
          this.loading.set(false);
        },
      });
  }

  private getInitials(name: string): string {
    if (!name) return 'S';
    const parts = name.trim().split(' ');
    if (parts.length >= 2)
      return `${parts[0][0]}${parts[parts.length - 1][0]}`.toUpperCase();
    return parts[0].substring(0, 2).toUpperCase();
  }
}
