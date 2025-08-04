import { Component } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { TieredMenuModule } from 'primeng/tieredmenu';
import { CardModule } from 'primeng/card';

@Component({
  selector: 'app-settings',
  standalone: true,
  imports: [
    RouterModule,
    TieredMenuModule,
    CardModule
  ],
  templateUrl: './settings.component.html',
  styleUrl: './settings.component.css'
})
export class SettingsComponent {

  settingsItems = [
    { label: 'Billing', icon: 'pi pi-credit-card', command: () => this.redirectTo('billing') },
    { label: 'Profile', icon: 'pi pi-user', command: () => this.redirectTo('profile') },
    { label: 'Configurations', icon: 'pi pi-cog', command: () => this.redirectTo('configurations') },
  ];

  constructor(private router: Router) {}

  redirectTo(page: string) {
    this.router.navigate([`/settings/${page}`]);
  }
}
