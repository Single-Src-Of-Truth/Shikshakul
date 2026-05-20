import { Component, EventEmitter, Input, Output, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import {
  CertificateService,
  TCResponse,
} from '@shikshakul/data-access/academic';

@Component({
  selector: 'shikshakul-tc-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './tc-form.component.html',
  styleUrl: './tc-form.component.scss',
})
export class TcFormComponent {
  private fb = inject(FormBuilder);
  private certificateService = inject(CertificateService);

  @Input() student: any;
  @Output() cancel = new EventEmitter<void>();
  @Output() success = new EventEmitter<TCResponse>();

  tcForm = this.fb.group({
    leaving_date: ['', Validators.required],
    reason: ['', Validators.required],
    conduct: ['Good', Validators.required],
    mark_inactive: [false],
  });

  onSubmit() {
    if (this.tcForm.valid && this.student) {
      const req = {
        student_id: this.student.id,
        leaving_date: this.tcForm.value.leaving_date!,
        reason: this.tcForm.value.reason!,
        conduct: this.tcForm.value.conduct!,
        mark_inactive: this.tcForm.value.mark_inactive!,
      };

      this.certificateService.generateTC(req).subscribe({
        next: (res: any) => {
          this.success.emit(res.data || res);
        },
        error: (err) => {
          alert(
            'Failed to generate TC: ' + (err.error?.message || err.message),
          );
        },
      });
    }
  }
}
