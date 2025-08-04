import { Component, OnInit } from '@angular/core';
import { LabService } from '../../services/lab.service';
import { ActivatedRoute } from '@angular/router';
import { GetLabResultResponse, LabFinishResultCriterion } from '../../types/lab';
import { NotificationService } from '../../services/notification.service';
import { DatePipe } from '@angular/common';
import { RatingModule } from 'primeng/rating';
import { TagModule } from 'primeng/tag';
import { CardModule } from 'primeng/card';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { AvatarModule } from 'primeng/avatar';
import { DividerModule } from 'primeng/divider';
import { TabsModule } from 'primeng/tabs';
import { DialogModule } from 'primeng/dialog';
import { FloatLabelModule } from 'primeng/floatlabel';
import { ButtonModule } from 'primeng/button';
import { ScrollPanelModule } from 'primeng/scrollpanel';
import { FilesDiffViewerComponent } from '../../components/files-diff-viewer/files-diff-viewer.component';
import { TableModule } from 'primeng/table';
import { UtilsService } from '../../services/utils.service';

@Component({
  selector: 'app-lab-session',
  standalone: true,
  imports: [RatingModule, TagModule,
     CardModule, FormsModule,
     CommonModule, AvatarModule,
     DividerModule, TabsModule,
     DialogModule, FloatLabelModule,
     ButtonModule, ScrollPanelModule,
     FilesDiffViewerComponent, TableModule,
     DatePipe],
  templateUrl: './lab-session.component.html',
  styleUrl: './lab-session.component.css'
})
export class LabSessionComponent implements OnInit {
  labResult: GetLabResultResponse = {} as GetLabResultResponse;
  criteriaResult: LabFinishResultCriterion[] = [] as LabFinishResultCriterion[];
  id: string = '';
  grade: string = '';
  tabs: any[] = [];


  avaliacaoDialog: boolean = false;
  avaliacaoRating: number = 0;
  avaliacaoFeedback: string = '';
  avaliacaoLabSlug: string = '';
  avaliacaoSessionId: string = '';
  avaliacaoFilesDiff: string = '';

  detailsDialog: boolean = false;
  detailsDialogData: any = {};
   
  constructor(
    private route: ActivatedRoute,
    private labService: LabService,
    private notificationService: NotificationService,
    public utilsService: UtilsService
  ) {}

  ngOnInit() {
    this.id = this.route.snapshot.params['id'];
    this.getLabResult();
  }

  getLabResult() {
    this.labService.getLabResult(this.id).subscribe({
      next: (result) => {
        this.labResult = result;
        this.grade = this.generateGrade(this.labResult.finishResult.totalScore);
        this.avaliacaoFilesDiff = result.finishResult.filesDiff || '';
        this.criteriaResult = result.finishResult?.criteriaResult || [];
      },
      error: (error) => {
        console.error(error);
        this.notificationService.showError('Error loading lab result');
      }
    });
  }

  generateGrade(score: number) {
    if (score >= 97) {
      return 'A+';
    } else if (score >= 93) {
      return 'A';
    } else if (score >= 90) {
      return 'A-';
    } else if (score >= 87) {
      return 'B+';
    } else if (score >= 83) {
      return 'B';
    } else if (score >= 80) {
      return 'B-';
    } else if (score >= 77) {
      return 'C+';
    } else if (score >= 73) {
      return 'C';
    } else if (score >= 70) {
      return 'C-';
    } else if (score >= 67) {
      return 'D+';
    } else {
      return 'D';
    }
  }

  onRate(event: any) {
    console.log(event);
    this.avaliacaoRating = event.value;
    this.avaliacaoLabSlug = this.labResult?.labSlug || '';
    this.avaliacaoSessionId = this.labResult?.namespace || '';
    this.avaliacaoDialog = true;
  }

  sendFeedback(sessionId: string) {
    this.labService.sendFeedback(sessionId, this.avaliacaoRating, this.avaliacaoFeedback).subscribe({
      next: () => {
        this.avaliacaoDialog = false;
        this.getLabResult();
        this.notificationService.showSuccess('Feedback sent successfully');
      },
      error: (error) => {
        console.error(error);
        this.notificationService.showError('Error sending feedback');
      }
    });
  }

  onCloseFeedbackDialog() {
    this.getLabResult();
  }

  formatMessage(msg: any): string {
    try {
      if (typeof msg === 'string') {
        const parsed = JSON.parse(msg);
        return JSON.stringify(parsed, null, 2);
      }
      return JSON.stringify(msg, null, 2);
    } catch {
      return msg?.toString?.() ?? '';
    }
  }

  onSeeCriterionDetails(criterion: LabFinishResultCriterion) {
    this.detailsDialogData = {
      name: criterion.name,
      score: criterion.score,
      weight: criterion.weight,
      required: criterion.required,
      message: criterion.message,
      status: criterion.status,
      rawOutput: this.formatMessage(criterion.rawOutput)
    };
    const html = `
      <!DOCTYPE html>
      <html>
        <head>
          <title>Raw Output</title>
          <meta name="color-scheme" content="light dark">
        </head>
        <body>
          <pre style="word-wrap: break-word; white-space: pre-wrap;">${this.escapeHtml(this.detailsDialogData.rawOutput)}</pre>
        </body>
      </html>
    `;
    const newWindow = window.open('', '_blank');
    if (newWindow) {
      newWindow.document.write(html);
      newWindow.document.close();
    }
  }


  private escapeHtml(unsafe: string): string {
  return unsafe
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
    }


  isValidDate(dateString: string): boolean {
    if (!dateString || dateString.trim() === '') {
      return false;
    }
    
    const date = new Date(dateString);
    return !isNaN(date.getTime()) && date.getFullYear() > 1900;
  }
}
