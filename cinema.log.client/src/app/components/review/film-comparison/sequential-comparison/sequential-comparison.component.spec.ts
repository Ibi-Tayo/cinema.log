import { ComponentFixture, TestBed } from '@angular/core/testing';
import { SequentialComparisonComponent } from './sequential-comparison.component';
import { Film } from '../../../../services/film.service';
import { describe, beforeEach, it, expect } from 'vitest';

describe('SequentialComparisonComponent', () => {
  let component: SequentialComparisonComponent;
  let fixture: ComponentFixture<SequentialComparisonComponent>;

  const mockTargetFilm: Film = {
    id: 'film1',
    externalId: 12345,
    title: 'Target Film',
    releaseYear: '2024',
    description: 'A target film',
    posterUrl: '/target-poster.jpg',
  };

  const mockChallengerFilm: Film = {
    id: 'film2',
    externalId: 67890,
    title: 'Challenger Film',
    releaseYear: '2023',
    description: 'A challenger film',
    posterUrl: '/challenger-poster.jpg',
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SequentialComparisonComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(SequentialComparisonComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilm', mockChallengerFilm);
    fixture.componentRef.setInput('progressText', '1 / 10');
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('should display target film title', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilm', mockChallengerFilm);
    fixture.componentRef.setInput('progressText', '1 / 10');
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    expect(compiled.textContent).toContain('Target Film');
  });

  it('should display challenger film title', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilm', mockChallengerFilm);
    fixture.componentRef.setInput('progressText', '1 / 10');
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    expect(compiled.textContent).toContain('Challenger Film');
  });

  it('should display progress text', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilm', mockChallengerFilm);
    fixture.componentRef.setInput('progressText', '3 / 10');
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const progressElement = compiled.querySelector('.comparison-progress');
    expect(progressElement?.textContent).toContain('Comparison 3 / 10');
  });

  it('should display "Which film do you prefer?" question', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilm', mockChallengerFilm);
    fixture.componentRef.setInput('progressText', '1 / 10');
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const questionElement = compiled.querySelector('.comparison-question p');
    expect(questionElement?.textContent).toContain('Which film do you prefer?');
  });

  it('should emit "better" when target film button is clicked', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilm', mockChallengerFilm);
    fixture.componentRef.setInput('progressText', '1 / 10');
    fixture.componentRef.setInput('isSubmitting', false);

    let emittedResult: string | undefined;
    fixture.componentInstance.comparisonResult.subscribe((result) => {
      emittedResult = result;
    });

    fixture.detectChanges();
    component.onSelectBetter();

    expect(emittedResult).toBe('better');
  });

  it('should emit "worse" when challenger film button is clicked', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilm', mockChallengerFilm);
    fixture.componentRef.setInput('progressText', '1 / 10');
    fixture.componentRef.setInput('isSubmitting', false);

    let emittedResult: string | undefined;
    fixture.componentInstance.comparisonResult.subscribe((result) => {
      emittedResult = result;
    });

    fixture.detectChanges();
    component.onSelectWorse();

    expect(emittedResult).toBe('worse');
  });

  it('should emit "same" when "About the Same" button is clicked', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilm', mockChallengerFilm);
    fixture.componentRef.setInput('progressText', '1 / 10');
    fixture.componentRef.setInput('isSubmitting', false);

    let emittedResult: string | undefined;
    fixture.componentInstance.comparisonResult.subscribe((result) => {
      emittedResult = result;
    });

    fixture.detectChanges();
    component.onSelectSame();

    expect(emittedResult).toBe('same');
  });

  it('should not emit when isSubmitting is true', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilm', mockChallengerFilm);
    fixture.componentRef.setInput('progressText', '1 / 10');
    fixture.componentRef.setInput('isSubmitting', true);

    let emittedResult: string | undefined;
    fixture.componentInstance.comparisonResult.subscribe((result) => {
      emittedResult = result;
    });

    fixture.detectChanges();
    component.onSelectBetter();

    expect(emittedResult).toBeUndefined();
  });

  it('should display poster placeholder when posterUrl is null', () => {
    const filmWithoutPoster = { ...mockTargetFilm, posterUrl: null };
    fixture.componentRef.setInput('targetFilm', filmWithoutPoster);
    fixture.componentRef.setInput('challengerFilm', mockChallengerFilm);
    fixture.componentRef.setInput('progressText', '1 / 10');
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const placeholders = compiled.querySelectorAll('.poster-placeholder');
    expect(placeholders.length).toBeGreaterThan(0);
  });

  it('should display both film release years', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilm', mockChallengerFilm);
    fixture.componentRef.setInput('progressText', '1 / 10');
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    expect(compiled.textContent).toContain('2024');
    expect(compiled.textContent).toContain('2023');
  });

  it('should have comparison-interface class on container', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilm', mockChallengerFilm);
    fixture.componentRef.setInput('progressText', '1 / 10');
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const container = compiled.querySelector('.comparison-interface');
    expect(container).toBeTruthy();
  });
});
