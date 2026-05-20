import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { CoScholasticListComponent } from '../components/co-scholastic-list/co-scholastic-list.component';
import { ExamTabsComponent } from '../components/exam-tabs/exam-tabs.component';
import { PerformanceStatsComponent } from '../components/performance-stats/performance-stats.component';
import { ScholasticTableComponent } from '../components/scholastic-table/scholastic-table.component';
import { StudentSummaryComponent } from '../components/student-summary/student-summary.component';
import { SubjectChartComponent } from '../components/subject-chart/subject-chart.component';

@Component({
  selector: 'shikshakul-academic-performance-page',
  standalone: true,
  imports: [
    CommonModule,
    CoScholasticListComponent,
    ExamTabsComponent,
    PerformanceStatsComponent,
    ScholasticTableComponent,
    StudentSummaryComponent,
    SubjectChartComponent,
  ],
  templateUrl: './academic-performance-page.component.html',
  styleUrl: './academic-performance-page.component.scss',
})
export class AcademicPerformancePageComponent {}
