import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { API_AUTH_TOKEN } from '../auth.token';

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const token = inject(API_AUTH_TOKEN);

  if (token) {
    const authReq = req.clone({
      setHeaders: {
        Authorization: `Bearer ${token}`,
      },
    });

    return next(authReq);
  }

  return next(req);
};
