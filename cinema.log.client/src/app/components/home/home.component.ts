import {
  Component,
  CUSTOM_ELEMENTS_SCHEMA,
  computed,
  effect,
  ChangeDetectionStrategy,
} from '@angular/core';
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
  isLoggedIn = computed(() => this.auth.currentUser() !== null);

  constructor(private auth: AuthService) {
    effect(() => {
      if (!this.auth.currentUser()) {
        this.auth.getCurrentUser().subscribe();
      }
    });
  }
}
