import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ReviewFormComponent } from './review-form.component';
import { Review } from '../../../services/review.service';

describe('ReviewFormComponent', () => {
  let component: ReviewFormComponent;
  let fixture: ComponentFixture<ReviewFormComponent>;

  const mockReview: Review = {
    id: 'review1',
    userId: 'user1',
    filmId: 'film1',
    title: 'Great movie!',
    rating: 5,
    date: '2024-01-23T00:00:00Z',
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ReviewFormComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(ReviewFormComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    fixture.componentRef.setInput('hasRating', false);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.componentRef.setInput('submitSuccess', false);
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('should show "Review" heading for new reviews', () => {
    fixture.componentRef.setInput('hasRating', false);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.componentRef.setInput('submitSuccess', false);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const heading = compiled.querySelector('h3');
    expect(heading?.textContent).toContain('Review');
    expect(heading?.textContent).not.toContain('Update');
  });

  it('should show "Update Review" heading for existing reviews', () => {
    fixture.componentRef.setInput('hasRating', true);
    fixture.componentRef.setInput('existingReview', mockReview);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.componentRef.setInput('submitSuccess', false);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const heading = compiled.querySelector('h3');
    expect(heading?.textContent).toContain('Update Review');
  });

  it('should show rating section for new reviews', () => {
    fixture.componentRef.setInput('hasRating', false);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.componentRef.setInput('submitSuccess', false);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const ratingSection = compiled.querySelector('.rating-section');
    expect(ratingSection).toBeTruthy();
  });

  it('should not show rating section for existing reviews', () => {
    fixture.componentRef.setInput('hasRating', true);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.componentRef.setInput('submitSuccess', false);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const ratingSection = compiled.querySelector('.rating-section');
    expect(ratingSection).toBeNull();
  });

  it('should render 5 stars', () => {
    fixture.componentRef.setInput('hasRating', false);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.componentRef.setInput('submitSuccess', false);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const stars = compiled.querySelectorAll('.star');
    expect(stars.length).toBe(5);
  });

  it('should update selected rating when star is clicked', () => {
    fixture.componentRef.setInput('hasRating', false);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.componentRef.setInput('submitSuccess', false);
    fixture.detectChanges();

    expect(component.selectedRating()).toBe(0);

    const compiled = fixture.nativeElement as HTMLElement;
    const stars = compiled.querySelectorAll('.star');
    (stars[2] as HTMLElement).click(); // Click 3rd star
    fixture.detectChanges();

    expect(component.selectedRating()).toBe(3);
  });

  it('should update review text when typing', () => {
    fixture.componentRef.setInput('hasRating', false);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.componentRef.setInput('submitSuccess', false);
    fixture.detectChanges();

    const textarea = fixture.nativeElement.querySelector(
      'textarea',
    ) as HTMLTextAreaElement;
    textarea.value = 'This is my review';
    textarea.dispatchEvent(new Event('input'));
    fixture.detectChanges();

    expect(component.reviewText()).toBe('This is my review');
  });

  it('should show submit button for new reviews', () => {
    fixture.componentRef.setInput('hasRating', false);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.componentRef.setInput('submitSuccess', false);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const submitButton = compiled.querySelector('.submit-button');
    expect(submitButton?.textContent).toContain('Submit Review');
  });

  it('should show update button for existing reviews', () => {
    fixture.componentRef.setInput('hasRating', true);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.componentRef.setInput('submitSuccess', false);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const submitButton = compiled.querySelector('.submit-button');
    expect(submitButton?.textContent).toContain('Update Review');
  });

  it('should disable submit button when submitting', () => {
    fixture.componentRef.setInput('hasRating', false);
    fixture.componentRef.setInput('isSubmitting', true);
    fixture.componentRef.setInput('submitSuccess', false);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const submitButton = compiled.querySelector(
      '.submit-button',
    ) as HTMLButtonElement;
    expect(submitButton.disabled).toBe(true);
  });

  it('should disable submit button when rating is 0 for new reviews', () => {
    fixture.componentRef.setInput('hasRating', false);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.componentRef.setInput('submitSuccess', false);
    fixture.detectChanges();

    component.reviewText.set('Some review text');
    component.selectedRating.set(0);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const submitButton = compiled.querySelector(
      '.submit-button',
    ) as HTMLButtonElement;
    expect(submitButton.disabled).toBe(true);
  });

  it('should emit reviewSubmit event when submitting new review', () => {
    fixture.componentRef.setInput('hasRating', false);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.componentRef.setInput('submitSuccess', false);

    let emittedData: { rating: number; content: string } | undefined;
    fixture.componentInstance.reviewSubmit.subscribe((data) => {
      emittedData = data;
    });

    component.selectedRating.set(5);
    component.reviewText.set('Great film!');
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const submitButton = compiled.querySelector(
      '.submit-button',
    ) as HTMLButtonElement;
    submitButton.click();

    expect(emittedData).toBeDefined();
    expect(emittedData!.rating).toBe(5);
    expect(emittedData!.content).toBe('Great film!');
  });

  it('should emit reviewUpdate event when updating existing review', () => {
    fixture.componentRef.setInput('hasRating', true);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.componentRef.setInput('submitSuccess', false);

    let emittedContent: string | undefined;
    fixture.componentInstance.reviewUpdate.subscribe((content) => {
      emittedContent = content;
    });

    component.reviewText.set('Updated review text');
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const submitButton = compiled.querySelector(
      '.submit-button',
    ) as HTMLButtonElement;
    submitButton.click();

    expect(emittedContent).toBeDefined();
    expect(emittedContent).toBe('Updated review text');
  });

  it('should show success message when submitSuccess is true', () => {
    fixture.componentRef.setInput('hasRating', false);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.componentRef.setInput('submitSuccess', true);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const successMessage = compiled.querySelector('.success-message');
    expect(successMessage).toBeTruthy();
    expect(successMessage?.textContent).toContain(
      'Review submitted successfully',
    );
  });

  it('should hide form when submitSuccess is true', () => {
    fixture.componentRef.setInput('hasRating', false);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.componentRef.setInput('submitSuccess', true);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const form = compiled.querySelector('.review-form');
    expect(form).toBeNull();
  });

  it('should show error message when submitError is provided', () => {
    fixture.componentRef.setInput('hasRating', false);
    fixture.componentRef.setInput('isSubmitting', false);
    fixture.componentRef.setInput('submitSuccess', false);
    fixture.componentRef.setInput('submitError', 'Failed to submit review');
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const errorMessage = compiled.querySelector('.error-message');
    expect(errorMessage).toBeTruthy();
    expect(errorMessage?.textContent).toContain('Failed to submit review');
  });
});
