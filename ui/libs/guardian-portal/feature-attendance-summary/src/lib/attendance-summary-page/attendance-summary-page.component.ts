import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { StudentHeaderComponent } from '../components/student-header/student-header.component';
import { AttendanceStatsComponent } from '../components/attendance-stats/attendance-stats.component';
import { CalendarHeatmapComponent } from '../components/calendar-heatmap/calendar-heatmap.component';
import { RecentAbsencesComponent } from '../components/recent-absences/recent-absences.component';
import { AttendanceTrendsComponent } from '../components/attendance-trends/attendance-trends.component';

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
export class AttendanceSummaryPageComponent {}
