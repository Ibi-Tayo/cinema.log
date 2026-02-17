import { Component, OnInit, CUSTOM_ELEMENTS_SCHEMA, signal, computed, ChangeDetectionStrategy, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { FilmService, Film } from '../../services/film.service';
import { getTMDBPosterUrl, TMDBPosterSize } from '../../utils/tmdb-image.util';

@Component({
  selector: 'app-films-to-review',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './films-to-review.component.html',
  styleUrl: './films-to-review.component.scss',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FilmsToReviewComponent implements OnInit {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private filmService = inject(FilmService);

  filmsToReview = signal<Film[]>([]);
  isLoading = signal(true);
  errorMessage = signal('');

  // Computed signals
  hasFilmsToReview = computed(() => this.filmsToReview().length > 0);

  ngOnInit(): void {
    const userId = this.route.snapshot.paramMap.get('userId');
    if (userId) {
      this.loadFilmsToReview(userId);
    } else {
      this.errorMessage.set('User ID is required');
      this.isLoading.set(false);
    }
  }

  loadFilmsToReview(userId: string): void {
    this.isLoading.set(true);
    this.filmService.getSeenUnratedFilms(userId).subscribe({
      next: (films) => {
        this.filmsToReview.set(films);
        this.isLoading.set(false);
      },
      error: (error) => {
        this.errorMessage.set(error.message || 'Failed to load films to review');
        this.isLoading.set(false);
      },
    });
  }

  selectFilm(filmId: string): void {
    this.router.navigate(['/review', filmId]);
  }

  /**
   * Gets the TMDB poster URL for a film with specified size
   * @param posterPath - The poster path from the film
   * @param size - The desired image size (default: 'w342' for display)
   */
  getPosterUrl(
    posterPath: string | null | undefined,
    size: TMDBPosterSize = 'w342',
  ): string {
    return getTMDBPosterUrl(posterPath, size);
  }
}
