import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { CardModule } from 'primeng/card';
import { InputTextModule } from 'primeng/inputtext';
import { PasswordModule } from 'primeng/password';
import { ButtonModule } from 'primeng/button';
import { DividerModule } from 'primeng/divider';
import { CommonModule } from '@angular/common';
import { AuthService } from '../../services/auth.service';
import { NotificationService } from '../../services/notification.service';
import { UtilsService } from '../../services/utils.service';
import { ConfigService } from '../../services/config.service';

@Component({
  selector: 'app-login-new',
  imports: [
    CommonModule,
    ReactiveFormsModule,
    CardModule,
    InputTextModule,
    PasswordModule,
    ButtonModule,
    DividerModule
  ],
  templateUrl: './login-new.component.html',
  styleUrl: './login-new.component.css'
})
export class LoginNewComponent implements OnInit {
  loginForm!: FormGroup;
  isLoading = false;
  isAuth0Loading = false;
  providerName = '';
  providerConfig: any;
  loginWithProviderLabel = 'Login with';

  isPasswordVisible = false;
  isProviderVisible = true;

  constructor(
    private fb: FormBuilder,
    private authService: AuthService,
    private notificationService: NotificationService,
    private utils: UtilsService,
    private configService: ConfigService
  ) { }

  ngOnInit() {
    const appConfig = this.configService.getAppConfig();
    if (!this.configService.getAppConfig()) {
      this.notificationService.showError('Application configuration is not loaded. Please try again later.');
      return;
    }
    this.isProviderVisible = appConfig?.providerLoginEnabled;
    this.isPasswordVisible = appConfig?.passwordAuthEnabled;
    this.providerName = appConfig?.providerConfig.name || 'External Provider';
    this.providerConfig = appConfig?.providerConfig || {};
    this.loginWithProviderLabel = `Login with ${this.providerName}`;


    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]]
    });
  }

  onLogin() {
    if (this.loginForm.valid) {
      this.isLoading = true;

      this.authService.login(this.loginForm.value.email, this.loginForm.value.password).subscribe({
        next: (response) => {
          console.log('Login successful:', response);
          this.isLoading = false;
          localStorage.setItem('access_token', response.access_token);
          localStorage.setItem('id_token', response.id_token);
          window.location.href = '/dashboard'; // Adjust the redirect path as needed
          this.notificationService.showSuccess('Login successful!');
        },
        error: (error) => {
          console.error('Login failed:', error);
          this.isLoading = false;
          this.notificationService.showError('Login failed. Please try again.');
        }
      });
    }
  }


  async loginWithProvider() {
    this.isAuth0Loading = true;

    try {
      const codeVerifier = this.utils.generateCodeVerifier();
      const codeChallenge = await this.utils.generateCodeChallenge(codeVerifier);

      sessionStorage.setItem('code_verifier', codeVerifier);
      const authorizationUrl = this.providerConfig.authorizationUrl;
      const client_id = this.providerConfig.clientId;
      const scope = this.providerConfig.scope;
      const redirect_uri = window.location.origin + '/callback' || 'http://localhost:4200/callback';
      const audience = this.providerConfig.audience;

      const providerAuthorizationUrl = authorizationUrl + '?' +
        'client_id=' + client_id + '&' +
        'response_type=code&' +
        'scope=' + scope + '&' +
        'audience=' + audience + '&' +
        'redirect_uri=' + redirect_uri + '&' +
        'code_challenge=' + codeChallenge + '&' +
        'code_challenge_method=S256';
      window.location.href = providerAuthorizationUrl;
    }
    catch (error) {
      console.error('Login with external provider failed:', error);
      this.isAuth0Loading = false;
      this.notificationService.showError('Login with external provider failed. Please try again.');
    }
  }
}
