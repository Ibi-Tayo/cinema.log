import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, catchError, throwError } from 'rxjs';
import { environment } from '../../environments/environment';

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
  rating: number;
}

@Injectable({
  providedIn: 'root',
})
export class ReviewService {
  constructor(private http: HttpClient) {}

  getAllReviewsByUserId(userId: string): Observable<Review[]> {
    return this.http
      .get<Review[]>(`${environment.apiUrl}/reviews/${userId}`, { withCredentials: true })
      .pipe(
        catchError((error) => {
          console.error('Error fetching reviews:', error);
          return throwError(() => new Error('Failed to fetch reviews. Please try again later.'));
        })
      );
  }

  createReview(review: CreateReviewRequest): Observable<Review> {
    return this.http
      .post<Review>(`${environment.apiUrl}/reviews`, review, { withCredentials: true })
      .pipe(
        catchError((error) => {
          console.error('Error creating review:', error);
          return throwError(() => new Error('Failed to create review. Please try again later.'));
        })
      );
  }

  updateReview(id: string, review: UpdateReviewRequest): Observable<Review> {
    return this.http
      .put<Review>(`${environment.apiUrl}/reviews/${id}`, review, { withCredentials: true })
      .pipe(
        catchError((error) => {
          console.error('Error updating review:', error);
          return throwError(() => new Error('Failed to update review. Please try again later.'));
        })
      );
  }

  deleteReview(id: string): Observable<void> {
    return this.http
      .delete<void>(`${environment.apiUrl}/reviews?id=${id}`, { withCredentials: true })
      .pipe(
        catchError((error) => {
          console.error('Error deleting review:', error);
          return throwError(() => new Error('Failed to delete review. Please try again later.'));
        })
      );
  }
}
