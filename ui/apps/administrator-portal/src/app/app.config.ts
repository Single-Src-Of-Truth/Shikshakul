import {
  ApplicationConfig,
  provideBrowserGlobalErrorListeners,
} from '@angular/core';
import { provideRouter } from '@angular/router';
import { appRoutes } from './app.routes';
import { provideCharts, withDefaultRegisterables } from 'ng2-charts';
import { environment } from '../environments/environment';
import { ACAD_API_URL } from '@shikshakul/data-access/academic';

export const appConfig: ApplicationConfig = {
  providers: [
    provideBrowserGlobalErrorListeners(),
    provideRouter(appRoutes),
    provideCharts(withDefaultRegisterables()),
    {
      provide: ACAD_API_URL,
      useValue: environment.acadServiceBaseUrl,
    },
  ],
};
