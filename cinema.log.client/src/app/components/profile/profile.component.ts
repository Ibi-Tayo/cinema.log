import { Component } from '@angular/core';

@Component({
  selector: 'app-profile',
  imports: [],
  templateUrl: './profile.component.html',
  styleUrl: './profile.component.scss',
})
export class ProfileComponent {
  // probably want to inject a user service
  // ngoninit an api call to userService.getUser
  // grab the data and wack it into the html
  // make sure you dont keep calling the api if you already have the data
}
