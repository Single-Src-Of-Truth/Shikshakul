import { Route } from '@angular/router';
import { CreateStudentPageComponent } from './create-student-page/create-student-page.component';
import { StudentListComponent } from './components/student-list/student-list.component';
import { StudentDetailComponent } from './components/student-detail/student-detail.component';

export const featureStudentsRoutes: Route[] = [
  {
    path: '',
    component: StudentListComponent,
  },
  {
    path: 'onboard',
    component: CreateStudentPageComponent,
  },
  {
    path: ':id',
    component: StudentDetailComponent,
  },
];
