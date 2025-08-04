import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../../../services/auth.service';
import { TableModule } from 'primeng/table';
import { RatingModule } from 'primeng/rating';
import { DatePipe } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ButtonModule } from 'primeng/button';
import { DialogModule } from 'primeng/dialog';
import { TextareaModule } from 'primeng/textarea';
import { FloatLabelModule } from 'primeng/floatlabel';
import { TagModule } from 'primeng/tag';
import { NotificationService } from '../../../services/notification.service';
import { AdminService } from '../../../services/admin.service';
import { UtilsService } from '../../../services/utils.service';

@Component({
  selector: 'app-admin-labs-sessions',
  imports: [TableModule, RatingModule,
    DatePipe, FormsModule,
    ButtonModule, DialogModule,
    TextareaModule, FloatLabelModule,
    TagModule],
  standalone: true,
  templateUrl: './admin-labs-sessions.component.html',
  styleUrl: './admin-labs-sessions.component.css'
})
export class AdminLabsSessionsComponent implements OnInit {

  public sessions: any[] = [];
  
  constructor(
    private authService: AuthService,
    private router: Router,
    private adminService: AdminService,
    private notificationService: NotificationService,
    public utilsService: UtilsService 
  ) {}

  ngOnInit(): void {
    this.adminService.getLabsSessions().subscribe({
      next: (sessions) => {
        this.sessions = sessions;
      },
      error: (error) => {
        console.error('Error loading sessions:', error);
        this.notificationService.showError('Failed to load sessions');
      }
    });
  }

    goToAdminLabSession(sessionId: string): void {
    this.router.navigate(['/admin/sessions', sessionId]);
  }

    isValidDate(dateString: string): boolean {
    if (!dateString || dateString.trim() === '') {
      return false;
    }
    
    const date = new Date(dateString);
    return !isNaN(date.getTime()) && date.getFullYear() > 1900;
  }
}
