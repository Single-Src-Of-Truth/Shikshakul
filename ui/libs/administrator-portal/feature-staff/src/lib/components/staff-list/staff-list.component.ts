import { CommonModule } from '@angular/common';
import { Component, inject, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { TeacherService } from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';

@Component({
  selector: 'shikshakul-staff-list',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './staff-list.component.html',
  styleUrl: './staff-list.component.scss',
})
export class StaffListComponent implements OnInit {
  private teacherService = inject(TeacherService);
  private router = inject(Router);
  private snackbar = inject(SnackbarService);

  staffList: any[] = [];
  loading = true;

  ngOnInit(): void {
    this.loadStaff();
  }

  loadStaff() {
    this.loading = true;
    this.teacherService.getTeachers().subscribe({
      next: (data) => {
        this.staffList = data;
        this.loading = false;
      },
      error: () => {
        this.loading = false;
        this.snackbar.error('Error', 'Failed to load staff directory.');
      },
    });
  }

  viewStaff(id: string) {
    this.router.navigate(['/staff', id]);
  }

  navigateToCreate() {
    this.router.navigate(['/staff/create']);
  }
}
