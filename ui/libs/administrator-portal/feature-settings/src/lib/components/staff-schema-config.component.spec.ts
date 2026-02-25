import { ComponentFixture, TestBed } from '@angular/core/testing';
import { StaffSchemaConfigComponent } from './staff-schema-config.component';

describe('StaffSchemaConfigComponent', () => {
  let component: StaffSchemaConfigComponent;
  let fixture: ComponentFixture<StaffSchemaConfigComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [StaffSchemaConfigComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(StaffSchemaConfigComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
