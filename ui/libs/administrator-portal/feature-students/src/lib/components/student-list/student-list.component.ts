import { CommonModule } from '@angular/common';
import { Component, inject, OnInit, ChangeDetectorRef } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import {
  ClassGrade,
  ClassManagementService,
  StudentService,
} from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';

@Component({
  selector: 'shikshakul-student-list',
  standalone: true,
  imports: [CommonModule, RouterModule, FormsModule],
  templateUrl: './student-list.component.html',
  styleUrl: './student-list.component.scss',
})
export class StudentListComponent implements OnInit {
  private studentService = inject(StudentService);
  private classService = inject(ClassManagementService);
  private router = inject(Router);
  private snackbar = inject(SnackbarService);
  private cdr = inject(ChangeDetectorRef);

  students: any[] = [];
  classes: ClassGrade[] = [];
  loading = true;

  selectedClassId = '';
  searchQuery = '';

  ngOnInit(): void {
    this.loadClasses();
    this.loadStudents();
  }

  loadClasses() {
    this.classService.getClasses().subscribe((res: any) => {
      const data = Array.isArray(res) ? res : (res?.data || []);
      this.classes = data.sort((a: any, b: any) => (a.sort_order || 0) - (b.sort_order || 0));
      this.cdr.detectChanges();
    });
  }

  loadStudents() {
    this.loading = true;
    this.cdr.detectChanges();
    const filters: any = {};
    if (this.selectedClassId) filters.class_id = this.selectedClassId;

    this.studentService.getStudents(filters).subscribe({
      next: (res: any) => {
        this.students = Array.isArray(res) ? res : (res?.data || []);
        this.loading = false;
        this.cdr.detectChanges();
      },
      error: () => {
        this.loading = false;
        this.snackbar.error('Error', 'Failed to load students.');
        this.cdr.detectChanges();
      },
    });
  }

  onClassFilterChange() {
    this.loadStudents();
  }

  viewStudent(id: string) {
    this.router.navigate(['/academics/students', id]);
  }

  navigateToCreate() {
    this.router.navigate(['/academics/students/onboard']);
  }
}
