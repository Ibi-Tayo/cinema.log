import { Component, ChangeDetectionStrategy } from '@angular/core';
import { ButtonModule } from 'primeng/button';
import { AuthService } from '../../services/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-login',
  imports: [ButtonModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss',
  schemas: [],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LoginComponent {
  constructor(public auth: AuthService, private router: Router) {}

  login() {
    this.auth.login();
  }

  devLogin() {
    this.auth.devLogin().subscribe({
      next: () => {
        // redirect to user profile
        this.auth.getCurrentUser().subscribe((user) => {
          this.router.navigate(['/profile', user.id]).then(() => {
            window.location.reload();
          });
        });
      },
      error: (err) => {
        console.error('Dev login error:', err);
        alert('Dev login failed. Please try again.');
      },
    });
  }
}
