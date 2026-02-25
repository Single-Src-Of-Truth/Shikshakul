import { ComponentFixture, TestBed } from '@angular/core/testing';
import { AdmissionSequenceConfigComponent } from './admission-sequence-config.component';

describe('AdmissionSequenceConfigComponent', () => {
  let component: AdmissionSequenceConfigComponent;
  let fixture: ComponentFixture<AdmissionSequenceConfigComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AdmissionSequenceConfigComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(AdmissionSequenceConfigComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
