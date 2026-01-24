import {
  Component,
  CUSTOM_ELEMENTS_SCHEMA,
  input,
  output,
  ChangeDetectionStrategy,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { Film } from '../../../../services/film.service';
import {
  getTMDBPosterUrl,
  TMDBPosterSize,
} from '../../../../utils/tmdb-image.util';

// PrimeNG imports
import { CardModule } from 'primeng/card';
import { ButtonModule } from 'primeng/button';
import { DividerModule } from 'primeng/divider';

@Component({
  selector: 'app-sequential-comparison',
  standalone: true,
  imports: [CommonModule, CardModule, ButtonModule, DividerModule],
  templateUrl: './sequential-comparison.component.html',
  styleUrl: './sequential-comparison.component.scss',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SequentialComparisonComponent {
  // Inputs
  targetFilm = input.required<Film>();
  challengerFilm = input.required<Film>();
  progressText = input.required<string>();
  isSubmitting = input.required<boolean>();

  // Outputs
  comparisonResult = output<'better' | 'worse' | 'same'>();

  onSelectBetter(): void {
    if (!this.isSubmitting()) {
      this.comparisonResult.emit('better');
    }
  }

  onSelectWorse(): void {
    if (!this.isSubmitting()) {
      this.comparisonResult.emit('worse');
    }
  }

  onSelectSame(): void {
    if (!this.isSubmitting()) {
      this.comparisonResult.emit('same');
    }
  }

  getPosterUrl(
    posterPath: string | null | undefined,
    size: TMDBPosterSize = 'w780',
  ): string {
    return getTMDBPosterUrl(posterPath, size);
  }
}
