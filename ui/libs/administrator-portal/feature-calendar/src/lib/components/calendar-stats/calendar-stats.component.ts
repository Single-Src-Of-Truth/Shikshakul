import {
  Component,
  Input,
  OnChanges,
  SimpleChanges,
  signal,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { EventResponse } from '@shikshakul/data-access/academic';

@Component({
  selector: 'shikshakul-calendar-stats',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './calendar-stats.component.html',
  styleUrl: './calendar-stats.component.scss',
})
export class CalendarStatsComponent implements OnChanges {
  @Input() events: EventResponse[] = [];
  @Input() loading = false;

  categoryCounts = signal({
    holidays: 0,
    exams: 0,
    ptm: 0,
    cultural: 0,
    other: 0,
  });

  ngOnChanges(changes: SimpleChanges) {
    if (changes['events']) {
      this.calculateCounts();
    }
  }

  private calculateCounts() {
    const counts = { holidays: 0, exams: 0, ptm: 0, cultural: 0, other: 0 };
    this.events.forEach((e) => {
      const type = e.event_type?.toLowerCase() || 'other';
      if (type === 'holiday' || e.is_holiday) counts.holidays++;
      else if (type === 'exam') counts.exams++;
      else if (type === 'ptm') counts.ptm++;
      else if (type === 'cultural') counts.cultural++;
      else counts.other++;
    });
    this.categoryCounts.set(counts);
  }
}
