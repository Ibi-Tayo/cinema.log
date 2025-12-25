import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, catchError, throwError } from 'rxjs';
import { environment } from '../../environments/environment';

export interface Film {
  id: string;
  externalId: number;
  title: string;
  description: string;
  posterUrl: string;
  releaseYear: string;
}

@Injectable({
  providedIn: 'root',
})
export class FilmService {
  constructor(private http: HttpClient) {}

  getFilmById(id: string): Observable<Film> {
    return this.http
      .get<Film>(`${environment.apiUrl}/films/${id}`, { withCredentials: true })
      .pipe(
        catchError((error) => {
          console.error('Error fetching film:', error);
          return throwError(() => new Error('Failed to fetch film. Please try again later.'));
        })
      );
  }

  searchFilms(query: string): Observable<Film[]> {
    return this.http
      .get<Film[]>(`${environment.apiUrl}/films/search?f=${encodeURIComponent(query)}`, {
        withCredentials: true,
      })
      .pipe(
        catchError((error) => {
          console.error('Error searching films:', error);
          return throwError(() => new Error('Failed to search films. Please try again later.'));
        })
      );
  }

  getCandidatesForComparison(): Observable<Film[]> {
    return this.http
      .get<Film[]>(`${environment.apiUrl}/films/candidates-for-comparison`, {
        withCredentials: true,
      })
      .pipe(
        catchError((error) => {
          console.error('Error fetching film candidates:', error);
          return throwError(
            () => new Error('Failed to fetch film candidates. Please try again later.')
          );
        })
      );
  }

  getFilmsForComparison(userId: string, filmId: string): Observable<Film[]> {
    return this.http
      .get<Film[]>(
        `${environment.apiUrl}/films/for-comparison?userId=${userId}&filmId=${filmId}`,
        { withCredentials: true }
      )
      .pipe(
        catchError((error) => {
          console.error('Error fetching films for comparison:', error);
          return throwError(
            () => new Error('Failed to fetch films for comparison. Please try again later.')
          );
        })
      );
  }
}
