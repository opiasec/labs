import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MarkdownClipboardButtonComponent } from './markdown-clipboard-button.component';

describe('MarkdownClipboardButtonComponent', () => {
  let component: MarkdownClipboardButtonComponent;
  let fixture: ComponentFixture<MarkdownClipboardButtonComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MarkdownClipboardButtonComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(MarkdownClipboardButtonComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
