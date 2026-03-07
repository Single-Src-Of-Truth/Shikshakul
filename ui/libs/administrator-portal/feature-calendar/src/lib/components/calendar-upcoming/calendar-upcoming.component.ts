import {
  Component,
  Input,
  OnChanges,
  SimpleChanges,
  signal,
  computed,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { EventResponse } from '@shikshakul/data-access/academic';
import { SelectComponent, SelectOption } from '@shikshakul/shared/ui/select';

@Component({
  selector: 'shikshakul-calendar-upcoming',
  standalone: true,
  imports: [CommonModule, FormsModule, SelectComponent],
  templateUrl: './calendar-upcoming.component.html',
  styleUrl: './calendar-upcoming.component.scss',
})
export class CalendarUpcomingComponent implements OnChanges {
  @Input() events: EventResponse[] = [];
  @Input() loading = false;

  selectedCategory = signal<string>('ALL');

  categoryOptions: SelectOption[] = [
    { label: 'All Categories', value: 'ALL' },
    { label: 'Holidays', value: 'HOLIDAY' },
    { label: 'Exams', value: 'EXAM' },
    { label: 'PTM', value: 'PTM' },
    { label: 'Cultural', value: 'CULTURAL' },
    { label: 'Other', value: 'OTHER' },
  ];

  upcomingEvents = signal<EventResponse[]>([]);

  groupedUpcomingEvents = computed(() => {
    const category = this.selectedCategory();

    const filtered = this.upcomingEvents().filter((ev) => {
      return (
        category === 'ALL' ||
        ev.event_type?.toUpperCase() === category ||
        (category === 'HOLIDAY' && ev.is_holiday)
      );
    });

    const groups: { month: string; events: EventResponse[] }[] = [];
    const monthMap = new Map<string, EventResponse[]>();

    filtered.forEach((ev) => {
      const date = new Date(ev.start_date);
      const monthYear = date.toLocaleString('default', {
        month: 'long',
        year: 'numeric',
      });

      if (!monthMap.has(monthYear)) {
        monthMap.set(monthYear, []);
        groups.push({ month: monthYear, events: monthMap.get(monthYear)! });
      }
      monthMap.get(monthYear)!.push(ev);
    });

    return groups;
  });

  ngOnChanges(changes: SimpleChanges) {
    if (changes['events']) {
      this.updateUpcomingEvents();
    }
  }

  private updateUpcomingEvents() {
    const now = new Date();
    now.setHours(0, 0, 0, 0);

    const upcoming = this.events
      .filter((e) => new Date(e.start_date).getTime() >= now.getTime())
      .sort(
        (a, b) =>
          new Date(a.start_date).getTime() - new Date(b.start_date).getTime(),
      );

    this.upcomingEvents.set(upcoming);
  }

  formatDateDay(isoDate: string): string {
    return new Date(isoDate).getDate().toString().padStart(2, '0');
  }

  formatDateDayName(isoDate: string): string {
    return new Date(isoDate).toLocaleString('default', { weekday: 'short' });
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
    switch (eventType?.toUpperCase()) {
      case 'HOLIDAY':
        return 'badge-red';
      case 'EXAM':
        return 'badge-amber';
      case 'PTM':
        return 'badge-purple';
      case 'CULTURAL':
        return 'badge-emerald';
      default:
        return 'badge-slate';
    }
  }

  getEventClassBg(eventType: string): string {
    switch (eventType?.toUpperCase()) {
      case 'HOLIDAY':
        return 'bg-red';
      case 'EXAM':
        return 'bg-amber';
      case 'PTM':
        return 'bg-purple';
      case 'CULTURAL':
        return 'bg-emerald';
      default:
        return 'bg-slate';
    }
  }
}
