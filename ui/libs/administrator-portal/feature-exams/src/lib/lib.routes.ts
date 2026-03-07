import { Route } from '@angular/router';
import { ExamSetupPageComponent } from './exam-setup-page/exam-setup-page.component';
import { EvaluationListComponent } from './evaluation-page/evaluation-page.component';
import { EvaluationMarkSheetComponent } from './components/evaluation-mark-sheet/evaluation-mark-sheet.component';
import { ResultManagementComponent } from './result-management/result-management.component';
import { ReportCardComponent } from './components/report-card/report-card.component';

export const featureExamsRoutes: Route[] = [
  {
    path: 'setup',
    component: ExamSetupPageComponent,
  },
  {
    path: 'evaluation',
    children: [
      {
        path: '',
        component: EvaluationListComponent,
      },
      {
        path: ':scheduleId',
        component: EvaluationMarkSheetComponent,
      },
    ],
  },
  {
    path: 'results',
    component: ResultManagementComponent,
  },
  {
    path: 'report-card/:studentId',
    component: ReportCardComponent,
  },
];
