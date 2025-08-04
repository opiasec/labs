import { ChangeDetectorRef, Component, OnInit, ViewChild } from '@angular/core';
import { ConfirmationService, MessageService } from 'primeng/api';
import { TableModule } from 'primeng/table';
import { Dialog } from 'primeng/dialog';
import { Ripple } from 'primeng/ripple';
import { ButtonModule } from 'primeng/button';
import { ToastModule } from 'primeng/toast';
import { ToolbarModule } from 'primeng/toolbar';
import { ConfirmDialog } from 'primeng/confirmdialog';
import { InputTextModule } from 'primeng/inputtext';
import { TextareaModule } from 'primeng/textarea';
import { CommonModule } from '@angular/common';
import { SelectModule } from 'primeng/select';
import { Tag } from 'primeng/tag';
import { FormsModule } from '@angular/forms';
import { InputNumber } from 'primeng/inputnumber';
import { IconFieldModule } from 'primeng/iconfield';
import { InputIconModule } from 'primeng/inputicon';
import { LabService } from '../../../services/lab.service';
import { AdminService } from '../../../services/admin.service';
import { UtilsService } from '../../../services/utils.service';
import { TabsModule } from 'primeng/tabs';
import { MultiSelectModule } from 'primeng/multiselect';
import { ToggleSwitch } from 'primeng/toggleswitch';
import { EditorModule } from 'primeng/editor';
import { ChipsModule } from 'primeng/chips';

@Component({
  selector: 'app-admin-labs-management',
  standalone: true,
  imports: [
    CommonModule, 
    FormsModule,
    ButtonModule,
    TableModule,
    Dialog, 
    Ripple,
    EditorModule,
    ToastModule,
    ToolbarModule, 
    ConfirmDialog,
    InputTextModule,
    ToggleSwitch,
    TextareaModule, 
    MultiSelectModule,
    SelectModule,
    TabsModule,
    Tag, 
    InputNumber,
    IconFieldModule, 
    InputIconModule,
    ChipsModule
  ],
  providers: [ConfirmationService, MessageService],
  templateUrl: './admin-labs-management.component.html',
  styleUrls: ['./admin-labs-management.component.css']
})
export class AdminLabsManagementComponent implements OnInit {

  selectedLabs: any[] = [];
  labs: any[] = [];
  lab: any = {};

  labDialog: boolean = false;
  submitted: boolean = false;
  possiblesDifficulties: any[] = [
    { label: 'Easy', value: 'easy' },
    { label: 'Medium', value: 'medium' },
    { label: 'Hard', value: 'hard' }
  ];
  possiblesVulnerabilities: any[] = [
  ];
  possiblesTechnologies: any[] = [
  ];
  possiblesLanguages: any[] = [
  ];
  
  possibleImages: any[] = [
    { label: 'Ubuntu 20.04', value: 'ubuntu:20.04' },
    { label: 'Ubuntu 22.04', value: 'ubuntu:22.04' },
    { label: 'Alpine Linux', value: 'alpine:latest' },
    { label: 'Node.js 18', value: 'node:18' },
    { label: 'Python 3.9', value: 'python:3.9' },
    { label: 'PHP 8.1', value: 'php:8.1' },
    { label: 'MySQL 8.0', value: 'mysql:8.0' },
    { label: 'PostgreSQL 14', value: 'postgres:14' },
  ];
  keyValuePairs: { key: string; value: string }[] = [];
  evaluators: any[] = [];
  possibleEvaluators: any[] = [
  ];

  ngOnInit(): void {
    this.loadLabs();
    this.getPossiblesVulnerabilities();
    this.getPossibleImages();
    this.getPossibleEvaluators();
    this.getPossiblesTechnologies();
    this.getPossiblesLanguages();
  }

  loadLabs(): void {
    this.adminService.getLabDefinitions().subscribe({
      next: (labs: any) => {
        this.labs = labs
        console.log('Labs loaded:', this.labs);
      },
      error: (error) => {
        console.error('Error loading labs:', error);
        this.messageService.add({ severity: 'error', summary: 'Error', detail: 'Failed to load labs' });
      }
    });
  }
  cols: any[] = []; 

  constructor(
    private confirmationService: ConfirmationService,
    private messageService: MessageService,
    private labService: LabService,
    public utilsService: UtilsService,
    private adminService: AdminService,
  ) {}


  openNew(): void {
    this.lab = {
      externalReferences: [],
      authors: [],
      active: true,
      requiresManualReview: false,
      config: {
        systemRequirements: {
          image: null,
          codeConfig: {
            gitUrl: '',
            gitBranch: 'main', 
            gitPath: ''
          },
          envVars: {}
        },
        evaluators: [],
      }
    };
    this.keyValuePairs = [];
    this.evaluators = [];
    this.labDialog = true;
  }
  
deleteSelectedLabs(): void {
  this.confirmationService.confirm({
    message: 'Are you sure you want to delete the selected labs?',
    header: 'Confirm',
    icon: 'pi pi-exclamation-triangle',
    accept: () => {
      if (this.selectedLabs.length > 0) {
        const deleteRequests = this.selectedLabs.map((lab) =>
          this.adminService.deleteLabDefinition(lab.slug).toPromise()
        );

        Promise.all(deleteRequests)
          .then(() => {
            this.messageService.add({
              severity: 'success',
              summary: 'Success',
              detail: 'All selected labs were deleted successfully',
            });
            this.loadLabs();
          })
          .catch((error) => {
            console.error('Error deleting labs:', error);
            this.messageService.add({
              severity: 'error',
              summary: 'Error',
              detail: 'Failed to delete some labs',
            });
          });
      } else {
        this.messageService.add({
          severity: 'warn',
          summary: 'Warning',
          detail: 'No labs selected for deletion',
        });
      }
    },
  });
}

getPossiblesVulnerabilities(): void {
    this.adminService.getPossiblesVulnerabilities().subscribe({
      next: (vulnerabilities: any[]) => {
        this.possiblesVulnerabilities = vulnerabilities
      },
      error: (error) => {
        console.error('Error loading possible vulnerabilities:', error);
        this.messageService.add({ severity: 'error', summary: 'Error', detail: 'Failed to load possible vulnerabilities' });
      }
    });
  }
  getPossibleImages(): void {
    this.adminService.getPossibleImages().subscribe({
      next: (images: any[]) => {
        this.possibleImages = images.map(image => ({
          label: image.name,
          value: image.value
        }));
      },
      error: (error) => {
        console.error('Error loading possible images:', error);
        this.messageService.add({ severity: 'error', summary: 'Error', detail: 'Failed to load possible images' });
      }
    });
  }
  getPossibleEvaluators(): void {
    this.adminService.getPossibleEvaluators().subscribe({
      next: (evaluators: any[]) => {
        this.possibleEvaluators = evaluators.map(evaluator => ({
          label: evaluator.name,
          value: evaluator.value
        }));
      },
      error: (error) => {
        console.error('Error loading possible evaluators:', error);
        this.messageService.add({ severity: 'error', summary: 'Error', detail: 'Failed to load possible evaluators' });
      }
    });
  }
  getPossiblesTechnologies(): void {
    this.adminService.getPossiblesTechnologies().subscribe({
      next: (technologies: any[]) => {
        this.possiblesTechnologies = technologies
      },
      error: (error) => {
        console.error('Error loading possible technologies:', error);
        this.messageService.add({ severity: 'error', summary: 'Error', detail: 'Failed to load possible technologies' });
      }
    });
  }
  getPossiblesLanguages(): void {
    this.adminService.getPossiblesLanguages().subscribe({
      next: (languages: any[]) => {
        this.possiblesLanguages = languages
      },
      error: (error) => {
        console.error('Error loading possible languages:', error);
        this.messageService.add({ severity: 'error', summary: 'Error', detail: 'Failed to load possible languages' });
      }
    });
  }

  getSeverity(isActive: boolean): string {
    switch (isActive) {
      case true:
        return 'success';
      case false:
          return 'danger';
        default:
          return 'info';
      }
  }

  getDifficultySeverity(difficulty: string) {
    switch (difficulty?.toLowerCase()) {
      case 'easy':
        return 'success';
      case 'medium':
        return 'warn';
      case 'hard':
        return 'danger';
      default:
        return null;
    }
  }

  hideDialog(): void {
    this.labDialog = false;
    this.lab = {};
  }

  saveLab(): void {
    this.submitted = true;
    
    if (!this.isFormValid()) {
      this.messageService.add({
        severity: 'error',
        summary: 'Error',
        detail: 'Please fill in all required fields'
      });
      return;
    }
    
    const validPairs = this.keyValuePairs.filter(pair => this.isKeyValuePairValid(pair));
    this.lab.config.systemRequirements.envVars = {};
    validPairs.forEach(pair => {
      this.lab.config.systemRequirements.envVars[pair.key] = pair.value;
    });
    
    this.lab.config.evaluators = this.evaluators.map(evaluator => {
      const processedEvaluator: any = {
        slug: evaluator.slug,
        weight: evaluator.weight,
        config: {}
      };
      
      if (evaluator.config) {
        if (Array.isArray(evaluator.config)) {
          processedEvaluator.config = this.convertConfigArrayToObject(evaluator.config);
        } else if (typeof evaluator.config === 'object') {
          processedEvaluator.config = evaluator.config;
        }
      }
      
      if (evaluator.slug === 'exploit' && evaluator.exploitTemplate) {
        processedEvaluator.exploitTemplate = evaluator.exploitTemplate;
      }
      
      return processedEvaluator;
    }).filter(evaluator => evaluator.slug);
    
    console.log('Saving lab:', this.lab);
    console.log('Environment Variables Object:', this.lab.config.systemRequirements.envVars);
    console.log('Code Config:', this.lab.config.systemRequirements.codeConfig);
    console.log('Evaluators:', this.lab.config.evaluators);
    
    if (this.isEditMode()) {
      this.updateLab();
    } else {
      this.createLab();
    }
  }

  private isEditMode(): boolean {
    return this.lab && (this.lab.id);
  }

  private createLab(): void {
    this.adminService.postLabDefinition(this.lab).subscribe({
      next: (response) => {
        this.messageService.add({ 
          severity: 'success', 
          summary: 'Success', 
          detail: 'Lab created successfully' 
        });
        this.loadLabs();
        this.hideDialog();
      },
      error: (error) => {
        console.error('Error creating lab:', error);
        this.messageService.add({ 
          severity: 'error', 
          summary: 'Error', 
          detail: 'Failed to create lab: ' + (error.error?.message || error.message)
        });
      }
    });
  }

  private updateLab(): void {
    const slug = this.lab.slug;
    this.adminService.putLabDefinition(slug, this.lab).subscribe({
      next: (response) => {
        this.messageService.add({ 
          severity: 'success', 
          summary: 'Success', 
          detail: 'Lab updated successfully' 
        });
        this.loadLabs();
        this.hideDialog();
      },
      error: (error) => {
        console.error('Error updating lab:', error);
        this.messageService.add({ 
          severity: 'error', 
          summary: 'Error', 
          detail: 'Failed to update lab: ' + (error.error?.message || error.message)
        });
      }
    });
  }

  private isFormValid(): boolean {
    if (!this.lab.slug || !this.lab.title || !this.lab.description || 
        !this.lab.difficulty || !this.lab.config?.systemRequirements?.image) {
      return false;
    }
    
    if (!this.lab.config.systemRequirements.codeConfig?.gitUrl || 
        !this.lab.config.systemRequirements.codeConfig?.gitBranch) {
      return false;
    }
    
    if (!this.lab.vulnerabilities || this.lab.vulnerabilities.length === 0) {
      return false;
    }
    
    return true;
  }

  editLab(lab: any): void {
    this.adminService.getLabDefinition(lab.slug).subscribe({
      next: (labData) => {
        this.lab = labData;
        this.lab.config.systemRequirements = labData.config.labSpec;
        this.ensureLabStructure(this.lab);
        
        this.evaluators = (this.lab.config.evaluators || []).map((evaluator: any) => ({
          slug: evaluator.slug,
          weight: evaluator.weight,
          exploitTemplate: evaluator.exploitTemplate || '',
          config: this.convertConfigObjectToArray(evaluator.config || {})
        }));
        
        this.keyValuePairs = this.convertConfigObjectToArray(
          this.lab.config.systemRequirements.envVars || {}
        );
        
        this.labDialog = true;
      },
      error: (error) => {
        console.error('Error loading lab for editing:', error);
        this.messageService.add({ severity: 'error', summary: 'Error', detail: 'Failed to load lab for editing' });
      }
    });
  }

  deleteLab(lab: any): void {
    this.confirmationService.confirm({
      message: 'Are you sure you want to delete this lab?',
      header: 'Confirm',
      icon: 'pi pi-exclamation-triangle',
      accept: () => {
        this.adminService.deleteLabDefinition(lab.slug).subscribe({
          next: () => {
            this.messageService.add({ severity: 'success', summary: 'Success', detail: 'Lab deleted successfully' });
            this.loadLabs();
          },
          error: (error) => {
            console.error('Error deleting lab:', error);
            this.messageService.add({ severity: 'error', summary: 'Error', detail: 'Failed to delete lab' });
          }
        });
      }
    });
  }

  onGlobalFilter(table: any, event: Event): void {
    const target = event.target as HTMLInputElement;
    table.filterGlobal(target.value, 'contains');
  }

  addKeyValuePair(): void {
    this.keyValuePairs.push({ key: '', value: '' });
  }

  removeKeyValuePair(index: number): void {
    this.keyValuePairs.splice(index, 1);
  }

  isKeyValuePairValid(pair: { key: string; value: string }): boolean {
    return pair.key.trim() !== '' && pair.value.trim() !== '';
  }

  addEvaluator(): void {
    this.evaluators.push({
      slug: null,
      weight: null,
      config: [],
      exploitTemplate: ''
    });
  }

  removeEvaluator(index: number): void {
    this.evaluators.splice(index, 1);
  }

  onEvaluatorTypeChange(evaluator: any): void {
    if (evaluator.slug !== 'exploit') {
      evaluator.exploitTemplate = '';
    }
  }

  addEvaluatorConfig(evaluatorIndex: number): void {
    if (!this.evaluators[evaluatorIndex].config) {
      this.evaluators[evaluatorIndex].config = [];
    }
    this.evaluators[evaluatorIndex].config.push({ key: '', value: '' });
  }

  removeEvaluatorConfig(evaluatorIndex: number, configIndex: number): void {
    this.evaluators[evaluatorIndex].config.splice(configIndex, 1);
  }

  private ensureLabStructure(lab: any): void {
    if (!lab.config) {
      lab.config = {};
    }
    if (!lab.config.systemRequirements) {
      lab.config.systemRequirements = {};
    }
    if (!lab.config.systemRequirements.codeConfig) {
      lab.config.systemRequirements.codeConfig = {
        gitUrl: '',
        gitBranch: 'main',
        gitPath: ''
      };
    }
  }

  private convertConfigObjectToArray(configObject: any): { key: string; value: string }[] {
    if (!configObject || typeof configObject !== 'object') {
      return [];
    }
    
    return Object.entries(configObject).map(([key, value]) => ({
      key: key,
      value: String(value)
    }));
  }

  private convertConfigArrayToObject(configArray: { key: string; value: string }[]): any {
    if (!configArray || !Array.isArray(configArray)) {
      return {};
    }
    
    const configObject: any = {};
    configArray.forEach(item => {
      if (item.key && item.key.trim() !== '') {
        configObject[item.key] = item.value || '';
      }
    });
    
    return configObject;
  }
}
