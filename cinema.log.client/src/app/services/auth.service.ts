import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { catchError, Observable, throwError } from 'rxjs';
import { environment } from '../../environments/environment';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  currentUser: User | null = null;

  constructor(private http: HttpClient, private router: Router) {}

  login(): void {
    // Redirect directly to the GitHub login endpoint
    window.location.href = `${environment.apiUrl}/auth/github-login`;
  }

  logout(): Observable<void> {
    return this.http.get<void>(`${environment.apiUrl}/auth/logout`, { withCredentials: true }).pipe(
      catchError((error) => {
        console.error('Logout failed:', error);
        return throwError(() => new Error('Logout failed. Please try again.'));
      })
    );
  }

  requestRefreshToken(): Observable<void> {
    return this.http.get<void>(`${environment.apiUrl}/auth/refresh-token`, { withCredentials: true }).pipe(
      catchError((error) => {
        console.error('Token refresh failed:', error);
        return throwError(
          () => new Error('Authentication session expired. Please log in again.')
        );
      })
    );
  }

  getCookie(name: string): string {
    try {
      let cookieStrings = document.cookie.split(';').map((s) => s.trim());
      let map = new Map();
      cookieStrings.forEach((cookiePair) => {
        let pair = cookiePair.split('=');
        map.set(pair[0], pair[1]);
      });
      return map.get(name) || '';
    } catch (error) {
      console.error('Error parsing cookies:', error);
      return '';
    }
  }
}

export interface User {
  id: string;
  githubId: number;
  name: string;
  username: string;
  profilePicUrl: string;
  createdAt: string;
  updatedAt: string;
}
