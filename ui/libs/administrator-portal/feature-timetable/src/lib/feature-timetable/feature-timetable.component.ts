import { CommonModule } from '@angular/common';
import {
  Component,
  inject,
  OnInit,
  ChangeDetectorRef,
  signal,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import {
  AcademicYear,
  AcademicYearService,
  ClassGrade,
  ClassManagementService,
  Section,
  TeacherService,
  TimetableService,
} from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';
import { SelectComponent } from '@shikshakul/shared/ui/select';
import { AddRoutineDrawerComponent } from '../components/add-routine-drawer/add-routine-drawer.component';

@Component({
  selector: 'shikshakul-feature-timetable',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    AddRoutineDrawerComponent,
    SelectComponent,
  ],
  templateUrl: './feature-timetable.component.html',
  styleUrl: './feature-timetable.component.scss',
})
export class FeatureTimetableComponent implements OnInit {
  private classService = inject(ClassManagementService);
  private teacherService = inject(TeacherService);
  private timetableService = inject(TimetableService);
  private acadYearService = inject(AcademicYearService);
  private snackbar = inject(SnackbarService);
  private cdr = inject(ChangeDetectorRef);

  classes = signal<ClassGrade[]>([]);
  sections = signal<Section[]>([]);
  teachers = signal<any[]>([]);
  routines = signal<any[]>([]);
  currentYear = signal<AcademicYear | undefined>(undefined);

  classOptions = signal<{ label: string; value: string }[]>([]);
  sectionOptions = signal<{ label: string; value: string }[]>([]);
  teacherOptions = signal<{ label: string; value: string }[]>([]);

  selectedClassId = '';
  selectedSectionId = '';
  selectedTeacherId = '';
  loading = signal(false);
  showAddDrawer = false;
  selectedRoutineForEdit: any = null;

  days = ['MONDAY', 'TUESDAY', 'WEDNESDAY', 'THURSDAY', 'FRIDAY', 'SATURDAY'];
  timeSlots = [
    '09:00 AM',
    '10:00 AM',
    '11:00 AM',
    '12:00 PM',
    '01:00 PM - BREAK',
    '02:00 PM',
    '03:00 PM',
    '04:00 PM',
    '05:00 PM',
  ];

  ngOnInit(): void {
    this.loadClasses();
    this.loadTeachers();
    this.loadCurrentYear();
  }

  loadCurrentYear() {
    this.acadYearService.getCurrentAcademicYear().subscribe((res: any) => {
      this.currentYear.set(res?.data || res);
      this.cdr.detectChanges();
    });
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

  onClassChange() {
    // Mutual exclusivity
    this.selectedTeacherId = '';

    this.selectedSectionId = '';
    this.sections.set([]);
    this.sectionOptions.set([]);

    if (this.selectedClassId && this.selectedClassId !== 'undefined') {
      this.classService
        .getSectionsByClass(this.selectedClassId)
        .subscribe((res: any) => {
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
    this.loadTimetable();
  }

  onSectionChange() {
    // Mutual exclusivity
    this.selectedTeacherId = '';
    this.loadTimetable();
  }

  onTeacherChange() {
    // Mutual exclusivity
    this.selectedClassId = '';
    this.selectedSectionId = '';
    this.sections.set([]);

    this.loadTimetable();
  }

  loadTimetable() {
    const filters: any = {};
    if (this.selectedSectionId) filters.section_id = this.selectedSectionId;
    if (this.selectedTeacherId) filters.teacher_id = this.selectedTeacherId;

    if (!filters.section_id && !filters.teacher_id) {
      this.routines.set([]);
      this.cdr.detectChanges();
      return;
    }

    this.loading.set(true);
    this.cdr.detectChanges();

    this.timetableService.getTimetable(filters).subscribe({
      next: (res: any) => {
        this.routines.set(Array.isArray(res) ? res : res?.data || []);
        this.loading.set(false);
        this.cdr.detectChanges();
      },
      error: () => {
        this.loading.set(false);
        this.snackbar.error('Error', 'Failed to load timetable.');
        this.cdr.detectChanges();
      },
    });
  }

  getRoutinesForDay(day: string) {
    return this.routines().filter((r) => r.day_of_week === day);
  }

  getRoutineForSlot(day: string, timeLabel: string) {
    if (timeLabel.includes('BREAK')) return null;

    // timeLabel is e.g., "09:00 AM"
    // Extract hour for matching. Internal format is "HH:mm"
    const [time, modifier] = timeLabel.split(' ');
    const [rawHours, minutes] = time.split(':');
    let hours = rawHours;
    if (hours === '12') {
      hours = modifier === 'AM' ? '00' : '12';
    } else {
      hours = modifier === 'PM' ? (parseInt(hours, 10) + 12).toString() : hours;
    }
    const matchHour = hours.padStart(2, '0');
    return this.routines().find(
      (r) => r.day_of_week === day && r.start_time.startsWith(matchHour),
    );
  }

  openAddRoutine() {
    this.selectedRoutineForEdit = null;
    this.showAddDrawer = true;
    this.cdr.detectChanges();
  }

  editRoutine(routine: any) {
    this.selectedRoutineForEdit = routine;
    this.showAddDrawer = true;
    this.cdr.detectChanges();
  }

  deleteRoutine(id: string) {
    if (confirm('Are you sure you want to delete this routine slot?')) {
      this.timetableService.deleteRoutine(id).subscribe({
        next: () => {
          this.snackbar.success('Success', 'Routine deleted successfully.');
          this.loadTimetable();
        },
        error: () => {
          this.snackbar.error('Error', 'Failed to delete routine.');
        },
      });
    }
  }
}
