import {
  Component,
  Input,
  Output,
  EventEmitter,
  OnChanges,
  SimpleChanges,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { DatepickerComponent } from '@shikshakul/shared/ui/datepicker';

@Component({
  selector: 'shikshakul-calendar-header',
  standalone: true,
  imports: [CommonModule, FormsModule, DatepickerComponent],
  templateUrl: './calendar-header.component.html',
  styleUrl: './calendar-header.component.scss',
})
export class CalendarHeaderComponent implements OnChanges {
  @Input() currentDate!: Date;
  @Output() dateChange = new EventEmitter<Date>();

  jumpDate = '';

  ngOnChanges(changes: SimpleChanges) {
    if (changes['currentDate']) {
      const cd = this.currentDate;
      const yyyy = cd.getFullYear();
      const mm = String(cd.getMonth() + 1).padStart(2, '0');
      const dd = String(cd.getDate()).padStart(2, '0');
      this.jumpDate = `${yyyy}-${mm}-${dd}`;
    }
  }

  onJumpDateChange(val: string) {
    if (val) {
      this.dateChange.emit(new Date(val));
    }
  }

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
