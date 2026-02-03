import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { FormGroup, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'shikshakul-staff-banking-payroll',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './staff-banking-payroll.component.html',
  styleUrl: './staff-banking-payroll.component.scss',
})
export class StaffBankingPayrollComponent {
  @Input() formGroup!: FormGroup;
}
