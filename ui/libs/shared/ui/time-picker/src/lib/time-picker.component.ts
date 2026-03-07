import { Component, forwardRef, Input, ElementRef, HostListener, signal, computed, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';

@Component({
    selector: 'shikshakul-time-picker',
    standalone: true,
    imports: [CommonModule],
    templateUrl: './time-picker.component.html',
    styleUrl: './time-picker.component.scss',
    providers: [
        {
            provide: NG_VALUE_ACCESSOR,
            useExisting: forwardRef(() => TimePickerComponent),
            multi: true
        }
    ]
})
export class TimePickerComponent implements ControlValueAccessor {
    @Input() placeholder = 'Select time';
    @Input() hasError = false;
    @Input() defaultHour = '09';
    @Input() defaultMinute = '00';
    @Input() defaultPeriod = 'AM';

    isOpen = signal(false);
    openUpwards = signal(false);
    selectedHour = signal<string | null>(null);
    selectedMinute = signal<string | null>(null);
    selectedPeriod = signal<string | null>(null);

    hours = Array.from({ length: 12 }, (_, i) => (i + 1).toString().padStart(2, '0'));
    minutes = Array.from({ length: 60 }, (_, i) => i.toString().padStart(2, '0'));
    periods = ['AM', 'PM'];

    formattedTime = computed(() => {
        const h = this.selectedHour();
        const m = this.selectedMinute();
        const p = this.selectedPeriod();
        if (!h || !m || !p) return '';
        return `${h}:${m} ${p}`;
    });

    constructor(private el: ElementRef, private cdr: ChangeDetectorRef) { }

    @HostListener('document:click', ['$event'])
    onClick(event: Event) {
        if (!this.el.nativeElement.contains(event.target)) {
            this.close();
        }
    }

    toggle() {
        this.isOpen.set(!this.isOpen());
        if (this.isOpen()) {
            this.checkPosition();
            // If empty when opening, set visual defaults from inputs
            if (!this.selectedHour()) this.selectedHour.set(this.defaultHour);
            if (!this.selectedMinute()) this.selectedMinute.set(this.defaultMinute);
            if (!this.selectedPeriod()) this.selectedPeriod.set(this.defaultPeriod);
            this.scrollToSelected();
        }
    }

    private checkPosition() {
        const rect = this.el.nativeElement.getBoundingClientRect();
        const spaceBelow = window.innerHeight - rect.bottom;
        const dropdownHeight = 300; // Estimated height of dropdown
        this.openUpwards.set(spaceBelow < dropdownHeight);
    }

    private scrollToSelected() {
        setTimeout(() => {
            const columns = this.el.nativeElement.querySelectorAll('.skl-column-scroll');
            columns.forEach((column: HTMLElement) => {
                const activeOption = column.querySelector('.active') as HTMLElement;
                if (activeOption) {
                    const scrollAmount = activeOption.offsetTop - (column.clientHeight / 2) + (activeOption.clientHeight / 2);
                    column.scrollTo({ top: scrollAmount, behavior: 'instant' });
                }
            });
        }, 50);
    }

    close() {
        if (this.isOpen()) {
            this.isOpen.set(false);
            this.updateValue();
            this.onTouched();
        }
    }

    selectHour(h: string, event: Event) {
        event.stopPropagation();
        this.selectedHour.set(h);
        this.updateValue();
    }

    selectMinute(m: string, event: Event) {
        event.stopPropagation();
        this.selectedMinute.set(m);
        this.updateValue();
    }

    selectPeriod(p: string, event: Event) {
        event.stopPropagation();
        this.selectedPeriod.set(p);
        this.updateValue();
    }

    private updateValue() {
        const time = this.formattedTime();
        this.onChange(time);
        this.cdr.detectChanges();
    }

    onChange = (value: any) => { };
    onTouched = () => { };

    writeValue(value: any): void {
        if (value) {
            // Expected format: "HH:mm AM/PM" or "HH:mm"
            const parts = value.split(' ');
            const timeParts = parts[0].split(':');

            if (timeParts.length === 2) {
                this.selectedHour.set(timeParts[0].padStart(2, '0'));
                this.selectedMinute.set(timeParts[1].padStart(2, '0'));
            }

            if (parts.length === 2) {
                this.selectedPeriod.set(parts[1].toUpperCase());
            } else {
                this.selectedPeriod.set('AM');
            }
        } else {
            this.selectedHour.set(null);
            this.selectedMinute.set(null);
            this.selectedPeriod.set(null);
        }
    }

    registerOnChange(fn: any): void {
        this.onChange = fn;
    }

    registerOnTouched(fn: any): void {
        this.onTouched = fn;
    }
}
