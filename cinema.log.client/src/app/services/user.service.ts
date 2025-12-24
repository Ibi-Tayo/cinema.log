import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, catchError, throwError } from 'rxjs';
import { environment } from '../../environments/environment';

export interface User {
  id: string;
  githubId: number;
  name: string;
  username: string;
  profilePicUrl: string;
  createdAt: string;
  updatedAt: string;
}

@Injectable({
  providedIn: 'root',
})
export class UserService {
  constructor(private http: HttpClient) {}

  getUserById(id: string): Observable<User> {
    return this.http
      .get<User>(`${environment.apiUrl}/users/${id}`, { withCredentials: true })
      .pipe(
        catchError((error) => {
          console.error('Error fetching user:', error);
          return throwError(() => new Error('Failed to fetch user. Please try again later.'));
        })
      );
  }

  getAllUsers(): Observable<User[]> {
    return this.http
      .get<User[]>(`${environment.apiUrl}/users`, { withCredentials: true })
      .pipe(
        catchError((error) => {
          console.error('Error fetching users:', error);
          return throwError(() => new Error('Failed to fetch users. Please try again later.'));
        })
      );
  }

  createUser(user: Partial<User>): Observable<User> {
    return this.http
      .post<User>(`${environment.apiUrl}/users`, user, { withCredentials: true })
      .pipe(
        catchError((error) => {
          console.error('Error creating user:', error);
          return throwError(() => new Error('Failed to create user. Please try again later.'));
        })
      );
  }

  updateUser(user: User): Observable<User> {
    return this.http
      .put<User>(`${environment.apiUrl}/users`, user, { withCredentials: true })
      .pipe(
        catchError((error) => {
          console.error('Error updating user:', error);
          return throwError(() => new Error('Failed to update user. Please try again later.'));
        })
      );
  }

  deleteUser(id: string): Observable<void> {
    return this.http
      .delete<void>(`${environment.apiUrl}/users/${id}`, { withCredentials: true })
      .pipe(
        catchError((error) => {
          console.error('Error deleting user:', error);
          return throwError(() => new Error('Failed to delete user. Please try again later.'));
        })
      );
  }
}
