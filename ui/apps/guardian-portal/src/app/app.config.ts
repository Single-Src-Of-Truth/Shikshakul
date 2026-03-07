import {
  ApplicationConfig,
} from '@angular/core';
import { provideRouter } from '@angular/router';
import { appRoutes } from './app.routes';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { ACAD_API_URL } from '@shikshakul/data-access/academic';
import { environment } from '@shikshakul/shared/environments';
import { authInterceptor, API_AUTH_TOKEN } from '@shikshakul/auth';

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(appRoutes),
    provideHttpClient(withInterceptors([authInterceptor])),
    {
      provide: ACAD_API_URL,
      useValue: environment.apiUrl,
    },
    {
      provide: API_AUTH_TOKEN,
      useValue: environment.authToken,
    }
  ],
};
