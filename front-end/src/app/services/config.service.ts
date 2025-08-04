import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';

@Injectable({
    providedIn: 'root',
})
export class ConfigService {
    private baseURL = environment.apiUrl + '/api';
    private config!: any;

    constructor(private http: HttpClient) { }

    load(): Promise<any> {
        return this.http.get(`${this.baseURL}/config/`).toPromise().then(config => {
            this.config = config;
            return config;
        }).catch(error => {
            console.error('Error loading configuration:', error);
        });
    }
    getAppConfig() {
        return this.config;
    }
}

