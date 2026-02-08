import { Route } from '@angular/router';

import { LayoutComponent } from '@shikshakul/admin/feature-shell';

export const appRoutes: Route[] = [
  {
    path: '',
    component: LayoutComponent,
    children: [
      {
        path: '',
        redirectTo: 'dashboard',
        pathMatch: 'full',
      },
      {
        path: 'dashboard',
        loadComponent: () =>
          import('@shikshakul/admin/feature-dashboard').then(
            (m) => m.DashboardComponent,
          ),
      },
      {
        path: 'students/create',
        loadComponent: () =>
          import('@shikshakul/admin/feature-students').then(
            (m) => m.CreateStudentPageComponent,
          ),
      },
      {
        path: 'academics/classes',
        loadComponent: () =>
          import('@shikshakul/admin/feature-academics').then(
            (m) => m.ClassManagementPageComponent,
          ),
      },
      {
        path: 'academics/subjects',
        loadComponent: () =>
          import('@shikshakul/admin/feature-academics').then(
            (m) => m.SubjectSetupPageComponent,
          ),
      },
      {
        path: 'staff/create',
        loadComponent: () =>
          import('@shikshakul/admin/feature-staff').then(
            (m) => m.CreateStaffPageComponent,
          ),
      },
      {
        path: 'exams/setup',
        loadComponent: () =>
          import('@shikshakul/admin/feature-exams').then(
            (m) => m.ExamSetupPageComponent,
          ),
      },
      {
        path: 'settings/roles',
        loadComponent: () =>
          import('@shikshakul/admin/feature-roles').then(
            (m) => m.RoleManagementPageComponent,
          ),
      },
      {
        path: 'settings/academic-years',
        loadComponent: () =>
          import('@shikshakul/admin/feature-settings').then(
            (m) => m.AcademicYearPageComponent,
          ),
      },
    ],
  },
];
