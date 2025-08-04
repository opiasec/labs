import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LabTimerComponent } from './lab-timer.component';

describe('LabTimerComponent', () => {
  let component: LabTimerComponent;
  let fixture: ComponentFixture<LabTimerComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [LabTimerComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(LabTimerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
