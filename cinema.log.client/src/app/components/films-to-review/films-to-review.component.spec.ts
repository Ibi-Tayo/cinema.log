import type { MockedObject } from 'vitest';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, Router } from '@angular/router';
import { of } from 'rxjs';
import { FilmsToReviewComponent } from './films-to-review.component';
import { FilmService } from '../../services/film.service';
import { describe, beforeEach, it, expect, vi } from 'vitest';

describe('FilmsToReviewComponent', () => {
  let component: FilmsToReviewComponent;
  let fixture: ComponentFixture<FilmsToReviewComponent>;
  let mockFilmService: Partial<MockedObject<FilmService>>;
  let mockRouter: Partial<MockedObject<Router>>;
  let mockActivatedRoute: Partial<ActivatedRoute>;

  beforeEach(async () => {
    mockFilmService = {
      getSeenUnratedFilms: vi.fn().mockName('FilmService.getSeenUnratedFilms'),
    };
    mockRouter = {
      navigate: vi.fn().mockName('Router.navigate'),
    };
    mockActivatedRoute = {
      snapshot: {
        paramMap: {
          get: vi.fn().mockReturnValue('test-user-id'),
        },
      } as any,
    };

    await TestBed.configureTestingModule({
      imports: [FilmsToReviewComponent],
      providers: [
        { provide: FilmService, useValue: mockFilmService },
        { provide: Router, useValue: mockRouter },
        { provide: ActivatedRoute, useValue: mockActivatedRoute },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(FilmsToReviewComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should load films to review on init when userId is present', () => {
    const mockFilms = [
      {
        id: '1',
        externalId: 123,
        title: 'Test Film',
        description: 'Test Description',
        posterUrl: '/test.jpg',
        releaseYear: '2023-01-01',
      },
    ];
    mockFilmService.getSeenUnratedFilms = vi.fn().mockReturnValue(of(mockFilms));

    fixture.detectChanges();

    expect(mockFilmService.getSeenUnratedFilms).toHaveBeenCalledWith('test-user-id');
    expect(component.filmsToReview()).toEqual(mockFilms);
    expect(component.isLoading()).toBe(false);
  });

  it('should navigate to review page when selectFilm is called', () => {
    const filmId = 'test-film-id';

    component.selectFilm(filmId);

    expect(mockRouter.navigate).toHaveBeenCalledWith(['/review', filmId]);
  });
});
