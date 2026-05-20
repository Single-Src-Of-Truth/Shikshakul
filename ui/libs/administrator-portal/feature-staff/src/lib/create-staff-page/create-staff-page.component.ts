import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
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
import { Router } from '@angular/router';
import {
  TeacherService,
  TeacherOnboardData,
} from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';

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
  private fb = inject(FormBuilder);
  private teacherService = inject(TeacherService);
  private snackbar = inject(SnackbarService);
  private router = inject(Router);

  staffForm: FormGroup;
  isSubmitting = false;

  constructor() {
    this.staffForm = this.fb.group({
      photo: [null],
      personal: this.fb.group({
        firstName: ['', Validators.required],
        lastName: ['', Validators.required],
        dob: ['', Validators.required],
        gender: ['MALE', Validators.required],
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
        employeeId: [''],
        department: ['', Validators.required],
        designation: ['', Validators.required],
        dateOfJoining: [
          new Date().toISOString().split('T')[0],
          Validators.required,
        ],
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
    return this.staffForm.get(name) as FormGroup;
  }

  onSubmit() {
    if (this.staffForm.invalid) {
      this.staffForm.markAllAsTouched();
      this.snackbar.error(
        'Validation Error',
        'Please check the form for missing fields.',
      );
      return;
    }

    this.isSubmitting = true;
    const formVal = this.staffForm.value;

    const payload: TeacherOnboardData = {
      first_name: formVal.personal.firstName,
      last_name: formVal.personal.lastName,
      email: formVal.contact.email,
      mobile: formVal.contact.mobile,

      profile_data: {
        dob: new Date(formVal.personal.dob).toISOString(),
        gender: formVal.personal.gender,
        father_spouse_name: formVal.personal.fatherSpouseName,
        aadhar_number: formVal.personal.aadharNumber,

        address: {
          line1: formVal.contact.currentAddress,
          city: formVal.contact.city,
          state: formVal.contact.state,
          pincode: formVal.contact.pincode,
        },

        employment: {
          department: formVal.employment.department,
          designation: formVal.employment.designation,
          date_of_joining: new Date(
            formVal.employment.dateOfJoining,
          ).toISOString(),
        },

        banking: {
          bank_name: formVal.banking.bankName,
          account_number: formVal.banking.accountNumber,
          ifsc: formVal.banking.ifscCode,
          pan: formVal.banking.panNumber,
        },
      },
      documents: {},
    };

    this.teacherService.onboardTeacher(payload).subscribe({
      next: () => {
        this.snackbar.success(
          'Success',
          'Staff member onboarded successfully!',
        );
        this.isSubmitting = false;
        this.router.navigate(['/staff/list']);
      },
      error: (err) => {
        console.error(err);
        this.snackbar.error('Error', 'Failed to onboard staff member.');
        this.isSubmitting = false;
      },
    });
  }

  onCancel() {
    this.router.navigate(['/staff/list']);
    console.log('Cancelled');
  }
}
