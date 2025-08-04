import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LabSessionComponent } from './lab-session.component';

describe('LabSessionComponent', () => {
  let component: LabSessionComponent;
  let fixture: ComponentFixture<LabSessionComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [LabSessionComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(LabSessionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
