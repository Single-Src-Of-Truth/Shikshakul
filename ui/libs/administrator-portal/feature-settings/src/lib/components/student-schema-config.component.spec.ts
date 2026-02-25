import { ComponentFixture, TestBed } from '@angular/core/testing';
import { StudentSchemaConfigComponent } from './student-schema-config.component';

describe('StudentSchemaConfigComponent', () => {
  let component: StudentSchemaConfigComponent;
  let fixture: ComponentFixture<StudentSchemaConfigComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [StudentSchemaConfigComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(StudentSchemaConfigComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
