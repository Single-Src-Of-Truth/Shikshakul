import { Component, EventEmitter, Input, Output, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import {
  CertificateService,
  BonafideResponse,
} from '@shikshakul/data-access/academic';

@Component({
  selector: 'shikshakul-bonafide-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './bonafide-form.component.html',
  styleUrl: './bonafide-form.component.scss',
})
export class BonafideFormComponent {
  private fb = inject(FormBuilder);
  private certificateService = inject(CertificateService);

  @Input() student: any;
  @Output() cancel = new EventEmitter<void>();
  @Output() success = new EventEmitter<BonafideResponse>();

  bonafideForm = this.fb.group({
    purpose: ['', Validators.required],
  });

  onSubmit() {
    if (this.bonafideForm.valid && this.student) {
      const req = {
        student_id: this.student.id,
        purpose: this.bonafideForm.value.purpose!,
      };

      this.certificateService.generateBonafide(req).subscribe({
        next: (res: any) => {
          this.success.emit(res.data || res);
        },
        error: (err) => {
          alert(
            'Failed to generate Bonafide: ' +
              (err.error?.message || err.message),
          );
        },
      });
    }
  }
}
