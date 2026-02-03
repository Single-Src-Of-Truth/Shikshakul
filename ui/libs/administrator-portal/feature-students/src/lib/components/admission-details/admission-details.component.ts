import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { FormGroup, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'shikshakul-admission-details',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './admission-details.component.html',
  styleUrl: './admission-details.component.scss',
})
export class AdmissionDetailsComponent {
  @Input() formGroup!: FormGroup;
}
