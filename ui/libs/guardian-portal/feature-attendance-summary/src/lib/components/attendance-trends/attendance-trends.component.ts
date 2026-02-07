import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'shikshakul-attendance-trends',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './attendance-trends.component.html',
  styleUrls: ['./attendance-trends.component.scss'],
})
export class AttendanceTrendsComponent {
  trends = [
    { month: 'Aug', value: 75, type: 'present' },
    { month: 'Sep', value: 82, type: 'present' },
    { month: 'Oct', value: 65, type: 'present' },
    { month: 'Nov', value: 100, type: 'future' },
  ];
}
