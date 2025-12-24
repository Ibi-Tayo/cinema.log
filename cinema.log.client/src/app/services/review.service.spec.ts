import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { ReviewService, Review, CreateReviewRequest, UpdateReviewRequest } from './review.service';
import { environment } from '../../environments/environment';

describe('ReviewService', () => {
  let service: ReviewService;
  let httpMock: HttpTestingController;

  const mockReview: Review = {
    id: '123e4567-e89b-12d3-a456-426614174000',
    title: 'Great movie!',
    date: '2024-01-01T00:00:00Z',
    rating: 9.5,
    filmId: '456e7890-e89b-12d3-a456-426614174001',
    userId: '789e0123-e89b-12d3-a456-426614174002',
  };

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [ReviewService],
    });

    service = TestBed.inject(ReviewService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should get all reviews by user id', () => {
    const userId = '789e0123-e89b-12d3-a456-426614174002';
    const mockReviews: Review[] = [mockReview];

    service.getAllReviewsByUserId(userId).subscribe((reviews) => {
      expect(reviews).toEqual(mockReviews);
      expect(reviews.length).toBe(1);
    });

    const req = httpMock.expectOne(`${environment.apiUrl}/reviews/${userId}`);
    expect(req.request.method).toBe('GET');
    expect(req.request.withCredentials).toBe(true);
    req.flush(mockReviews);
  });

  it('should handle error when getting reviews by user id', () => {
    const userId = 'invalid-user-id';

    service.getAllReviewsByUserId(userId).subscribe({
      next: () => fail('should have failed'),
      error: (error) => {
        expect(error.message).toContain('Failed to fetch reviews');
      },
    });

    const req = httpMock.expectOne(`${environment.apiUrl}/reviews/${userId}`);
    req.error(new ProgressEvent('error'));
  });

  it('should create review', () => {
    const createRequest: CreateReviewRequest = {
      content: 'Great movie!',
      rating: 9.5,
      filmId: '456e7890-e89b-12d3-a456-426614174001',
    };

    service.createReview(createRequest).subscribe((review) => {
      expect(review).toEqual(mockReview);
    });

    const req = httpMock.expectOne(`${environment.apiUrl}/reviews`);
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual(createRequest);
    expect(req.request.withCredentials).toBe(true);
    req.flush(mockReview);
  });

  it('should handle error when creating review', () => {
    const createRequest: CreateReviewRequest = {
      content: 'Great movie!',
      rating: 9.5,
      filmId: 'invalid-film-id',
    };

    service.createReview(createRequest).subscribe({
      next: () => fail('should have failed'),
      error: (error) => {
        expect(error.message).toContain('Failed to create review');
      },
    });

    const req = httpMock.expectOne(`${environment.apiUrl}/reviews`);
    req.error(new ProgressEvent('error'));
  });

  it('should update review', () => {
    const reviewId = '123e4567-e89b-12d3-a456-426614174000';
    const updateRequest: UpdateReviewRequest = {
      content: 'Updated content',
      rating: 8.5,
    };

    const updatedReview: Review = { ...mockReview, title: 'Updated content', rating: 8.5 };

    service.updateReview(reviewId, updateRequest).subscribe((review) => {
      expect(review).toEqual(updatedReview);
    });

    const req = httpMock.expectOne(`${environment.apiUrl}/reviews/${reviewId}`);
    expect(req.request.method).toBe('PUT');
    expect(req.request.body).toEqual(updateRequest);
    expect(req.request.withCredentials).toBe(true);
    req.flush(updatedReview);
  });

  it('should handle error when updating review', () => {
    const reviewId = 'invalid-review-id';
    const updateRequest: UpdateReviewRequest = {
      content: 'Updated content',
      rating: 8.5,
    };

    service.updateReview(reviewId, updateRequest).subscribe({
      next: () => fail('should have failed'),
      error: (error) => {
        expect(error.message).toContain('Failed to update review');
      },
    });

    const req = httpMock.expectOne(`${environment.apiUrl}/reviews/${reviewId}`);
    req.error(new ProgressEvent('error'));
  });

  it('should delete review', () => {
    const reviewId = '123e4567-e89b-12d3-a456-426614174000';

    service.deleteReview(reviewId).subscribe(() => {
      expect(true).toBe(true);
    });

    const req = httpMock.expectOne(`${environment.apiUrl}/reviews?id=${reviewId}`);
    expect(req.request.method).toBe('DELETE');
    expect(req.request.withCredentials).toBe(true);
    req.flush(null);
  });

  it('should handle error when deleting review', () => {
    const reviewId = 'invalid-review-id';

    service.deleteReview(reviewId).subscribe({
      next: () => fail('should have failed'),
      error: (error) => {
        expect(error.message).toContain('Failed to delete review');
      },
    });

    const req = httpMock.expectOne(`${environment.apiUrl}/reviews?id=${reviewId}`);
    req.error(new ProgressEvent('error'));
  });
});
