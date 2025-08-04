import { Component, OnInit } from '@angular/core';
import { LabService } from '../../services/lab.service';
import { TableModule } from 'primeng/table';
import { RatingModule } from 'primeng/rating';
import { GetAllLabsByUserAndStatusResponse } from '../../types/lab';
import { DatePipe } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ButtonModule } from 'primeng/button';
import { Router } from '@angular/router';
import { DialogModule } from 'primeng/dialog';
import { TextareaModule } from 'primeng/textarea';
import { FloatLabelModule } from 'primeng/floatlabel';
import { TagModule } from 'primeng/tag';
import { NotificationService } from '../../services/notification.service';
import { UtilsService } from '../../services/utils.service';

@Component({
  selector: 'app-labs-sessions',
  standalone: true,
  templateUrl: './labs-sessions.component.html',
  styleUrls: ['./labs-sessions.component.css'],
  imports: [TableModule, RatingModule, FormsModule, ButtonModule, DialogModule, TextareaModule, FloatLabelModule, TagModule, DatePipe]
})
export class LabsSessionsComponent implements OnInit {
  sessions: GetAllLabsByUserAndStatusResponse[] = [];
  assessmentDialog: boolean = false;


  assessmentSessionId: string = '';
  assessmentRating: number = 0;
  assessmentLabSlug: string = '';
  assessmentFeedback: string = '';

  constructor(private labService: LabService,
    private router: Router,
    private notificationService: NotificationService,
    public utilsService: UtilsService) { }

  ngOnInit(): void {
    this.loadSessions();
  }

  loadSessions(): void {
    this.labService.getAllLabsByUserAndStatus('finished').subscribe({
      next: (sessions) => {
        this.sessions = sessions;
      },
      error: (error) => {
        console.error('Error loading sessions:', error);
        this.notificationService.showError('failed to load sessions');
      }
    });
  }

  onRate(event: any, sessionId: string, labSlug: string): void {
    this.assessmentDialog = true;
    this.assessmentRating = event.value;
    this.assessmentSessionId = sessionId;
    this.assessmentLabSlug = labSlug;
  }

  goToLabSession(sessionId: string): void {
    this.router.navigate(['/sessions', sessionId]);
  }

  sendFeedback(sessionId: string): void {
    this.labService.sendFeedback(sessionId, this.assessmentRating, this.assessmentFeedback).subscribe({
      next: () => {
        this.assessmentDialog = false;
        this.loadSessions();
        this.notificationService.showSuccess('feedback sent successfully');
      },
      error: (error) => {
        console.error('Error sending feedback:', error);
        this.notificationService.showError('failed to load feedback');
      }
    });
  }

    isValidDate(dateString: string): boolean {
    if (!dateString || dateString.trim() === '') {
      return false;
    }
    
    const date = new Date(dateString);
    return !isNaN(date.getTime()) && date.getFullYear() > 1900;
  }

}

