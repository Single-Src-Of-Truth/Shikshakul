import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';

@Component({
  selector: 'shikshakul-class-stats',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './class-stats.component.html',
  styleUrl: './class-stats.component.scss',
})
export class ClassStatsComponent {
  @Input() totalClasses = 14;
  @Input() totalSections = 42;
  @Input() totalStudents = 1245;
}
