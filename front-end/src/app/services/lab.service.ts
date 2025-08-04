import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { GetIDEStatusResponse,
  CreateLabResponse, FinishLabResponse,
  GetLabResultResponse,
  LabDefinition,
  GetAllLabsByUserAndStatusResponse } from '../types/lab';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class LabService {


  private baseURL = environment.apiUrl + '/api';

  constructor(private http: HttpClient) { }

  getLabsDefinition() {
    const response = this.http.get(`${this.baseURL}/lab-definition/`);
    return response;
  }


  createLab(labSlug: string) {
    console.log("Criando laboratório com slug:", labSlug);
    const response = this.http.post<CreateLabResponse>(`${this.baseURL}/labs/`, { labSlug: labSlug });

    return response;
  }

  getLabStatus (namespace: string): Observable<any> {
    const response = this.http.get(`${this.baseURL}/labs/${namespace}/status`);
    return response;
  }

  getLabCodeServerStatus (labCodeServerUrl: string): Observable<any> {
    const response = this.http.get(`${labCodeServerUrl}`);
    return response;
  }

  getLabResult (namespace: string): Observable<any> {
    const response = this.http.get(`${this.baseURL}/labs/${namespace}`);
    return response;
  }

  getAllLabsByUserAndStatus (status: string): Observable<any[]> {
    const response = this.http.get<any[]>(`${this.baseURL}/labs/?status=${status}`);
    return response;
  }
  
  redeployLab (namespace: string) {
    const response = this.http.post(`${this.baseURL}/labs/${namespace}/redeploy`, {});
    return response;
  }

  finishLab (namespace: string) {
    const response = this.http.post<FinishLabResponse>(`${this.baseURL}/labs/${namespace}/finish`, {});
    return response;
  }
  leaveLab (namespace: string) {
    const response = this.http.post(`${this.baseURL}/labs/${namespace}/leave`, {});
    return response;
  }

  getLabDefinition (labSlug: string) {
    const response = this.http.get<LabDefinition>(`${this.baseURL}/lab-definition/${labSlug}`);
    return response;
  }

  sendFeedback (namespace: string, rating: number, feedback: string) {
    const response = this.http.post(`${this.baseURL}/labs/${namespace}/feedback`, { rating: rating, feedback: feedback });
    return response;
  }
}
