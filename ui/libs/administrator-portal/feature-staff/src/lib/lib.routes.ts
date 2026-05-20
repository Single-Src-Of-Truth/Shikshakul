import { Route } from '@angular/router';
import { StaffListComponent } from './components/staff-list/staff-list.component';
import { CreateStaffPageComponent } from './create-staff-page/create-staff-page.component';
import { StaffDetailComponent } from './components/staff-detail/staff-detail.component';

export const featureStaffRoutes: Route[] = [
  { path: '', component: StaffListComponent },
  { path: 'create', component: CreateStaffPageComponent },
  { path: ':id', component: StaffDetailComponent },
];
