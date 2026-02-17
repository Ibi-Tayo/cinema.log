import { Component, OnInit, CUSTOM_ELEMENTS_SCHEMA, signal, computed, ChangeDetectionStrategy, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { FilmService, Film } from '../../services/film.service';
import {
  ReviewService,
  CreateReviewRequest,
  UpdateReviewRequest,
  Review,
} from '../../services/review.service';
import { AuthService } from '../../services/auth.service';
import {
  RatingService,
  UserFilmRating,
  ComparisonRequest,
} from '../../services/rating.service';
import { ComparisonStateService } from '../../services/comparison-state.service';

// Child components
import { LoadingStateComponent } from '../shared/loading-state/loading-state.component';
import { ErrorStateComponent } from '../shared/error-state/error-state.component';
import { FilmDisplayComponent } from './film-display/film-display.component';
import { ReviewFormComponent } from './review-form/review-form.component';
import { SequentialComparisonComponent } from './film-comparison/sequential-comparison/sequential-comparison.component';
import { BulkComparisonComponent } from './film-comparison/bulk-comparison/bulk-comparison.component';

// PrimeNG imports
import { CardModule } from 'primeng/card';
import { ButtonModule } from 'primeng/button';
import { ProgressSpinnerModule } from 'primeng/progressspinner';
import { MessageModule } from 'primeng/message';

@Component({
  selector: 'app-review',
  standalone: true,
  imports: [
    CommonModule,
    LoadingStateComponent,
    ErrorStateComponent,
    FilmDisplayComponent,
    ReviewFormComponent,
    SequentialComparisonComponent,
    BulkComparisonComponent,
    CardModule,
    ButtonModule,
    ProgressSpinnerModule,
    MessageModule,
  ],
  templateUrl: './review.component.html',
  styleUrl: './review.component.scss',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ReviewComponent implements OnInit {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private filmService = inject(FilmService);
  private reviewService = inject(ReviewService);
  private authService = inject(AuthService);
  private ratingService = inject(RatingService);
  comparisonState = inject(ComparisonStateService);

  private static readonly REDIRECT_DELAY_MS = 1500;

  // Core data signals
  film = signal<Film | null>(null);
  review = signal<Review | null>(null);
  isLoading = signal(true);
  errorMessage = signal('');

  // Review form data signals
  isSubmitting = signal(false);
  submitSuccess = signal(false);
  submitError = signal('');

  // ELO rating data signals
  filmRating = signal<UserFilmRating | null>(null);

  // Comparison data signals
  filmsForComparison = signal<Film[]>([]);
  isLoadingComparisons = signal(false);
  comparisonError = signal('');
  currentComparisonIndex = signal(0);
  isSubmittingComparison = signal(false);
  showComparisons = signal(false);
  loadedFilmsCount = signal(0);
  maxBulkFilms = signal(50);

  // Computed signals
  hasRating = computed(() => this.filmRating() !== null);
  currentComparisonFilm = computed(() => {
    const index = this.currentComparisonIndex();
    const films = this.filmsForComparison();
    return index < films.length ? films[index] : null;
  });
  comparisonProgress = computed(() => {
    const index = this.currentComparisonIndex();
    const total = this.filmsForComparison().length;
    return `${index + 1} / ${total}`;
  });
  canLoadMore = computed(() => {
    const loaded = this.loadedFilmsCount();
    const max = this.maxBulkFilms();
    return loaded < max && !this.isLoadingComparisons();
  });

  constructor() {
    const navigation = this.router.currentNavigation?.();
    const filmFromState = navigation?.extras.state?.['film'];
    if (filmFromState) {
      this.film.set(filmFromState);
    }
  }

  ngOnInit(): void {
    const filmId = this.route.snapshot.paramMap.get('filmId');
    if (filmId) {
      this.loadFilm(filmId);
    } else {
      this.errorMessage.set('Film ID is required');
      this.isLoading.set(false);
    }

    // Load comparison mode preference
    this.comparisonState.loadModePreference();
  }

  loadFilm(filmId: string): void {
    this.isLoading.set(true);
    if (this.film()) {
      // Film was passed via navigation state - user probably came from search and thus we already have the film object
      this.loadFilmFromObject(this.film()!);
    } else {
      // Film not passed via navigation, fetch from server directly
      this.loadFilmFromId(filmId);
    }
  }

  checkForExistingRating(filmId: string): void {
    const userId = this.authService.currentUser()?.id;
    if (!userId) return;

    this.ratingService.getRating(userId, filmId).subscribe({
      next: (rating) => {
        this.filmRating.set(rating);
        // Load comparisons if user already has a rating for this film
        this.loadFilmsForComparison();
      },
      error: () => {
        // Rating doesn't exist yet, which is fine
        this.filmRating.set(null);
      },
    });
  }

  checkForExistingReview(): void {
    const userId = this.authService.currentUser()?.id;
    if (!userId || !this.film()) return;

    this.reviewService.getAllReviewsByUserId(userId).subscribe({
      next: (reviews) => {
        const existingReview = reviews.find(
          (review) => review.filmId === this.film()?.id,
        );
        if (existingReview) {
          this.review.set(existingReview);
        }
      },
      error: (error) => {
        console.error('Error checking for existing review:', error);
      },
    });
  }

  // Handler for review form submission (new review)
  onReviewSubmit(data: { rating: number; content: string }): void {
    if (!this.film()) {
      this.submitError.set('Film information is missing.');
      return;
    }

    this.isSubmitting.set(true);
    this.submitError.set('');
    this.submitSuccess.set(false);

    const reviewRequest: CreateReviewRequest = {
      filmId: this.film()!.id,
      rating: data.rating,
      content: data.content,
    };

    this.reviewService.createReview(reviewRequest).subscribe({
      next: () => {
        this.submitSuccess.set(true);
        this.isSubmitting.set(false);
        // Load the rating and comparisons
        const film = this.film();
        if (film) {
          this.checkForExistingRating(film.id);
          this.loadFilmsForComparison();
        }
      },
      error: (error) => {
        console.error('Error submitting review:', error);
        this.submitError.set('Failed to submit review. Please try again.');
        this.isSubmitting.set(false);
      },
    });
  }

  // Handler for review form update (existing review)
  onReviewUpdate(content: string): void {
    if (!this.film() || !content.trim()) {
      this.submitError.set('Please provide a review.');
      return;
    }

    this.isSubmitting.set(true);
    this.submitError.set('');
    this.submitSuccess.set(false);

    const updateReviewRequest: UpdateReviewRequest = {
      reviewId: this.review()?.id || '',
      content: content.trim(),
    };

    this.reviewService.updateReview(updateReviewRequest).subscribe({
      next: () => {
        this.submitSuccess.set(true);
        this.isSubmitting.set(false);
        // Load the rating and comparisons
        const film = this.film();
        if (film) {
          this.checkForExistingRating(film.id);
          this.loadFilmsForComparison();
        }
        // wait then refresh the page
        setTimeout(() => {
          window.location.reload();
        }, ReviewComponent.REDIRECT_DELAY_MS);
      },
      error: (error) => {
        console.error('Error submitting review update:', error);
        this.submitError.set(
          'Failed to submit review update. Please try again.',
        );
        this.isSubmitting.set(false);
      },
    });
  }

  goBackToSearch(): void {
    this.router.navigate(['/search']);
  }

  loadFilmsForComparison(): void {
    const userId = this.authService.currentUser()?.id;
    const film = this.film();
    if (!userId || !film) return;

    this.isLoadingComparisons.set(true);
    this.comparisonError.set('');
    this.showComparisons.set(true);

    this.filmService.getFilmsForComparison(userId, film.id).subscribe({
      next: (films) => {
        console.log('Loaded films for comparison:', films);
        const filmList = films ?? [];
        this.filmsForComparison.set(filmList);
        this.loadedFilmsCount.set(filmList.length);
        this.isLoadingComparisons.set(false);
        this.currentComparisonIndex.set(0);
        console.log('Reset comparison index to 0, total films:', films?.length);
      },
      error: (error) => {
        console.error('Error loading films for comparison:', error);
        this.comparisonError.set('Failed to load films for comparison.');
        this.isLoadingComparisons.set(false);
      },
    });
  }

  // Handler for sequential comparison result
  onComparisonResult(result: 'better' | 'worse' | 'same'): void {
    const userId = this.authService.currentUser()?.id;
    const film = this.film();
    if (!userId || !film) return;

    const comparisonFilm = this.currentComparisonFilm();
    if (!comparisonFilm) return;

    this.isSubmittingComparison.set(true);
    this.comparisonError.set('');

    let winningFilmId: string;
    let wasEqual = false;

    if (result === 'better') {
      winningFilmId = film.id;
    } else if (result === 'worse') {
      winningFilmId = comparisonFilm.id;
    } else {
      // same
      winningFilmId = film.id; // For equal, we'll use the current film's ID
      wasEqual = true;
    }

    const comparisonRequest: ComparisonRequest = {
      userId: userId,
      filmAId: film.id,
      filmBId: comparisonFilm.id,
      winningFilmId: winningFilmId,
      wasEqual: wasEqual,
    };

    this.ratingService.compareFilms(comparisonRequest).subscribe({
      next: () => {
        this.isSubmittingComparison.set(false);
        // Refresh the rating (but don't reload comparisons)
        const currentFilm = this.film();
        if (currentFilm) {
          const userId = this.authService.currentUser()?.id;
          if (userId) {
            this.ratingService.getRating(userId, currentFilm.id).subscribe({
              next: (rating) => {
                this.filmRating.set(rating);
              },
              error: () => {
                this.filmRating.set(null);
              },
            });
          }
        }
        // Move to next comparison or finish
        const currentIndex = this.currentComparisonIndex();
        const totalFilms = this.filmsForComparison().length;
        if (currentIndex < totalFilms - 1) {
          this.currentComparisonIndex.set(currentIndex + 1);
        } else {
          // All comparisons done, redirect to profile
          this.finishComparisons();
        }
      },
      error: (error) => {
        console.error('Error submitting comparison:', error);
        this.comparisonError.set(
          'Failed to submit comparison. Please try again.',
        );
        this.isSubmittingComparison.set(false);
      },
    });
  }

  finishComparisons(): void {
    const userId = this.authService.currentUser()?.id;
    if (userId) {
      setTimeout(() => {
        this.router.navigate(['/profile', userId]);
      }, ReviewComponent.REDIRECT_DELAY_MS);
    }
  }

  // Handler for mode toggle from bulk comparison component
  onToggleMode(): void {
    this.comparisonState.toggleMode();

    // Reset and reload comparisons when switching modes
    this.filmsForComparison.set([]);
    this.comparisonState.resetSelections();
    this.currentComparisonIndex.set(0);
    this.loadedFilmsCount.set(0);

    if (this.film()) {
      this.loadFilmsForComparison();
    }
  }

  // Handler for load more from bulk comparison component
  onLoadMore(): void {
    const userId = this.authService.currentUser()?.id;
    const film = this.film();
    if (!userId || !film || this.isLoadingComparisons()) return;

    const currentCount = this.loadedFilmsCount();
    if (currentCount >= this.maxBulkFilms()) return;

    this.isLoadingComparisons.set(true);
    this.comparisonError.set('');

    // Get IDs of already loaded films to exclude them
    const currentFilms = this.filmsForComparison();
    const excludeFilmIds = currentFilms.map(f => f.id);

    this.filmService.getFilmsForComparison(userId, film.id, excludeFilmIds).subscribe({
      next: (newFilms) => {
        if (newFilms && newFilms.length > 0) {
          const currentFilms = this.filmsForComparison();
          const allFilms = [...currentFilms, ...newFilms];
          this.filmsForComparison.set(allFilms);
          this.loadedFilmsCount.set(allFilms.length);
        }
        this.isLoadingComparisons.set(false);
      },
      error: (error) => {
        console.error('Error loading more films:', error);
        this.comparisonError.set('Failed to load more films.');
        this.isLoadingComparisons.set(false);
      },
    });
  }

  // Handler for batch submit from bulk comparison component
  onBatchSubmit(): void {
    const userId = this.authService.currentUser()?.id;
    const film = this.film();
    if (!userId || !film) return;

    const selections = this.comparisonState.getAllSelections();
    if (selections.length === 0) {
      this.comparisonError.set('Please select at least one comparison.');
      return;
    }

    this.isSubmittingComparison.set(true);
    this.comparisonError.set('');

    // Convert selections to comparison items
    const comparisons = selections.map((selection) => ({
      challengerFilmId: selection.filmId,
      result: selection.result,
    }));

    const request = {
      userId,
      targetFilmId: film.id,
      comparisons,
    };

    this.ratingService.compareBatch(request).subscribe({
      next: () => {
        this.isSubmittingComparison.set(false);
        // Clear selections
        this.comparisonState.resetSelections();
        // Refresh the rating
        this.ratingService.getRating(userId, film.id).subscribe({
          next: (rating) => this.filmRating.set(rating),
          error: () => this.filmRating.set(null),
        });
        // Redirect to profile
        this.finishComparisons();
      },
      error: (error) => {
        console.error('Error submitting batch comparisons:', error);
        this.comparisonError.set(
          'Failed to submit comparisons. Please try again.',
        );
        this.isSubmittingComparison.set(false);
      },
    });
  }

  private loadFilmFromId(filmId: string) {
    this.filmService.getFilmById(filmId).subscribe({
      next: (film) => {
        this.film.set(film);
        this.isLoading.set(false);
        // Check if the user has already rated this film
        this.checkForExistingRating(filmId);
        // Check if the user has already reviewed this film
        this.checkForExistingReview();
      },
      error: (error) => {
        console.error('Error loading film:', error);
        this.errorMessage.set('Failed to load film details. Please try again.');
        this.isLoading.set(false);
      },
    });
  }

  private loadFilmFromObject(film: Film) {
    this.filmService.createFilm(film).subscribe({
      next: (film) => {
        this.film.set(film);
        this.isLoading.set(false);
        // Check if the user has already rated this film
        this.checkForExistingRating(film.id);
        // Check if the user has already reviewed this film
        this.checkForExistingReview();
      },
      error: (error) => {
        console.error('Error loading film:', error);
        this.errorMessage.set('Failed to load film details. Please try again.');
        this.isLoading.set(false);
      },
    });
  }
}
