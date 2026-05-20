import { CommonModule } from '@angular/common';
import {
  Component,
  inject,
  ChangeDetectorRef,
  Output,
  EventEmitter,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import {
  FeeService,
  StudentService,
  StudentStatus,
} from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';

@Component({
  selector: 'shikshakul-fee-collection',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './fee-collection.component.html',
  styleUrl: './fee-collection.component.scss',
})
export class FeeCollectionComponent {
  StudentStatus = StudentStatus;
  private feeService = inject(FeeService);
  private studentService = inject(StudentService);
  private snackbar = inject(SnackbarService);

  @Output() loadingState = new EventEmitter<boolean>();

  viewMode: 'DIRECTORY' | 'COLLECTOR' = 'DIRECTORY';
  searchQuery = '';
  students: any[] = [];
  selectedStudent: any = null;
  dues: any[] = [];

  payingDue: any = null;
  paymentAmount = 0;
  todayDate = new Date().toISOString().split('T')[0];
  paymentMode: 'CASH' | 'ONLINE' | 'CHEQUE' | 'UPI' = 'CASH';
  remarks = '';

  totalDue = 0;
  transactionHistory: any[] = [];

  constructor(private cdr: ChangeDetectorRef) {}

  searchStudent() {
    if (!this.searchQuery) return;
    this.loadingState.emit(true);

    this.studentService.getActiveStudents().subscribe({
      next: (res: any) => {
        const list = Array.isArray(res) ? res : res?.data || [];
        this.students = list.filter(
          (s: any) =>
            s.first_name
              .toLowerCase()
              .includes(this.searchQuery.toLowerCase()) ||
            s.last_name
              .toLowerCase()
              .includes(this.searchQuery.toLowerCase()) ||
            s.admission_no
              ?.toLowerCase()
              .includes(this.searchQuery.toLowerCase()),
        );
        this.loadingState.emit(false);
        this.cdr.detectChanges();
      },
      error: () => {
        this.snackbar.error('Error', 'Failed to fetch students');
        this.loadingState.emit(false);
      },
    });
  }

  selectStudent(student: any) {
    this.selectedStudent = student;
    this.viewMode = 'COLLECTOR';
    this.searchQuery = '';

    this.loadingState.emit(true);
    this.feeService.getStudentDues(student.id).subscribe({
      next: (data: any) => {
        this.dues = Array.isArray(data) ? data : data?.data || [];
        // Sum up all unpaid dues
        this.totalDue = this.dues.reduce(
          (sum, item) => sum + (item.amount_due || 0),
          0,
        );
        this.paymentAmount = this.totalDue;

        // Auto-select the first due for total collection if available
        if (this.dues.length > 0) {
          this.payingDue = this.dues[0];
        }

        this.loadingState.emit(false);
        this.cdr.detectChanges();
      },
      error: () => {
        this.snackbar.error('Error', 'Could not load dues');
        this.loadingState.emit(false);
        this.cdr.detectChanges();
      },
    });

    this.feeService.getTransactionHistory(student.id).subscribe({
      next: (res: any) => {
        const history = Array.isArray(res) ? res : res?.data || [];
        this.transactionHistory = history.slice(0, 3);
      },
      error: () => (this.transactionHistory = []),
    });
  }

  openPayment(due: any) {
    this.payingDue = due;
    this.paymentAmount = due.amount_due;
  }

  submitPayment() {
    if (!this.selectedStudent) {
      this.snackbar.error('Error', 'No student selected');
      return;
    }

    if (!this.payingDue && this.dues.length > 0) {
      this.payingDue = this.dues[0];
    }

    if (!this.payingDue) {
      this.snackbar.error('Error', 'No outstanding dues to collect');
      return;
    }

    if (this.paymentAmount <= 0) {
      this.snackbar.error(
        'Validation Error',
        'Please enter a valid amount greater than 0',
      );
      return;
    }

    const payload = {
      student_fee_id: this.payingDue.id,
      amount: this.paymentAmount,
      payment_mode: this.paymentMode,
      remarks: this.remarks || 'Collected via Fee Management Portal',
    };

    this.loadingState.emit(true);
    this.feeService.collectFee(payload).subscribe({
      next: () => {
        this.snackbar.success('Success', 'Payment collected successfully!');
        this.remarks = '';
        this.selectStudent(this.selectedStudent); // Refresh context
      },
      error: (err) => {
        console.error('Payment Error:', err);
        this.snackbar.error('Error', 'Payment failed. Please try again.');
        this.loadingState.emit(false);
        this.cdr.detectChanges();
      },
    });
  }

  backToDirectory() {
    this.viewMode = 'DIRECTORY';
    this.selectedStudent = null;
  }
}
