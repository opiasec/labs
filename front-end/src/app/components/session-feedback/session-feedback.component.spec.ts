import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SessionFeedbackComponent } from './session-feedback.component';

describe('SessionFeedbackComponent', () => {
  let component: SessionFeedbackComponent;
  let fixture: ComponentFixture<SessionFeedbackComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SessionFeedbackComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SessionFeedbackComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
