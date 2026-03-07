import { Route } from '@angular/router';
import { CertificateManagementPageComponent } from './certificate-management-page/certificate-management-page.component';

export const featureCertificatesRoutes: Route[] = [
  {
    path: '',
    component: CertificateManagementPageComponent,
  },
];
