import { Route } from '@angular/router';
import { LayoutComponent } from '@shikshakul/teacher/feature-shell';

export const appRoutes: Route[] = [
  {
    path: '',
    component: LayoutComponent,
    children: [
      { path: '', redirectTo: 'dashboard', pathMatch: 'full' },
      {
        path: 'dashboard',
        loadComponent: () =>
          import('@shikshakul/teacher/feature-dashboard').then(
            (m) => m.DashboardComponent,
          ),
      },
      {
        path: 'classes',
        loadComponent: () =>
          import('@shikshakul/teacher/feature-classes').then(
            (m) => m.MyClassesPageComponent,
          ),
      },
      {
        path: 'attendance',
        loadComponent: () =>
          import('@shikshakul/teacher/feature-attendance').then(
            (m) => m.AttendanceEntryPageComponent,
          ),
      },
      {
        path: 'marks',
        loadComponent: () =>
          import('@shikshakul/teacher/feature-marks-entry').then(
            (m) => m.MarksEntryPageComponent,
          ),
      },
      {
        path: 'students',
        loadComponent: () =>
          import('@shikshakul/teacher/feature-my-students').then(
            (m) => m.MyStudentsPageComponent,
          ),
      },
      {
        path: 'student/:id',
        loadComponent: () =>
          import('@shikshakul/teacher/feature-student-details').then(
            (m) => m.StudentDetailsPageComponent,
          ),
      },
    ],
  },
];
