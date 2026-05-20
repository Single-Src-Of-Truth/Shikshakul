import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'shikshakul-student-list',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './student-list.component.html',
  styleUrl: './student-list.component.scss',
})
export class StudentListComponent {
  @Input() students: any[] = [];
  @Output() statusChange = new EventEmitter<void>();

  setStatus(student: any, status: 'PRESENT' | 'ABSENT' | 'LATE') {
    student.status = status;
    this.statusChange.emit();
  }

  markAllPresent(event: any) {
    const isChecked = event.target.checked;
    if (isChecked) {
      this.students.forEach((s) => (s.status = 'PRESENT'));
    }
    this.statusChange.emit();
  }
}
