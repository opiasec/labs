import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';

@Injectable({
    providedIn: 'root', 
})
export class UtilsService {
    private baseURL = environment.apiUrl + '/api';
     constructor(private http: HttpClient) { }
    getStatusSeverity(status: string) {
        switch (status?.toLowerCase()) {
            case 'success':
                return 'success';
            case 'passed':
                return 'success';
            case 'approved':
                return 'success';
            case 'timeout':
                return 'danger';
            case 'rejected':
                return 'danger';
            case 'abbandoned':
                return 'danger';
            case 'failed':
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
                return 'warning';
            case 'hard':
                return 'danger';
            default:
                return 'info';
        }
    }
    base64URLEncode(array: Uint8Array): string {
      return btoa(String.fromCharCode(...array))
        .replace(/\+/g, '-')
        .replace(/\//g, '_')
        .replace(/=/g, '');
    }
    generateCodeVerifier(): string {
    const array = new Uint8Array(32);
    crypto.getRandomValues(array);
    return this.base64URLEncode(array);
  }
  async generateCodeChallenge(codeVerifier: string): Promise<string> {
    const encoder = new TextEncoder();
    const data = encoder.encode(codeVerifier);
    const digest = await crypto.subtle.digest('SHA-256', data);
    return this.base64URLEncode(new Uint8Array(digest));
  }
  getAppConfig() {
    return this.http.get(`${this.baseURL}/config`);
  }
}
