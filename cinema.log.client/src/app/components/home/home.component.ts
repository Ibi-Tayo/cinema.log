import { Component, CUSTOM_ELEMENTS_SCHEMA, OnInit } from '@angular/core';
import { AuthService } from '../../services/auth.service';


@Component({
  selector: 'app-home',
  standalone: true,
  imports: [],
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
})
export class HomeComponent implements OnInit {

  isLoggedIn = false;

  constructor(private auth: AuthService){}

  ngOnInit(): void {
    this.auth.getCurrentUser().subscribe({
      next: (_) => {
        this.isLoggedIn = true;
      },
      error: (_) => {
        this.isLoggedIn = false;
      }
    });
  }
}
