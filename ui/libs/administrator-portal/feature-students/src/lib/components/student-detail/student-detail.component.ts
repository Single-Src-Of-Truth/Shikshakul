import { CommonModule } from '@angular/common';
import { Component, inject, OnInit } from '@angular/core';
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

  studentId: string | null = null;
  student: any = null;
  loading = true;
  isEditing = false;

  editForm!: FormGroup;

  ngOnInit(): void {
    this.studentId = this.route.snapshot.paramMap.get('id');
    if (this.studentId) {
      this.loadStudent(this.studentId);
    }
  }

  loadStudent(id: string) {
    this.loading = true;
    this.studentService.getStudentById(id).subscribe({
      next: (data) => {
        this.student = data;
        this.initForm(data);
        this.loading = false;
      },
      error: () => {
        this.loading = false;
        this.snackbar.error('Error', 'Failed to load student details.');
        this.router.navigate(['/academics/students']);
      },
    });
  }

  initForm(data: any) {
    this.editForm = this.fb.group({
      first_name: [data.first_name, Validators.required],
      last_name: [data.last_name, Validators.required],
      email: [data.email, Validators.email],
      mobile: [data.mobile, [Validators.pattern('^[0-9]{10}$')]],
      father_name: [data.profile_data?.father_name || ''],
    });
  }

  toggleEdit() {
    this.isEditing = !this.isEditing;
    if (!this.isEditing) {
      this.initForm(this.student);
    }
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
        ...this.student.profile_data,
        father_name: formVal.father_name,
      },
    };

    this.studentService.updateStudent(this.studentId, payload).subscribe({
      next: (updated) => {
        this.student = updated;
        this.isEditing = false;
        this.snackbar.success(
          'Updated',
          'Student details updated successfully.',
        );
      },
      error: () => this.snackbar.error('Error', 'Failed to update student.'),
    });
  }
}
