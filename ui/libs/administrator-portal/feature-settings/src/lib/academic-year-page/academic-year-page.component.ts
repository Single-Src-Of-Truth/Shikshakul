import { CommonModule } from '@angular/common';
import { Component, inject, OnInit, signal, ChangeDetectorRef } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { AcademicYearService, AcademicYear } from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';
import { AcademicYearStatsComponent } from '../components/academic-year-stats/academic-year-stats.component';
import { StudentSchemaConfigComponent } from '../components/student-schema-config/student-schema-config.component';
import { StaffSchemaConfigComponent } from '../components/staff-schema-config/staff-schema-config.component';
import { AdmissionSequenceConfigComponent } from '../components/admission-sequence-config/admission-sequence-config.component';

@Component({
  selector: 'shikshakul-academic-year-page',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    AcademicYearStatsComponent,
    StudentSchemaConfigComponent,
    StaffSchemaConfigComponent,
    AdmissionSequenceConfigComponent,
  ],
  templateUrl: './academic-year-page.component.html',
  styleUrl: './academic-year-page.component.scss',
})
export class AcademicYearPageComponent implements OnInit {
  private acadYearService = inject(AcademicYearService);
  private snackbar = inject(SnackbarService);
  private fb = inject(FormBuilder);
  private cdr = inject(ChangeDetectorRef);

  years = signal<AcademicYear[]>([]);
  loading = signal(false);
  showModal = signal(false);
  submitting = signal(false);
  
  yearForm: FormGroup;

  constructor() {
    this.yearForm = this.fb.group({
      name: ['', [Validators.required, Validators.pattern(/^\d{4}-\d{4}$/)]],
      start_date: ['', Validators.required],
      end_date: ['', Validators.required],
    });
  }

  ngOnInit(): void {
    this.loadYears();
  }

  loadYears() {
    this.loading.set(true);
    this.acadYearService.getAcademicYears().subscribe({
      next: (res: any) => {
        const data = Array.isArray(res) ? res : res?.data || [];
        this.years.set(data);
        this.loading.set(false);
        this.cdr.detectChanges();
      },
      error: () => {
        this.loading.set(false);
        this.snackbar.error('Error', 'Failed to load academic years.');
        this.cdr.detectChanges();
      }
    });
  }

  toggleModal() {
    this.showModal.set(!this.showModal());
    if (!this.showModal()) this.yearForm.reset();
  }

  onSubmit() {
    if (this.yearForm.invalid) return;

    this.submitting.set(true);
    this.acadYearService.createAcademicYear(this.yearForm.value).subscribe({
      next: () => {
        this.snackbar.success('Success', 'Academic year created successfully.');
        this.submitting.set(false);
        this.toggleModal();
        this.loadYears();
      },
      error: (err) => {
        this.submitting.set(false);
        this.snackbar.error('Error', err.error?.message || 'Failed to create academic year.');
      }
    });
  }

  activateYear(id: string) {
    this.acadYearService.updateAcademicYear(id, { is_current: true }).subscribe({
      next: () => {
        this.snackbar.success('Success', 'Academic year activated successfully.');
        this.loadYears();
      },
      error: (err) => {
        this.snackbar.error('Error', err.error?.message || 'Failed to activate year.');
      }
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
        this.snackbar.error('Error', err.error?.message || 'Failed to delete year.');
      }
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
