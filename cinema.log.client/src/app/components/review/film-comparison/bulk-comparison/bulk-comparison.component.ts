import { Component, CUSTOM_ELEMENTS_SCHEMA, input, output, computed, ChangeDetectionStrategy, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Film } from '../../../../services/film.service';
import { ComparisonStateService } from '../../../../services/comparison-state.service';
import {
  getTMDBPosterUrl,
  TMDBPosterSize,
} from '../../../../utils/tmdb-image.util';

// PrimeNG imports
import { CardModule } from 'primeng/card';
import { ButtonModule } from 'primeng/button';

@Component({
  selector: 'app-bulk-comparison',
  standalone: true,
  imports: [CommonModule, CardModule, ButtonModule],
  templateUrl: './bulk-comparison.component.html',
  styleUrl: './bulk-comparison.component.scss',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class BulkComparisonComponent {
  comparisonState = inject(ComparisonStateService);

  // Inputs
  targetFilm = input.required<Film>();
  challengerFilms = input.required<Film[]>();
  loadedFilmsCount = input.required<number>();
  maxFilms = input.required<number>();
  canLoadMore = input.required<boolean>();
  isSubmitting = input.required<boolean>();

  // Outputs
  batchSubmit = output<void>();
  loadMore = output<void>();

  // Computed
  selectedCount = computed(() => this.comparisonState.selectionCount());

  setSelection(filmId: string, result: 'better' | 'worse' | 'same'): void {
    this.comparisonState.setSelection(filmId, result);
  }

  onLoadMore(): void {
    this.loadMore.emit();
  }

  onSubmitBatch(): void {
    if (this.selectedCount() > 0 && !this.isSubmitting()) {
      this.batchSubmit.emit();
    }
  }

  getPosterUrl(
    posterPath: string | null | undefined,
    size: TMDBPosterSize = 'w780',
  ): string {
    return getTMDBPosterUrl(posterPath, size);
  }
}
