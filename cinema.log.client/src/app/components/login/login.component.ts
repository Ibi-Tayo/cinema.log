import {
  Component,
  ChangeDetectionStrategy,
  signal,
  CUSTOM_ELEMENTS_SCHEMA,
  OnInit,
} from '@angular/core';
import { ButtonModule } from 'primeng/button';
import { AuthService } from '../../services/auth.service';
import { Router, ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-login',
  imports: [ButtonModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LoginComponent implements OnInit {
  errorMessage = signal<string>('');
  branchName = signal<string>(
    import.meta.env.NG_APP_BRANCH_NAME || 'production',
  );

  constructor(
    public auth: AuthService,
    private router: Router,
    private route: ActivatedRoute,
  ) {}

  ngOnInit() {
    // Check for error query parameter
    this.route.queryParams.subscribe((params) => {
      const error = params['error'];
      if (error === 'github_auth_failed') {
        this.errorMessage.set(
          'GitHub authentication failed. Please try again.',
        );
      } else if (error === 'google_auth_failed') {
        this.errorMessage.set(
          'Google authentication failed. Please try again.',
        );
      }
    });
  }

  login() {
    this.auth.login();
  }

  googleLogin() {
    this.auth.googleLogin();
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
        this.errorMessage.set('Dev login failed. Please try again.');
      },
    });
  }

  devGoogleLogin() {
    this.auth.devGoogleLogin().subscribe({
      next: () => {
        // redirect to user profile
        this.auth.getCurrentUser().subscribe((user) => {
          this.router.navigate(['/profile', user.id]).then(() => {
            window.location.reload();
          });
        });
      },
      error: (err) => {
        console.error('Dev Google login error:', err);
        this.errorMessage.set('Dev Google login failed. Please try again.');
      },
    });
  }
}
