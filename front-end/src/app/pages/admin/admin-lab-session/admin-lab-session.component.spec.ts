import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AdminLabSessionComponent } from './admin-lab-session.component';

describe('AdminLabSessionComponent', () => {
  let component: AdminLabSessionComponent;
  let fixture: ComponentFixture<AdminLabSessionComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AdminLabSessionComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(AdminLabSessionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
