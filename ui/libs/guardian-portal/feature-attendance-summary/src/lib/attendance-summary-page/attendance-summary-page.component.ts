import { CommonModule } from '@angular/common';
import { Component, inject, OnInit, signal } from '@angular/core';
import { StudentHeaderComponent } from '../components/student-header/student-header.component';
import { AttendanceStatsComponent } from '../components/attendance-stats/attendance-stats.component';
import { CalendarHeatmapComponent } from '../components/calendar-heatmap/calendar-heatmap.component';
import { RecentAbsencesComponent } from '../components/recent-absences/recent-absences.component';
import { AttendanceTrendsComponent } from '../components/attendance-trends/attendance-trends.component';
import { AttendanceService, StudentService, StudentAttendanceHistory } from '@shikshakul/data-access/academic';

@Component({
  selector: 'shikshakul-attendance-summary-page',
  standalone: true,
  imports: [
    CommonModule,
    StudentHeaderComponent,
    AttendanceStatsComponent,
    AttendanceTrendsComponent,
    CalendarHeatmapComponent,
    RecentAbsencesComponent,
  ],
  templateUrl: './attendance-summary-page.component.html',
  styleUrl: './attendance-summary-page.component.scss',
})
export class AttendanceSummaryPageComponent implements OnInit {
  private studentService = inject(StudentService);
  private attendanceService = inject(AttendanceService);

  studentId = signal<string>('');
  student = signal<any>(null);
  loading = signal(true);

  // Stats signals
  totalDays = signal(0);
  presentDays = signal(0);
  absentDays = signal(0);
  lateDays = signal(0);

  // Raw records
  records = signal<StudentAttendanceHistory[]>([]);

  ngOnInit(): void {
    // Determine context. In guardian portal, usually we'd fetch the guardian's linked students.
    // Assuming backend returns related students for the current token when no filters are sent.
    this.studentService.getStudents().subscribe({
      next: (res: any) => {
        const data = res?.data || res;
        const stds = Array.isArray(data) ? data : (data?.students || []);
        if (stds.length > 0) {
          // Use the first student linked to the guardian
          this.student.set(stds[0]);
          const id = stds[0].id || stds[0].student_id || stds[0]._id;
          this.studentId.set(id);
          this.fetchAttendanceHistory(id);
        } else {
          this.loading.set(false);
          console.warn('No student profiles linked to this user context.');
        }
      },
      error: () => {
        this.loading.set(false);
      }
    });
  }

  fetchAttendanceHistory(studentId: string) {
    this.attendanceService.getStudentHistory(studentId).subscribe({
      next: (res: any) => {
        const history: StudentAttendanceHistory[] = Array.isArray(res?.data) ? res.data : [];
        this.records.set(history);

        // Compute stats
        const total = history.length;
        let present = 0, absent = 0, late = 0;

        history.forEach(r => {
          const s = r.status.toLowerCase();
          if (s === 'present' || s === 'p') present++;
          if (s === 'absent' || s === 'a') absent++;
          if (s === 'late' || s === 'l') late++;
        });

        this.totalDays.set(total);
        this.presentDays.set(present);
        this.absentDays.set(absent);
        this.lateDays.set(late);
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
      }
    });
  }
}
