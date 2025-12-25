import { Component, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { FilmService, Film } from '../../services/film.service';
import { debounceTime, Subject } from 'rxjs';

@Component({
  selector: 'app-search',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './search.component.html',
  styleUrl: './search.component.scss',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
})
export class SearchComponent {
  searchQuery = '';
  searchResults: Film[] = [];
  isLoading = false;
  errorMessage = '';
  hasSearched = false;
  private searchSubject = new Subject<string>();

  constructor(
    private filmService: FilmService,
    private router: Router
  ) {
    // Debounce search input to avoid excessive API calls
    this.searchSubject.pipe(debounceTime(500)).subscribe((query) => {
      this.performSearch(query);
    });
  }

  onSearchInput(): void {
    if (this.searchQuery.trim().length > 0) {
      this.searchSubject.next(this.searchQuery.trim());
    } else {
      this.searchResults = [];
      this.hasSearched = false;
      this.errorMessage = '';
    }
  }

  performSearch(query: string): void {
    if (!query || query.trim().length === 0) {
      return;
    }

    this.isLoading = true;
    this.errorMessage = '';
    this.hasSearched = true;

    this.filmService.searchFilms(query).subscribe({
      next: (results) => {
        this.searchResults = results;
        this.isLoading = false;
        if (results.length === 0) {
          this.errorMessage = 'No films found matching your search.';
        }
      },
      error: (error) => {
        console.error('Error searching films:', error);
        this.errorMessage = 'Failed to search films. Please try again.';
        this.isLoading = false;
        this.searchResults = [];
      },
    });
  }

  selectFilm(film: Film): void {
    // Navigate to review page with film ID
    this.router.navigate(['/review', film.id]);
  }
}
