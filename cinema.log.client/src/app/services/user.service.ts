import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, catchError } from 'rxjs';
import { handleHttpError } from '../utils/error-handler.util';

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
      .get<User>(`${import.meta.env.NG_APP_API_URL}/users/${id}`, {
        withCredentials: true,
      })
      .pipe(
        catchError(
          handleHttpError(
            'fetching user',
            'Failed to fetch user. Please try again later.',
          ),
        ),
      );
  }

  getAllUsers(): Observable<User[]> {
    return this.http
      .get<
        User[]
      >(`${import.meta.env.NG_APP_API_URL}/users`, { withCredentials: true })
      .pipe(
        catchError(
          handleHttpError(
            'fetching users',
            'Failed to fetch users. Please try again later.',
          ),
        ),
      );
  }

  createUser(user: Partial<User>): Observable<User> {
    return this.http
      .post<User>(`${import.meta.env.NG_APP_API_URL}/users`, user, {
        withCredentials: true,
      })
      .pipe(
        catchError(
          handleHttpError(
            'creating user',
            'Failed to create user. Please try again later.',
          ),
        ),
      );
  }

  updateUser(user: User): Observable<User> {
    return this.http
      .put<User>(`${import.meta.env.NG_APP_API_URL}/users`, user, {
        withCredentials: true,
      })
      .pipe(
        catchError(
          handleHttpError(
            'updating user',
            'Failed to update user. Please try again later.',
          ),
        ),
      );
  }

  deleteUser(id: string): Observable<void> {
    return this.http
      .delete<void>(`${import.meta.env.NG_APP_API_URL}/users/${id}`, {
        withCredentials: true,
      })
      .pipe(
        catchError(
          handleHttpError(
            'deleting user',
            'Failed to delete user. Please try again later.',
          ),
        ),
      );
  }
}
