import { Component, Input, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CardModule } from 'primeng/card';

@Component({
  selector: 'app-simple-card',
  imports: [CommonModule, CardModule],
  templateUrl: './simple-card.component.html',
  styleUrl: './simple-card.component.css'
})
export class SimpleCardComponent implements OnInit {
  @Input() title: string = '';
  @Input() details: string = '';
  @Input() icon: string = '';

  ngOnInit() {

  }

  hasIcon(): boolean {
    return !!this.icon && this.icon.trim().length > 0;
  }
}
