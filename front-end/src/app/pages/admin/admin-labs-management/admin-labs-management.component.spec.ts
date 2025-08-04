import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AdminLabsManagementComponent } from './admin-labs-management.component';

describe('AdminLabsManagementComponent', () => {
  let component: AdminLabsManagementComponent;
  let fixture: ComponentFixture<AdminLabsManagementComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AdminLabsManagementComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(AdminLabsManagementComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
