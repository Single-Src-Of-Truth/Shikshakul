import { CommonModule } from '@angular/common';
import {
  Component,
  inject,
  OnInit,
  ChangeDetectorRef,
  Output,
  EventEmitter,
} from '@angular/core';
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
import { SelectComponent, SelectOption } from '@shikshakul/shared/ui/select';
import { forkJoin } from 'rxjs';

@Component({
  selector: 'shikshakul-fee-setup',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, FormsModule, SelectComponent],
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

  @Output() loadingState = new EventEmitter<boolean>();

  classes: any[] = [];
  feeHeads: any[] = [];
  activeYearId = '';
  currentStructure: any[] = [];

  selectedClassId = '';
  loading = false;

  // Select Options Helpers
  get classOptions(): SelectOption[] {
    return this.classes
      .map((c) => ({
        label: c.name,
        value: this.getEntityId(c, ['id', 'class_id', '_id']),
      }))
      .filter((option) => !!option.value);
  }

  get feeHeadOptions(): SelectOption[] {
    return this.feeHeads
      .map((h) => ({
        label: h.name,
        value: this.getEntityId(h, ['id', 'fee_head_id', '_id']),
      }))
      .filter((option) => !!option.value);
  }

  get paymentTypeOptions(): SelectOption[] {
    return [
      { label: 'Recurring', value: 'RECURRING' },
      { label: 'One Time', value: 'ONE_TIME' },
      { label: 'Term Wise', value: 'TERM_WISE' },
      { label: 'Admission Fee', value: 'ADMISSION' },
    ];
  }

  get frequencyOptions(): SelectOption[] {
    return [
      { label: 'Monthly', value: 'MONTHLY' },
      { label: 'Quarterly', value: 'QUARTERLY' },
      { label: 'Half Yearly', value: 'HALF_YEARLY' },
      { label: 'Yearly', value: 'YEARLY' },
      { label: 'One Time', value: 'ONE_TIME' },
    ];
  }

  getSelectedClassName(): string {
    const cls = this.classes.find(
      (c) =>
        this.getEntityId(c, ['id', 'class_id', '_id']) === this.selectedClassId,
    );
    return cls ? cls.name : '';
  }

  showAddModal = false;
  showHeadModal = false;
  editingStructure: any = null;
  editingHead: any = null;

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
    this.loadingState.emit(true);
    this.cdr.detectChanges();
    forkJoin({
      year: this.yearService.getCurrentAcademicYear(),
      classes: this.classService.getClasses(),
      heads: this.feeService.getFeeHeads(),
    }).subscribe({
      next: (res: any) => {
        const yearData = res.year?.data || res.year;
        const normalizedYear = Array.isArray(yearData) ? yearData[0] : yearData;
        this.activeYearId = this.getEntityId(normalizedYear, [
          'id',
          'academic_year_id',
          '_id',
        ]);
        this.classes = Array.isArray(res.classes)
          ? res.classes
          : res.classes?.data || [];
        this.feeHeads = Array.isArray(res.heads)
          ? res.heads
          : res.heads?.data || [];
        this.loading = false;
        this.loadingState.emit(false);
        this.cdr.detectChanges();
      },
      error: () => {
        this.loading = false;
        this.loadingState.emit(false);
        this.cdr.detectChanges();
      },
    });
  }

  addFeeHead() {
    if (this.headForm.invalid) {
      this.headForm.markAllAsTouched();
      return;
    }

    if (this.editingHead) {
      this.loadingState.emit(true);
      this.feeService
        .updateFeeHead(this.editingHead.id, this.headForm.value)
        .subscribe({
          next: (head) => {
            const index = this.feeHeads.findIndex(
              (h) => h.id === this.editingHead.id,
            );
            if (index !== -1) {
              this.feeHeads[index] = head;
            }
            this.headForm.reset({ type: 'RECURRING' });
            this.editingHead = null;
            this.snackbar.success('Success', 'Fee Head Updated');
            this.loadingState.emit(false);
            this.cdr.detectChanges();
          },
          error: () => {
            this.snackbar.error('Error', 'Failed to update fee head');
            this.loadingState.emit(false);
            this.cdr.detectChanges();
          },
        });
    } else {
      this.loadingState.emit(true);
      this.feeService.createFeeHead(this.headForm.value).subscribe({
        next: (head) => {
          this.feeHeads.push(head);
          this.headForm.reset({ type: 'RECURRING' });
          this.snackbar.success('Success', 'Fee Head Added');
          this.loadingState.emit(false);
          this.cdr.detectChanges();
        },
        error: () => {
          this.snackbar.error('Error', 'Failed to add fee head');
          this.loadingState.emit(false);
          this.cdr.detectChanges();
        },
      });
    }
  }

  editFeeHead(head: any) {
    this.editingHead = head;
    this.headForm.patchValue({
      name: head.name,
      type: head.type,
    });
  }

  cancelEditFeeHead() {
    this.editingHead = null;
    this.headForm.reset({ type: 'RECURRING' });
  }

  onClassChange(classId?: string) {
    const resolvedClassId = this.sanitizeId(classId ?? this.selectedClassId);
    this.selectedClassId = resolvedClassId;

    this.cdr.detectChanges();

    if (!this.selectedClassId) {
      this.currentStructure = [];
      return;
    }

    if (!this.activeYearId) {
      this.snackbar.error(
        'Configuration Error',
        'No active academic year found. Please set one up first.',
      );
      return;
    }

    this.loading = true;
    this.loadingState.emit(true);
    this.cdr.detectChanges();
    this.feeService
      .getFeeStructure(this.selectedClassId, this.activeYearId)
      .subscribe({
        next: (res: any) => {
          this.currentStructure = Array.isArray(res) ? res : res?.data || [];
          this.loading = false;
          this.loadingState.emit(false);
          this.cdr.detectChanges();
        },
        error: () => {
          this.loading = false;
          this.loadingState.emit(false);
          this.cdr.detectChanges();
        },
      });
  }

  addStructure() {
    const classId = this.sanitizeId(this.selectedClassId);
    const academicYearId = this.sanitizeId(this.activeYearId);

    if (!classId) {
      this.snackbar.error('Error', 'Please select a class first');
      return;
    }

    if (!academicYearId) {
      this.snackbar.error(
        'Configuration Error',
        'No active academic year found. Please set one up first.',
      );
      return;
    }

    if (this.structureForm.invalid) {
      this.structureForm.markAllAsTouched();
      return;
    }

    const payload = {
      ...this.structureForm.value,
      class_id: classId,
      academic_year_id: academicYearId,
    };

    this.loadingState.emit(true);
    if (this.editingStructure) {
      this.feeService
        .updateFeeStructure(this.editingStructure.id, payload)
        .subscribe({
          next: () => {
            this.snackbar.success('Success', 'Fee updated');
            this.onClassChange();
            this.closeStructureModal();
          },
          error: () => {
            this.snackbar.error('Error', 'Failed to update fee');
            this.loadingState.emit(false);
            this.cdr.detectChanges();
          },
        });
    } else {
      this.feeService.createFeeStructure(payload).subscribe({
        next: () => {
          this.snackbar.success('Success', 'Fee assigned to class');
          this.onClassChange();
          this.closeStructureModal();
        },
        error: () => {
          this.snackbar.error('Error', 'Failed to assign fee');
          this.loadingState.emit(false);
          this.cdr.detectChanges();
        },
      });
    }
  }

  openEditStructureModal(structure: any) {
    this.editingStructure = structure;
    this.structureForm.patchValue({
      fee_head_id:
        this.getEntityId(structure, ['fee_head_id']) ||
        this.getEntityId(structure?.fee_head, ['id', 'fee_head_id', '_id']),
      amount: structure.amount,
      frequency: structure.frequency,
      due_date_day: structure.due_date_day || 10,
    });
    this.showAddModal = true;
  }

  closeStructureModal() {
    this.showAddModal = false;
    this.editingStructure = null;
    this.structureForm.reset({ frequency: 'MONTHLY', due_date_day: 10 });
  }

  deleteStructure(id: string) {
    if (confirm('Remove this fee?')) {
      this.loadingState.emit(true);
      this.feeService.deleteFeeStructure(id).subscribe({
        next: () => {
          this.snackbar.success('Success', 'Fee removed');
          this.onClassChange();
        },
        error: () => {
          this.snackbar.error('Error', 'Failed to remove fee');
          this.loadingState.emit(false);
          this.cdr.detectChanges();
        },
      });
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

  private sanitizeId(value: unknown): string {
    const normalized = String(value ?? '').trim();
    if (!normalized || normalized === 'undefined' || normalized === 'null') {
      return '';
    }
    return normalized;
  }

  private getEntityId(entity: any, keys: string[]): string {
    if (!entity) return '';
    for (const key of keys) {
      const value = this.sanitizeId(entity?.[key]);
      if (value) return value;
    }
    return '';
  }
}
