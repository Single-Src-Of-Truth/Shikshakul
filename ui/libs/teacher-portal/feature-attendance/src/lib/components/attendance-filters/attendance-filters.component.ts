import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';

@Component({
  selector: 'shikshakul-attendance-filters',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './attendance-filters.component.html',
  styleUrl: './attendance-filters.component.scss',
})
export class AttendanceFiltersComponent {
  selectedDate = '2023-10-16';
}
