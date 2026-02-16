// app-routing.module.ts
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { authGuard } from './auth/auth-guard';
import { LoginComponent } from './login/login.component';



export const routes: Routes = [
  // 1. The Entry Point: Loaded immediately
  { 
    path: 'login', 
    component: LoginComponent 
  },

  // 2. The Forum/Posts Feature: Lazy Loaded & Protected
  { 
    path: 'posts', 
    canActivate: [authGuard],
    loadComponent: () => import('./posts/posts.component').then(m => m.PostsComponent)
  },

  // 3. The Default Redirect
  { 
    path: '', 
    redirectTo: '/posts', 
    pathMatch: 'full' 
  },

  // 4. The "Catch-All" (Optional: Redirect to login or a 404 page)
  { 
    path: '**', 
    redirectTo: '/login' 
  }
]



