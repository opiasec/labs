import { Injectable } from '@angular/core';
import { MessageService } from 'primeng/api';

@Injectable({
  providedIn: 'root'
})
export class NotificationService {

  constructor(private _messageService: MessageService) { }

  showSuccess(message: string, title: string = 'success!') {
    this._messageService.add({ severity: 'success', summary: title, detail: message });
  }

  showError(message: string, title: string = 'Oh no! D:') {
    this._messageService.add({ severity: 'error', summary: title, detail: message });
  }

  showInfo(message: string, title: string = 'hey!') {
    this._messageService.add({ severity: 'info', summary: title, detail: message });
  }

  showWarn(message: string, title: string = 'caution!') {
    this._messageService.add({ severity: 'warn', summary: title, detail: message });
  }

  showCustom(message: string, title: string = 'Custom', severity: 'success' | 'error' | 'info' | 'warn' = 'info') {
    this._messageService.add({ severity: severity, summary: title, detail: message });
  }

  showSecondary(message: string, title: string = 'Info') {
    this._messageService.add({ severity: 'secondary', summary: title, detail: message });
  }

  showContrast(message: string, title: string = 'Info') {
    this._messageService.add({ severity: 'contrast', summary: title, detail: message });
  }
}
