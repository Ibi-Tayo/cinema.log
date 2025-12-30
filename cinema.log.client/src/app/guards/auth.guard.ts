import { Injectable } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { Router } from '@angular/router';
import { Observable, of } from 'rxjs';
import { map, catchError } from 'rxjs/operators';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class AuthGuard {
  constructor(private authService: AuthService, private router: Router) {}

  canActivate(): Observable<boolean> {
    // If we are in development mode, bypass all auth restrictions
    if (!environment.production) {
      return of(true);
    }
    // If we already have the current user cached, allow access
    if (this.authService.currentUser) {
      return of(true);
    }

    // Otherwise, fetch the current user from /auth/me
    return this.authService.getCurrentUser().pipe(
      map((user) => {
        // Cache the user and allow access
        this.authService.currentUser = user;
        return true;
      }),
      catchError(() => {
        // If authentication fails, redirect to login
        this.router.navigate(['/login']);
        return of(false);
      })
    );
  }
}
