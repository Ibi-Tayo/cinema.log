import { Component, OnInit, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { FilmService, Film } from '../../services/film.service';
import { ReviewService, CreateReviewRequest } from '../../services/review.service';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-review',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './review.component.html',
  styleUrl: './review.component.scss',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
})
export class ReviewComponent implements OnInit {
  film: Film | null = null;
  isLoading = true;
  errorMessage = '';
  
  // Review form data
  selectedRating = 0;
  reviewText = '';
  isSubmitting = false;
  submitSuccess = false;
  submitError = '';

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private filmService: FilmService,
    private reviewService: ReviewService,
    private authService: AuthService
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
      },
      error: (error) => {
        console.error('Error loading film:', error);
        this.errorMessage = 'Failed to load film details. Please try again.';
        this.isLoading = false;
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
        // Redirect to profile after a short delay
        setTimeout(() => {
          const userId = this.authService.currentUser?.id;
          if (userId) {
            this.router.navigate(['/profile', userId]);
          } else {
            // Fallback to home if user is somehow not available
            this.router.navigate(['/home']);
          }
        }, 2000);
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
}
