import { CommonModule } from '@angular/common';
import { Component, inject, OnInit } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import {
  TeacherOnboardData,
  TeacherService,
} from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';

@Component({
  selector: 'shikshakul-staff-detail',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './staff-detail.component.html',
  styleUrl: './staff-detail.component.scss',
})
export class StaffDetailComponent implements OnInit {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private teacherService = inject(TeacherService);
  private snackbar = inject(SnackbarService);
  private fb = inject(FormBuilder);

  staffId: string | null = null;
  staff: any = null;
  loading = true;
  isEditing = false;
  editForm!: FormGroup;

  ngOnInit(): void {
    this.staffId = this.route.snapshot.paramMap.get('id');
    if (this.staffId) {
      this.loadStaff(this.staffId);
    }
  }

  loadStaff(id: string) {
    this.loading = true;
    this.teacherService.getTeacherById(id).subscribe({
      next: (data) => {
        this.staff = data;
        this.initForm(data);
        this.loading = false;
      },
      error: () => {
        this.snackbar.error('Error', 'Failed to load staff details.');
        this.router.navigate(['/staff']);
      },
    });
  }

  initForm(data: any) {
    this.editForm = this.fb.group({
      first_name: [data.first_name, Validators.required],
      last_name: [data.last_name, Validators.required],
      email: [data.email, [Validators.required, Validators.email]],
      mobile: [data.mobile, Validators.required],
      department: [
        data.profile_data?.employment?.department || '',
        Validators.required,
      ],
      designation: [
        data.profile_data?.employment?.designation || '',
        Validators.required,
      ],
    });
  }

  toggleEdit() {
    this.isEditing = !this.isEditing;
    if (!this.isEditing) this.initForm(this.staff);
  }

  saveChanges() {
    if (this.editForm.invalid || !this.staffId) return;

    const formVal = this.editForm.value;

    const payload: Partial<TeacherOnboardData> = {
      first_name: formVal.first_name,
      last_name: formVal.last_name,
      email: formVal.email,
      mobile: formVal.mobile,
      profile_data: {
        ...this.staff.profile_data,
        employment: {
          ...this.staff.profile_data?.employment,
          department: formVal.department,
          designation: formVal.designation,
        },
      },
    };

    this.teacherService.updateTeacher(this.staffId, payload).subscribe({
      next: (updated) => {
        this.staff = updated;
        this.isEditing = false;
        this.snackbar.success('Success', 'Profile updated successfully.');
      },
      error: () => this.snackbar.error('Error', 'Failed to update profile.'),
    });
  }

  deleteStaff() {
    if (
      confirm(
        'Are you sure you want to delete this staff record? This cannot be undone.',
      )
    ) {
      this.teacherService.deleteTeacher(this.staffId!).subscribe({
        next: () => {
          this.snackbar.success(
            'Deleted',
            'Staff record deleted successfully.',
          );
          this.router.navigate(['/staff']);
        },
        error: () => this.snackbar.error('Error', 'Failed to delete staff.'),
      });
    }
  }

  updateStatus(action: 'APPROVE' | 'REJECT') {
    this.teacherService.approveTeacher(this.staffId!, action).subscribe({
      next: (updatedStaff) => {
        this.staff = updatedStaff;
        this.snackbar.success(
          'Success',
          `Staff member ${action === 'APPROVE' ? 'Approved' : 'Rejected'}.`,
        );
      },
      error: () =>
        this.snackbar.error(
          'Error',
          `Failed to ${action.toLowerCase()} staff.`,
        ),
    });
  }
}
