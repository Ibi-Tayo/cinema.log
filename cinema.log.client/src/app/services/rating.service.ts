import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, catchError } from 'rxjs';
import {
  handleHttpError,
  handleExpectedError,
} from '../utils/error-handler.util';
import { EnvService } from './env.service';

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

export interface UserFilmRatingDetail {
  rating: UserFilmRating;
  filmTitle: string;
  filmReleaseYear: number;
  filmPosterURL: string;
  eloRank?: number;
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

export interface ComparisonItem {
  challengerFilmId: string;
  result: 'better' | 'worse' | 'same';
}

export interface BatchComparisonRequest {
  userId: string;
  targetFilmId: string;
  comparisons: ComparisonItem[];
}

export interface BatchComparisonResponse {
  success: boolean;
  message: string;
}

@Injectable({
  providedIn: 'root',
})
export class RatingService {
  constructor(
    private http: HttpClient,
    private envService: EnvService,
  ) {}

  getRating(userId: string, filmId: string): Observable<UserFilmRating> {
    return this.http
      .get<UserFilmRating>(
        `${this.envService.apiUrl}/ratings?userId=${userId}&filmId=${filmId}`,
        { withCredentials: true },
      )
      .pipe(
        // Use handleExpectedError because 404 is expected when rating doesn't exist yet
        catchError(
          handleExpectedError(
            'Failed to fetch rating. Please try again later.',
          ),
        ),
      );
  }

  getRatingsByUserId(userId: string): Observable<UserFilmRatingDetail[]> {
    return this.http
      .get<UserFilmRatingDetail[]>(
        `${this.envService.apiUrl}/ratings/${userId}`,
        {
          withCredentials: true,
        },
      )
      .pipe(
        catchError(
          handleHttpError(
            'fetching ratings by user ID',
            'Failed to fetch ratings. Please try again later.',
          ),
        ),
      );
  }

  compareFilms(comparison: ComparisonRequest): Observable<ComparisonPair> {
    return this.http
      .post<ComparisonPair>(
        `${this.envService.apiUrl}/ratings/compare-films`,
        comparison,
        {
          withCredentials: true,
        },
      )
      .pipe(
        catchError(
          handleHttpError(
            'comparing films',
            'Failed to compare films. Please try again later.',
          ),
        ),
      );
  }

  compareBatch(
    request: BatchComparisonRequest,
  ): Observable<BatchComparisonResponse> {
    return this.http
      .post<BatchComparisonResponse>(
        `${this.envService.apiUrl}/ratings/compare-films-batch`,
        request,
        {
          withCredentials: true,
        },
      )
      .pipe(
        catchError(
          handleHttpError(
            'comparing films in batch',
            'Failed to compare films. Please try again later.',
          ),
        ),
      );
  }
}
