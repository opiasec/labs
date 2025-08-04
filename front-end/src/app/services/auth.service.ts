import { Injectable, inject } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { ConfigService } from './config.service';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  private baseURL = environment.apiUrl + '/api/auth';
  private loaded = false;
  private isAdminSubject = new BehaviorSubject<boolean>(false);

  constructor(
    private http: HttpClient,
    private configService: ConfigService
  ) { }


  isAdmin(): Observable<boolean> {

    const access_token = localStorage.getItem('access_token');
    if (!access_token) {
      this.isAdminSubject.next(false);
      return this.isAdminSubject.asObservable();
    }
    const payload = JSON.parse(atob(access_token.split('.')[1]));

    for (const role of payload.permissions) {
      if (role === 'admin') {
        this.isAdminSubject.next(true);
        return this.isAdminSubject.asObservable();
      }
    }

    return this.isAdminSubject.asObservable();
  }
  
  isAuthenticated(): boolean {
    const access_token = localStorage.getItem('access_token');
    const id_token = localStorage.getItem('id_token');
    if (access_token && id_token) {
      const payload = JSON.parse(atob(access_token.split('.')[1]));
      const currentTime = Math.floor(Date.now() / 1000);
      return payload.exp > currentTime;
    }
    return false;
  }

  getUser(): Promise<any> {
    return new Promise((resolve, reject) => {
      const id_token = localStorage.getItem('id_token');
      if (!id_token) {
        reject('No ID token found');
        return;
      }

      const payload = JSON.parse(atob(id_token.split('.')[1]));
      resolve(payload);
    });
  }

  login(email: string, password: string): Observable<any> {
    const formData = new FormData();
    formData.append('email', email);
    formData.append('password', password);
    
    return this.http.post(`${this.baseURL}/login`, formData);
    
  }
  exchangeCodeForToken(code: string, codeVerifier: string): Observable<any> {
    const config = this.configService.getAppConfig();
    const headers = new HttpHeaders({
      'Content-Type': 'application/x-www-form-urlencoded'
    });
    const baseURL = window.location.origin || 'http://localhost:4200';
    const body = new URLSearchParams();
    body.set('code', code);
    body.set('grant_type', 'authorization_code');
    body.set('client_id', config.providerConfig.clientId);
    body.set('redirect_uri', baseURL + '/callback');
    body.set('code_verifier', codeVerifier);
    return this.http.post(`${config.providerConfig.tokenUrl}`, body.toString(), { headers });
  }
  logout(): void {

    localStorage.removeItem('access_token');
    localStorage.removeItem('id_token');
    sessionStorage.removeItem('code_verifier');
  }
}
