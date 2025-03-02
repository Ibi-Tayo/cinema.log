import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { catchError, Observable, throwError } from 'rxjs';
import { environment } from '../../environments/environment';
import { Response } from '../models/Response';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  userId: string;
  currentUser: User | null = null;

  constructor(private http: HttpClient, private router: Router) {
    this.userId = this.getCookie('userId');
    this.getUser(this.userId).subscribe({
      next: (res) => {
        this.currentUser = res.data;
        console.log(this.currentUser);
      },
      error: (error) => {
        console.error('Failed to retrieve user data:', error);
      },
    });
  }

  getUser(id: string): Observable<Response<User>> {
    return this.http
      .get<Response<User>>(`${environment.apiUrl}/user/${id}`)
      .pipe(
        catchError((error) => {
          console.error('Error fetching user data:', error);
          return throwError(
            () =>
              new Error('Failed to fetch user data. Please try again later.')
          );
        })
      );
  }

  login(): void {
    this.http
      .get<Response<string>>(`${environment.apiUrl}/auth/github-login`)
      .subscribe({
        next: (res) => {
          window.location.href = res.data; // redirect url within the request data
        },
        error: (error) => {
          console.error('GitHub login failed:', error);
          // Optionally redirect to error page or show error message
        },
      });
  }

  logout(): void {
    this.http.get(`${environment.apiUrl}/auth/logout`).subscribe({
      next: () => {
        this.currentUser = null;
        this.userId = '';
        this.router.navigate(['/login']);
      },
      error: (error) => {
        console.error('Logout failed:', error);
        // Still clear the user data on client side even if server logout fails
        this.currentUser = null;
        this.userId = '';
        this.router.navigate(['/login']);
      },
    });
  }

  requestRefreshToken(): Observable<any> {
    return this.http.get(`${environment.apiUrl}/auth/refresh-token`).pipe(
      catchError((error) => {
        console.error('Token refresh failed:', error);
        return throwError(
          () =>
            new Error('Authentication session expired. Please log in again.')
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
  userId: string;
  name: string;
  username: string;
  profilePicUrl: string;
}
