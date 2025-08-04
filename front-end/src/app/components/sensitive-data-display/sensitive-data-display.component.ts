// sensitive-data-display.component.ts
import { Component, Input } from '@angular/core';
import { Clipboard } from '@angular/cdk/clipboard';
import { MessageService } from 'primeng/api';
import { TooltipModule } from 'primeng/tooltip';
import { NotificationService } from '../../services/notification.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-sensitive-data-display',
  templateUrl: './sensitive-data-display.component.html',
  providers: [MessageService],
  standalone: true,
  imports: [TooltipModule, CommonModule]
})
export class SensitiveDataDisplayComponent {
  @Input() value: string = '';
  visible: boolean = false;

  constructor(private clipboard: Clipboard, private notificationService: NotificationService) {}

  toggleVisibility() {
    this.visible = !this.visible;
  }

  get hiddenValue(): string {
  return '•'.repeat(this.value.length);
}

  copy() {
    this.clipboard.copy(this.value);
    this.notificationService.showSuccess("Password copied to clipboard!");
  }
}
