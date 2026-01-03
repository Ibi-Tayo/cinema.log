import {
  Component,
  OnInit,
  CUSTOM_ELEMENTS_SCHEMA,
  signal,
  computed,
  ChangeDetectionStrategy,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { UserService, User } from '../../services/user.service';
import { ReviewService, Review } from '../../services/review.service';
import { FilmService, Film } from '../../services/film.service';
import { forkJoin, of, switchMap } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { getTMDBPosterUrl, TMDBPosterSize } from '../../utils/tmdb-image.util';
import { RatingService, UserFilmRating, UserFilmRatingDetail } from '../../services/rating.service';

interface ReviewWithFilm extends Review {
  film?: Film;
}

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './profile.component.html',
  styleUrl: './profile.component.scss',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ProfileComponent implements OnInit {
  user = signal<User | null>(null);
  recentReviews = signal<ReviewWithFilm[]>([]);
  userRatings = signal<UserFilmRatingDetail[]>([]);
  isLoading = signal(true);
  errorMessage = signal('');
  currentCarouselIndex = signal(0);

  // Computed signals
  hasReviews = computed(() => this.recentReviews().length > 0);
  currentReview = computed(() => {
    const reviews = this.recentReviews();
    const index = this.currentCarouselIndex();
    return reviews[index];
  });

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private userService: UserService,
    private reviewService: ReviewService,
    private filmService: FilmService,
    private ratingsService: RatingService
  ) {}

  ngOnInit(): void {
    const userId = this.route.snapshot.paramMap.get('id');
    if (userId) {
      this.loadUserProfile(userId);
      this.loadUserRatings(userId);
    } else {
      this.errorMessage.set('User ID is required');
      this.isLoading.set(false);
    }
  }

  loadUserProfile(userId: string): void {
    this.isLoading.set(true);

    this.userService
      .getUserById(userId)
      .pipe(
        switchMap((user) => {
          this.user.set(user);
          return this.reviewService.getAllReviewsByUserId(userId);
        }),
        switchMap((reviews) => {
          if (!reviews || reviews.length === 0) {
            return of([]);
          }
          // Get the most recent 5 reviews
          const recentReviews = reviews.slice(0, 5);

          // Fetch film details for each review
          const filmRequests = recentReviews.map((review) =>
            this.filmService.getFilmById(review.filmId).pipe(
              map((film) => ({ ...review, film } as ReviewWithFilm)),
              catchError(() => of({ ...review } as ReviewWithFilm))
            )
          );

          return forkJoin(filmRequests);
        }),
        catchError((error) => {
          this.errorMessage.set(error.message || 'Failed to load profile');
          return of([]);
        })
      )
      .subscribe((reviewsWithFilms) => {
        this.recentReviews.set(reviewsWithFilms);
        this.isLoading.set(false);
      });
  }

  loadUserRatings(userId: string): void {
    this.ratingsService.getRatingsByUserId(userId).subscribe({
      next: (ratings) => {
        this.userRatings.set(ratings);
      },
      error: (error) => {
        this.errorMessage.set(error.message || 'Failed to load user ratings');
      },
    });
  }

  nextCarouselSlide(): void {
    const reviews = this.recentReviews();
    if (reviews.length > 0) {
      this.currentCarouselIndex.set(
        (this.currentCarouselIndex() + 1) % reviews.length
      );
    }
  }

  previousCarouselSlide(): void {
    const reviews = this.recentReviews();
    if (reviews.length > 0) {
      this.currentCarouselIndex.set(
        (this.currentCarouselIndex() - 1 + reviews.length) % reviews.length
      );
    }
  }

  goToSlide(index: number): void {
    this.currentCarouselIndex.set(index);
  }

  getStars(rating: number): string[] {
    return Array(5)
      .fill('star')
      .map((_, i) => (i < rating ? 'full' : 'empty'));
  }

  selectFilm(filmId: string): void {
    // Navigate to review page with film ID
    this.router.navigate(['/review', filmId]);
  }

  /**
   * Gets the TMDB poster URL for a film with specified size
   * @param posterPath - The poster path from the film
   * @param size - The desired image size (default: 'w342' for carousel display)
   */
  getPosterUrl(
    posterPath: string | null | undefined,
    size: TMDBPosterSize = 'w342'
  ): string {
    return getTMDBPosterUrl(posterPath, size);
  }
}
