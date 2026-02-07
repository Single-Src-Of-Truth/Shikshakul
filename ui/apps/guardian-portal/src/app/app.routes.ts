import { Route } from '@angular/router';
import { LayoutComponent } from '@shikshakul/guardian/feature-shell';

export const appRoutes: Route[] = [
  {
    path: '',
    component: LayoutComponent,
    children: [
      { path: '', redirectTo: 'dashboard', pathMatch: 'full' },
      // {
      //   path: 'dashboard',
      //   loadComponent: () =>
      //     import('@shikshakul/guardian/feature-dashboard').then(
      //       (m) => m.DashboardComponent,
      //     ),
      // },
      {
        path: 'dashboard', // We treat the profile as the dashboard for now based on your flow
        loadComponent: () =>
          import('@shikshakul/guardian/feature-student-profile').then(
            (m) => m.StudentProfilePageComponent,
          ),
      },
      {
        path: 'attendance',
        loadComponent: () =>
          import('@shikshakul/guardian/feature-attendance-summary').then(
            (m) => m.AttendanceSummaryPageComponent,
          ),
      },
      {
        path: 'reports',
        loadComponent: () =>
          import('@shikshakul/guardian/feature-academic-performance').then(
            (m) => m.AcademicPerformancePageComponent,
          ),
      },
      {
        path: 'exams',
        loadComponent: () =>
          import('@shikshakul/guardian/feature-exam-schedule').then(
            (m) => m.ExamSchedulePageComponent,
          ),
      },
      {
        path: 'teacher',
        loadComponent: () =>
          import('@shikshakul/guardian/feature-class-teacher').then(
            (m) => m.ClassTeacherPageComponent,
          ),
      },
    ],
  },
];
