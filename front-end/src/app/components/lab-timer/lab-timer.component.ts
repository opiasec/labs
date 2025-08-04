import { Component, Input, OnInit, OnDestroy, EventEmitter, Output } from '@angular/core';
import { TagModule } from 'primeng/tag';
import { LabSessionService } from '../../services/lab-session.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-lab-timer',
  templateUrl: './lab-timer.component.html',
  imports: [TagModule]
})
export class LabTimerComponent implements OnInit, OnDestroy {
  @Input() endTime!: string;
  @Input() labSlug!: string;
  @Input() namespace!: string;

  @Output() onEnd = new EventEmitter<void>();

  remainingTime: string = '';
  private intervalId: any;
  constructor(private labSessionService: LabSessionService, private router: Router) {}
  ngOnInit() {
    this.updateTimer();
    this.intervalId = setInterval(() => this.updateTimer(), 1000);
  }

  ngOnDestroy() {
    clearInterval(this.intervalId);
  }

  get isCritical(): boolean {
    const parts = this.remainingTime.split(':');
    return parseInt(parts[0]) === 0 && parseInt(parts[1]) <= 59;
  }
  

  private updateTimer() {
    const end = new Date(this.endTime).getTime();
    const now = Date.now();
    const diff = end - now;

    if (diff <= 0) {
      this.remainingTime = '00:00';
      clearInterval(this.intervalId);
      this.onEnd.emit();
      this.labSessionService.clearSession();
      return;
    }

    const minutes = Math.floor(diff / 60000);
    const seconds = Math.floor((diff % 60000) / 1000);
    this.remainingTime = `${this.pad(minutes)}:${this.pad(seconds)}`;
  }

  private pad(value: number): string {
    return value < 10 ? `0${value}` : `${value}`;
  }

  goToLab() {
    this.router.navigate(['/lab'], { queryParams: { slug: this.labSlug} });
  }
}
