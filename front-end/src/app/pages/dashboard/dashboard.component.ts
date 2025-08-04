import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { DashboardService } from '../../services/dashboard.service';
import { ButtonModule } from 'primeng/button';
import { ChartModule } from 'primeng/chart';
import { ProgressBarModule } from 'primeng/progressbar';
import { SelectModule } from 'primeng/select';
import { Router } from '@angular/router';
import { CardModule } from 'primeng/card';
import { SimpleCardComponent } from '../../components/simple-card/simple-card.component';
import { NotificationService } from '../../services/notification.service';


@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ButtonModule,
    ChartModule,
    ProgressBarModule,
    SelectModule,
    CardModule,
    SimpleCardComponent
  ],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.css'
})
export class DashboardComponent implements OnInit {

  dashboardData: any;
  
  selectedTime: string = '7 days';

  weeklyActivityChartData: any;
  weeklyActivityChartOptions: any;

  scoreEvolutionChartData: any;
  scoreEvolutionChartOptions: any;



  ngOnInit(): void {
    this.loadDashboardData();
  }

  constructor(
    private dashboardService: DashboardService, 
    private notificationService: NotificationService,
    private router: Router) {
  }

  loadDashboardData() {
    this.dashboardService.getDashboardData().subscribe({
      next: (data: any) => {
        this.dashboardData = data;
        this.initWeeklyActivityChart(data.weeklyActivity);
        this.initScoreEvolutionChart(data.scoreEvolution);
      },
      error: (error: any) => {
        console.error('Error loading dashboard data', error);
        this.initWeeklyActivityChart(null);
        this.initScoreEvolutionChart(null);
      }
    });
  }

  initWeeklyActivityChart(weeklyActivity: any) {
    let labels: string[] = [];
    let data: number[] = [];
    
    if (weeklyActivity && weeklyActivity.length > 0) {
      weeklyActivity.forEach((item: any) => {
        const date = new Date(item.day);
        const dayName = date.toLocaleDateString('en-US', { weekday: 'short' });
        labels.push(dayName);
        data.push(item.count || 0);
      });
    } else {
      labels = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
      data = [0, 0, 0, 0, 0, 0, 0];
    }

    this.weeklyActivityChartData = {
      labels: labels,
      datasets: [
        {
          label: 'Activity',
          data: data,
          fill: true,
          tension: 0.4,
          backgroundColor: 'rgb(210, 210, 210)',
        }
      ]
    };

    this.weeklyActivityChartOptions = {
      responsive: true,
      plugins: {
        legend: {
          position: 'top'
        }
      }
    };
  }
  initScoreEvolutionChart(scoreEvolution: any) {
    let labels: string[] = [];
    let data: number[] = [];
    
    if (scoreEvolution && scoreEvolution.length > 0) {
      scoreEvolution.forEach((item: any) => {
        const date = new Date(item.day);
        const dayName = date.toLocaleDateString('en-US', { weekday: 'short' });
        labels.push(dayName);
        data.push(item.score || 0);
      });
    } else {
      labels = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
      data = [0, 0, 0, 0, 0, 0, 0];
    }
    
    this.scoreEvolutionChartData = {
      labels: labels,
      datasets: [
        {
          label: 'Score Evolution',
          data: data,
          fill: false,
          borderColor: 'rgb(210, 210, 210)',
          tension: 0.2
        }
      ]
    };

    this.scoreEvolutionChartOptions = {
      responsive: true,
      plugins: {
        legend: {
          position: 'top'
        }
      },
      scales: {
        y: {
          beginAtZero: true,
          min: 0,
          max: 100,
          ticks: {
            stepSize: 10
          }
        }
      }
    };
  }
}