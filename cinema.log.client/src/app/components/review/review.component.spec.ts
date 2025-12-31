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
  let mockFilmService: jasmine.SpyObj<FilmService>;
  let mockReviewService: jasmine.SpyObj<ReviewService>;
  let mockAuthService: jasmine.SpyObj<AuthService>;
  let mockRatingService: jasmine.SpyObj<RatingService>;
  let mockRouter: jasmine.SpyObj<Router>;
  let mockActivatedRoute: any;

  beforeEach(async () => {
    mockFilmService = jasmine.createSpyObj('FilmService', ['getFilmById']);
    mockReviewService = jasmine.createSpyObj('ReviewService', ['createReview']);
    mockAuthService = jasmine.createSpyObj('AuthService', ['getCurrentUser']);
    mockRatingService = jasmine.createSpyObj('RatingService', [
      'getRating',
      'getFilmsForComparison',
      'compareFilms',
    ]);
    mockRouter = jasmine.createSpyObj('Router', ['navigate'], {
      currentNavigation: jasmine
        .createSpy('currentNavigation')
        .and.returnValue(null),
    });
    mockActivatedRoute = {
      snapshot: {
        paramMap: {
          get: jasmine.createSpy('get').and.returnValue('test-film-id'),
        },
      },
    };

    // Setup default return values
    mockFilmService.getFilmById.and.returnValue(
      of({
        id: 'test-film-id',
        externalId: 123,
        title: 'Test Film',
        releaseYear: '2023',
        description: 'Test description',
        posterUrl: 'test-poster.jpg',
      })
    );
    mockAuthService.getCurrentUser.and.returnValue(
      of({
        id: 'test-user-id',
        githubId: 456,
        name: 'Test User',
        username: 'testuser',
        profilePicUrl: 'test-profile.jpg',
        createdAt: '2023-01-01T00:00:00Z',
        updatedAt: '2023-01-01T00:00:00Z',
      })
    );
    mockRatingService.getRating.and.returnValue(of(null as any));

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
