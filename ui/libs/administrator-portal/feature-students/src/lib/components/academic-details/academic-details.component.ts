import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { FormGroup, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'shikshakul-academic-details',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './academic-details.component.html',
  styleUrl: './academic-details.component.scss',
})
export class AcademicDetailsComponent {
  @Input() formGroup!: FormGroup;
}
