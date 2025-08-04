import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LabsSessionsComponent } from './labs-sessions.component';

describe('LabsSessionsComponent', () => {
  let component: LabsSessionsComponent;
  let fixture: ComponentFixture<LabsSessionsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [LabsSessionsComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(LabsSessionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
