import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DialogModule } from 'primeng/dialog';
import { ButtonModule } from 'primeng/button';
import { AvatarModule } from 'primeng/avatar';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { NotificationService } from '../../services/notification.service';

@Component({
  selector: 'profile-modal-component',
  templateUrl: './profile-modal.component.html',
  styleUrls: ['./profile-modal.component.css'],
  standalone: true,
  imports: [CommonModule, DialogModule, ButtonModule, AvatarModule]
})
export class ProfileModalComponent implements OnInit {
  @Input() visible = false;
  @Output() visibleChange = new EventEmitter<boolean>();
  user: any;

  constructor(
    private authService: AuthService,
    private notificationService: NotificationService,
    private router: Router
  ) {}

  async ngOnInit() {
    console.log('ProfileModalComponent initialized');
    this.user = await this.authService.getUser();
    console.log('User data:', this.user);
  }

  hideDialog() {
    this.visible = false;
    this.visibleChange.emit(false);
  }
  

  redirectTo(page: string) {

    this.hideDialog();
    this.router.navigate([page]);
  }

  async signOut() {
    try {
      await this.authService.logout();
      this.notificationService.showInfo('You have been signed out successfully.', 'Hey! See you soon! :D');
      this.hideDialog();
      this.router.navigate(['/login']);
    } catch (error) {
      console.error('Sign out failed:', error);
    }
  }
}
