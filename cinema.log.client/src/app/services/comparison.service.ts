import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, catchError, throwError } from 'rxjs';
import { environment } from '../../environments/environment';
import { ComparisonPair } from './rating.service';

export interface ComparisonRequest {
  userId: string;
  filmAId: string;
  filmBId: string;
  winningFilmId: string;
  wasEqual: boolean;
}

@Injectable({
  providedIn: 'root',
})
export class ComparisonService {
  constructor(private http: HttpClient) {}

  compareFilms(comparison: ComparisonRequest): Observable<ComparisonPair> {
    return this.http
      .post<ComparisonPair>(`${environment.apiUrl}/comparisons/compare`, comparison, {
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
