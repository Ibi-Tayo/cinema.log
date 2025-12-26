import { Component, OnInit, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { UserService, User } from '../../services/user.service';
import { ReviewService, Review } from '../../services/review.service';
import { FilmService, Film } from '../../services/film.service';
import { AuthService } from '../../services/auth.service';
import { forkJoin, of, switchMap } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { getTMDBPosterUrl, TMDBPosterSize } from '../../utils/tmdb-image.util';

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
})
export class ProfileComponent implements OnInit {
  user: User | null = null;
  recentReviews: ReviewWithFilm[] = [];
  isLoading = true;
  errorMessage = '';
  currentCarouselIndex = 0;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private userService: UserService,
    private reviewService: ReviewService,
    private filmService: FilmService,
    private authService: AuthService
  ) {}

  ngOnInit(): void {
    const userId = this.route.snapshot.paramMap.get('id');
    if (userId) {
      this.loadUserProfile(userId);
    } else {
      this.errorMessage = 'User ID is required';
      this.isLoading = false;
    }
  }

  loadUserProfile(userId: string): void {
    this.isLoading = true;

    this.userService
      .getUserById(userId)
      .pipe(
        switchMap((user) => {
          this.user = user;
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
          this.errorMessage = error.message || 'Failed to load profile';
          return of([]);
        })
      )
      .subscribe((reviewsWithFilms) => {
        this.recentReviews = reviewsWithFilms;
        this.isLoading = false;
      });
  }

  nextCarouselSlide(): void {
    if (this.recentReviews.length > 0) {
      this.currentCarouselIndex =
        (this.currentCarouselIndex + 1) % this.recentReviews.length;
    }
  }

  previousCarouselSlide(): void {
    if (this.recentReviews.length > 0) {
      this.currentCarouselIndex =
        (this.currentCarouselIndex - 1 + this.recentReviews.length) %
        this.recentReviews.length;
    }
  }

  goToSlide(index: number): void {
    this.currentCarouselIndex = index;
  }

  getStars(rating: number): string[] {
    return Array(5)
      .fill('star')
      .map((_, i) => (i < rating ? 'full' : 'empty'));
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
