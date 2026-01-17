import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { StaffPhotoComponent } from '../components/staff-photo/staff-photo.component';
import { StaffPersonalInfoComponent } from '../components/staff-personal-info/staff-personal-info.component';
import { StaffContactDetailsComponent } from '../components/staff-contact-details/staff-contact-details.component';
import { StaffEmploymentInfoComponent } from '../components/staff-employment-info/staff-employment-info.component';
import { StaffBankingPayrollComponent } from '../components/staff-banking-payroll/staff-banking-payroll.component';

@Component({
  selector: 'shikshakul-create-staff-page',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    StaffPhotoComponent,
    StaffPersonalInfoComponent,
    StaffContactDetailsComponent,
    StaffEmploymentInfoComponent,
    StaffBankingPayrollComponent,
  ],
  templateUrl: './create-staff-page.component.html',
  styleUrl: './create-staff-page.component.scss',
})
export class CreateStaffPageComponent {
  staffForm: FormGroup;

  constructor(private fb: FormBuilder) {
    this.staffForm = this.fb.group({
      photo: [null],
      personal: this.fb.group({
        firstName: ['', Validators.required],
        lastName: ['', Validators.required],
        dob: ['', Validators.required],
        gender: ['', Validators.required],
        fatherSpouseName: [''],
        aadharNumber: [
          '',
          [Validators.required, Validators.pattern('^[0-9]{12}$')],
        ],
      }),
      contact: this.fb.group({
        mobile: ['', [Validators.required, Validators.pattern('^[0-9]{10}$')]],
        email: ['', [Validators.required, Validators.email]],
        currentAddress: ['', Validators.required],
        state: ['', Validators.required],
        city: ['', Validators.required],
        pincode: ['', Validators.required],
        isPermanentSame: [false],
      }),
      employment: this.fb.group({
        employeeId: [''], // Auto-generated
        department: ['', Validators.required],
        designation: ['', Validators.required],
        dateOfJoining: ['', Validators.required],
      }),
      banking: this.fb.group({
        bankName: ['', Validators.required],
        accountNumber: ['', Validators.required],
        ifscCode: ['', Validators.required],
        panNumber: ['', Validators.required],
      }),
    });
  }

  getGroup(name: string): FormGroup {
    const control = this.staffForm.get(name);
    if (!control) throw new Error(`Form group '${name}' not found`);
    return control as FormGroup;
  }

  onSubmit() {
    if (this.staffForm.valid) {
      console.log('Staff Data:', this.staffForm.value);
    } else {
      this.staffForm.markAllAsTouched();
    }
  }

  onCancel() {
    console.log('Cancelled');
  }
}
