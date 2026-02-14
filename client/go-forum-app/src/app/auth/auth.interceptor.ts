import { HttpInterceptorFn } from '@angular/common/http';

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const token = localStorage.getItem('access_token');

  // 1. Don't add headers to the login/signup calls 
  const isAuthRequest = req.url.includes('/auth/login') || req.url.includes('/auth/register');

  // 2. If we have a token and it's not an auth request, clone and stamp it
  if (token && !isAuthRequest) {
    const authReq = req.clone({
      setHeaders: {
        Authorization: `Bearer ${token}`
      }
    });
    return next(authReq);
  }

  // 3. Otherwise, let it pass through unchanged
  return next(req);
};