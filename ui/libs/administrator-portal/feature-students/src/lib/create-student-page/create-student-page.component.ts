import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
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
export class CreateStudentPageComponent {
  studentForm: FormGroup;

  constructor(private fb: FormBuilder) {
    this.studentForm = this.fb.group({
      photo: [null],
      admission: this.fb.group({
        admissionNumber: ['ADM-2024-0892', Validators.required],
        dateOfAdmission: [new Date(), Validators.required],
        academicYear: ['2024-2025', Validators.required],
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
        board: ['', Validators.required],
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
        email: ['', Validators.email],
      }),
      address: this.fb.group({
        addressLine1: ['', Validators.required],
        city: ['', Validators.required],
        state: ['', Validators.required],
        pincode: ['', [Validators.required, Validators.pattern('^[0-9]{6}$')]],
      }),
    });
  }

  getGroup(name: string): FormGroup {
    const control = this.studentForm.get(name);

    // Safety check (optional but good practice)
    if (!control) {
      throw new Error(`Form group '${name}' not found`);
    }

    return control as FormGroup;
  }

  onSubmit() {
    if (this.studentForm.valid) {
      console.log('Form Submitted!', this.studentForm.value);
      // Here you would call a service to send data to the API
    } else {
      console.log('Form is invalid. Please check the fields.');
      this.studentForm.markAllAsTouched(); // Trigger validation messages
    }
  }

  onCancel() {
    // Navigate back or reset form
    console.log('Cancelled');
  }
}
