import { Routes } from '@angular/router';
import { HomeComponent } from './components/home/home.component';
import { SearchComponent } from './components/search/search.component';
import { LoginComponent } from './components/login/login.component';
import { ProfileComponent } from './components/profile/profile.component';
import { ReviewComponent } from './components/review/review.component';
import { AuthGuard } from './guards/auth.guard';

export const routes: Routes = [
  { path: '', redirectTo: 'home', pathMatch: 'full' },
  { path: 'home', component: HomeComponent },
  { path: 'search', component: SearchComponent, canActivate: [AuthGuard] },
  { path: 'review', component: ReviewComponent, canActivate: [AuthGuard] },
  { path: 'login', component: LoginComponent },
  {
    path: 'profile/:id',
    component: ProfileComponent,
    canActivate: [AuthGuard],
  },
  { path: '**', component: LoginComponent }, // this will be a dedicated 404 component
];
