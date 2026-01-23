import {
  Component,
  CUSTOM_ELEMENTS_SCHEMA,
  input,
  output,
  signal,
  ChangeDetectionStrategy,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { Review } from '../../../services/review.service';

@Component({
  selector: 'app-review-form',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './review-form.component.html',
  styleUrl: './review-form.component.scss',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ReviewFormComponent {
  // Inputs
  hasRating = input.required<boolean>();
  existingReview = input<Review | null>(null);
  isSubmitting = input.required<boolean>();
  submitSuccess = input.required<boolean>();
  submitError = input<string>('');

  // Outputs
  reviewSubmit = output<{ rating: number; content: string }>();
  reviewUpdate = output<string>();

  // Internal state
  selectedRating = signal(0);
  reviewText = signal('');

  /**
   * Select a star rating
   */
  selectRating(rating: number): void {
    this.selectedRating.set(rating);
  }

  /**
   * Update review text
   */
  updateReviewText(value: string): void {
    this.reviewText.set(value);
  }

  /**
   * Submit a new review
   */
  onSubmitReview(): void {
    if (this.selectedRating() === 0 || !this.reviewText().trim()) {
      return;
    }
    this.reviewSubmit.emit({
      rating: this.selectedRating(),
      content: this.reviewText().trim(),
    });
  }

  /**
   * Update existing review
   */
  onUpdateReview(): void {
    if (!this.reviewText().trim()) {
      return;
    }
    this.reviewUpdate.emit(this.reviewText().trim());
  }

  /**
   * Get array of star numbers for rendering
   */
  getStars(): number[] {
    return [1, 2, 3, 4, 5];
  }

  /**
   * Check if submit button should be disabled
   */
  isSubmitDisabled(): boolean {
    if (this.isSubmitting()) return true;
    if (this.hasRating()) {
      return !this.reviewText().trim();
    }
    return this.selectedRating() === 0 || !this.reviewText().trim();
  }
}
