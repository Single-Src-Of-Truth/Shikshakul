import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { AcademicYearStatsComponent } from '../components/academic-year-stats/academic-year-stats.component';

@Component({
  selector: 'shikshakul-academic-year-page',
  standalone: true,
  imports: [CommonModule, AcademicYearStatsComponent],
  templateUrl: './academic-year-page.component.html',
  styleUrl: './academic-year-page.component.scss',
})
export class AcademicYearPageComponent {
  years = [
    {
      name: '2024-2025',
      startDate: 'April 1, 2024',
      endDate: 'March 31, 2025',
      status: 'Active',
      statusClass: 'active',
      sub: 'Current Session',
    },
    {
      name: '2025-2026',
      startDate: 'April 1, 2025',
      endDate: 'March 31, 2026',
      status: 'Upcoming',
      statusClass: 'upcoming',
      sub: 'Next Session',
    },
    {
      name: '2023-2024',
      startDate: 'April 1, 2023',
      endDate: 'March 31, 2024',
      status: 'Archived',
      statusClass: 'archived',
      sub: '',
    },
    {
      name: '2022-2023',
      startDate: 'April 1, 2022',
      endDate: 'March 31, 2023',
      status: 'Archived',
      statusClass: 'archived',
      sub: '',
    },
  ];
}
