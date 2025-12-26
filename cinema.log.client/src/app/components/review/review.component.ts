import { Component, OnInit, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { FilmService, Film } from '../../services/film.service';
import {
  ReviewService,
  CreateReviewRequest,
} from '../../services/review.service';
import { AuthService } from '../../services/auth.service';
import { RatingService, UserFilmRating } from '../../services/rating.service';
import {
  ComparisonService,
  ComparisonRequest,
} from '../../services/comparison.service';
import { getTMDBPosterUrl, TMDBPosterSize } from '../../utils/tmdb-image.util';

@Component({
  selector: 'app-review',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './review.component.html',
  styleUrl: './review.component.scss',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
})
export class ReviewComponent implements OnInit {
  private static readonly REDIRECT_DELAY_MS = 1500;

  film: Film | null = null;
  isLoading = true;
  errorMessage = '';

  // Review form data
  selectedRating = 0;
  reviewText = '';
  isSubmitting = false;
  submitSuccess = false;
  submitError = '';

  // ELO rating data
  filmRating: UserFilmRating | null = null;
  hasRating = false;

  // Comparison data
  filmsForComparison: Film[] = [];
  isLoadingComparisons = false;
  comparisonError = '';
  currentComparisonIndex = 0;
  isSubmittingComparison = false;
  showComparisons = false;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private filmService: FilmService,
    private reviewService: ReviewService,
    private authService: AuthService,
    private ratingService: RatingService,
    private comparisonService: ComparisonService
  ) {}

  ngOnInit(): void {
    const filmId = this.route.snapshot.paramMap.get('filmId');
    if (filmId) {
      this.loadFilm(filmId);
    } else {
      this.errorMessage = 'Film ID is required';
      this.isLoading = false;
    }
  }

  loadFilm(filmId: string): void {
    this.isLoading = true;
    this.filmService.getFilmById(filmId).subscribe({
      next: (film) => {
        this.film = film;
        this.isLoading = false;
        // Check if the user has already rated this film
        this.checkForExistingRating(filmId);
      },
      error: (error) => {
        console.error('Error loading film:', error);
        this.errorMessage = 'Failed to load film details. Please try again.';
        this.isLoading = false;
      },
    });
  }

  checkForExistingRating(filmId: string): void {
    const userId = this.authService.currentUser?.id;
    if (!userId) return;

    this.ratingService.getRating(userId, filmId).subscribe({
      next: (rating) => {
        this.filmRating = rating;
        this.hasRating = true;
      },
      error: () => {
        // Rating doesn't exist yet, which is fine
        this.hasRating = false;
      },
    });
  }

  selectRating(rating: number): void {
    this.selectedRating = rating;
  }

  submitReview(): void {
    if (!this.film || this.selectedRating === 0 || !this.reviewText.trim()) {
      this.submitError = 'Please provide both a rating and a review.';
      return;
    }

    this.isSubmitting = true;
    this.submitError = '';
    this.submitSuccess = false;

    const reviewRequest: CreateReviewRequest = {
      filmId: this.film.id,
      rating: this.selectedRating,
      content: this.reviewText.trim(),
    };

    this.reviewService.createReview(reviewRequest).subscribe({
      next: () => {
        this.submitSuccess = true;
        this.isSubmitting = false;
        // Reset form
        this.selectedRating = 0;
        this.reviewText = '';
        // Load the rating and comparisons
        if (this.film) {
          this.checkForExistingRating(this.film.id);
          this.loadFilmsForComparison();
        }
      },
      error: (error) => {
        console.error('Error submitting review:', error);
        this.submitError = 'Failed to submit review. Please try again.';
        this.isSubmitting = false;
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
    const userId = this.authService.currentUser?.id;
    if (!userId || !this.film) return;

    this.isLoadingComparisons = true;
    this.comparisonError = '';
    this.showComparisons = true;

    this.filmService.getFilmsForComparison(userId, this.film.id).subscribe({
      next: (films) => {
        this.filmsForComparison = films;
        this.isLoadingComparisons = false;
        this.currentComparisonIndex = 0;
      },
      error: (error) => {
        console.error('Error loading films for comparison:', error);
        this.comparisonError = 'Failed to load films for comparison.';
        this.isLoadingComparisons = false;
      },
    });
  }

  getCurrentComparisonFilm(): Film | null {
    if (this.currentComparisonIndex < this.filmsForComparison.length) {
      return this.filmsForComparison[this.currentComparisonIndex];
    }
    return null;
  }

  submitComparison(result: 'better' | 'worse' | 'same'): void {
    const userId = this.authService.currentUser?.id;
    if (!userId || !this.film) return;

    const comparisonFilm = this.getCurrentComparisonFilm();
    if (!comparisonFilm) return;

    this.isSubmittingComparison = true;
    this.comparisonError = '';

    let winningFilmId: string;
    let wasEqual = false;

    if (result === 'better') {
      winningFilmId = this.film.id;
    } else if (result === 'worse') {
      winningFilmId = comparisonFilm.id;
    } else {
      // same
      winningFilmId = this.film.id; // For equal, we'll use the current film's ID
      wasEqual = true;
    }

    const comparisonRequest: ComparisonRequest = {
      userId: userId,
      filmAId: this.film.id,
      filmBId: comparisonFilm.id,
      winningFilmId: winningFilmId,
      wasEqual: wasEqual,
    };

    this.comparisonService.compareFilms(comparisonRequest).subscribe({
      next: () => {
        this.isSubmittingComparison = false;
        // Refresh the rating
        if (this.film) {
          this.checkForExistingRating(this.film.id);
        }
        // Move to next comparison or finish
        if (this.currentComparisonIndex < this.filmsForComparison.length - 1) {
          this.currentComparisonIndex++;
        } else {
          // All comparisons done, redirect to profile
          this.finishComparisons();
        }
      },
      error: (error) => {
        console.error('Error submitting comparison:', error);
        this.comparisonError = 'Failed to submit comparison. Please try again.';
        this.isSubmittingComparison = false;
      },
    });
  }

  skipComparison(): void {
    if (this.currentComparisonIndex < this.filmsForComparison.length - 1) {
      this.currentComparisonIndex++;
    } else {
      this.finishComparisons();
    }
  }

  finishComparisons(): void {
    const userId = this.authService.currentUser?.id;
    if (userId) {
      setTimeout(() => {
        this.router.navigate(['/profile', userId]);
      }, ReviewComponent.REDIRECT_DELAY_MS);
    }
  }

  getComparisonProgress(): string {
    return `${this.currentComparisonIndex + 1} / ${
      this.filmsForComparison.length
    }`;
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
}
