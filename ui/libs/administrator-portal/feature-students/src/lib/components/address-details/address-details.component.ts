import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { FormGroup, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'shikshakul-address-details',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './address-details.component.html',
  styleUrl: './address-details.component.scss',
})
export class AddressDetailsComponent {
  @Input() formGroup!: FormGroup;
}
