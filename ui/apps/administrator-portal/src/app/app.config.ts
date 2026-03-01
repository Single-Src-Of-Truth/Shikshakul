import {
  ApplicationConfig,
  provideZonelessChangeDetection,
  provideBrowserGlobalErrorListeners,
} from '@angular/core';
import { provideRouter } from '@angular/router';
import { appRoutes } from './app.routes';
import { provideCharts, withDefaultRegisterables } from 'ng2-charts';
import { environment } from '@shikshakul/shared/environments';
import { ACAD_API_URL } from '@shikshakul/data-access/academic';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { authInterceptor, API_AUTH_TOKEN } from '@shikshakul/auth';

export const appConfig: ApplicationConfig = {
  providers: [
    provideZonelessChangeDetection(),
    provideBrowserGlobalErrorListeners(),
    provideRouter(appRoutes),
    provideCharts(withDefaultRegisterables()),
    provideHttpClient(withInterceptors([authInterceptor])),
    {
      provide: ACAD_API_URL,
      useValue: environment.acadServiceBaseUrl,
    },
    {
      provide: API_AUTH_TOKEN,
      useValue: environment.authToken,
    },
  ],
};
