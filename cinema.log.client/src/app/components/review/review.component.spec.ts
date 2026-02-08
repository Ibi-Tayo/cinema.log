import { beforeEach, describe, expect, it, vi, type MockedObject } from 'vitest';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, Router } from '@angular/router';
import { of } from 'rxjs';
import { ReviewComponent } from './review.component';
import { FilmService } from '../../services/film.service';
import { ReviewService } from '../../services/review.service';
import { AuthService } from '../../services/auth.service';
import { RatingService } from '../../services/rating.service';

describe('ReviewComponent', () => {
  let component: ReviewComponent;
  let fixture: ComponentFixture<ReviewComponent>;
  let mockFilmService: Partial<MockedObject<FilmService>>;
  let mockReviewService: Partial<MockedObject<ReviewService>>;
  let mockAuthService: Partial<MockedObject<AuthService>>;
  let mockRatingService: Partial<MockedObject<RatingService>>;
  let mockRouter: Partial<MockedObject<Router>>;
  let mockActivatedRoute: any;

  beforeEach(async () => {
    mockFilmService = {
      getFilmById: vi.fn().mockName('FilmService.getFilmById'),
    };
    mockReviewService = {
      createReview: vi.fn().mockName('ReviewService.createReview'),
    };
    mockAuthService = {
      getCurrentUser: vi.fn().mockName('AuthService.getCurrentUser'),
    };
    mockRatingService = {
      getRating: vi.fn().mockName('RatingService.getRating'),
      compareFilms: vi.fn().mockName('RatingService.compareFilms'),
      compareBatch: vi.fn().mockName('RatingService.compareBatch'),
    };
    mockRouter = {
      navigate: vi.fn().mockName('Router.navigate'),
    };
    mockActivatedRoute = {
      snapshot: {
        paramMap: {
          get: vi.fn().mockReturnValue('test-film-id'),
        },
      },
    };

    // Setup default return values
    mockFilmService.getFilmById!.mockReturnValue(
      of({
        id: 'test-film-id',
        externalId: 123,
        title: 'Test Film',
        releaseYear: '2023',
        description: 'Test description',
        posterUrl: 'test-poster.jpg',
      }),
    );
    mockAuthService.getCurrentUser!.mockReturnValue(
      of({
        id: 'test-user-id',
        githubId: 456,
        name: 'Test User',
        username: 'testuser',
        profilePicUrl: 'test-profile.jpg',
        createdAt: '2023-01-01T00:00:00Z',
        updatedAt: '2023-01-01T00:00:00Z',
      }),
    );
    mockRatingService.getRating!.mockReturnValue(of(null as any));

    await TestBed.configureTestingModule({
      imports: [ReviewComponent],
      providers: [
        { provide: FilmService, useValue: mockFilmService },
        { provide: ReviewService, useValue: mockReviewService },
        { provide: AuthService, useValue: mockAuthService },
        { provide: RatingService, useValue: mockRatingService },
        { provide: Router, useValue: mockRouter },
        { provide: ActivatedRoute, useValue: mockActivatedRoute },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(ReviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
