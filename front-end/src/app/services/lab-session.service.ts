import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { environment } from '../../environments/environment';

@Injectable({ providedIn: 'root' })
export class LabSessionService {
  private labEndTime$ = new BehaviorSubject<string | null>(null);

  constructor() {
    const stored = sessionStorage.getItem('activeLab');
    if (stored) {
      const { expiresAt } = JSON.parse(stored);
      this.labEndTime$.next(expiresAt);
    }
   }

  setLabEndTime(endTime: string) {
    this.labEndTime$.next(endTime);
    const current = localStorage.getItem('activeLab');
    const data = current ? JSON.parse(current) : {};
    localStorage.setItem('activeLab', JSON.stringify({
      ...data,
      expiresAt: endTime
    }));
  }

  getLabEndTimeObservable() {
    return this.labEndTime$.asObservable();
  }
  clearSession() {
    this.labEndTime$.next(null);
    const activeLab = localStorage.getItem('activeLab');
    if(activeLab) {
      const { labSlug } = JSON.parse(activeLab);
      localStorage.removeItem('activeLab');
      localStorage.removeItem(labSlug);
    }
  }
}
