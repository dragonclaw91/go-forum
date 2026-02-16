import { CanActivateFn } from '@angular/router';
import { inject } from '@angular/core';
import { Router, UrlTree } from '@angular/router';

export const authGuard: CanActivateFn = (route, state) => {
  const router = inject(Router);


  const token = localStorage.getItem('access_token')


    // If token exists, allow them through (return true)
  // If not, return the "Map" to the login page
  return token ? true : router.parseUrl('login');

};
