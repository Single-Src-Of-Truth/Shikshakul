import {
  Component,
  Input,
  Output,
  EventEmitter,
  OnChanges,
  SimpleChanges,
  computed,
  signal,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { EventResponse } from '@shikshakul/data-access/academic';

interface CalendarDay {
  date: Date;
  isCurrentMonth: boolean;
  isToday: boolean;
  events: EventResponse[];
}

@Component({
  selector: 'shikshakul-calendar-grid',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './calendar-grid.component.html',
  styleUrl: './calendar-grid.component.scss',
})
export class CalendarGridComponent implements OnChanges {
  @Input() currentDate!: Date;
  @Input() events: EventResponse[] = [];
  @Output() eventClick = new EventEmitter<EventResponse>();

  weekDays = ['MON', 'TUE', 'WED', 'THU', 'FRI', 'SAT', 'SUN'];
  calendarDays = signal<CalendarDay[]>([]);

  ngOnChanges(changes: SimpleChanges) {
    if (changes['currentDate'] || changes['events']) {
      this.generateGrid();
    }
  }

  generateGrid() {
    const year = this.currentDate.getFullYear();
    const month = this.currentDate.getMonth();

    const firstDayOfMonth = new Date(year, month, 1);
    const lastDayOfMonth = new Date(year, month + 1, 0);

    // JS getDay() returns 0 for Sunday, 1 for Monday. We want Monday to be 0
    let startDayOfWeek = firstDayOfMonth.getDay() - 1;
    if (startDayOfWeek === -1) startDayOfWeek = 6; // Sunday becomes 6

    const daysInMonth = lastDayOfMonth.getDate();
    const daysFromPrevMonth = startDayOfWeek;

    const totalCells = Math.ceil((daysInMonth + daysFromPrevMonth) / 7) * 7;
    const days: CalendarDay[] = [];

    const today = new Date();
    today.setHours(0, 0, 0, 0);

    for (let i = 0; i < totalCells; i++) {
      const cellDate = new Date(year, month, i - daysFromPrevMonth + 1);
      cellDate.setHours(0, 0, 0, 0);

      const isCurrentMonth = cellDate.getMonth() === month;
      const isToday = cellDate.getTime() === today.getTime();

      // Filter events for this day
      const dayEvents = this.events.filter((e) => {
        const eStart = new Date(e.start_date);
        eStart.setHours(0, 0, 0, 0);
        // Note: End date logic might span multiple days, simplifying to start date for now
        return eStart.getTime() === cellDate.getTime();
      });

      days.push({
        date: cellDate,
        isCurrentMonth,
        isToday,
        events: dayEvents,
      });
    }

    this.calendarDays.set(days);
  }

  getEventClass(eventType: string): string {
    switch (eventType.toUpperCase()) {
      case 'HOLIDAY':
        return 'event-holiday';
      case 'EXAM':
        return 'event-exam';
      case 'PTM':
        return 'event-ptm';
      case 'CULTURAL':
        return 'event-cultural';
      default:
        return 'event-other';
    }
  }
}
