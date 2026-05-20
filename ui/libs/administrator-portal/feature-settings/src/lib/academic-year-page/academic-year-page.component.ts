import { CommonModule } from '@angular/common';
import {
  Component,
  inject,
  OnInit,
  signal,
  ChangeDetectorRef,
} from '@angular/core';
import {
  AcademicYearService,
  AcademicYear,
} from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';
import { AcademicYearStatsComponent } from '../components/academic-year-stats/academic-year-stats.component';
import { StudentSchemaConfigComponent } from '../components/student-schema-config/student-schema-config.component';
import { StaffSchemaConfigComponent } from '../components/staff-schema-config/staff-schema-config.component';
import { AdmissionSequenceConfigComponent } from '../components/admission-sequence-config/admission-sequence-config.component';
import { AcademicYearFormComponent } from '../components/academic-year-form/academic-year-form.component';

@Component({
  selector: 'shikshakul-academic-year-page',
  standalone: true,
  imports: [
    CommonModule,
    AcademicYearStatsComponent,
    StudentSchemaConfigComponent,
    StaffSchemaConfigComponent,
    AdmissionSequenceConfigComponent,
    AcademicYearFormComponent,
  ],
  templateUrl: './academic-year-page.component.html',
  styleUrl: './academic-year-page.component.scss',
})
export class AcademicYearPageComponent implements OnInit {
  private acadYearService = inject(AcademicYearService);
  private snackbar = inject(SnackbarService);
  private cdr = inject(ChangeDetectorRef);

  years = signal<AcademicYear[]>([]);
  loading = signal(false);
  showDrawer = signal(false);

  ngOnInit(): void {
    this.loadYears();
  }

  loadYears() {
    this.loading.set(true);
    this.acadYearService.getAcademicYears().subscribe({
      next: (res: any) => {
        let data = Array.isArray(res) ? res : res?.data || [];
        // Normalize IDs if backend returns academic_year_id instead of id
        data = data.map((year: any) => ({
          ...year,
          id: year.id || year.academic_year_id,
        }));
        this.years.set(data);
        this.loading.set(false);
        this.cdr.detectChanges();
      },
      error: () => {
        this.loading.set(false);
        this.snackbar.error('Error', 'Failed to load academic years.');
        this.cdr.detectChanges();
      },
    });
  }

  toggleDrawer() {
    this.showDrawer.set(!this.showDrawer());
  }

  activateYear(id: string) {
    if (!id) {
      this.snackbar.error('Error', 'Invalid Academic Year ID');
      return;
    }
    this.acadYearService
      .updateAcademicYear(id, { is_current: true })
      .subscribe({
        next: () => {
          this.snackbar.success(
            'Success',
            'Academic year activated successfully.',
          );
          this.loadYears();
        },
        error: (err) => {
          this.snackbar.error(
            'Error',
            err.error?.message || 'Failed to activate year.',
          );
        },
      });
  }

  deleteYear(id: string) {
    if (!confirm('Are you sure you want to delete this academic year?')) return;

    this.acadYearService.deleteAcademicYear(id).subscribe({
      next: () => {
        this.snackbar.success('Success', 'Academic year deleted.');
        this.loadYears();
      },
      error: (err) => {
        this.snackbar.error(
          'Error',
          err.error?.message || 'Failed to delete year.',
        );
      },
    });
  }

  getStatusClass(year: AcademicYear): string {
    if (year.is_current) return 'active';
    const today = new Date();
    const startDate = new Date(year.start_date);
    if (startDate > today) return 'upcoming';
    return 'archived';
  }

  getStatusLabel(year: AcademicYear): string {
    if (year.is_current) return 'Active';
    const today = new Date();
    const startDate = new Date(year.start_date);
    if (startDate > today) return 'Upcoming';
    return 'Archived';
  }
}
