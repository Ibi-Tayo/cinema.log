import {
  Component,
  OnInit,
  CUSTOM_ELEMENTS_SCHEMA,
  signal,
  computed,
  ChangeDetectionStrategy,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
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
import { getTMDBPosterUrl, TMDBPosterSize } from '../../utils/tmdb-image.util';

// PrimeNG imports
import { CardModule } from 'primeng/card';
import { ButtonModule } from 'primeng/button';
import { ProgressSpinnerModule } from 'primeng/progressspinner';
import { MessageModule } from 'primeng/message';
import { DividerModule } from 'primeng/divider';

@Component({
  selector: 'app-review',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    CardModule,
    ButtonModule,
    ProgressSpinnerModule,
    MessageModule,
    DividerModule,
  ],
  templateUrl: './review.component.html',
  styleUrl: './review.component.scss',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ReviewComponent implements OnInit {
  private static readonly REDIRECT_DELAY_MS = 1500;

  // Core data signals
  film = signal<Film | null>(null);
  review = signal<Review | null>(null);
  isLoading = signal(true);
  errorMessage = signal('');

  // Review form data signals
  selectedRating = signal(0);
  reviewText = signal('');
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

  // Computed signals
  hasRating = computed(() => this.filmRating() !== null);
  canSubmit = computed(
    () => this.selectedRating() > 0 && this.reviewText().trim().length > 0
  );
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

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private filmService: FilmService,
    private reviewService: ReviewService,
    private authService: AuthService,
    private ratingService: RatingService
  ) {
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
          (review) => review.filmId === this.film()?.id
        );
        if (existingReview) {
          this.reviewText.set(existingReview.title);
          this.review.set(existingReview);
        }
      },
      error: (error) => {
        console.error('Error checking for existing review:', error);
      },
    });
  }

  selectRating(rating: number): void {
    this.selectedRating.set(rating);
  }

  submitReview(): void {
    if (
      !this.film() ||
      this.selectedRating() === 0 ||
      !this.reviewText().trim()
    ) {
      this.submitError.set('Please provide both a rating and a review.');
      return;
    }

    this.isSubmitting.set(true);
    this.submitError.set('');
    this.submitSuccess.set(false);

    const reviewRequest: CreateReviewRequest = {
      filmId: this.film()!.id,
      rating: this.selectedRating(),
      content: this.reviewText().trim(),
    };

    this.reviewService.createReview(reviewRequest).subscribe({
      next: () => {
        this.submitSuccess.set(true);
        this.isSubmitting.set(false);
        // Reset form
        this.selectedRating.set(0);
        this.reviewText.set('');
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

  updateReview(): void {
    if (!this.film() || !this.reviewText().trim()) {
      this.submitError.set('Please provide a review.');
      return;
    }

    this.isSubmitting.set(true);
    this.submitError.set('');
    this.submitSuccess.set(false);

    const updateReviewRequest: UpdateReviewRequest = {
      reviewId: this.review()?.id || '',
      content: this.reviewText().trim(),
    };

    this.reviewService.updateReview(updateReviewRequest).subscribe({
      next: () => {
        this.submitSuccess.set(true);
        this.isSubmitting.set(false);
        // Reset form
        this.reviewText.set('');
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
          'Failed to submit review update. Please try again.'
        );
        this.isSubmitting.set(false);
      },
    });
  }

  getStars(): number[] {
    return [1, 2, 3, 4, 5];
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
        this.filmsForComparison.set(films ?? []); // This is because the server might return null
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

  submitComparison(result: 'better' | 'worse' | 'same'): void {
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
          'Failed to submit comparison. Please try again.'
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

  /**
   * Gets the TMDB poster URL for a film with specified size
   * @param posterPath
   * @param size
   */
  getPosterUrl(
    posterPath: string | null | undefined,
    size: TMDBPosterSize = 'original'
  ): string {
    return getTMDBPosterUrl(posterPath, size);
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
