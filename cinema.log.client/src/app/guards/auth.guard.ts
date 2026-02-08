import { Injectable, inject } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { Router } from '@angular/router';
import { Observable, of } from 'rxjs';
import { map, catchError } from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class AuthGuard {
  private authService = inject(AuthService);
  private router = inject(Router);


  canActivate(): Observable<boolean> {
    // If we already have the current user cached, allow access
    if (this.authService.currentUser()) {
      return of(true);
    }

    // Otherwise, fetch the current user from /auth/me
    return this.authService.getCurrentUser().pipe(
      map((user) => {
        // User is cached via tap in the service, allow access
        return true;
      }),
      catchError(() => {
        // If authentication fails, redirect to login
        this.router.navigate(['/login']);
        return of(false);
      }),
    );
  }
}
