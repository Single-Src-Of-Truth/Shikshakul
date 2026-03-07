import { CommonModule } from '@angular/common';
import { Component, Input, computed } from '@angular/core';

@Component({
  selector: 'shikshakul-attendance-stats',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './attendance-stats.component.html',
  styleUrl: './attendance-stats.component.scss',
})
export class AttendanceStatsComponent {
  @Input() totalDays = 0;
  @Input() presentDays = 0;
  @Input() absentDays = 0;
  @Input() lateDays = 0;

  get overallPercentage(): number {
    if (this.totalDays === 0) return 0;
    return Math.round(((this.presentDays + this.lateDays) / this.totalDays) * 100);
  }
}
