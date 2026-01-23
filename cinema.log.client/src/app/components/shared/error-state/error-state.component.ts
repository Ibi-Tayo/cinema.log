import {
  Component,
  CUSTOM_ELEMENTS_SCHEMA,
  input,
  output,
  ChangeDetectionStrategy,
} from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-error-state',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './error-state.component.html',
  styleUrl: './error-state.component.scss',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ErrorStateComponent {
  message = input<string>('An error occurred. Please try again.');
  showRetryButton = input<boolean>(false);
  retryButtonLabel = input<string>('Back to Search');

  retry = output<void>();

  onRetry(): void {
    this.retry.emit();
  }
}
