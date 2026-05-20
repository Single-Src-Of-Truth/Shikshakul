import { HttpInterceptorFn } from '@angular/common/http';

/**
 * Credentials interceptor — attaches `withCredentials: true` to every HTTP
 * request so the browser automatically sends the `skl_session` HttpOnly cookie
 * that the IAM service sets on login.
 *
 * This replaces the previous static bearer-token pattern. The token is no
 * longer read from the environment; it is managed entirely by the browser's
 * cookie jar after a successful login.
 */
export const credentialsInterceptor: HttpInterceptorFn = (req, next) => {
  return next(req.clone({ withCredentials: true }));
};
