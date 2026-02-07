import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';

@Component({
  selector: 'shikshakul-academic-year-stats',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './academic-year-stats.component.html',
  styleUrl: './academic-year-stats.component.scss',
})
export class AcademicYearStatsComponent {
  @Input() currentYear = '2024-2025';
  @Input() upcomingYear = '2025-2026';
  @Input() totalArchives = 8;
}
