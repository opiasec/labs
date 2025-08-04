import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FilesDiffViewerComponent } from './files-diff-viewer.component';

describe('FilesDiffViewerComponent', () => {
  let component: FilesDiffViewerComponent;
  let fixture: ComponentFixture<FilesDiffViewerComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [FilesDiffViewerComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(FilesDiffViewerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
