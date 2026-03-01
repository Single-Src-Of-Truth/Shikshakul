import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'shikshakul-calendar-header',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './calendar-header.component.html',
  styleUrl: './calendar-header.component.scss',
})
export class CalendarHeaderComponent {
  @Input() currentDate!: Date;
  @Output() dateChange = new EventEmitter<Date>();
  @Output() addEvent = new EventEmitter<void>();

  get monthYear(): string {
    return this.currentDate.toLocaleString('default', {
      month: 'long',
      year: 'numeric',
    });
  }

  previousMonth() {
    const d = new Date(this.currentDate);
    d.setMonth(d.getMonth() - 1);
    this.dateChange.emit(d);
  }

  nextMonth() {
    const d = new Date(this.currentDate);
    d.setMonth(d.getMonth() + 1);
    this.dateChange.emit(d);
  }

  goToday() {
    this.dateChange.emit(new Date());
  }
}
