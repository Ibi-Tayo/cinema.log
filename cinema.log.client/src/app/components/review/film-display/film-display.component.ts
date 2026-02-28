import {
  Component,
  CUSTOM_ELEMENTS_SCHEMA,
  input,
  ChangeDetectionStrategy,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { Film } from '../../../services/film.service';
import { UserFilmRating } from '../../../services/rating.service';
import {
  getTMDBPosterUrl,
  TMDBPosterSize,
} from '../../../utils/tmdb-image.util';

@Component({
  selector: 'app-film-display',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './film-display.component.html',
  styleUrl: './film-display.component.scss',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FilmDisplayComponent {
  film = input.required<Film>();
  filmRating = input<UserFilmRating | null>(null);

  getStars(rating: number): string[] {
    return Array(5)
      .fill('star')
      .map((_, index) => {
        const starIndex = index + 1;
        if (rating >= starIndex) return 'full';
        if (rating >= starIndex - 0.5) return 'half';
        return 'empty';
      });
  }

  getPosterUrl(
    posterPath: string | null | undefined,
    size: TMDBPosterSize = 'original',
  ): string {
    return getTMDBPosterUrl(posterPath, size);
  }
}
