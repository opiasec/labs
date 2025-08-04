import { TestBed } from '@angular/core/testing';

import { LabSessionService } from './lab-session.service';

describe('LabSessionService', () => {
  let service: LabSessionService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(LabSessionService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
