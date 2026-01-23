import { Injectable, signal, computed } from '@angular/core';

export type ComparisonResult = 'better' | 'worse' | 'same';

@Injectable({
  providedIn: 'root',
})
export class ComparisonStateService {
  // Comparison mode (bulk or sequential)
  private _mode = signal<'bulk' | 'sequential'>('bulk');
  mode = this._mode.asReadonly();

  // Bulk selections: Map<filmId, result>
  private _selections = signal<Map<string, ComparisonResult>>(new Map());
  selections = this._selections.asReadonly();

  // Computed: number of selections made
  selectionCount = computed(() => this._selections().size);

  /**
   * Load comparison mode preference from localStorage
   */
  loadModePreference(): void {
    const savedMode = localStorage.getItem('comparisonMode');
    if (savedMode === 'bulk' || savedMode === 'sequential') {
      this._mode.set(savedMode);
    }
  }

  /**
   * Save comparison mode preference to localStorage
   */
  saveModePreference(mode: 'bulk' | 'sequential'): void {
    this._mode.set(mode);
    localStorage.setItem('comparisonMode', mode);
  }

  /**
   * Toggle between bulk and sequential modes
   */
  toggleMode(): void {
    const newMode = this._mode() === 'bulk' ? 'sequential' : 'bulk';
    this.saveModePreference(newMode);
  }

  /**
   * Set a comparison selection for a specific film
   */
  setSelection(filmId: string, result: ComparisonResult): void {
    const selections = new Map(this._selections());
    selections.set(filmId, result);
    this._selections.set(selections);
  }

  /**
   * Remove a comparison selection for a specific film
   */
  removeSelection(filmId: string): void {
    const selections = new Map(this._selections());
    selections.delete(filmId);
    this._selections.set(selections);
  }

  /**
   * Get the selection for a specific film
   */
  getSelection(filmId: string): ComparisonResult | undefined {
    return this._selections().get(filmId);
  }

  /**
   * Clear all selections
   */
  resetSelections(): void {
    this._selections.set(new Map());
  }

  /**
   * Get all selections as an array
   */
  getAllSelections(): Array<{ filmId: string; result: ComparisonResult }> {
    return Array.from(this._selections().entries()).map(([filmId, result]) => ({
      filmId,
      result,
    }));
  }
}
