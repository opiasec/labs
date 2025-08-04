import { Component, OnInit, inject } from '@angular/core';
import { ToolbarModule } from 'primeng/toolbar';
import { ButtonModule } from 'primeng/button';
import { IconFieldModule } from 'primeng/iconfield';
import { InputIconModule } from 'primeng/inputicon';
import { SplitButtonModule } from 'primeng/splitbutton';
import { AvatarModule } from 'primeng/avatar';
import { AvatarGroupModule } from 'primeng/avatargroup';
import { ProfileModalComponent } from '../profile-modal/profile-modal.component';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-topbar',
  imports: [
    ToolbarModule,
    ButtonModule,
    IconFieldModule,
    InputIconModule,
    ProfileModalComponent,
    SplitButtonModule,
    AvatarGroupModule,
    AvatarModule],
  templateUrl: './topbar.component.html',
  styleUrl: './topbar.component.css'
})
export class TopbarComponent implements OnInit {

  constructor(
    private authService: AuthService,
  ) {}

  avatarUrl = '';
  name = '';
  displayProfileModal = false;
  ngOnInit() {
    this.authService.getUser().then(user => {
      this.avatarUrl = user.picture || '';
      this.name = user.name || '';
    });
  }


  toggleProfileModal() {
    this.displayProfileModal = !this.displayProfileModal
  }
}

