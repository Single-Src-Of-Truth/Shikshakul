import {
  ApplicationConfig,
  APP_INITIALIZER,
  inject,
} from '@angular/core';
import { provideRouter } from '@angular/router';
import { appRoutes } from './app.routes';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { ACAD_API_URL } from '@shikshakul/data-access/academic';
import { IAM_API_URL } from '@shikshakul/data-access/iam';
import { environment } from '@shikshakul/shared/environments';
import {
  credentialsInterceptor,
  gatewayContextInterceptor,
  AuthStore,
} from '@shikshakul/auth';
import { firstValueFrom } from 'rxjs';

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(appRoutes),
    provideHttpClient(
      withInterceptors([credentialsInterceptor, gatewayContextInterceptor])
    ),
    {
      provide: ACAD_API_URL,
      useValue: environment.apiUrl,
    },
    {
      provide: IAM_API_URL,
      useValue: environment.iamServiceBaseUrl,
    },
    {
      provide: APP_INITIALIZER,
      useFactory: () => {
        const store = inject(AuthStore);
        return () => firstValueFrom(store.loadProfile());
      },
      multi: true,
    },
  ],
};
