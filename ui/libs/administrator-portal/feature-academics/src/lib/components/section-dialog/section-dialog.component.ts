import { CommonModule } from '@angular/common';
import { Component, EventEmitter, inject, Input, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ClassManagementService } from '@shikshakul/data-access/academic';

@Component({
  selector: 'shikshakul-section-dialog',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './section-dialog.component.html',
  styleUrl: './section-dialog.component.scss',
})
export class SectionDialogComponent {
  private classService = inject(ClassManagementService);

  @Input() classId!: string;
  @Input() className = '';

  @Output() close = new EventEmitter<void>();
  @Output() save = new EventEmitter<void>();

  sectionData = {
    name: '',
    capacity: 40,
  };

  isSubmitting = false;

  onSubmit() {
    if (!this.sectionData.name || !this.classId) return;

    this.isSubmitting = true;

    this.classService
      .createSection(this.classId, {
        name: this.sectionData.name,
        capacity: this.sectionData.capacity,
      })
      .subscribe({
        next: () => {
          this.isSubmitting = false;
          this.save.emit();
          this.close.emit();
        },
        error: () => {
          this.isSubmitting = false;
          alert('Failed to create section.');
        },
      });
  }

  onCancel() {
    this.close.emit();
  }
}
