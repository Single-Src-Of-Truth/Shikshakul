import { Component, EventEmitter, Input, Output } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ClassGrade } from '@shikshakul/data-access/academic';

@Component({
  selector: 'shikshakul-class-dialog',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './class-dialog.component.html',
  styleUrl: './class-dialog.component.scss',
})
export class ClassDialogComponent {
  @Input() classData: Partial<ClassGrade> = { name: '', sort_order: 0 };
  @Input() isEditMode = false;

  @Output() close = new EventEmitter<void>();
  @Output() save = new EventEmitter<Partial<ClassGrade>>();

  onSubmit() {
    if (this.classData.name && this.classData.sort_order) {
      this.save.emit(this.classData);
    }
  }

  onCancel() {
    this.close.emit();
  }
}
