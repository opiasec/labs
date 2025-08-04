import { Component, OnInit, OnDestroy } from '@angular/core';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';
import { SplitterModule } from 'primeng/splitter';
import { ButtonModule } from 'primeng/button';
import { CommonModule } from '@angular/common';
import { ProgressSpinnerModule } from 'primeng/progressspinner';
import { AvatarModule } from 'primeng/avatar';
import { TabsModule } from 'primeng/tabs';
import { LabService } from '../../services/lab.service';
import { interval } from 'rxjs';
import { ActivatedRoute } from '@angular/router';
import { ScrollPanelModule } from 'primeng/scrollpanel';
import { GetIDEStatusResponse, CreateLabResponse, FinishLabResponse } from '../../types/lab';
import { CardModule } from 'primeng/card';
import { TooltipModule } from 'primeng/tooltip';
import { LabSessionService } from '../../services/lab-session.service';
import { NotificationService } from '../../services/notification.service';
import { SensitiveDataDisplayComponent } from '../../components/sensitive-data-display/sensitive-data-display.component';
import { MarkdownModule } from 'ngx-markdown';
import { MarkdownClipboardButtonComponent } from '../../components/markdown-clipboard-button/markdown-clipboard-button.component';
import { ScrollTop } from 'primeng/scrolltop';
import { ConfirmDialog } from 'primeng/confirmdialog';
import { ConfirmationService } from 'primeng/api';

@Component({
  selector: 'app-lab',
  standalone: true,
  imports: [SplitterModule, ButtonModule,
    CommonModule, ProgressSpinnerModule,
    AvatarModule, TabsModule,
    ScrollPanelModule,
    CardModule,
    SensitiveDataDisplayComponent,
    MarkdownModule,
    ScrollTop,
    TooltipModule,
    ConfirmDialog,
  ],
  providers: [ConfirmationService],
  templateUrl: './lab.component.html',
  styleUrl: './lab.component.css'
})


export class LabComponent implements OnInit, OnDestroy {
  readonly buttonComponent = MarkdownClipboardButtonComponent;
  statusPolling: any;
  labSlug: string = '';
  labNamespace: string = '';
  labActiveSession: any = null;
  labIsDeploying = false;
  labIsReady = false;
  


  labIsRedeploying = false;

  labPassword: string = '';
  labIntroMarkdown: string = '';
  ideUrl: any = '';
  ideUrlUnsafe: any = '';

  ideReady = false;
  appReady = false;

  codeServerUrl: any;
  labAppUrl: any;
  labStarted = false;
  getIDEStatusResponse: GetIDEStatusResponse | null = null;
  finishLabResponse: FinishLabResponse | null = null;
  introHTML: SafeHtml | null = null;

  labPasswordVisible = false;

  constructor(
    private sanitizer: DomSanitizer, 
    private labService: LabService,
    private route: ActivatedRoute,
    private labSessionService: LabSessionService,
    private notificationService: NotificationService,
    private confirmationService: ConfirmationService
  ) {}

  ngOnInit() {

    this.route.queryParams.subscribe(params => {
      this.labSlug = params['slug'];
      if (this.labSlug) {
        this.loadLabContent();
        if(this.checkForActiveLab(this.labSlug)) {
          return;
        }
      }
    });
  }

  checkForActiveLab(labSlug: string): boolean {
    const activeLab = JSON.parse(localStorage.getItem('activeLab') || '{}');
    if (!activeLab || !activeLab.labSlug || activeLab.labSlug !== labSlug) {
      console.log("No active lab found or lab slug does not match.");
      return false;
    }
      this.ideUrlUnsafe = activeLab!.ideUrl!;
      this.ideUrl = this.sanitizer.bypassSecurityTrustResourceUrl(activeLab!.ideUrl!);
      this.labSessionService.setLabEndTime(activeLab!.expiresAt!);
      this.labNamespace = activeLab!.namespace!;
      this.labIsReady = true;
      this.labStarted = true;
      this.labIsDeploying = false;
      console.log("new date", new Date(activeLab?.expiresAt));
      if (activeLab?.expiresAt && new Date(activeLab?.expiresAt) < new Date()) {
        localStorage.removeItem('activeLab');
        return false;
      }
      return true;
    }

  deployLab() {
    if (!this.labSlug) return;
    this.createLab();
    console.log("Deploying lab (deploy lab function):", this.labSlug);
  }

  createLab() {
    if (!this.labSlug) return;
    
    if (this.labIsDeploying) {
      console.log("Lab already being deployed, skipping...");
      return;
    }

    if (this.hasActiveLab()) {
      console.log("Active lab found, skipping deployment.");
      this.notificationService.showInfo('You already have an active lab.');
      return;
    }
    
    this.labStarted = false;
    this.labIsReady = false;
    this.ideReady = false;
    this.labIsDeploying = true;
    
    if (this.statusPolling) {
      this.statusPolling.unsubscribe();
    }
    
    this.labService.createLab(this.labSlug).subscribe({
      next: (startResponse: CreateLabResponse) => {
        if (!startResponse.namespace) {
          this.labIsDeploying = false;
            this.labStarted = false;
            this.labIsReady = false;
            this.ideReady = false;
            this.notificationService.showError('error creating lab');
          return;
        }
        this.notificationService.showInfo('Lab is being deployed', "It's quick, I promise! :)");

        this.labNamespace = startResponse.namespace!;
        this.labPassword = startResponse.labPassword!;
        localStorage.setItem('activeLab', JSON.stringify({
          labSlug: this.labSlug,
          ...startResponse
        }));
        this.setActiveLab({
          labSlug: this.labSlug,
          ...startResponse
        }, 'deploying');
        this.ideUrlUnsafe = startResponse!.ideUrl!;
        this.ideUrl = this.sanitizer.bypassSecurityTrustResourceUrl(startResponse!.ideUrl!);

        this.checkActiveLabStatus();

      },
      error: (error) => {
        this.labIsDeploying = false;
        console.error("Error creating lab:", error);
        this.notificationService.showError('Error creating lab');
      }
    });
  }

  hasActiveLab(): boolean {
    const activeLab = localStorage.getItem('activeLab');
    return !!activeLab;
  }

  checkActiveLabStatus() {
    if (!this.labSlug){
      console.error("Lab slug is not defined, cannot check lab status.");
      this.notificationService.showError('Lab slug is not defined');
      return;
    }
    const activeLab = JSON.parse(localStorage.getItem('activeLab')!);

    if (!activeLab) {
      console.error("No active lab found in local storage.");
      this.notificationService.showError('No active lab found');
      return;
    }

    this.statusPolling = interval(5000).subscribe(() => {
          this.labService.getLabStatus(this.labNamespace).subscribe({
            next: (labStatusResponse) => {
              this.getIDEStatusResponse = labStatusResponse;
              if (labStatusResponse.status === 'alive' || labStatusResponse.status === 'expired') {
                console.log("Active lab status:", activeLab);

                this.setActiveLab(activeLab, labStatusResponse.status);
                this.ideReady = true;
                this.labIsReady = true;
                this.labStarted = true;
                this.labIsDeploying = false;
                console.log("Lab is ready");
                this.notificationService.showSuccess('Lab is ready.', "LETS GO!" );
                this.statusPolling.unsubscribe();
                return;
              }
            },
            error: (error: any) => {
              console.error("lab is not ready", error);
            }
          });
        });
  }

  setActiveLab(activeLabData: any, status: string = 'deploying') {
    const sessionData = {
      ...activeLabData,
      status: status,
    }
    
    localStorage.setItem("activeLab", JSON.stringify(sessionData));
    this.labSessionService.setLabEndTime(activeLabData.expiresAt!);
  }
  

  finishLab() {
    if (!this.labSlug) return;
    this.labStarted = true;
    this.labService.finishLab(this.labNamespace).subscribe({
      next: (finishResponse : FinishLabResponse) => {
        if (finishResponse.message === 'finish lab started') {
          this.labSessionService.clearSession();
          localStorage.removeItem(this.labSlug);
          this.labIsRedeploying = false;
          this.labIsReady = false;
          this.labStarted = false;
        }
        this.notificationService.showSuccess('lab finished successfully.');
      },
      error: (error) => {
        console.error("Error finishing lab:", error);
        this.notificationService.showError('Error finishing lab');
      }
    });
  }
  leaveLab() {
    console.log("Leaving lab:", this.labSlug);
    console.log("Lab namespace:", this.labNamespace);
    this.confirmationService.confirm({
    message: 'Are you sure you want to leave the lab?',
    header: 'Confirm',
    icon: 'pi pi-exclamation-triangle',
    accept: () => {
    if (!this.labSlug) return;
    this.labService.leaveLab(this.labNamespace).subscribe({
      next: () => {
        this.labSessionService.clearSession();
        localStorage.removeItem(this.labSlug);
        this.labIsRedeploying = false;
        this.labIsReady = false;
        this.labStarted = false;
        this.notificationService.showSuccess('You left the lab successfully.');
      },
      error: (error) => {
        console.error("Error leaving lab:", error);
        this.notificationService.showError('Error leaving lab');
      }
    });
  }
  });
  }

  toggleLabPasswordVisibility() {
    this.labPasswordVisible = !this.labPasswordVisible;
  }

  getMaskedLabPassword(): string {
    return this.labPassword ? '•'.repeat(this.labPassword.length) : '';
  }

  openIdeInNewTab() {
    if (this.ideUrlUnsafe) {
      window.open(this.ideUrlUnsafe, '_blank');
    }
  }

  async loadLabContent() {
    if (!this.labSlug) return;
    this.labService.getLabDefinition(this.labSlug).subscribe({
      next: (labContent) => {
        this.labIntroMarkdown = labContent.readme || '';
      },
      error: (error) => {
        console.error("Error loading lab content:", error);
        this.notificationService.showError('Error loading lab content');
      }
    });
  }

  ngOnDestroy() {
    if (this.statusPolling) {
      this.statusPolling.unsubscribe();
    }
  }
}