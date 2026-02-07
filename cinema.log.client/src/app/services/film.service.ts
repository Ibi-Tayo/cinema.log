import { HttpClient } from '@angular/common/http';
import { Injectable, signal } from '@angular/core';
import { Observable, catchError, throwError, tap, of, map } from 'rxjs';
import { handleHttpError } from '../utils/error-handler.util';
import { EnvService } from './env.service';

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
  private filmCache = signal<Map<string, Film>>(new Map());
  private searchCache = signal<Map<string, Film[]>>(new Map());

  constructor(
    private http: HttpClient,
    private envService: EnvService,
  ) {}

  createFilm(film: Film): Observable<Film> {
    return this.http
      .post<Film>(`${this.envService.apiUrl}/films`, film, {
        withCredentials: true,
      })
      .pipe(
        catchError((error) => {
          console.error('Error creating film:', error);
          return throwError(
            () => new Error('Failed to create film. Please try again later.'),
          );
        }),
      );
  }

  getFilmById(id: string): Observable<Film> {
    // Check cache first
    const cached = this.filmCache().get(id);
    if (cached) {
      return of(cached);
    }

    return this.http
      .get<Film>(`${this.envService.apiUrl}/films/${id}`, {
        withCredentials: true,
      })
      .pipe(
        tap((film) => {
          // Update cache
          const currentCache = new Map(this.filmCache());
          currentCache.set(id, film);
          this.filmCache.set(currentCache);
        }),
        catchError(
          handleHttpError(
            'fetching film',
            'Failed to fetch film. Please try again later.',
          ),
        ),
      );
  }

  searchFilms(query: string): Observable<Film[]> {
    // Check cache first
    const cacheKey = query.toLowerCase().trim();
    const cached = this.searchCache().get(cacheKey);
    if (cached) {
      return of(cached);
    }

    return this.http
      .get<Film[]>(
        `${this.envService.apiUrl}/films/search?f=${encodeURIComponent(query)}`,
        {
          withCredentials: true,
        },
      )
      .pipe(
        tap((films) => {
          // Update cache
          const currentCache = new Map(this.searchCache());
          currentCache.set(cacheKey, films);
          this.searchCache.set(currentCache);

          // Also cache individual films
          const filmCache = new Map(this.filmCache());
          films.forEach((film) => filmCache.set(film.id, film));
          this.filmCache.set(filmCache);
        }),
        catchError(
          handleHttpError(
            'searching films',
            'Failed to search films. Please try again later.',
          ),
        ),
      );
  }

  getFilmsForComparison(userId: string, filmId: string): Observable<Film[]> {
    return this.http
      .get<
        Film[]
      >(`${this.envService.apiUrl}/films/for-comparison?userId=${userId}&filmId=${filmId}`, { withCredentials: true })
      .pipe(
        catchError(
          handleHttpError(
            'fetching films for comparison',
            'Failed to fetch films for comparison. Please try again later.',
          ),
        ),
      );
  }

  generateRecommendations(userId: string, films: Film[]): Observable<Film[]> {
    return this.http
      .post<
        Film[]
      >(`${this.envService.apiUrl}/films/generate-recommendations?userId=${userId}`, films, { withCredentials: true })
      .pipe(
        tap((recommendedFilms) => {
          // Cache recommended films
          const filmCache = new Map(this.filmCache());
          recommendedFilms.forEach((film) => filmCache.set(film.id, film));
          this.filmCache.set(filmCache);
        }),
        catchError(
          handleHttpError(
            'generating recommendations',
            'Failed to generate recommendations. Please try again later.',
          ),
        ),
      );
  }

  getSeenUnratedFilms(userId: string): Observable<Film[]> {
    return this.http
      .get<Film[] | null>(
        `${this.envService.apiUrl}/films/seen-unrated/${userId}`,
        {
          withCredentials: true,
        },
      )
      .pipe(
        map((films) => films || []), // Handle null response from API
        tap((films) => {
          // Cache films
          const filmCache = new Map(this.filmCache());
          films.forEach((film) => filmCache.set(film.id, film));
          this.filmCache.set(filmCache);
        }),
        catchError(
          handleHttpError(
            'fetching films to review',
            'Failed to fetch films to review. Please try again later.',
          ),
        ),
      );
  }

  /**
   * Clear all caches - useful when data needs to be refreshed
   */
  clearCache(): void {
    this.filmCache.set(new Map());
    this.searchCache.set(new Map());
  }

  /**
   * Clear specific film from cache
   */
  clearFilmCache(filmId: string): void {
    const currentCache = new Map(this.filmCache());
    currentCache.delete(filmId);
    this.filmCache.set(currentCache);
  }
}
