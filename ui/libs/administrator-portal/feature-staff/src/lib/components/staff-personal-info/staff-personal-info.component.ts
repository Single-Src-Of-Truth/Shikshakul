import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { FormGroup, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'shikshakul-staff-personal-info',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './staff-personal-info.component.html',
  styleUrl: './staff-personal-info.component.scss',
})
export class StaffPersonalInfoComponent {
  @Input() formGroup!: FormGroup;
}
