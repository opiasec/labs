import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LoginNewComponent } from './login-new.component';

describe('LoginNewComponent', () => {
  let component: LoginNewComponent;
  let fixture: ComponentFixture<LoginNewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [LoginNewComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(LoginNewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
