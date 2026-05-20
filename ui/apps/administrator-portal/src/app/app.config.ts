import {
  ApplicationConfig,
  provideZonelessChangeDetection,
  provideBrowserGlobalErrorListeners,
  APP_INITIALIZER,
  inject,
} from '@angular/core';
import { provideRouter } from '@angular/router';
import { appRoutes } from './app.routes';
import { provideCharts, withDefaultRegisterables } from 'ng2-charts';
import { environment } from '@shikshakul/shared/environments';
import { ACAD_API_URL } from '@shikshakul/data-access/academic';
import { IAM_API_URL } from '@shikshakul/data-access/iam';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import {
  credentialsInterceptor,
  gatewayContextInterceptor,
  AuthStore,
} from '@shikshakul/auth';
import { firstValueFrom } from 'rxjs';

export const appConfig: ApplicationConfig = {
  providers: [
    provideZonelessChangeDetection(),
    provideBrowserGlobalErrorListeners(),
    provideRouter(appRoutes),
    provideCharts(withDefaultRegisterables()),
    provideHttpClient(
      withInterceptors([credentialsInterceptor, gatewayContextInterceptor])
    ),
    {
      provide: ACAD_API_URL,
      useValue: environment.acadServiceBaseUrl,
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
