import { Routes } from '@angular/router';
import { DashboardComponent } from './pages/dashboard/dashboard.component';
import { authGuard } from './guards/auth.guard';
import { adminGuard } from './guards/admin.guard';
import { ApplayoutComponent } from './components/layout/layout.component';
import { LabComponent } from './pages/lab/lab.component';
import { LabsComponent } from './pages/labs/labs.component';
import { ProfileComponent } from './components/settings/profile/profile.component';
import { SettingsComponent } from './components/settings/settings.component';
import { BillingComponent } from './components/settings/billing/billing.component';
import { ConfigurationsComponent } from './components/settings/configurations/configurations.component';
import { LabsSessionsComponent } from './pages/labs-sessions/labs-sessions.component';
import { LabSessionComponent } from './pages/lab-session/lab-session.component';
import { AdminLabsSessionsComponent } from './pages/admin/admin-labs-sessions/admin-labs-sessions.component';
import { AdminLabSessionComponent } from './pages/admin/admin-lab-session/admin-lab-session.component';
import { AdminLabsManagementComponent } from './pages/admin/admin-labs-management/admin-labs-management.component';
import { LoginNewComponent } from './pages/login-new/login-new.component';
import { CallbackComponent } from './pages/callback/callback.component';

export const routes: Routes = [
  {
    path: 'login',
    component: LoginNewComponent
  },
  {
    path: 'callback',
    component: CallbackComponent
  },
  {
    path: '',
    component: ApplayoutComponent,
    canActivate: [authGuard],
    children: [
      {
        path: 'dashboard',
        component: DashboardComponent,
        canActivate: [authGuard]
      },
      {
        path: 'sessions',
        component: LabsSessionsComponent,
        canActivate: [authGuard]
      },
      {
        path: 'admin',
        canActivate: [adminGuard],
        children: [
          {
            path: 'sessions',
            component: AdminLabsSessionsComponent
          },
          {
            path: 'sessions/:id',
            component: AdminLabSessionComponent
          },
          {
            path: 'labs/management',
            component: AdminLabsManagementComponent
          },
          { path: '', redirectTo: 'dashboard', pathMatch: 'full' }
        ]
      },

      {
        path: 'sessions/:id',
        component: LabSessionComponent,
        canActivate: [authGuard]
      },
      { path: 'labs', component: LabsComponent, canActivate: [authGuard] },
      { path: 'lab', component: LabComponent, canActivate: [authGuard] },
      { path: '', redirectTo: '/dashboard', pathMatch: 'full' }
    ]
  }
];
