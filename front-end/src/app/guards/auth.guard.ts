import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { NotificationService } from '../services/notification.service';

export const authGuard = async (): Promise<boolean> => {
  const router = inject(Router);
  const authService = inject(AuthService);
  const notificationService = inject(NotificationService);

  try {
    const isAuthenticated = authService.isAuthenticated();

    if (!isAuthenticated) {
      await router.navigate(['/login']);
      localStorage.removeItem('access_token');
      localStorage.removeItem('id_token');
      return false;
    }
  
    return true;
  } catch (error) {
    console.error('Error at authGuard:', error);
    return false;
  }
};
