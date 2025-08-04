import { Component } from '@angular/core';
import { MenuItem } from 'primeng/api';
import { MenuModule } from 'primeng/menu';
import { ButtonModule } from 'primeng/button';
import { PanelMenuModule } from 'primeng/panelmenu';
import { DividerModule } from 'primeng/divider';
import { CardModule } from 'primeng/card';
import { TieredMenuModule } from 'primeng/tieredmenu';
import { Router } from '@angular/router';
import { OnInit, OnDestroy } from '@angular/core';
import { LabSessionService } from '../../services/lab-session.service';
import { Subscription } from 'rxjs';
import { LabTimerComponent } from '../lab-timer/lab-timer.component';
import { TagModule } from 'primeng/tag';
import { CommonModule } from '@angular/common';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-navbar',
  imports: [MenuModule, ButtonModule, PanelMenuModule, DividerModule, CardModule, TieredMenuModule, LabTimerComponent, TagModule, CommonModule],
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.css'
})
export class NavbarComponent implements OnInit, OnDestroy {
  private sub: Subscription = new Subscription();
  endTime: string | null = null;
  labSlug: string = '';
  namespace: string = '';
  showSessionTimer: boolean = false;
  isAdmin: boolean = false;

  navBarItems: MenuItem[] = [
    {
      label: 'Dashboard',
      routerLink: '/dashboard'
    },
    {
      label: 'Labs',
      routerLink: '/labs'
    },
    {
      label: 'Sessions',
      routerLink: '/sessions'
    }
  ];
  adminNavBarItems: MenuItem[] = [
    {
      label: 'Sessions View',
      routerLink: '/admin/sessions'
    },
    {
      label: 'Labs Management',
      routerLink: '/admin/labs/management'
    }
  ];

  constructor(
    private router: Router,
    private labSessionService: LabSessionService,
    private authService: AuthService
  ) { }

  ngOnInit() {
    this.updateActiveLab();
    this.sub = this.labSessionService.getLabEndTimeObservable()
      .subscribe(end => {
        this.endTime = end;
        const activeLab = JSON.parse(localStorage.getItem('activeLab') || '{}');
        if (activeLab) {
          console.log("Active lab found:", activeLab);
          this.labSlug = activeLab.labSlug || '';
          this.namespace = activeLab.namespace || '';
        }
        if (this.endTime) {
          this.showSessionTimer = true;
        }
      });

    this.sub.add(
      this.authService.isAdmin().subscribe(isAdmin => {
        this.isAdmin = isAdmin;
      })
    );
  }

  private updateAdminMenu(isAdmin: boolean) {
    this.navBarItems = this.navBarItems.filter(item => item.label !== 'Admin');

    if (isAdmin) {
      this.navBarItems.push({
        label: 'Admin',
        items: [
          {
            label: 'Sessions View',
            routerLink: '/admin/sessions'
          }
        ]
      });
    }
  }

  updateActiveLab() {
    const activeLab = JSON.parse(localStorage.getItem('activeLab') || '{}');
    if (!activeLab || !activeLab.labSlug) {
      console.log("No active lab found or lab slug does not match.");
      return
    }
    this.labSessionService.setLabEndTime(activeLab.expiresAt);
    return
  }

  isUserAdmin(): boolean {
    return this.isAdmin;
  }

  onSessionEnd() {
    this.showSessionTimer = false;
  }

  ngOnDestroy() {
    this.sub.unsubscribe();
  }

  navigateToHome() {
    this.router.navigate(['/']);
  }
}