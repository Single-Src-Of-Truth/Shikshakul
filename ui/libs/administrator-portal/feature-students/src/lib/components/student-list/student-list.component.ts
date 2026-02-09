import { CommonModule } from '@angular/common';
import { Component, inject, OnInit } from '@angular/core';
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
    this.classService.getClasses().subscribe((data) => {
      this.classes = data.sort((a, b) => a.sort_order - b.sort_order);
    });
  }

  loadStudents() {
    this.loading = true;
    const filters: any = {};
    if (this.selectedClassId) filters.class_id = this.selectedClassId;

    this.studentService.getStudents(filters).subscribe({
      next: (data) => {
        this.students = data;
        this.loading = false;
      },
      error: () => {
        this.loading = false;
        this.snackbar.error('Error', 'Failed to load students.');
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
