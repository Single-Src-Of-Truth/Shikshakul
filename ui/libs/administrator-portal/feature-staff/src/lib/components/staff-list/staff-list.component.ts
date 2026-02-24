import { CommonModule } from '@angular/common';
import { Component, inject, OnInit, ChangeDetectorRef } from '@angular/core';
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
  private cdr = inject(ChangeDetectorRef);

  staffList: any[] = [];
  loading = true;

  ngOnInit(): void {
    this.loadStaff();
  }

  loadStaff() {
    this.loading = true;
    this.cdr.detectChanges();
    this.teacherService.getTeachers().subscribe({
      next: (res: any) => {
        this.staffList = Array.isArray(res) ? res : (res?.data || []);
        this.loading = false;
        this.cdr.detectChanges();
      },
      error: () => {
        this.loading = false;
        this.snackbar.error('Error', 'Failed to load staff directory.');
        this.cdr.detectChanges();
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
