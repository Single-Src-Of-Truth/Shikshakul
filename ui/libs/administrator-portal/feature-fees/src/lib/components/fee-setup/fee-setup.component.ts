import { CommonModule } from '@angular/common';
import { Component, inject, OnInit, ChangeDetectorRef } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import {
  FeeService,
  ClassManagementService,
  AcademicYearService,
} from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';
import { forkJoin } from 'rxjs';

@Component({
  selector: 'shikshakul-fee-setup',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, FormsModule],
  templateUrl: './fee-setup.component.html',
  styleUrl: './fee-setup.component.scss',
})
export class FeeSetupComponent implements OnInit {
  private feeService = inject(FeeService);
  private classService = inject(ClassManagementService);
  private yearService = inject(AcademicYearService);
  private fb = inject(FormBuilder);
  private snackbar = inject(SnackbarService);
  private cdr = inject(ChangeDetectorRef);

  classes: any[] = [];
  feeHeads: any[] = [];
  activeYearId = '';
  currentStructure: any[] = [];

  selectedClassId = '';
  loading = false;

  showAddModal = false;
  showHeadModal = false;

  headForm: FormGroup;
  structureForm: FormGroup;

  constructor() {
    this.headForm = this.fb.group({
      name: ['', Validators.required],
      type: ['RECURRING', Validators.required],
    });

    this.structureForm = this.fb.group({
      fee_head_id: ['', Validators.required],
      amount: [0, [Validators.required, Validators.min(1)]],
      frequency: ['MONTHLY', Validators.required],
      due_date_day: [10],
    });
  }

  ngOnInit(): void {
    this.loadInitialData();
  }

  loadInitialData() {
    this.loading = true;
    this.cdr.detectChanges();
    forkJoin({
      year: this.yearService.getCurrentAcademicYear(),
      classes: this.classService.getClasses(),
      heads: this.feeService.getFeeHeads(),
    }).subscribe({
      next: (res: any) => {
        const yearData = res.year?.data || res.year;
        this.activeYearId = yearData?.id || '';
        this.classes = Array.isArray(res.classes) ? res.classes : (res.classes?.data || []);
        this.feeHeads = Array.isArray(res.heads) ? res.heads : (res.heads?.data || []);
        this.loading = false;
        this.cdr.detectChanges();
      },
      error: () => {
        this.loading = false;
        this.cdr.detectChanges();
      }
    });
  }

  addFeeHead() {
    if (this.headForm.invalid) return;
    this.feeService.createFeeHead(this.headForm.value).subscribe({
      next: (head) => {
        this.feeHeads.push(head);
        this.headForm.reset({ type: 'RECURRING' });
        this.snackbar.success('Success', 'Fee Head Added');
      },
      error: () => this.snackbar.error('Error', 'Failed to add fee head'),
    });
  }

  onClassChange() {
    if (!this.selectedClassId || !this.activeYearId) return;
    this.loading = true;
    this.cdr.detectChanges();
    this.feeService
      .getFeeStructure(this.selectedClassId, this.activeYearId)
      .subscribe({
        next: (res: any) => {
          this.currentStructure = Array.isArray(res) ? res : (res?.data || []);
          this.loading = false;
          this.cdr.detectChanges();
        },
        error: () => {
          this.loading = false;
          this.cdr.detectChanges();
        }
      });
  }

  addStructure() {
    if (this.structureForm.invalid || !this.selectedClassId) return;

    const payload = {
      ...this.structureForm.value,
      class_id: this.selectedClassId,
      academic_year_id: this.activeYearId,
    };

    this.feeService.createFeeStructure(payload).subscribe({
      next: () => {
        this.snackbar.success('Success', 'Fee assigned to class');
        this.onClassChange();
        this.structureForm.reset({ frequency: 'MONTHLY', due_date_day: 10 });
        this.showAddModal = false;
      },
      error: () => this.snackbar.error('Error', 'Failed to assign fee'),
    });
  }

  deleteStructure(id: string) {
    if (confirm('Remove this fee?')) {
      this.feeService
        .deleteFeeStructure(id)
        .subscribe(() => this.onClassChange());
    }
  }

  getIconForFee(name = ''): string {
    const n = name.toLowerCase();
    if (n.includes('bus') || n.includes('transport')) return 'directions_bus';
    if (n.includes('lab') || n.includes('computer')) return 'computer';
    if (n.includes('exam')) return 'assignment';
    if (n.includes('library')) return 'menu_book';
    if (n.includes('sport')) return 'sports_soccer';
    return 'account_balance_wallet'; // Default
  }

  calculateTotal(): number {
    return this.currentStructure.reduce((sum, item) => sum + item.amount, 0);
  }
}
