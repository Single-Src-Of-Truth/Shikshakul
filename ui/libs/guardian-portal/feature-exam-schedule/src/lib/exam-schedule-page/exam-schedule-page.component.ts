import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { AdmitCardBannerComponent } from '../components/admit-card-banner/admit-card-banner.component';
import { ExamStatsComponent } from '../components/exam-stats/exam-stats.component';
import { ScheduleTableComponent } from '../components/schedule-table/schedule-table.component';
import { ExamInstructionsComponent } from '../components/exam-instructions/exam-instructions.component';

@Component({
  selector: 'shikshakul-exam-schedule-page',
  imports: [
    CommonModule,
    AdmitCardBannerComponent,
    ExamStatsComponent,
    ScheduleTableComponent,
    ExamInstructionsComponent,
  ],
  templateUrl: './exam-schedule-page.component.html',
  styleUrl: './exam-schedule-page.component.scss',
})
export class ExamSchedulePageComponent {}
