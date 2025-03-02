import { Component } from '@angular/core';
import { ButtonModule } from 'primeng/button';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-login',
  imports: [ButtonModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss',
  schemas: [],
})
export class LoginComponent {
  constructor(private auth: AuthService) {}

  login() {
    this.auth.login();
  }
}
