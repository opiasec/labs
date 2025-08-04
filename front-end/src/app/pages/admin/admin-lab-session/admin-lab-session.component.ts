import { Component, OnInit } from '@angular/core';
import { LabService } from '../../../services/lab.service';
import { ActivatedRoute, Router } from '@angular/router';
import { GetLabResultResponse, LabFinishResultCriterion } from '../../../types/lab';
import { NotificationService } from '../../../services/notification.service';
import { RatingModule } from 'primeng/rating';
import { TagModule } from 'primeng/tag';
import { CardModule } from 'primeng/card';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';
import { AvatarModule } from 'primeng/avatar';
import { DividerModule } from 'primeng/divider';
import { TabsModule } from 'primeng/tabs';
import { DialogModule } from 'primeng/dialog';
import { FloatLabelModule } from 'primeng/floatlabel';
import { ButtonModule } from 'primeng/button';
import { ScrollPanelModule } from 'primeng/scrollpanel';
import { FilesDiffViewerComponent } from '../../../components/files-diff-viewer/files-diff-viewer.component';
import { TableModule } from 'primeng/table';
import { AdminService } from '../../../services/admin.service';
import { AuthService } from '../../../services/auth.service';
import { MessageModule } from 'primeng/message';
import { Select } from 'primeng/select';
import { EditorModule } from 'primeng/editor';
import { UtilsService } from '../../../services/utils.service';


@Component({
  selector: 'app-admin-lab-session',
  imports: [CommonModule, FormsModule,
    CardModule, AvatarModule,
    DividerModule, TabsModule,
    DialogModule, FloatLabelModule,
    ButtonModule, ScrollPanelModule,
    ReactiveFormsModule, Select, MessageModule,
    TableModule, RatingModule, TagModule,
    EditorModule, 
    FilesDiffViewerComponent],
  templateUrl: './admin-lab-session.component.html',
  standalone: true,
  styleUrl: './admin-lab-session.component.css'
})
export class AdminLabSessionComponent implements OnInit {

  constructor(private adminService: AdminService,
              private authService: AuthService,
              private route: ActivatedRoute,
              private notificationService: NotificationService,
              private router: Router,
              private formBuilder: FormBuilder,
              public utilsService: UtilsService) {
    this.statusChangeForm = this.formBuilder.group({
      status: ['', Validators.required],
      comment: ['']
    });
  }

  namespaceId: string = '';
  labResult: any = {};
  criteriaResult: any[] = [];
  detailsDialogData: any = {};
  avaliacaoDialog: boolean = false;
  avaliacaoRating: number = 0;
  avaliacaoFeedback: string = '';
  avaliacaoLabSlug: string = '';
  avaliacaoSessionId: string = '';
  avaliacaoFilesDiff: string = '';

  statusChangeForm: FormGroup;
  statusComment: string = '';
  statusChangeFormSubmitted = false;
  statusOptions: { name: string, value: string }[] = [
  ];

  labSessionLogs: any[] = [];

  ngOnInit() {
    this.authService.isAdmin().subscribe(isAdmin => {
      if (!isAdmin) {
        this.router.navigate(['/']);
      }
    });
    this.namespaceId = this.route.snapshot.params['id'];
    if (this.namespaceId) {
      this.loadLabSession(this.namespaceId);
    } else {
      this.notificationService.showError('Invalid lab session ID');
      this.router.navigate(['/']);
    }
  }

  loadPossibleStatuses(statusFrom: string) {
    this.adminService.getPossiblesStatus(statusFrom).subscribe({
      next: (data) => {
        this.statusOptions = data.map((status: any) => ({
          name: status.name,
          value: status.value
        }));
      },
      error: (error) => {
        console.error('Error loading possible statuses:', error);
        this.notificationService.showError('Failed to load possible statuses');
      }
    });
  }

  loadLabSession(namespaceId: string) {
    console.log('Loading lab session for namespace:', namespaceId);
    this.adminService.getLabSession(namespaceId).subscribe({
      next: (data) => {
        this.labResult = data;
        this.criteriaResult = data.finishResult?.criteriaResult || [];
        this.avaliacaoFilesDiff = data.finishResult?.filesDiff || '';
        
        const statusFrom: string = data.status || 'unknown';
        this.loadPossibleStatuses(statusFrom);

        this.labSessionLogs = data.logs || [];
      },
      error: (error) => {
        console.error('Error loading lab session:', error);
        this.notificationService.showError('Failed to load lab session');
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
  }

  isInvalid(fieldName: string): boolean {
    const field = this.statusChangeForm.get(fieldName);
    return !!(field && field.invalid && (field.dirty || field.touched));
  }

  onStatusChangeSubmit() {
    this.statusChangeFormSubmitted = true;
    if (this.statusChangeForm.valid) {
      const statusId = this.statusChangeForm.value.status.value;
      const comment = this.statusChangeForm.value.comment || '';
      this.adminService.postChangeLabStatus(this.namespaceId, statusId, comment).subscribe({
        next: () => {
          this.notificationService.showSuccess('Lab session status updated successfully');
          this.loadLabSession(this.namespaceId);
        },
        error: (error) => {
          console.error('Error updating lab session status:', error);
          this.notificationService.showError('Failed to update lab session status');
        }
      });
    } else {
      Object.keys(this.statusChangeForm.controls).forEach(key => {
        this.statusChangeForm.get(key)?.markAsTouched();
      });
    }
    console.log('Status Change Form Submitted:', this.statusChangeForm.value);
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

  getSeverity(status: string) {
    switch (status?.toLowerCase()) {
      case 'success':
        return 'success';
      case 'passed':
        return 'success';
      case 'timeout':
        return 'danger';
      case 'failed':
        return 'danger';
      default:
        return 'info';
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
