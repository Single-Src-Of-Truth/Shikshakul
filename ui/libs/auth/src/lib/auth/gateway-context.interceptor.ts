import { HttpErrorResponse, HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { throwError } from 'rxjs';
import { AuthStore } from '../auth.store';

const UUID_REGEX =
  /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;

const isUuid = (value: string | null | undefined): boolean =>
  Boolean(value && UUID_REGEX.test(value));

/**
 * Prevents gateway-bound requests from being sent with invalid auth context.
 * This avoids noisy backend failures like "Invalid UUID format in gateway headers".
 */
export const gatewayContextInterceptor: HttpInterceptorFn = (req, next) => {
  const isAcademicRequest =
    req.url.startsWith('/api/v1') || req.url.includes('/api/v1/');

  if (!isAcademicRequest) {
    return next(req);
  }

  const store = inject(AuthStore);
  const profile = store.profile();
  const userId = profile?.id;
  const tenantId = profile?.tenant_id;

  if (!isUuid(userId) || !isUuid(tenantId)) {
    console.error('[Gateway Context] Blocked request due to invalid UUID context', {
      url: req.url,
      userId: userId ?? null,
      tenantId: tenantId ?? null,
    });

    return throwError(
      () =>
        new HttpErrorResponse({
          status: 0,
          statusText: 'Gateway Context Invalid',
          url: req.url,
          error: {
            error:
              'Tenant/user context not ready. Please sign in again or select a tenant.',
          },
        })
    );
  }

  return next(req);
};
