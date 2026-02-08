import { Component, OnInit, OnDestroy, CUSTOM_ELEMENTS_SCHEMA, signal, computed, ChangeDetectionStrategy, inject } from '@angular/core';

import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { FilmService, Film } from '../../services/film.service';
import { debounceTime, Subject } from 'rxjs';
import { getTMDBPosterUrl, TMDBPosterSize } from '../../utils/tmdb-image.util';
import { DataViewModule } from 'primeng/dataview';
import { InputTextModule } from 'primeng/inputtext';
import { IconFieldModule } from 'primeng/iconfield';
import { InputIconModule } from 'primeng/inputicon';
import { CardModule } from 'primeng/card';
import { SkeletonModule } from 'primeng/skeleton';

@Component({
  selector: 'app-search',
  standalone: true,
  imports: [
    FormsModule,
    CommonModule,
    DataViewModule,
    InputTextModule,
    IconFieldModule,
    InputIconModule,
    CardModule,
    SkeletonModule,
  ],
  templateUrl: './search.component.html',
  styleUrl: './search.component.scss',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SearchComponent implements OnInit, OnDestroy {
  private filmService = inject(FilmService);
  private router = inject(Router);

  searchQuery = signal('');
  searchResults = signal<Film[]>([]);
  isLoading = signal(false);
  errorMessage = signal('');
  hasSearched = signal(false);
  private searchSubject = new Subject<string>();

  // Computed signals
  hasResults = computed(() => this.searchResults().length > 0);

  ngOnInit(): void {
    // Debounce search input to avoid excessive API calls
    this.searchSubject.pipe(debounceTime(500)).subscribe((query) => {
      this.performSearch(query);
    });
  }

  ngOnDestroy(): void {
    this.searchSubject.complete();
  }

  onSearchInput(): void {
    const query = this.searchQuery();
    if (query.trim().length > 0) {
      this.searchSubject.next(query.trim());
    } else {
      this.searchResults.set([]);
      this.hasSearched.set(false);
      this.errorMessage.set('');
    }
  }

  performSearch(query: string): void {
    if (!query || query.trim().length === 0) {
      return;
    }

    this.isLoading.set(true);
    this.errorMessage.set('');
    this.hasSearched.set(true);

    this.filmService.searchFilms(query).subscribe({
      next: (results) => {
        this.searchResults.set(results);
        this.isLoading.set(false);
        if (results.length === 0) {
          this.errorMessage.set('No films found matching your search.');
        }
      },
      error: (error) => {
        console.error('Error searching films:', error);
        this.errorMessage.set('Failed to search films. Please try again.');
        this.isLoading.set(false);
        this.searchResults.set([]);
      },
    });
  }

  selectFilm(film: Film): void {
    // Navigate to review page with film ID
    this.router.navigate(['/review', film.id], { state: { film } });
  }

  /**
   * Gets the TMDB poster URL for a film with specified size
   * @param posterPath
   * @param size - note that it should be a bit smaller than detail view
   */
  getPosterUrl(
    posterPath: string | null | undefined,
    size: TMDBPosterSize = 'w780'
  ): string {
    return getTMDBPosterUrl(posterPath, size);
  }
}
