import { Component, OnInit, inject } from '@angular/core';
import { AvatarModule } from 'primeng/avatar';
import { DialogModule } from 'primeng/dialog';

import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { ButtonModule } from 'primeng/button';

interface Profile {
  avatarUrl: string;
  name: string;
  email: string;
  createdAt: string;
}

@Component({
  selector: 'profile',
  templateUrl: './profile.component.html',
  imports: [ AvatarModule, DialogModule, CommonModule, ButtonModule ],
  standalone: true,
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit {
  profile: Profile = {
    avatarUrl: '',
    name: '',
    email: '',
    createdAt: ''
  };
  displayProfileModal = false;

  ngOnInit() {

  }

  constructor(private router: Router) {}

  redirectTo(page: string) {
    this.router.navigate([page])
  }

  signOut() {

  }

}