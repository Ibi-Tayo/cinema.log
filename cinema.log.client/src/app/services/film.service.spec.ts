import { TestBed } from '@angular/core/testing';
import {
  HttpClientTestingModule,
  HttpTestingController,
} from '@angular/common/http/testing';
import { FilmService, Film } from './film.service';
import { describe, beforeEach, vi, afterEach, it, expect } from 'vitest';

describe('FilmService', () => {
  let service: FilmService;
  let httpMock: HttpTestingController;

  const mockFilm: Film = {
    id: '123e4567-e89b-12d3-a456-426614174000',
    externalId: 550,
    title: 'Fight Club',
    description:
      'An insomniac office worker and a devil-may-care soap maker form an underground fight club.',
    posterUrl: 'https://example.com/poster.jpg',
    releaseYear: '1999',
  };

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [FilmService],
    });

    service = TestBed.inject(FilmService);
    httpMock = TestBed.inject(HttpTestingController);
    vi.spyOn(console, 'error');
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should get film by id', () => {
    service
      .getFilmById('123e4567-e89b-12d3-a456-426614174000')
      .subscribe((film) => {
        expect(film).toEqual(mockFilm);
      });

    const req = httpMock.expectOne(
      `${import.meta.env.NG_APP_API_URL}/films/123e4567-e89b-12d3-a456-426614174000`,
    );
    expect(req.request.method).toBe('GET');
    expect(req.request.withCredentials).toBe(true);
    req.flush(mockFilm);
  });

  it('should handle error when getting film by id', () => {
    service.getFilmById('invalid-id').subscribe({
      next: () => expect.fail('should have failed'),
      error: (error) => {
        expect(error.message).toContain('Failed to fetch film');
      },
    });

    const req = httpMock.expectOne(
      `${import.meta.env.NG_APP_API_URL}/films/invalid-id`,
    );
    req.error(new ProgressEvent('error'));
  });

  it('should search films', () => {
    const mockFilms: Film[] = [mockFilm];
    const query = 'Fight Club';

    service.searchFilms(query).subscribe((films) => {
      expect(films).toEqual(mockFilms);
      expect(films.length).toBe(1);
    });

    const req = httpMock.expectOne(
      `${import.meta.env.NG_APP_API_URL}/films/search?f=${encodeURIComponent(query)}`,
    );
    expect(req.request.method).toBe('GET');
    expect(req.request.withCredentials).toBe(true);
    req.flush(mockFilms);
  });

  it('should handle error when searching films', () => {
    service.searchFilms('test').subscribe({
      next: () => expect.fail('should have failed'),
      error: (error) => {
        expect(error.message).toContain('Failed to search films');
      },
    });

    const req = httpMock.expectOne(
      `${import.meta.env.NG_APP_API_URL}/films/search?f=test`,
    );
    req.error(new ProgressEvent('error'));
  });
});
