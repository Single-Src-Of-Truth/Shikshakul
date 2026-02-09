import { CommonModule } from '@angular/common';
import { Component, inject, OnDestroy, OnInit } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { StudentPhotoComponent } from '../components/student-photo/student-photo.component';
import { AddressDetailsComponent } from '../components/address-details/address-details.component';
import { AdmissionDetailsComponent } from '../components/admission-details/admission-details.component';
import { PersonalDetailsComponent } from '../components/personal-details/personal-details.component';
import { AcademicDetailsComponent } from '../components/academic-details/academic-details.component';
import { GuardianDetailsComponent } from '../components/guardian-details/guardian-details.component';
import { Router } from '@angular/router';
import {
  StudentService,
  StudentOnboardData,
  ClassGrade,
  Section,
  ClassManagementService,
  AcademicYearService,
} from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';
import { Subject, takeUntil } from 'rxjs';

@Component({
  selector: 'shikshakul-create-student-page',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    StudentPhotoComponent,
    AdmissionDetailsComponent,
    PersonalDetailsComponent,
    AcademicDetailsComponent,
    GuardianDetailsComponent,
    AddressDetailsComponent,
  ],
  templateUrl: './create-student-page.component.html',
  styleUrl: './create-student-page.component.scss',
})
export class CreateStudentPageComponent implements OnInit, OnDestroy {
  private fb = inject(FormBuilder);
  private router = inject(Router);

  // Services
  private studentService = inject(StudentService);
  private classService = inject(ClassManagementService);
  private academicYearService = inject(AcademicYearService);
  private snackbar = inject(SnackbarService);

  studentForm: FormGroup;
  isSubmitting = false;
  destroy$ = new Subject<void>();

  classes: ClassGrade[] = [];
  sections: Section[] = [];

  activeAcademicYearId: string | null = null;

  constructor() {
    this.studentForm = this.fb.group({
      photo: [null],
      admission: this.fb.group({
        admissionNumber: ['', Validators.required],
        dateOfAdmission: [
          new Date().toISOString().split('T')[0],
          Validators.required,
        ],
        academicYear: [{ value: '', disabled: true }, Validators.required], // Disabled so user can't break it
      }),
      personal: this.fb.group({
        firstName: ['', Validators.required],
        lastName: ['', Validators.required],
        dob: ['', Validators.required],
        gender: ['', Validators.required],
        bloodGroup: [''],
        category: ['General'],
        aadharNumber: [''],
      }),
      academic: this.fb.group({
        board: ['CBSE', Validators.required],
        class: ['', Validators.required],
        section: ['', Validators.required],
        rollNumber: [''],
      }),
      guardian: this.fb.group({
        fatherName: ['', Validators.required],
        motherName: ['', Validators.required],
        primaryMobile: [
          '',
          [Validators.required, Validators.pattern('^[0-9]{10}$')],
        ],
        email: ['', [Validators.email]],
      }),
      address: this.fb.group({
        addressLine1: ['', Validators.required],
        city: ['', Validators.required],
        state: ['', Validators.required],
        pincode: ['', [Validators.required, Validators.pattern('^[0-9]{6}$')]],
      }),
    });
  }

  ngOnInit() {
    this.loadInitialData();
    this.setupClassChangeListener();
  }

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }

  private loadInitialData() {
    this.academicYearService
      .getCurrentAcademicYear()
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: (year) => {
          if (year) {
            this.activeAcademicYearId = year.id;
            this.studentForm.get('admission.academicYear')?.setValue(year.name);
          }
        },
        error: () =>
          this.snackbar.error('Error', 'Could not fetch active academic year.'),
      });

    this.classService
      .getClasses()
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: (data) => {
          this.classes = data.sort((a, b) => a.sort_order - b.sort_order);
        },
        error: () => this.snackbar.error('Error', 'Could not load classes.'),
      });
  }

  private setupClassChangeListener() {
    const classControl = this.studentForm.get('academic.class');
    const sectionControl = this.studentForm.get('academic.section');

    classControl?.valueChanges
      .pipe(takeUntil(this.destroy$))
      .subscribe((classId) => {
        // Reset section when class changes
        sectionControl?.setValue('');
        this.sections = [];

        if (classId) {
          this.fetchSections(classId);
        }
      });
  }

  private fetchSections(classId: string) {
    // Disable section dropdown while loading
    this.studentForm.get('academic.section')?.disable();

    this.classService
      .getSectionsByClass(classId)
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: (data) => {
          this.sections = data;
          this.studentForm.get('academic.section')?.enable();
        },
        error: () => {
          this.snackbar.error(
            'Error',
            'Failed to load sections for selected class.',
          );
          this.studentForm.get('academic.section')?.enable();
        },
      });
  }

  getGroup(name: string): FormGroup {
    return this.studentForm.get(name) as FormGroup;
  }

  onSubmit() {
    if (this.studentForm.invalid) {
      this.studentForm.markAllAsTouched();
      this.snackbar.error(
        'Validation Error',
        'Please check the form for missing fields.',
      );
      return;
    }

    this.isSubmitting = true;
    const formValue = this.studentForm.getRawValue(); // getRawValue includes disabled fields (like Year)

    const payload: StudentOnboardData = {
      // Use the stored ID, not the display name from the disabled input
      academic_year_id: this.activeAcademicYearId!,
      class_id: formValue.academic.class,
      section_id: formValue.academic.section,
      first_name: formValue.personal.firstName,
      last_name: formValue.personal.lastName,
      dob: new Date(formValue.personal.dob).toISOString(),
      gender: formValue.personal.gender,
      email: formValue.guardian.email,
      mobile: formValue.guardian.primaryMobile,
      profile_data: {
        admission_number: formValue.admission.admissionNumber,
        date_of_admission: formValue.admission.dateOfAdmission,
        blood_group: formValue.personal.bloodGroup,
        category: formValue.personal.category,
        aadhar_number: formValue.personal.aadharNumber,
        father_name: formValue.guardian.fatherName,
        mother_name: formValue.guardian.motherName,
        board: formValue.academic.board,
        roll_number: formValue.academic.rollNumber,
        address: {
          line1: formValue.address.addressLine1,
          city: formValue.address.city,
          state: formValue.address.state,
          pincode: formValue.address.pincode,
        },
      },
      documents: {},
    };

    this.studentService.onboardStudent(payload).subscribe({
      next: () => {
        this.snackbar.success(
          'Success',
          `Student ${payload.first_name} admitted successfully!`,
        );
        this.router.navigate(['/academics/students']);
      },
      error: (err) => {
        console.error(err);
        this.snackbar.error('Error', 'Failed to admit student.');
        this.isSubmitting = false;
      },
    });
  }

  onCancel() {
    this.router.navigate(['/academics/students']);
  }
}
