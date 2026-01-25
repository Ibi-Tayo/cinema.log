import {
  Component,
  CUSTOM_ELEMENTS_SCHEMA,
  effect,
  ChangeDetectionStrategy,
  ViewChild,
  OnInit,
} from '@angular/core';
import { RouterLink, RouterLinkActive, Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { MenuModule } from 'primeng/menu';
import { MenuItem } from 'primeng/api';
import { Menu } from 'primeng/menu';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [RouterLink, RouterLinkActive, MenuModule],
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.scss',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class NavbarComponent implements OnInit {
  @ViewChild('profileMenu') profileMenu!: Menu;
  menuItems: MenuItem[] = [];

  constructor(
    public auth: AuthService,
    private router: Router,
  ) {
    effect(() => {
      if (!auth.currentUser()) {
        auth.getCurrentUser().subscribe();
      }
    });
  }

  ngOnInit() {
    this.menuItems = [
      {
        label: 'Profile',
        icon: 'pi pi-user',
        command: () =>
          this.router.navigate(['/profile', this.auth.currentUser()?.id]),
      },
      {
        label: 'Film Graph',
        icon: 'pi pi-chart-bar',
        command: () => this.router.navigate(['/film-graph']),
      },
      {
        label: 'Sign Out',
        icon: 'pi pi-sign-out',
        command: () => this.onSignOut(),
      },
    ];
  }

  toggleProfileMenu(event: Event) {
    this.profileMenu.toggle(event);
  }

  onSignOut() {
    this.auth.logout().subscribe(() => {
      this.router.navigate(['/login']);
    });
  }
}
