import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { FormGroup, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'shikshakul-guardian-details',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './guardian-details.component.html',
  styleUrl: './guardian-details.component.scss',
})
export class GuardianDetailsComponent {
  @Input() formGroup!: FormGroup;
}
