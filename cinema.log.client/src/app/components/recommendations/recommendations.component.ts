import {
  Component,
  OnInit,
  OnDestroy,
  CUSTOM_ELEMENTS_SCHEMA,
  signal,
  computed,
  ChangeDetectionStrategy,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { FilmService, Film } from '../../services/film.service';
import { getTMDBPosterUrl, TMDBPosterSize } from '../../utils/tmdb-image.util';
import { debounceTime, Subject } from 'rxjs';

// PrimeNG imports
import { ButtonModule } from 'primeng/button';
import { CardModule } from 'primeng/card';
import { InputTextModule } from 'primeng/inputtext';
import { IconFieldModule } from 'primeng/iconfield';
import { InputIconModule } from 'primeng/inputicon';
import { DataViewModule } from 'primeng/dataview';
import { SkeletonModule } from 'primeng/skeleton';
import { ChipModule } from 'primeng/chip';
import { ToastModule } from 'primeng/toast';
import { MessageService } from 'primeng/api';

enum RecommendationStep {
  SEED_SELECTION = 'seed',
  RECOMMENDATIONS = 'recommendations',
}

interface FilmSelection {
  film: Film;
  hasSeen: boolean;
}

@Component({
  selector: 'app-recommendations',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ButtonModule,
    CardModule,
    InputTextModule,
    IconFieldModule,
    InputIconModule,
    DataViewModule,
    SkeletonModule,
    ChipModule,
    ToastModule,
  ],
  templateUrl: './recommendations.component.html',
  styleUrl: './recommendations.component.scss',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  changeDetection: ChangeDetectionStrategy.OnPush,
  providers: [MessageService],
})
export class RecommendationsComponent implements OnInit, OnDestroy {
  userId = signal<string>('');
  currentStep = signal<RecommendationStep>(RecommendationStep.SEED_SELECTION);

  // Seed selection
  seedFilms = signal<Film[]>([]);
  searchQuery = signal('');
  searchResults = signal<Film[]>([]);
  isSearching = signal(false);
  searchError = signal('');
  private searchSubject = new Subject<string>();

  // Recommendations
  recommendedFilms = signal<FilmSelection[]>([]);
  isLoadingRecommendations = signal(false);
  recommendationsError = signal('');
  currentRound = signal(1);

  // Computed signals
  canStartGenerating = computed(() => this.seedFilms().length === 3);
  selectedSeenCount = computed(
    () => this.recommendedFilms().filter((f) => f.hasSeen).length,
  );
  canProceedToNextRound = computed(() => {
    const count = this.selectedSeenCount();
    return count > 0 && count <= 10;
  });
  isInSeedStep = computed(
    () => this.currentStep() === RecommendationStep.SEED_SELECTION,
  );
  isInRecommendationStep = computed(
    () => this.currentStep() === RecommendationStep.RECOMMENDATIONS,
  );

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private filmService: FilmService,
    private messageService: MessageService,
  ) {}

  ngOnInit(): void {
    const userId = this.route.snapshot.paramMap.get('userId');
    if (userId) {
      this.userId.set(userId);
    } else {
      this.messageService.add({
        severity: 'error',
        summary: 'Error',
        detail: 'User ID is required',
      });
      this.router.navigate(['/home']);
    }

    // Setup search debounce
    this.searchSubject.pipe(debounceTime(500)).subscribe((query) => {
      this.performSearch(query);
    });
  }

  ngOnDestroy(): void {
    this.searchSubject.complete();
  }

  // Seed selection methods
  onSearchInput(): void {
    const query = this.searchQuery();
    if (query.trim().length > 0) {
      this.searchSubject.next(query.trim());
    } else {
      this.searchResults.set([]);
      this.searchError.set('');
    }
  }

  performSearch(query: string): void {
    if (!query || query.trim().length === 0) {
      return;
    }

    this.isSearching.set(true);
    this.searchError.set('');

    this.filmService.searchFilms(query).subscribe({
      next: (results) => {
        this.searchResults.set(results);
        this.isSearching.set(false);
        if (results.length === 0) {
          this.searchError.set('No films found matching your search.');
        }
      },
      error: (error) => {
        console.error('Error searching films:', error);
        this.searchError.set('Failed to search films. Please try again.');
        this.isSearching.set(false);
        this.searchResults.set([]);
      },
    });
  }

  addSeedFilm(film: Film): void {
    const seeds = this.seedFilms();
    if (seeds.length >= 3) {
      this.messageService.add({
        severity: 'warn',
        summary: 'Maximum Reached',
        detail: 'You can only select 3 seed films',
      });
      return;
    }

    if (seeds.find((f) => f.id === film.id)) {
      this.messageService.add({
        severity: 'info',
        summary: 'Already Added',
        detail: 'This film is already in your seed list',
      });
      return;
    }

    this.seedFilms.set([...seeds, film]);
    this.searchQuery.set('');
    this.searchResults.set([]);
  }

  removeSeedFilm(film: Film): void {
    this.seedFilms.set(this.seedFilms().filter((f) => f.id !== film.id));
  }

  startGeneratingRecommendations(): void {
    if (!this.canStartGenerating()) {
      return;
    }

    this.isLoadingRecommendations.set(true);
    this.recommendationsError.set('');

    this.filmService
      .generateRecommendations(this.userId(), this.seedFilms())
      .subscribe({
        next: (recommendations) => {
          this.recommendedFilms.set(
            recommendations.map((film) => ({ film, hasSeen: false })),
          );
          this.currentStep.set(RecommendationStep.RECOMMENDATIONS);
          this.isLoadingRecommendations.set(false);
          this.messageService.add({
            severity: 'success',
            summary: 'Success',
            detail: `Generated ${recommendations.length} recommendations`,
          });
        },
        error: (error) => {
          console.error('Error generating recommendations:', error);
          this.recommendationsError.set(
            'Failed to generate recommendations. Please try again.',
          );
          this.isLoadingRecommendations.set(false);
          this.messageService.add({
            severity: 'error',
            summary: 'Error',
            detail: 'Failed to generate recommendations',
          });
        },
      });
  }

  // Recommendation methods
  toggleFilmSeen(selection: FilmSelection): void {
    const count = this.selectedSeenCount();

    // If trying to select and already at limit
    if (!selection.hasSeen && count >= 10) {
      this.messageService.add({
        severity: 'warn',
        summary: 'Maximum Reached',
        detail: 'You can only select up to 10 films per round',
      });
      return;
    }

    selection.hasSeen = !selection.hasSeen;
    // Trigger change detection
    this.recommendedFilms.set([...this.recommendedFilms()]);
  }

  nextRound(): void {
    if (!this.canProceedToNextRound()) {
      return;
    }

    const seenFilms = this.recommendedFilms()
      .filter((s) => s.hasSeen)
      .map((s) => s.film);

    this.isLoadingRecommendations.set(true);
    this.recommendationsError.set('');

    this.filmService
      .generateRecommendations(this.userId(), seenFilms)
      .subscribe({
        next: (recommendations) => {
          this.recommendedFilms.set(
            recommendations.map((film) => ({ film, hasSeen: false })),
          );
          this.currentRound.update((r) => r + 1);
          this.isLoadingRecommendations.set(false);
          this.messageService.add({
            severity: 'success',
            summary: 'Success',
            detail: `Round ${this.currentRound()} - ${recommendations.length} new recommendations`,
          });
        },
        error: (error) => {
          console.error('Error generating next round:', error);
          this.recommendationsError.set(
            'Failed to generate next round. Please try again.',
          );
          this.isLoadingRecommendations.set(false);
          this.messageService.add({
            severity: 'error',
            summary: 'Error',
            detail: 'Failed to generate next round',
          });
        },
      });
  }

  finish(): void {
    if (!this.canProceedToNextRound()) {
      this.messageService.add({
        severity: 'warn',
        summary: 'No Films Selected',
        detail: 'Please select at least one film you have seen',
      });
      return;
    }

    const seenFilms = this.recommendedFilms()
      .filter((s) => s.hasSeen)
      .map((s) => s.film);

    this.isLoadingRecommendations.set(true);

    // Send final batch to log them
    this.filmService
      .generateRecommendations(this.userId(), seenFilms)
      .subscribe({
        next: () => {
          this.messageService.add({
            severity: 'success',
            summary: 'Success',
            detail: 'Films added to your profile for review',
          });
          setTimeout(() => {
            this.router.navigate(['/profile', this.userId()]);
          }, 1500);
        },
        error: (error) => {
          console.error('Error finishing recommendations:', error);
          this.isLoadingRecommendations.set(false);
          this.messageService.add({
            severity: 'error',
            summary: 'Error',
            detail: 'Failed to save your selections',
          });
        },
      });
  }

  getPosterUrl(
    posterPath: string | null | undefined,
    size: TMDBPosterSize = 'w342',
  ): string {
    return getTMDBPosterUrl(posterPath, size);
  }

  goBack(): void {
    this.router.navigate(['/profile', this.userId()]);
  }
}
