import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Subject } from '@shikshakul/data-access/academic';

@Component({
  selector: 'shikshakul-subject-dialog',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './subject-dialog.component.html',
  styleUrl: './subject-dialog.component.scss',
})
export class SubjectDialogComponent {
  @Input() subject: Partial<Subject> = { name: '', code: '', type: 'THEORY' };
  @Input() isEditMode = false;

  @Output() close = new EventEmitter<void>();
  @Output() save = new EventEmitter<Partial<Subject>>();

  isSubmitting = false;

  onSubmit() {
    if (this.subject.name && this.subject.code) {
      this.isSubmitting = true;
      this.save.emit(this.subject);
    }
  }

  onCancel() {
    this.close.emit();
  }
}
