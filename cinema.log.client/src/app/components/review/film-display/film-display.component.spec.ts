import { ComponentFixture, TestBed } from '@angular/core/testing';
import { FilmDisplayComponent } from './film-display.component';
import { Film } from '../../../services/film.service';
import { UserFilmRating } from '../../../services/rating.service';
import { describe, beforeEach, it, expect } from 'vitest';

describe('FilmDisplayComponent', () => {
  let component: FilmDisplayComponent;
  let fixture: ComponentFixture<FilmDisplayComponent>;

  const mockFilm: Film = {
    id: 'film1',
    externalId: 12345,
    title: 'Test Film',
    releaseYear: '2024',
    description: 'A test film description',
    posterUrl: '/test-poster.jpg',
  };

  const mockFilmRating: UserFilmRating = {
    id: 'rating1',
    userId: 'user1',
    filmId: 'film1',
    eloRating: 1500,
    numberOfComparisons: 10,
    lastUpdated: '2024-01-23T00:00:00Z',
    initialRating: 1500,
    kConstantValue: 32,
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [FilmDisplayComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(FilmDisplayComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    fixture.componentRef.setInput('film', mockFilm);
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('should display film title', () => {
    fixture.componentRef.setInput('film', mockFilm);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const titleElement = compiled.querySelector('h2');
    expect(titleElement?.textContent).toContain('Test Film');
  });

  it('should display film release year', () => {
    fixture.componentRef.setInput('film', mockFilm);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const yearElement = compiled.querySelector('.film-year');
    expect(yearElement?.textContent).toContain('2024');
  });

  it('should display film description', () => {
    fixture.componentRef.setInput('film', mockFilm);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const descriptionElement = compiled.querySelector('.film-description');
    expect(descriptionElement?.textContent).toContain(
      'A test film description',
    );
  });

  it('should display film poster when posterUrl exists', () => {
    fixture.componentRef.setInput('film', mockFilm);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const posterImg = compiled.querySelector('.film-poster img');
    expect(posterImg).toBeTruthy();
    expect(posterImg?.getAttribute('alt')).toBe('Test Film');
  });

  it('should display poster placeholder when posterUrl is null', () => {
    const filmWithoutPoster = { ...mockFilm, posterUrl: null };
    fixture.componentRef.setInput('film', filmWithoutPoster);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const placeholder = compiled.querySelector('.poster-placeholder');
    expect(placeholder).toBeTruthy();
  });

  it('should not display ELO rating section when filmRating is null', () => {
    fixture.componentRef.setInput('film', mockFilm);
    fixture.componentRef.setInput('filmRating', null);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const eloSection = compiled.querySelector('.elo-rating-section');
    expect(eloSection).toBeNull();
  });

  it('should display ELO rating section when filmRating is provided', () => {
    fixture.componentRef.setInput('film', mockFilm);
    fixture.componentRef.setInput('filmRating', mockFilmRating);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const eloSection = compiled.querySelector('.elo-rating-section');
    expect(eloSection).toBeTruthy();
  });

  it('should display ELO rating value', () => {
    fixture.componentRef.setInput('film', mockFilm);
    fixture.componentRef.setInput('filmRating', mockFilmRating);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const eloValue = compiled.querySelector('.elo-value');
    expect(eloValue?.textContent).toContain('1,500');
  });

  it('should display number of comparisons', () => {
    fixture.componentRef.setInput('film', mockFilm);
    fixture.componentRef.setInput('filmRating', mockFilmRating);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const eloHint = compiled.querySelector('.elo-hint');
    expect(eloHint?.textContent).toContain('Based on 10 comparisons');
  });

  it('should call getPosterUrl with correct parameters', () => {
    fixture.componentRef.setInput('film', mockFilm);
    fixture.detectChanges();

    const posterUrl = component.getPosterUrl('/test-poster.jpg');
    expect(posterUrl).toBeTruthy();
    expect(posterUrl).toContain('test-poster.jpg');
  });

  it('should have film-section class on container', () => {
    fixture.componentRef.setInput('film', mockFilm);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const container = compiled.querySelector('.film-section');
    expect(container).toBeTruthy();
  });
});
