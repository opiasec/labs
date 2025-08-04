import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { NotificationService } from '../services/notification.service';
import { firstValueFrom } from 'rxjs/internal/firstValueFrom';

export const adminGuard = async (): Promise<boolean> => {
    const router = inject(Router);
    const access_token = localStorage.getItem('access_token');
    const authService = inject(AuthService);
    const notificationService = inject(NotificationService);

    try {
        const isAdmin = await firstValueFrom(authService.isAdmin());
        if (isAdmin) {
            return true
        }
        notificationService.showError('You do not have permission to access this resource.');
        console.warn('Admin guard: User is not an admin');
        await router.navigate(['/dashboard']);
        return false;
    } catch (error) {
        console.error('Admin guard error', error);
        return false;
    }
};
