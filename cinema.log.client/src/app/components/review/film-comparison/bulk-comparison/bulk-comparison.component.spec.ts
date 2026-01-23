import { ComponentFixture, TestBed } from '@angular/core/testing';
import { BulkComparisonComponent } from './bulk-comparison.component';
import { ComparisonStateService } from '../../../../services/comparison-state.service';
import { Film } from '../../../../services/film.service';

describe('BulkComparisonComponent', () => {
  let component: BulkComparisonComponent;
  let fixture: ComponentFixture<BulkComparisonComponent>;
  let comparisonStateService: ComparisonStateService;

  const mockTargetFilm: Film = {
    id: 'film1',
    externalId: 12345,
    title: 'Target Film',
    releaseYear: '2024',
    description: 'A target film',
    posterUrl: '/target-poster.jpg',
  };

  const mockChallengerFilms: Film[] = [
    {
      id: 'film2',
      externalId: 67890,
      title: 'Challenger Film 1',
      releaseYear: '2023',
      description: 'First challenger',
      posterUrl: '/challenger1.jpg',
    },
    {
      id: 'film3',
      externalId: 11111,
      title: 'Challenger Film 2',
      releaseYear: '2022',
      description: 'Second challenger',
      posterUrl: '',
    },
  ];

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [BulkComparisonComponent],
      providers: [ComparisonStateService],
    }).compileComponents();

    fixture = TestBed.createComponent(BulkComparisonComponent);
    component = fixture.componentInstance;
    comparisonStateService = TestBed.inject(ComparisonStateService);
    comparisonStateService.resetSelections();
  });

  afterEach(() => {
    comparisonStateService.resetSelections();
  });

  it('should create', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilms', mockChallengerFilms);
    fixture.componentRef.setInput('loadedFilmsCount', 2);
    fixture.componentRef.setInput('maxFilms', 50);
    fixture.componentRef.setInput('canLoadMore', true);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('should display target film title', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilms', mockChallengerFilms);
    fixture.componentRef.setInput('loadedFilmsCount', 2);
    fixture.componentRef.setInput('maxFilms', 50);
    fixture.componentRef.setInput('canLoadMore', true);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    expect(compiled.textContent).toContain('Target Film');
  });

  it('should display films loaded count', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilms', mockChallengerFilms);
    fixture.componentRef.setInput('loadedFilmsCount', 2);
    fixture.componentRef.setInput('maxFilms', 50);
    fixture.componentRef.setInput('canLoadMore', true);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const filmsCount = compiled.querySelector('.films-count');
    expect(filmsCount?.textContent).toContain('2 / 50 films loaded');
  });

  it('should display all challenger films', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilms', mockChallengerFilms);
    fixture.componentRef.setInput('loadedFilmsCount', 2);
    fixture.componentRef.setInput('maxFilms', 50);
    fixture.componentRef.setInput('canLoadMore', true);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const filmItems = compiled.querySelectorAll('.bulk-film-item');
    expect(filmItems.length).toBe(2);
  });

  it('should display challenger film titles', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilms', mockChallengerFilms);
    fixture.componentRef.setInput('loadedFilmsCount', 2);
    fixture.componentRef.setInput('maxFilms', 50);
    fixture.componentRef.setInput('canLoadMore', true);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    expect(compiled.textContent).toContain('Challenger Film 1');
    expect(compiled.textContent).toContain('Challenger Film 2');
  });

  it('should set selection when button is clicked', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilms', mockChallengerFilms);
    fixture.componentRef.setInput('loadedFilmsCount', 2);
    fixture.componentRef.setInput('maxFilms', 50);
    fixture.componentRef.setInput('canLoadMore', true);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.detectChanges();

    component.setSelection('film2', 'better');

    expect(comparisonStateService.getSelection('film2')).toBe('better');
    expect(comparisonStateService.selectionCount()).toBe(1);
  });

  it('should update selectedCount when selections change', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilms', mockChallengerFilms);
    fixture.componentRef.setInput('loadedFilmsCount', 2);
    fixture.componentRef.setInput('maxFilms', 50);
    fixture.componentRef.setInput('canLoadMore', true);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.detectChanges();

    expect(component.selectedCount()).toBe(0);

    comparisonStateService.setSelection('film2', 'better');
    fixture.detectChanges();

    expect(component.selectedCount()).toBe(1);

    comparisonStateService.setSelection('film3', 'worse');
    fixture.detectChanges();

    expect(component.selectedCount()).toBe(2);
  });

  it('should emit batchSubmit when submit button is clicked with selections', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilms', mockChallengerFilms);
    fixture.componentRef.setInput('loadedFilmsCount', 2);
    fixture.componentRef.setInput('maxFilms', 50);
    fixture.componentRef.setInput('canLoadMore', true);
    fixture.componentRef.setInput('isSubmitting', false);

    let batchSubmitted = false;
    fixture.componentInstance.batchSubmit.subscribe(() => {
      batchSubmitted = true;
    });

    comparisonStateService.setSelection('film2', 'better');
    fixture.detectChanges();

    component.onSubmitBatch();

    expect(batchSubmitted).toBe(true);
  });

  it('should not emit batchSubmit when no selections', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilms', mockChallengerFilms);
    fixture.componentRef.setInput('loadedFilmsCount', 2);
    fixture.componentRef.setInput('maxFilms', 50);
    fixture.componentRef.setInput('canLoadMore', true);
    fixture.componentRef.setInput('isSubmitting', false);

    let batchSubmitted = false;
    fixture.componentInstance.batchSubmit.subscribe(() => {
      batchSubmitted = true;
    });

    fixture.detectChanges();
    component.onSubmitBatch();

    expect(batchSubmitted).toBe(false);
  });

  it('should emit loadMore when load more button is clicked', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilms', mockChallengerFilms);
    fixture.componentRef.setInput('loadedFilmsCount', 10);
    fixture.componentRef.setInput('maxFilms', 50);
    fixture.componentRef.setInput('canLoadMore', true);
    fixture.componentRef.setInput('isSubmitting', false);

    let loadMoreEmitted = false;
    fixture.componentInstance.loadMore.subscribe(() => {
      loadMoreEmitted = true;
    });

    fixture.detectChanges();
    component.onLoadMore();

    expect(loadMoreEmitted).toBe(true);
  });

  it('should show load more button when canLoadMore is true and has 10+ films', () => {
    const manyFilms = Array.from({ length: 10 }, (_, i) => ({
      ...mockChallengerFilms[0],
      id: `film${i}`,
    }));

    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilms', manyFilms);
    fixture.componentRef.setInput('loadedFilmsCount', 10);
    fixture.componentRef.setInput('maxFilms', 50);
    fixture.componentRef.setInput('canLoadMore', true);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    expect(compiled.textContent).toContain('Load More Films');
  });

  it('should not show load more button when less than 10 films', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilms', mockChallengerFilms);
    fixture.componentRef.setInput('loadedFilmsCount', 2);
    fixture.componentRef.setInput('maxFilms', 50);
    fixture.componentRef.setInput('canLoadMore', true);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    expect(compiled.textContent).not.toContain('Load More Films');
  });

  it('should display poster placeholder when posterUrl is null', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilms', mockChallengerFilms);
    fixture.componentRef.setInput('loadedFilmsCount', 2);
    fixture.componentRef.setInput('maxFilms', 50);
    fixture.componentRef.setInput('canLoadMore', true);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const placeholders = compiled.querySelectorAll('.poster-placeholder-small');
    expect(placeholders.length).toBeGreaterThan(0);
  });

  it('should highlight selected button with selected class', () => {
    fixture.componentRef.setInput('targetFilm', mockTargetFilm);
    fixture.componentRef.setInput('challengerFilms', mockChallengerFilms);
    fixture.componentRef.setInput('loadedFilmsCount', 2);
    fixture.componentRef.setInput('maxFilms', 50);
    fixture.componentRef.setInput('canLoadMore', true);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.detectChanges();

    comparisonStateService.setSelection('film2', 'better');
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const betterButtons = compiled.querySelectorAll('.bulk-button.better');
    const firstBetterButton = betterButtons[0];

    expect(firstBetterButton.classList.contains('selected')).toBe(true);
  });
});
