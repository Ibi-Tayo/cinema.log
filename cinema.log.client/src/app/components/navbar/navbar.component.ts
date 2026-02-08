import { Component, CUSTOM_ELEMENTS_SCHEMA, effect, ChangeDetectionStrategy, ViewChild, OnInit, inject } from '@angular/core';
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
  auth = inject(AuthService);
  private router = inject(Router);

  @ViewChild('profileMenu') profileMenu!: Menu;
  menuItems: MenuItem[] = [];

  constructor() {
    const auth = this.auth;

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
        label: 'Find Films To Review',
        icon: 'pi pi-search',
        command: () =>
          this.router.navigate([
            '/recommendations',
            this.auth.currentUser()?.id,
          ]),
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

  isProfileAreaActive(): boolean {
    return (
      this.router.url.includes('/profile') ||
      this.router.url.includes('/recommendations') ||
      this.router.url.includes('/film-graph')
    );
  }

  onSignOut() {
    this.auth.logout().subscribe(() => {
      this.router.navigate(['/login']);
    });
  }
}
