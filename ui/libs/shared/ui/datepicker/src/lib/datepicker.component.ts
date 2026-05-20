import { Component, forwardRef, Input, ElementRef, HostListener, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';

@Component({
    selector: 'shikshakul-datepicker',
    standalone: true,
    imports: [CommonModule],
    templateUrl: './datepicker.component.html',
    styleUrls: ['./datepicker.component.scss'],
    providers: [
        {
            provide: NG_VALUE_ACCESSOR,
            useExisting: forwardRef(() => DatepickerComponent),
            multi: true
        }
    ]
})
export class DatepickerComponent implements ControlValueAccessor {
    @Input() placeholder = 'Select date';
    @Input() hasError = false;

    isOpen = signal(false);
    currentDate = signal<Date | null>(null);
    viewDate = signal<Date>(new Date());

    weekDays = ['Su', 'Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa'];

    calendarGrid = computed(() => {
        const year = this.viewDate().getFullYear();
        const month = this.viewDate().getMonth();

        const firstDay = new Date(year, month, 1);
        const startingDayOfWeek = firstDay.getDay();
        const daysInMonth = new Date(year, month + 1, 0).getDate();
        const daysInPrevMonth = new Date(year, month, 0).getDate();
        const grid: { date: Date, isCurrentMonth: boolean, isToday: boolean, isSelected: boolean }[] = [];

        const today = new Date();
        today.setHours(0, 0, 0, 0);
        const selected = this.currentDate();

        for (let i = startingDayOfWeek - 1; i >= 0; i--) {
            const d = new Date(year, month - 1, daysInPrevMonth - i);
            grid.push({ date: d, isCurrentMonth: false, isToday: false, isSelected: false });
        }

        for (let i = 1; i <= daysInMonth; i++) {
            const d = new Date(year, month, i);
            const isToday = d.getTime() === today.getTime();
            const isSelected = selected ? d.getTime() === new Date(selected).setHours(0, 0, 0, 0) : false;
            grid.push({ date: d, isCurrentMonth: true, isToday, isSelected });
        }

        const padding = 42 - grid.length;
        for (let i = 1; i <= padding; i++) {
            const d = new Date(year, month + 1, i);
            grid.push({ date: d, isCurrentMonth: false, isToday: false, isSelected: false });
        }

        return grid;
    });

    formattedDate = computed(() => {
        const d = this.currentDate();
        if (!d) return '';
        return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
    });

    currentMonthYear = computed(() => {
        return this.viewDate().toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
    });

    constructor(private el: ElementRef) { }

    // Close when clicking outside
    @HostListener('document:click', ['$event'])
    onClick(event: Event) {
        if (!this.el.nativeElement.contains(event.target)) {
            this.close();
        }
    }

    toggle() {
        this.isOpen.set(!this.isOpen());
        if (this.isOpen() && this.currentDate()) {
            this.viewDate.set(new Date(this.currentDate()!));
        }
        this.onTouched();
    }

    close() {
        this.isOpen.set(false);
    }

    prevMonth(event: Event) {
        event.stopPropagation();
        const d = this.viewDate();
        this.viewDate.set(new Date(d.getFullYear(), d.getMonth() - 1, 1));
    }

    nextMonth(event: Event) {
        event.stopPropagation();
        const d = this.viewDate();
        this.viewDate.set(new Date(d.getFullYear(), d.getMonth() + 1, 1));
    }

    selectDate(cell: any, event: Event) {
        event.stopPropagation();
        this.currentDate.set(cell.date);

        const d = new Date(cell.date);
        const value = d.getFullYear() + '-' + String(d.getMonth() + 1).padStart(2, '0') + '-' + String(d.getDate()).padStart(2, '0');

        this.onChange(value);
        this.close();
    }

    onChange = (value: any) => { };
    onTouched = () => { };

    writeValue(value: any): void {
        if (value) {
            this.currentDate.set(new Date(value));
            this.viewDate.set(new Date(value));
        } else {
            this.currentDate.set(null);
        }
    }

    registerOnChange(fn: any): void {
        this.onChange = fn;
    }

    registerOnTouched(fn: any): void {
        this.onTouched = fn;
    }
}
