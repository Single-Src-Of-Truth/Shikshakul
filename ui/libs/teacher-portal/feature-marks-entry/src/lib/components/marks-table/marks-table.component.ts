import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'shikshakul-marks-table',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './marks-table.component.html',
  styleUrl: './marks-table.component.scss',
})
export class MarksTableComponent {
  @Input() students: any[] = [];

  // Config
  maxTheory = 80;
  maxPractical = 20;

  // Calculate total automatically
  calculateTotal(student: any) {
    if (student.isAbsent) {
      student.total = 0;
      student.grade = 'AB';
      return;
    }

    const t = Number(student.theory) || 0;
    const p = Number(student.practical) || 0;

    // Simple Validation
    if (t > this.maxTheory || p > this.maxPractical) {
      student.grade = 'Error';
      return;
    }

    student.total = t + p;
    this.assignGrade(student);
  }

  assignGrade(student: any) {
    const score = student.total;
    if (score >= 90) student.grade = 'A1';
    else if (score >= 80) student.grade = 'A2';
    else if (score >= 70) student.grade = 'B1';
    else if (score >= 60) student.grade = 'B2';
    else if (score >= 33) student.grade = 'C';
    else student.grade = 'F';
  }

  toggleAbsent(student: any) {
    student.isAbsent = !student.isAbsent;
    if (student.isAbsent) {
      student.theory = null; // Clear inputs
      student.practical = null;
    }
    this.calculateTotal(student);
  }
}
