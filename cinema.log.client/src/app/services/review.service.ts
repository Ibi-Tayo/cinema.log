import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, catchError } from 'rxjs';
import { handleHttpError } from '../utils/error-handler.util';
import { EnvService } from './env.service';

export interface Review {
  id: string;
  title: string;
  date: string;
  rating: number;
  filmId: string;
  userId: string;
}

export interface CreateReviewRequest {
  content: string;
  rating: number;
  filmId: string;
}

export interface UpdateReviewRequest {
  content: string;
  reviewId: string;
}

@Injectable({
  providedIn: 'root',
})
export class ReviewService {
  constructor(
    private http: HttpClient,
    private envService: EnvService,
  ) {}

  getAllReviewsByUserId(userId: string): Observable<Review[]> {
    return this.http
      .get<
        Review[]
      >(`${this.envService.apiUrl}/reviews/${userId}`, { withCredentials: true })
      .pipe(
        catchError(
          handleHttpError(
            'fetching reviews',
            'Failed to fetch reviews. Please try again later.',
          ),
        ),
      );
  }

  createReview(review: CreateReviewRequest): Observable<Review> {
    return this.http
      .post<Review>(`${this.envService.apiUrl}/reviews`, review, {
        withCredentials: true,
      })
      .pipe(
        catchError(
          handleHttpError(
            'creating review',
            'Failed to create review. Please try again later.',
          ),
        ),
      );
  }

  updateReview(review: UpdateReviewRequest): Observable<Review> {
    return this.http
      .put<Review>(
        `${this.envService.apiUrl}/reviews/${review.reviewId}`,
        review,
        { withCredentials: true },
      )
      .pipe(
        catchError(
          handleHttpError(
            'updating review',
            'Failed to update review. Please try again later.',
          ),
        ),
      );
  }

  deleteReview(id: string): Observable<void> {
    return this.http
      .delete<void>(`${this.envService.apiUrl}/reviews?id=${id}`, {
        withCredentials: true,
      })
      .pipe(
        catchError(
          handleHttpError(
            'deleting review',
            'Failed to delete review. Please try again later.',
          ),
        ),
      );
  }
}
