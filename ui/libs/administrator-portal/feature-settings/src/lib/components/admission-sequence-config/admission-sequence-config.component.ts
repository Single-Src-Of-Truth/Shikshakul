import { CommonModule } from '@angular/common';
import { Component, inject, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { StudentService } from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';

@Component({
  selector: 'lib-admission-sequence-config',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './admission-sequence-config.component.html',
  styleUrl: './admission-sequence-config.component.scss',
})
export class AdmissionSequenceConfigComponent implements OnInit {
  private fb = inject(FormBuilder);
  private studentService = inject(StudentService);
  private snackbar = inject(SnackbarService);

  sequenceForm!: FormGroup;
  saving = false;

  ngOnInit(): void {
    this.sequenceForm = this.fb.group({
      prefix: [''],
      suffix: [''],
      start_sequence: [1, [Validators.required, Validators.min(1)]]
    });
  }

  getPreview(): string {
    const prefix = this.sequenceForm.get('prefix')?.value || '';
    const suffix = this.sequenceForm.get('suffix')?.value || '';
    const start = this.sequenceForm.get('start_sequence')?.value || 1;

    // Format start sequence to at least 4 digits if needed
    const paddedStart = start.toString().padStart(4, '0');

    return `${prefix}${paddedStart}${suffix}`;
  }

  saveSequence() {
    if (this.sequenceForm.invalid) {
      this.snackbar.error('Error', 'Please provide a valid start sequence number.');
      return;
    }

    this.saving = true;
    this.studentService.configureSequence(this.sequenceForm.value).subscribe({
      next: () => {
        this.snackbar.success('Success', 'Admission sequence configured successfully.');
        this.saving = false;
      },
      error: (err) => {
        this.snackbar.error('Error', 'Failed to save admission sequence configuration.');
        console.error(err);
        this.saving = false;
      }
    });
  }
}
