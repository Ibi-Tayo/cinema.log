import { TestBed } from '@angular/core/testing';
import {
  HttpClientTestingModule,
  HttpTestingController,
} from '@angular/common/http/testing';
import {
  RatingService,
  UserFilmRating,
  ComparisonPair,
  ComparisonRequest,
} from './rating.service';
import { environment } from '../../environments/environment';

describe('RatingService', () => {
  let service: RatingService;
  let httpMock: HttpTestingController;

  const mockRating: UserFilmRating = {
    id: '123e4567-e89b-12d3-a456-426614174000',
    userId: '456e7890-e89b-12d3-a456-426614174001',
    filmId: '789e0123-e89b-12d3-a456-426614174002',
    eloRating: 1500,
    numberOfComparisons: 10,
    lastUpdated: '2024-01-01T00:00:00Z',
    initialRating: 8.5,
    kConstantValue: 32,
  };

  const mockComparisonPair: ComparisonPair = {
    filmA: mockRating,
    filmB: { ...mockRating, id: '123e4567-e89b-12d3-a456-426614174999' },
  };

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [RatingService],
    });

    service = TestBed.inject(RatingService);
    httpMock = TestBed.inject(HttpTestingController);
    spyOn(console, 'error');
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should get rating', () => {
    const userId = '456e7890-e89b-12d3-a456-426614174001';
    const filmId = '789e0123-e89b-12d3-a456-426614174002';

    service.getRating(userId, filmId).subscribe((rating) => {
      expect(rating).toEqual(mockRating);
    });

    const req = httpMock.expectOne(
      `${environment.apiUrl}/ratings?userId=${userId}&filmId=${filmId}`
    );
    expect(req.request.method).toBe('GET');
    expect(req.request.withCredentials).toBe(true);
    req.flush(mockRating);
  });

  it('should handle error when getting rating', () => {
    const userId = 'invalid-user-id';
    const filmId = 'invalid-film-id';

    service.getRating(userId, filmId).subscribe({
      next: () => fail('should have failed'),
      error: (error) => {
        expect(error.message).toContain('Failed to fetch rating');
      },
    });

    const req = httpMock.expectOne(
      `${environment.apiUrl}/ratings?userId=${userId}&filmId=${filmId}`
    );
    req.error(new ProgressEvent('error'));
  });

  it('should get ratings for comparison', () => {
    const userId = '456e7890-e89b-12d3-a456-426614174001';
    const mockRatings: UserFilmRating[] = [mockRating];

    service.getRatingsForComparison(userId).subscribe((ratings) => {
      expect(ratings).toEqual(mockRatings);
      expect(ratings.length).toBe(1);
    });

    const req = httpMock.expectOne(
      `${environment.apiUrl}/ratings/for-comparison?userId=${userId}`
    );
    expect(req.request.method).toBe('GET');
    expect(req.request.withCredentials).toBe(true);
    req.flush(mockRatings);
  });

  it('should handle error when getting ratings for comparison', () => {
    const userId = 'invalid-user-id';

    service.getRatingsForComparison(userId).subscribe({
      next: () => fail('should have failed'),
      error: (error) => {
        expect(error.message).toContain(
          'Failed to fetch ratings for comparison'
        );
      },
    });

    const req = httpMock.expectOne(
      `${environment.apiUrl}/ratings/for-comparison?userId=${userId}`
    );
    req.error(new ProgressEvent('error'));
  });

  it('should compare films', () => {
    const comparisonRequest: ComparisonRequest = {
      userId: '456e7890-e89b-12d3-a456-426614174001',
      filmAId: '789e0123-e89b-12d3-a456-426614174002',
      filmBId: '789e0123-e89b-12d3-a456-426614174003',
      winningFilmId: '789e0123-e89b-12d3-a456-426614174002',
      wasEqual: false,
    };

    service.compareFilms(comparisonRequest).subscribe((pair) => {
      expect(pair).toEqual(mockComparisonPair);
    });

    const req = httpMock.expectOne(
      `${environment.apiUrl}/ratings/compare-films`
    );
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual(comparisonRequest);
    expect(req.request.withCredentials).toBe(true);
    req.flush(mockComparisonPair);
  });

  it('should compare films when equal (no winner)', () => {
    const comparisonRequest: ComparisonRequest = {
      userId: '456e7890-e89b-12d3-a456-426614174001',
      filmAId: '789e0123-e89b-12d3-a456-426614174002',
      filmBId: '789e0123-e89b-12d3-a456-426614174003',
      wasEqual: true,
    };

    service.compareFilms(comparisonRequest).subscribe((pair) => {
      expect(pair).toEqual(mockComparisonPair);
    });

    const req = httpMock.expectOne(
      `${environment.apiUrl}/ratings/compare-films`
    );
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual(comparisonRequest);
    expect(req.request.withCredentials).toBe(true);
    req.flush(mockComparisonPair);
  });

  it('should handle error when comparing films', () => {
    const comparisonRequest: ComparisonRequest = {
      userId: 'invalid-user-id',
      filmAId: 'invalid-film-a',
      filmBId: 'invalid-film-b',
      winningFilmId: 'invalid-winner',
      wasEqual: false,
    };

    service.compareFilms(comparisonRequest).subscribe({
      next: () => fail('should have failed'),
      error: (error) => {
        expect(error.message).toContain('Failed to compare films');
      },
    });

    const req = httpMock.expectOne(
      `${environment.apiUrl}/ratings/compare-films`
    );
    req.error(new ProgressEvent('error'));
  });
});
