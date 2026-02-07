import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { FormGroup, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'shikshakul-staff-employment-info',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './staff-employment-info.component.html',
  styleUrl: './staff-employment-info.component.scss',
})
export class StaffEmploymentInfoComponent {
  @Input() formGroup!: FormGroup;
}
