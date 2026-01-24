import { TestBed } from '@angular/core/testing';
import {
  ComparisonStateService,
  ComparisonResult,
} from './comparison-state.service';

describe('ComparisonStateService', () => {
  let service: ComparisonStateService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ComparisonStateService);
    // Clear localStorage before each test
    localStorage.clear();
  });

  afterEach(() => {
    // Clean up localStorage after each test
    localStorage.clear();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe('Mode Management', () => {
    it('should initialize with bulk mode by default', () => {
      expect(service.mode()).toBe('bulk');
    });

    it('should load mode preference from localStorage', () => {
      localStorage.setItem('comparisonMode', 'sequential');
      service.loadModePreference();
      expect(service.mode()).toBe('sequential');
    });

    it('should save mode preference to localStorage', () => {
      service.saveModePreference('sequential');
      expect(localStorage.getItem('comparisonMode')).toBe('sequential');
      expect(service.mode()).toBe('sequential');
    });

    it('should toggle mode between bulk and sequential', () => {
      service.saveModePreference('bulk');
      expect(service.mode()).toBe('bulk');

      service.toggleMode();
      expect(service.mode()).toBe('sequential');
      expect(localStorage.getItem('comparisonMode')).toBe('sequential');

      service.toggleMode();
      expect(service.mode()).toBe('bulk');
      expect(localStorage.getItem('comparisonMode')).toBe('bulk');
    });

    it('should ignore invalid mode in localStorage', () => {
      localStorage.setItem('comparisonMode', 'invalid');
      service.loadModePreference();
      expect(service.mode()).toBe('bulk'); // Should remain default
    });
  });

  describe('Selection Management', () => {
    it('should start with empty selections', () => {
      expect(service.selectionCount()).toBe(0);
      expect(service.getAllSelections()).toEqual([]);
    });

    it('should set a selection', () => {
      service.setSelection('film1', 'better');
      expect(service.getSelection('film1')).toBe('better');
      expect(service.selectionCount()).toBe(1);
    });

    it('should update an existing selection', () => {
      service.setSelection('film1', 'better');
      service.setSelection('film1', 'worse');
      expect(service.getSelection('film1')).toBe('worse');
      expect(service.selectionCount()).toBe(1);
    });

    it('should handle multiple selections', () => {
      service.setSelection('film1', 'better');
      service.setSelection('film2', 'worse');
      service.setSelection('film3', 'same');

      expect(service.selectionCount()).toBe(3);
      expect(service.getSelection('film1')).toBe('better');
      expect(service.getSelection('film2')).toBe('worse');
      expect(service.getSelection('film3')).toBe('same');
    });

    it('should remove a selection', () => {
      service.setSelection('film1', 'better');
      service.setSelection('film2', 'worse');
      expect(service.selectionCount()).toBe(2);

      service.removeSelection('film1');
      expect(service.selectionCount()).toBe(1);
      expect(service.getSelection('film1')).toBeUndefined();
      expect(service.getSelection('film2')).toBe('worse');
    });

    it('should return all selections as an array', () => {
      service.setSelection('film1', 'better');
      service.setSelection('film2', 'worse');
      service.setSelection('film3', 'same');

      const allSelections = service.getAllSelections();
      expect(allSelections.length).toBe(3);
      expect(allSelections).toContain({ filmId: 'film1', result: 'better' });
      expect(allSelections).toContain({ filmId: 'film2', result: 'worse' });
      expect(allSelections).toContain({ filmId: 'film3', result: 'same' });
    });

    it('should reset all selections', () => {
      service.setSelection('film1', 'better');
      service.setSelection('film2', 'worse');
      expect(service.selectionCount()).toBe(2);

      service.resetSelections();
      expect(service.selectionCount()).toBe(0);
      expect(service.getAllSelections()).toEqual([]);
    });

    it('should return undefined for non-existent selection', () => {
      expect(service.getSelection('nonexistent')).toBeUndefined();
    });
  });

  describe('Computed Properties', () => {
    it('should reactively update selection count', () => {
      expect(service.selectionCount()).toBe(0);

      service.setSelection('film1', 'better');
      expect(service.selectionCount()).toBe(1);

      service.setSelection('film2', 'worse');
      expect(service.selectionCount()).toBe(2);

      service.removeSelection('film1');
      expect(service.selectionCount()).toBe(1);

      service.resetSelections();
      expect(service.selectionCount()).toBe(0);
    });
  });
});
