import { CommonModule, KeyValuePipe } from '@angular/common';
import { Component, inject, OnInit, ChangeDetectorRef } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { StudentService } from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';

@Component({
  selector: 'shikshakul-student-detail',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './student-detail.component.html',
  styleUrl: './student-detail.component.scss',
})
export class StudentDetailComponent implements OnInit {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private studentService = inject(StudentService);
  private snackbar = inject(SnackbarService);
  private fb = inject(FormBuilder);
  private cdr = inject(ChangeDetectorRef);

  studentId: string | null = null;
  student: any = null;
  loading = true;
  isEditing = false;

  editForm!: FormGroup;

  ngOnInit(): void {
    this.studentId = this.route.snapshot.paramMap.get('id');
    if (this.studentId) {
      this.loadStudent(this.studentId);
    } else {
      this.loading = false;
      this.snackbar.error('Error', 'No student ID provided.');
    }
  }

  loadStudent(id: string) {
    this.loading = true;
    this.cdr.detectChanges();

    this.studentService.getStudentById(id).subscribe({
      next: (res) => {
        const studentData = res?.data || res;
        this.student = studentData;

        if (studentData) {
          this.initForm(studentData);
        }

        this.loading = false;
        this.cdr.detectChanges();
      },
      error: (err) => {
        console.error("Failed to load student:", err);
        this.loading = false;
        this.snackbar.error('Error', 'Failed to load student details.');
        this.cdr.detectChanges();
        this.router.navigate(['/academics/students']);
      },
    });
  }

  initForm(data: any) {
    const profileKeys = Object.keys(data?.profile_data || {});
    const profileControls: any = {};

    for (const key of profileKeys) {
      profileControls[key] = [data?.profile_data[key]];
    }

    if (Object.keys(profileControls).length === 0) {
      profileControls['father_name'] = [''];
    }

    this.editForm = this.fb.group({
      first_name: [data?.first_name || '', Validators.required],
      last_name: [data?.last_name || '', Validators.required],
      email: [data?.email || '', Validators.email],
      mobile: [data?.mobile || '', [Validators.pattern('^[0-9]{10}$')]],
      profile_data: this.fb.group(profileControls),
    });
  }

  get profileDataControls() {
    return (this.editForm.get('profile_data') as FormGroup).controls;
  }

  objectKeys(obj: any): string[] {
    return obj ? Object.keys(obj) : [];
  }

  formatKey(key: any): string {
    return String(key).split('_').join(' ').toUpperCase();
  }

  toggleEdit() {
    this.isEditing = !this.isEditing;
    if (!this.isEditing && this.student) {
      this.initForm(this.student);
    }
    this.cdr.detectChanges();
  }

  saveChanges() {
    if (this.editForm.invalid || !this.studentId) return;

    const formVal = this.editForm.value;
    const payload = {
      ...this.student,
      first_name: formVal.first_name,
      last_name: formVal.last_name,
      email: formVal.email,
      mobile: formVal.mobile,
      profile_data: {
        ...(this.student?.profile_data || {}),
        ...formVal.profile_data,
      },
    };

    this.studentService.updateStudent(this.studentId, payload).subscribe({
      next: (updated) => {
        this.student = updated?.data || updated;
        this.isEditing = false;
        this.snackbar.success(
          'Updated',
          'Student details updated successfully.',
        );
        this.cdr.detectChanges();
      },
      error: (err) => {
        console.error("Failed to update student:", err);
        this.snackbar.error('Error', 'Failed to update student.');
        this.cdr.detectChanges();
      },
    });
  }
}
