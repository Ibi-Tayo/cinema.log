import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, catchError } from 'rxjs';
import { environment } from '../../environments/environment';
import { handleHttpError, handleExpectedError } from '../utils/error-handler.util';

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
  userId: string;
  filmAId: string;
  filmBId: string;
  winningFilmId?: string;
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
        // Use handleExpectedError because 404 is expected when rating doesn't exist yet
        catchError(handleExpectedError('Failed to fetch rating. Please try again later.'))
      );
  }

  compareFilms(comparison: ComparisonRequest): Observable<ComparisonPair> {
    return this.http
      .post<ComparisonPair>(`${environment.apiUrl}/ratings/compare-films`, comparison, {
        withCredentials: true,
      })
      .pipe(
        catchError(
          handleHttpError(
            'comparing films',
            'Failed to compare films. Please try again later.'
          )
        )
      );
  }
}
