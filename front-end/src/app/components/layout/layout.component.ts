import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { NavbarComponent } from '../navbar/navbar.component';
import { TopbarComponent } from '../topbar/topbar.component';


@Component({
  selector: 'app-layout',
  standalone: true,
  imports: [NavbarComponent, RouterOutlet, TopbarComponent],
  templateUrl: './layout.component.html',
  styles: []
})
export class ApplayoutComponent {

}