// app-routing.module.ts
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';



export const routes: Routes = [
  { path: '', pathMatch:"full", loadComponent:() => {
    return import('./login/login.component').then(m => m.LoginComponent)
  }
},    { path: 'home', pathMatch:"full", loadComponent:() => {
  return import('./home/home.component').then(m => m.HomeComponent)
}
}, 
];



