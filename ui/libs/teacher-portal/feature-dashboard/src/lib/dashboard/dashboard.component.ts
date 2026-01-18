import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { StatsCardsComponent } from '../components/stats-cards/stats-cards.component';
import { ScheduleWidgetComponent } from '../components/schedule-widget/schedule-widget.component';
import { PendingTasksComponent } from '../components/pending-tasks/pending-tasks.component';

@Component({
  selector: 'shikshakul-teacher-dashboard',
  standalone: true,
  imports: [
    CommonModule,
    StatsCardsComponent,
    ScheduleWidgetComponent,
    PendingTasksComponent,
  ],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.scss',
})
export class DashboardComponent {
  currentDate = new Date();
}
