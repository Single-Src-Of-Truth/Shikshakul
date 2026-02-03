import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { ExamFormComponent } from '../components/exam-form/exam-form.component';
import { RecentExamsWidgetComponent } from '../components/recent-exams-widget/recent-exams-widget.component';
import { ExamGuidelinesWidgetComponent } from '../components/exam-guidelines-widget/exam-guidelines-widget.component';

@Component({
  selector: 'shikshakul-exam-setup-page',
  standalone: true,
  imports: [
    CommonModule,
    ExamFormComponent,
    RecentExamsWidgetComponent,
    ExamGuidelinesWidgetComponent,
  ],
  templateUrl: './exam-setup-page.component.html',
  styleUrl: './exam-setup-page.component.scss',
})
export class ExamSetupPageComponent {}
