import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';

@Component({
  selector: 'shikshakul-marks-history',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './marks-history.component.html',
  styleUrl: './marks-history.component.scss',
})
export class MarksHistoryComponent {
  @Input() exams: any[] = [];
  @Input() observation = '';

  getGradeClass(grade: string): string {
    if (grade?.startsWith('A')) return 'green';
    if (grade?.startsWith('B')) return 'yellow';
    return '';
  }
}
