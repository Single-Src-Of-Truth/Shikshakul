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
  @Input() totalClasses = 5;
  @Input() totalStudents = 186;
  @Input() classesToday = 3;
}
