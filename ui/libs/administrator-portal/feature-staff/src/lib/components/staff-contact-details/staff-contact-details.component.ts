import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { FormGroup, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'shikshakul-staff-contact-details',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './staff-contact-details.component.html',
  styleUrl: './staff-contact-details.component.scss',
})
export class StaffContactDetailsComponent {
  @Input() formGroup!: FormGroup;
}
