import {
  Component,
  Input,
  OnChanges,
  SimpleChanges,
  signal,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { EventResponse } from '@shikshakul/data-access/academic';

@Component({
  selector: 'shikshakul-calendar-sidebar',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './calendar-sidebar.component.html',
  styleUrl: './calendar-sidebar.component.scss',
})
export class CalendarSidebarComponent implements OnChanges {
  @Input() events: EventResponse[] = [];
  @Input() currentDate!: Date;

  categoryCounts = signal({
    holidays: 0,
    exams: 0,
    ptm: 0,
    cultural: 0,
    other: 0,
  });

  upcomingEvents = signal<EventResponse[]>([]);

  ngOnChanges(changes: SimpleChanges) {
    if (changes['events']) {
      this.calculateCounts();
      this.filterUpcomingEvents();
    }
  }

  private calculateCounts() {
    const counts = { holidays: 0, exams: 0, ptm: 0, cultural: 0, other: 0 };
    this.events.forEach((e) => {
      const type = e.event_type.toLowerCase();
      if (type === 'holiday' || e.is_holiday) counts.holidays++;
      else if (type === 'exam') counts.exams++;
      else if (type === 'ptm') counts.ptm++;
      else if (type === 'cultural') counts.cultural++;
      else counts.other++;
    });
    this.categoryCounts.set(counts);
  }

  private filterUpcomingEvents() {
    const now = new Date();
    now.setHours(0, 0, 0, 0);

    const upcoming = this.events
      .filter((e) => new Date(e.start_date).getTime() >= now.getTime())
      .sort(
        (a, b) =>
          new Date(a.start_date).getTime() - new Date(b.start_date).getTime(),
      )
      .slice(0, 5);

    this.upcomingEvents.set(upcoming);
  }

  formatDateShort(isoDate: string): string {
    const d = new Date(isoDate);
    return d
      .toLocaleString('default', { month: 'short', day: '2-digit' })
      .toUpperCase();
  }

  formatTimeRange(start: string, end: string): string {
    const s = new Date(start);
    const e = new Date(end);
    if (s.getHours() === 0 && e.getHours() === 0) return 'All Day';

    const sStr = s.toLocaleTimeString([], {
      hour: '2-digit',
      minute: '2-digit',
    });
    const eStr = e.toLocaleTimeString([], {
      hour: '2-digit',
      minute: '2-digit',
    });
    return `${sStr} - ${eStr}`;
  }

  getEventClassBadge(eventType: string): string {
    switch (eventType.toUpperCase()) {
      case 'HOLIDAY':
        return 'badge-holiday';
      case 'EXAM':
        return 'badge-exam';
      case 'PTM':
        return 'badge-ptm';
      case 'CULTURAL':
        return 'badge-cultural';
      default:
        return 'badge-other';
    }
  }
}
