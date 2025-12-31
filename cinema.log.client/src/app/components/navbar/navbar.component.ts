import {
  Component,
  CUSTOM_ELEMENTS_SCHEMA,
  effect,
  ChangeDetectionStrategy,
} from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-navbar',
  imports: [RouterLink, RouterLinkActive],
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.scss',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class NavbarComponent {
  constructor(public auth: AuthService) {
    effect(() => {
      if (!auth.currentUser()) {
        auth.getCurrentUser().subscribe();
      }
    });
  }
}
