import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AdminLabsSessionsComponent } from './admin-labs-sessions.component';

describe('AdminLabsSessionsComponent', () => {
  let component: AdminLabsSessionsComponent;
  let fixture: ComponentFixture<AdminLabsSessionsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AdminLabsSessionsComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(AdminLabsSessionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
