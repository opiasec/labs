import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class AdminService {
  private baseURL = environment.apiUrl + '/api/admin';

  constructor(private http: HttpClient) { }

  // Lab Sessions Management
  getLabsSessions(): Observable<any> {
    return this.http.get(`${this.baseURL}/labs`);
  }
  getLabSession(sessionId: string): Observable<any> {
    return this.http.get(`${this.baseURL}/labs/${sessionId}`);
  }
  getPossiblesStatus(statusFrom: string): Observable<any> {
    return this.http.get(`${this.baseURL}/labs/status?from=${statusFrom}`);
  }
  postChangeLabStatus(namespace: string, statusId: string, comment: string): Observable<any> {
    return this.http.post(`${this.baseURL}/labs/${namespace}/status`, { statusId, comment });
  }

  // Lab Definition Management
  getLabDefinitions(): Observable<any> {
    return this.http.get(`${this.baseURL}/lab-definition/`);
  }
  getLabDefinition(slug: string): Observable<any> {
    return this.http.get(`${this.baseURL}/lab-definition/${slug}`);
  }
  postLabDefinition(labDefinition: any): Observable<any> {
    return this.http.post(`${this.baseURL}/lab-definition/`, labDefinition);
  }
  putLabDefinition(slug: string, labDefinition: any): Observable<any> {
    return this.http.put(`${this.baseURL}/lab-definition/${slug}`, labDefinition);
  }
  deleteLabDefinition(slug: string): Observable<any> {
    return this.http.delete(`${this.baseURL}/lab-definition/${slug}`);
  }
  // Lab Definition Management Utils
  getPossiblesVulnerabilities(): Observable<any> {
    return this.http.get(`${this.baseURL}/lab-definition/vulnerabilities`);
  }
  getPossibleImages(): Observable<any> {
    return this.http.get(`${this.baseURL}/lab-definition/images`);
  }
  getPossibleEvaluators(): Observable<any> {
    return this.http.get(`${this.baseURL}/lab-definition/evaluators`);
  }
  getPossiblesLanguages(): Observable<any> {
    return this.http.get(`${this.baseURL}/lab-definition/languages`);
  }
  getPossiblesTechnologies(): Observable<any> {
    return this.http.get(`${this.baseURL}/lab-definition/technologies`);
  }


}