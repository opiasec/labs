import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { CardModule } from 'primeng/card';
import { ButtonModule } from 'primeng/button';
import { LabService } from '../../services/lab.service';
import { TableModule, Table } from 'primeng/table';
import {  Lab } from '../../types/lab';
import { DropdownModule } from 'primeng/dropdown';
import { InputTextModule } from 'primeng/inputtext';
import { MultiSelectModule } from 'primeng/multiselect';
import { SelectModule } from 'primeng/select';
import { FormsModule } from '@angular/forms';
import { TagModule } from 'primeng/tag';
import { IconFieldModule } from 'primeng/iconfield';
import { InputIconModule } from 'primeng/inputicon';
import { FilterService } from 'primeng/api';
import { NotificationService } from '../../services/notification.service';

interface DifficultyOption {
  name: string;
}

@Component({
  selector: 'app-labs',
  standalone: true,
  imports: [
    CommonModule,
    CardModule,
    ButtonModule,
    SelectModule,
    TableModule,
    DropdownModule,
    InputTextModule,
    MultiSelectModule,
    FormsModule,
    TagModule,
    IconFieldModule,
    InputIconModule,
  ],
  templateUrl: './labs.component.html',
  styleUrls: ['./labs.component.css']
})
export class LabsComponent implements OnInit {
  labs!: Lab[];
  filteredLabs!: Lab[];

  difficultyOptions: DifficultyOption[] = [
    { name: 'easy' },
    { name: 'medium' },
    { name: 'hard' }
  ];
  selectedDifficulty: any = null;
  selectedVulnerabilities: any[] = [];
  selectedLanguages: any[] = [];
  selectedTechnologies: any[] = [];


  constructor(
    private labService: LabService,
    private router: Router,
    private filterService: FilterService,
    private notificationService: NotificationService
  ) {
    // Custom filter for vulnerabilities
  this.filterService.register('vulnerabilityContains', (value: any, filter: any): boolean => {
    if (!filter || filter.trim() === '') {
      return true;
    }
    
    if (!value || !Array.isArray(value)) {
      return false;
    }
    
    return value.some((vuln: any) => 
      vuln.name && vuln.name.toLowerCase().includes(filter.toLowerCase())
    );
  });

  // Custom filter for languages
  this.filterService.register('languageContains', (value: any, filter: any): boolean => {
    if (!filter || filter.trim() === '') {
      return true;
    }
    
    if (!value || !Array.isArray(value)) {
      return false;
    }
    
    return value.some((lang: any) => 
      lang.name && lang.name.toLowerCase().includes(filter.toLowerCase())
    );
  });

  // Custom filter for technologies
  this.filterService.register('technologyContains', (value: any, filter: any): boolean => {
    if (!filter || filter.trim() === '') {
      return true;
    }
    
    if (!value || !Array.isArray(value)) {
      return false;
    }
    
    return value.some((tech: any) => 
      tech.name && tech.name.toLowerCase().includes(filter.toLowerCase())
    );
  });
}

  ngOnInit() {
    this.loadLabs();
  }

  loadLabs() {
    this.labService.getLabsDefinition().subscribe(
      (labs: any) => {
        this.labs = labs;
        this.filteredLabs = labs;
      },
      error => {
        console.error('Error loading labs:', error);
        this.notificationService.showError('failed to load labs');
      }
    );
  }

  startLab(labSlug: string) {
    this.router.navigate(['/lab'], { queryParams: { slug: labSlug } });
  }

  getSeverity(difficulty: string) {
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
} 