import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormsModule, NgForm } from '@angular/forms';
import { Subject } from '@shikshakul/data-access/academic';
import { SelectComponent } from '@shikshakul/shared/ui/select';

@Component({
  selector: 'shikshakul-subject-drawer',
  standalone: true,
  imports: [CommonModule, FormsModule, SelectComponent],
  templateUrl: './subject-drawer.component.html',
  styleUrl: './subject-drawer.component.scss',
})
export class SubjectDrawerComponent {
  subjectTypeOptions = [
    { label: 'Theory', value: 'THEORY' },
    { label: 'Practical', value: 'PRACTICAL' },
    { label: 'Lab', value: 'LAB' }
  ];
  @Input() isOpen = false;
  @Input() subject: Partial<Subject> = { name: '', code: '', type: 'THEORY' };
  @Input() isEditMode = false;

  @Output() close = new EventEmitter<void>();
  @Output() save = new EventEmitter<Partial<Subject>>();

  isSubmitting = false;
  isFormSubmitted = false;

  onSubmit(form: NgForm) {
    this.isFormSubmitted = true;
    
    if (form.valid && this.subject.name && this.subject.code) {
      this.isSubmitting = true;
      this.save.emit(this.subject);
      setTimeout(() => {
        this.isSubmitting = false;
        this.isFormSubmitted = false;
      }, 500); // Reset for next open
    }
  }

  onCancel() {
    this.isFormSubmitted = false;
    this.close.emit();
  }

  resetForm(form: NgForm) {
    if (confirm('Are you sure you want to clear the form?')) {
      form.resetForm();
      this.subject = { name: '', code: '', type: 'THEORY' };
      this.isFormSubmitted = false;
    }
  }
}
