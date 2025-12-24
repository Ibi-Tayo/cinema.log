import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, catchError, throwError } from 'rxjs';
import { environment } from '../../environments/environment';

export interface UserFilmRating {
  id: string;
  userId: string;
  filmId: string;
  eloRating: number;
  numberOfComparisons: number;
  lastUpdated: string;
  initialRating: number;
  kConstantValue: number;
}

export interface ComparisonPair {
  filmA: UserFilmRating;
  filmB: UserFilmRating;
}

export interface ComparisonRequest {
  id?: string;
  userId: string;
  filmAId: string;
  filmBId: string;
  winningFilmId: string;
  comparisonDate?: string;
  wasEqual: boolean;
}

@Injectable({
  providedIn: 'root',
})
export class RatingService {
  constructor(private http: HttpClient) {}

  getRating(userId: string, filmId: string): Observable<UserFilmRating> {
    return this.http
      .get<UserFilmRating>(
        `${environment.apiUrl}/ratings?userId=${userId}&filmId=${filmId}`,
        { withCredentials: true }
      )
      .pipe(
        catchError((error) => {
          console.error('Error fetching rating:', error);
          return throwError(() => new Error('Failed to fetch rating. Please try again later.'));
        })
      );
  }

  getRatingsForComparison(userId: string): Observable<UserFilmRating[]> {
    return this.http
      .get<UserFilmRating[]>(
        `${environment.apiUrl}/ratings/for-comparison?userId=${userId}`,
        { withCredentials: true }
      )
      .pipe(
        catchError((error) => {
          console.error('Error fetching ratings for comparison:', error);
          return throwError(
            () => new Error('Failed to fetch ratings for comparison. Please try again later.')
          );
        })
      );
  }

  compareFilms(comparison: ComparisonRequest): Observable<ComparisonPair> {
    return this.http
      .post<ComparisonPair>(`${environment.apiUrl}/ratings/compare-films`, comparison, {
        withCredentials: true,
      })
      .pipe(
        catchError((error) => {
          console.error('Error comparing films:', error);
          return throwError(() => new Error('Failed to compare films. Please try again later.'));
        })
      );
  }
}
