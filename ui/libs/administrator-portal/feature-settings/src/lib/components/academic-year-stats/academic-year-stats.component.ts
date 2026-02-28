import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { AcademicYear } from '@shikshakul/data-access/academic';

@Component({
  selector: 'shikshakul-academic-year-stats',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './academic-year-stats.component.html',
  styleUrl: './academic-year-stats.component.scss',
})
export class AcademicYearStatsComponent {
  @Input() currentYear: AcademicYear[] = [];

  get activeYearName(): string {
    return this.currentYear.find(y => y.is_current)?.name || 'None';
  }

  get upcomingYearName(): string {
    const today = new Date();
    return this.currentYear.find(y => !y.is_current && new Date(y.start_date) > today)?.name || 'N/A';
  }

  get archivedCount(): number {
    const today = new Date();
    return this.currentYear.filter(y => !y.is_current && new Date(y.start_date) <= today).length;
  }
}
