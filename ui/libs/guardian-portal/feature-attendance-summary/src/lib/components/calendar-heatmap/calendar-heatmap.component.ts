import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';

@Component({
  selector: 'shikshakul-calendar-heatmap',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './calendar-heatmap.component.html',
  styleUrl: './calendar-heatmap.component.scss',
})
export class CalendarHeatmapComponent {
  days = ['SUN', 'MON', 'TUE', 'WED', 'THU', 'FRI', 'SAT'];

  calendarGrid = [
    { day: '', status: 0 },
    { day: '1', status: 5 },
    { day: '2', status: 4 },
    { day: '3', status: 1 },
    { day: '4', status: 1 },
    { day: '5', status: 3 },
    { day: '6', status: 1 },
    { day: '7', status: 5 },
    { day: '8', status: 2 },
    { day: '9', status: 1 },
    { day: '10', status: 1 },
    { day: '11', status: 1 },
    { day: '12', status: 2 },
    { day: '13', status: 1 },
    { day: '14', status: 5 },
    { day: '15', status: 1 },
    { day: '16', status: 1 },
    { day: '17', status: 1 },
    { day: '18', status: 1 },
    { day: '19', status: 1 },
    { day: '20', status: 1 },
    { day: '21', status: 5 },
    { day: '22', status: 1 },
    { day: '23', status: 1 },
    { day: '24', status: 5 },
    { day: '25', status: 1, isToday: true },
    { day: '26', status: 0 },
    { day: '27', status: 0 },
    { day: '28', status: 0 },
    { day: '29', status: 0 },
    { day: '30', status: 0 },
    { day: '31', status: 0 },
  ];

  getStatusClass(status: number): string {
    switch (status) {
      case 1:
        return 'present';
      case 2:
        return 'absent';
      case 3:
        return 'late';
      case 4:
        return 'holiday';
      case 5:
        return 'off';
      default:
        return '';
    }
  }
}
