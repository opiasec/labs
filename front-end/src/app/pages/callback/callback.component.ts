import { Component, OnInit } from '@angular/core';
import { AuthService } from '../../services/auth.service';
import { NotificationService } from '../../services/notification.service';
import { ProgressSpinnerModule } from 'primeng/progressspinner';
@Component({
  selector: 'app-callback',
  imports: [ProgressSpinnerModule],
  templateUrl: './callback.component.html',
  styleUrl: './callback.component.css'
})
export class CallbackComponent implements OnInit {

  constructor(
    private authService: AuthService,
    private notificationService: NotificationService
  ) {

  }

  ngOnInit() {
    this.tokenExchange();
  }

  tokenExchange() {
    const code = this.getCodeFromUrl();
    const codeVerifier = sessionStorage.getItem('code_verifier');
    if (code && codeVerifier) {
      this.authService.exchangeCodeForToken(code, codeVerifier).subscribe({
        next: (response) => {
          console.log('Token exchange successful:', response);
          localStorage.setItem('access_token', response.access_token);
          localStorage.setItem('id_token', response.id_token);
          sessionStorage.removeItem('code_verifier');
          this.notificationService.showSuccess('Login successful!');
          window.location.href = '/dashboard';
        },
        error: (error) => {
          this.notificationService.showError('Token exchange failed. Please try again.');
          console.error('Token exchange failed:', error);
        }
      });
    }
  }

  private getCodeFromUrl(): string | null {
    const urlParams = new URLSearchParams(window.location.search);
    return urlParams.get('code');
  }

}
