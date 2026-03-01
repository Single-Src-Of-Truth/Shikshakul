import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

export type SnackbarType = 'success' | 'error' | 'info' | 'warning';

export interface SnackbarConfig {
  title?: string;
  message: string;
  type: SnackbarType;
  duration?: number;
}

@Injectable({
  providedIn: 'root',
})
export class SnackbarService {
  private snackbarState = new BehaviorSubject<SnackbarConfig | null>(null);
  snackbar$ = this.snackbarState.asObservable();

  private timer: any;

  show(
    message: string,
    title?: string,
    type: SnackbarType = 'info',
    duration = 4000,
  ) {
    if (this.timer) clearTimeout(this.timer);

    this.snackbarState.next({ message, title, type, duration });

    this.timer = setTimeout(() => {
      this.hide();
    }, duration);
  }

  hide() {
    this.snackbarState.next(null);
  }

  success(title: string, message: string, duration = 4000) {
    this.show(message, title, 'success', duration);
  }

  error(title: string, message: string, duration = 5000) {
    this.show(message, title, 'error', duration);
  }

  info(title: string, message: string, duration = 4000) {
    this.show(message, title, 'info', duration);
  }

  warning(title: string, message: string, duration = 4000) {
    this.show(message, title, 'warning', duration);
  }
}
