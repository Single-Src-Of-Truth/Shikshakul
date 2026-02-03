import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { FormGroup, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'shikshakul-staff-photo',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './staff-photo.component.html',
  styleUrl: './staff-photo.component.scss',
})
export class StaffPhotoComponent {
  @Input() parentForm!: FormGroup;
  imagePreview: string | ArrayBuffer | null = null;

  onFileSelected(event: Event) {
    const file = (event.target as HTMLInputElement).files?.[0];
    if (file) {
      this.parentForm.patchValue({ photo: file });
      this.parentForm.get('photo')?.markAsTouched();

      const reader = new FileReader();
      reader.onload = () => {
        this.imagePreview = reader.result;
      };
      reader.readAsDataURL(file);
    }
  }
}
