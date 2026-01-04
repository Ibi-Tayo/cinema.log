import { Routes } from '@angular/router';
import { AuthGuard } from './guards/auth.guard';

export const routes: Routes = [
  { path: '', redirectTo: 'home', pathMatch: 'full' },
  { 
    path: 'home', 
    loadComponent: () => import('./components/home/home.component').then(m => m.HomeComponent)
  },
  { 
    path: 'search', 
    loadComponent: () => import('./components/search/search.component').then(m => m.SearchComponent),
    canActivate: [AuthGuard] 
  },
  { 
    path: 'review/:filmId', 
    loadComponent: () => import('./components/review/review.component').then(m => m.ReviewComponent),
    canActivate: [AuthGuard] 
  },
  { 
    path: 'login', 
    loadComponent: () => import('./components/login/login.component').then(m => m.LoginComponent)
  },
  {
    path: 'profile/:id',
    loadComponent: () => import('./components/profile/profile.component').then(m => m.ProfileComponent),
    canActivate: [AuthGuard],
  },
  { 
    path: '**', 
    loadComponent: () => import('./components/login/login.component').then(m => m.LoginComponent)
  },
];
