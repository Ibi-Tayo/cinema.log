import { HttpClient } from '@angular/common/http';
import { Injectable, signal } from '@angular/core';
import { catchError, Observable, tap } from 'rxjs';
import { EnvService } from './env.service';
import {
  handleHttpError,
  handleExpectedError,
} from '../utils/error-handler.util';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  currentUser = signal<User | null>(null);
  env: Env;

  constructor(
    private http: HttpClient,
    private envService: EnvService,
  ) {
    this.env = this.envService.environment === Env.PROD ? Env.PROD : Env.DEV;
  }

  getCurrentUser(): Observable<User> {
    return this.http
      .get<User>(`${this.envService.apiUrl}/auth/me`, {
        withCredentials: true,
      })
      .pipe(
        tap((user) => this.currentUser.set(user)),
        // Use handleExpectedError because 401/403 is expected when user is not logged in
        catchError(
          handleExpectedError('Failed to authenticate. Please log in.'),
        ),
      );
  }

  login(): void {
    // Redirect directly to the GitHub login endpoint
    window.location.href = `${this.envService.apiUrl}/auth/github-login`;
  }

  devLogin(): Observable<void> {
    return this.http
      .get<void>(`${this.envService.apiUrl}/auth/dev/login`, {
        withCredentials: true,
      })
      .pipe(
        catchError(
          handleHttpError(
            'during dev login',
            'Dev login failed. Please try again.',
          ),
        ),
      );
  }

  googleLogin(): void {
    // Redirect directly to the Google login endpoint
    window.location.href = `${this.envService.apiUrl}/auth/google-login`;
  }

  devGoogleLogin(): Observable<void> {
    return this.http
      .post<void>(
        `${this.envService.apiUrl}/auth/dev/google-login`,
        {},
        {
          withCredentials: true,
        },
      )
      .pipe(
        catchError(
          handleHttpError(
            'during dev Google login',
            'Dev Google login failed. Please try again.',
          ),
        ),
      );
  }

  logout(): Observable<void> {
    return this.http
      .get<void>(`${this.envService.apiUrl}/auth/logout`, {
        withCredentials: true,
      })
      .pipe(
        tap(() => this.currentUser.set(null)),
        catchError(
          handleHttpError('during logout', 'Logout failed. Please try again.'),
        ),
      );
  }

  requestRefreshToken(): Observable<void> {
    return this.http
      .get<void>(`${this.envService.apiUrl}/auth/refresh-token`, {
        withCredentials: true,
      })
      .pipe(
        catchError(
          handleHttpError(
            'during token refresh',
            'Authentication session expired. Please log in again.',
          ),
        ),
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
  githubId?: number;
  googleId?: string;
  name: string;
  username: string;
  profilePicUrl: string;
  createdAt: string;
  updatedAt: string;
}

export enum Env {
  DEV = 'development',
  PROD = 'production',
}
