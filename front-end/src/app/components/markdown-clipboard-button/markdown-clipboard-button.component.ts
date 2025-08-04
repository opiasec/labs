import { Component } from '@angular/core';
import { ButtonModule } from 'primeng/button';
import { NotificationService } from '../../services/notification.service';

@Component({
  selector: 'app-markdown-clipboard-button',
  imports: [ButtonModule],
  templateUrl: './markdown-clipboard-button.component.html',
  styleUrl: './markdown-clipboard-button.component.css'
})
export class MarkdownClipboardButtonComponent {

  constructor(
    private notificationService: NotificationService
  ) { }

  onClick() {
    this.notificationService.showSuccess('copied to clipboard!');
  }
}
