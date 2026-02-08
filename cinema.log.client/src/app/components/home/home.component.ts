import { Component, CUSTOM_ELEMENTS_SCHEMA, computed, effect, ChangeDetectionStrategy, inject } from '@angular/core';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [],
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class HomeComponent {
  private auth = inject(AuthService);

  isLoggedIn = computed(() => this.auth.currentUser() !== null);

  constructor() {
    effect(() => {
      if (!this.auth.currentUser()) {
        this.auth.getCurrentUser().subscribe();
      }
    });
  }
}
