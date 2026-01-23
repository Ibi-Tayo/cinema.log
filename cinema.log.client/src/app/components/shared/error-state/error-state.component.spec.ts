import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ErrorStateComponent } from './error-state.component';

describe('ErrorStateComponent', () => {
  let component: ErrorStateComponent;
  let fixture: ComponentFixture<ErrorStateComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ErrorStateComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(ErrorStateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should display default error message when no input provided', () => {
    const compiled = fixture.nativeElement as HTMLElement;
    const messageElement = compiled.querySelector('p');
    expect(messageElement?.textContent).toContain(
      'An error occurred. Please try again.',
    );
  });

  it('should display custom error message when input provided', () => {
    fixture = TestBed.createComponent(ErrorStateComponent);
    fixture.componentRef.setInput(
      'message',
      'Failed to load film details. Please try again.',
    );
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const messageElement = compiled.querySelector('p');
    expect(messageElement?.textContent).toContain(
      'Failed to load film details. Please try again.',
    );
  });

  it('should render lord-icon element', () => {
    const compiled = fixture.nativeElement as HTMLElement;
    const lordIcon = compiled.querySelector('lord-icon');
    expect(lordIcon).toBeTruthy();
    expect(lordIcon?.getAttribute('trigger')).toBe('loop');
  });

  it('should have error-state class on container', () => {
    const compiled = fixture.nativeElement as HTMLElement;
    const container = compiled.querySelector('.error-state');
    expect(container).toBeTruthy();
  });

  it('should not show retry button by default', () => {
    const compiled = fixture.nativeElement as HTMLElement;
    const retryButton = compiled.querySelector('.retry-button');
    expect(retryButton).toBeNull();
  });

  it('should show retry button when showRetryButton is true', () => {
    fixture = TestBed.createComponent(ErrorStateComponent);
    fixture.componentRef.setInput('showRetryButton', true);
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const retryButton = compiled.querySelector('.retry-button');
    expect(retryButton).toBeTruthy();
  });

  it('should display custom retry button label', () => {
    fixture = TestBed.createComponent(ErrorStateComponent);
    fixture.componentRef.setInput('showRetryButton', true);
    fixture.componentRef.setInput('retryButtonLabel', 'Try Again');
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const retryButton = compiled.querySelector('.retry-button');
    expect(retryButton?.textContent).toContain('Try Again');
  });

  it('should emit retry event when button is clicked', () => {
    fixture = TestBed.createComponent(ErrorStateComponent);
    fixture.componentRef.setInput('showRetryButton', true);

    let retryEmitted = false;
    fixture.componentInstance.retry.subscribe(() => {
      retryEmitted = true;
    });

    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const retryButton = compiled.querySelector(
      '.retry-button',
    ) as HTMLButtonElement;
    retryButton.click();

    expect(retryEmitted).toBe(true);
  });
});
