import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { FeeService, StudentService } from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';

@Component({
  selector: 'shikshakul-fee-collection',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './fee-collection.component.html',
  styleUrl: './fee-collection.component.scss',
})
export class FeeCollectionComponent {
  private feeService = inject(FeeService);
  private studentService = inject(StudentService);
  private snackbar = inject(SnackbarService);

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

  searchStudent() {
    if (!this.searchQuery) return;
    // Assuming you have a search API or we fetch list and filter
    this.studentService.getStudents().subscribe((res: any) => {
      const list = Array.isArray(res) ? res : (res?.data || []);
      this.students = list.filter(
        (s: any) =>
          s.first_name.toLowerCase().includes(this.searchQuery.toLowerCase()) ||
          s.profile_data?.admission_number?.includes(this.searchQuery),
      );
    });
  }

  selectStudent(student: any) {
    this.selectedStudent = student;
    this.students = [];
    this.searchQuery = '';

    this.feeService.getStudentDues(student.id).subscribe({
      next: (data) => {
        this.dues = data;
        // Sum up all unpaid dues
        this.totalDue = this.dues.reduce(
          (sum, item) => sum + (item.amount_due || 0),
          0,
        );
        this.paymentAmount = this.totalDue;
      },
      error: () => this.snackbar.error('Error', 'Could not load dues'),
    });

    this.feeService.getTransactionHistory(student.id).subscribe({
      next: (data) => (this.transactionHistory = data.slice(0, 3)),
    });
  }

  openPayment(due: any) {
    this.payingDue = due;
    this.paymentAmount = due.amount_due;
  }

  submitPayment() {
    if (!this.payingDue) return;

    const payload = {
      student_fee_id: this.payingDue.id,
      amount: this.paymentAmount,
      payment_mode: this.paymentMode,
      remarks: 'Collected via Admin Portal',
    };

    this.feeService.collectFee(payload).subscribe({
      next: () => {
        this.snackbar.success('Success', 'Payment collected!');
        this.payingDue = null;
        this.selectStudent(this.selectedStudent); // Refresh dues
      },
      error: () => this.snackbar.error('Error', 'Payment failed'),
    });
  }
}
