import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, catchError } from 'rxjs';
import { environment } from '../../environments/environment';
import { handleHttpError } from '../utils/error-handler.util';

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

  createFilm(film: Film): Observable<Film> {
    return this.http
      .post<Film>(`${environment.apiUrl}/films`, film, {
        withCredentials: true,
      })
      .pipe(
        catchError((error) => {
          console.error('Error creating film:', error);
          return throwError(
            () => new Error('Failed to create film. Please try again later.')
          );
        })
      );
  }

  getFilmById(id: string): Observable<Film> {
    return this.http
      .get<Film>(`${environment.apiUrl}/films/${id}`, { withCredentials: true })
      .pipe(
        catchError(
          handleHttpError(
            'fetching film',
            'Failed to fetch film. Please try again later.'
          )
        )
      );
  }

  searchFilms(query: string): Observable<Film[]> {
    return this.http
      .get<Film[]>(
        `${environment.apiUrl}/films/search?f=${encodeURIComponent(query)}`,
        {
          withCredentials: true,
        }
      )
      .pipe(
        catchError(
          handleHttpError(
            'searching films',
            'Failed to search films. Please try again later.'
          )
        )
      );
  }

  getCandidatesForComparison(): Observable<Film[]> {
    return this.http
      .get<Film[]>(`${environment.apiUrl}/films/candidates-for-comparison`, {
        withCredentials: true,
      })
      .pipe(
        catchError(
          handleHttpError(
            'fetching film candidates',
            'Failed to fetch film candidates. Please try again later.'
          )
        )
      );
  }

  getFilmsForComparison(userId: string, filmId: string): Observable<Film[]> {
    return this.http
      .get<Film[]>(
        `${environment.apiUrl}/films/for-comparison?userId=${userId}&filmId=${filmId}`,
        { withCredentials: true }
      )
      .pipe(
        catchError(
          handleHttpError(
            'fetching films for comparison',
            'Failed to fetch films for comparison. Please try again later.'
          )
        )
      );
  }
}
