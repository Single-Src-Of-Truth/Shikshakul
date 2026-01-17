import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { DashboardStatsComponent } from './components/dashboard-stats/dashboard-stats.component';
import { AttendanceChartComponent } from './components/attendance-chart/attendance-chart.component';
import { FeeCollectionChartComponent } from './components/fee-collection-chart/fee-collection-chart.component';
import { NoticeBoardComponent } from './components/notice-board/notice-board.component';
import { TransportTrackerComponent } from './components/transport-tracker/transport-tracker.component';

@Component({
  selector: 'shikshakul-dashboard',
  standalone: true,
  imports: [
    CommonModule,
    DashboardStatsComponent,
    AttendanceChartComponent,
    FeeCollectionChartComponent,
    NoticeBoardComponent,
    TransportTrackerComponent,
  ],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.scss',
})
export class DashboardComponent {}
