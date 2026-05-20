import {
  Component,
  Input,
  Output,
  EventEmitter,
  inject,
  OnInit,
  OnChanges,
  SimpleChanges,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import {
  EventResponse,
  CalendarEventType,
} from '@shikshakul/data-access/academic';
import { DatepickerComponent } from '@shikshakul/shared/ui/datepicker';
import { TimePickerComponent } from '@shikshakul/shared/ui/time-picker';
import { SelectComponent, SelectOption } from '@shikshakul/shared/ui/select';

@Component({
  selector: 'shikshakul-event-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    DatepickerComponent,
    TimePickerComponent,
    SelectComponent,
  ],
  templateUrl: './event-form.component.html',
  styleUrl: './event-form.component.scss',
})
export class EventFormComponent implements OnInit, OnChanges {
  @Input() mode: 'view' | 'add' | 'edit' = 'add';
  @Input() event: EventResponse | null = null;
  @Input() initialDate: Date | null = null;

  @Output() save = new EventEmitter<any>();
  @Output() delete = new EventEmitter<string>();
  @Output() cancel = new EventEmitter<void>();

  private fb = inject(FormBuilder);

  eventForm!: FormGroup;
  @Input() submitting = false;
  private initialized = false;

  eventTypeOptions: SelectOption[] = [
    { label: 'Holiday', value: 'HOLIDAY' },
    { label: 'Exam', value: 'EXAM' },
    { label: 'PTM', value: 'PTM' },
    { label: 'Cultural', value: 'CULTURAL' },
    { label: 'Other', value: 'OTHER' },
  ];

  ngOnInit() {
    this.createForm();
    this.initialized = true;
  }

  private createForm() {
    this.eventForm = this.fb.group({
      title: ['', Validators.required],
      description: [''],
      event_type: ['OTHER' as CalendarEventType, Validators.required],
      start_date: ['', Validators.required],
      start_time: ['09:00 AM'],
      end_date: ['', Validators.required],
      end_time: ['10:00 AM'],
      is_holiday: [false],
    });
  }

  ngOnChanges(changes: SimpleChanges) {
    if (!this.initialized) return;

    // Only re-initialize if the specific event being edited has changed
    // or if we are switching from View to Add/Edit and the form is fresh.
    if (changes['event'] && this.event && this.mode === 'edit') {
      this.initFormWithEvent();
    } else if (
      changes['mode'] &&
      this.mode === 'add' &&
      changes['initialDate']
    ) {
      // If we jumped to a new date, we should probably update the form
      this.initFormForAdd();
    } else if (
      changes['mode'] &&
      this.mode === 'add' &&
      !this.eventForm.get('title')?.value
    ) {
      // If opening fresh add mode
      this.initFormForAdd();
    }
  }

  private initFormWithEvent() {
    if (!this.event) return;
    this.eventForm.patchValue({
      title: this.event.title,
      description: this.event.description,
      event_type: this.event.event_type as CalendarEventType,
      start_date: this.formatDateForPicker(this.event.start_date),
      start_time: this.formatTimeForPicker(this.event.start_date),
      end_date: this.formatDateForPicker(this.event.end_date),
      end_time: this.formatTimeForPicker(this.event.end_date),
      is_holiday: this.event.is_holiday,
    });
  }

  private initFormForAdd() {
    let start_date = '';
    let end_date = '';
    const start_time = '09:00 AM';
    const end_time = '10:00 AM';

    if (this.initialDate) {
      start_date = this.formatDateForPicker(this.initialDate.toISOString());
      end_date = start_date;
    }

    this.eventForm.reset({
      event_type: 'OTHER',
      is_holiday: false,
      start_date,
      start_time,
      end_date,
      end_time,
    });
  }

  resetForm() {
    if (confirm('Are you sure you want to clear the form?')) {
      if (this.mode === 'edit') {
        this.initFormWithEvent();
      } else {
        this.initFormForAdd();
      }
    }
  }

  private formatDateForPicker(isoString: string): string {
    if (!isoString) return '';
    try {
      const d = new Date(isoString);
      const yyyy = d.getFullYear();
      const mm = String(d.getMonth() + 1).padStart(2, '0');
      const dd = String(d.getDate()).padStart(2, '0');
      return `${yyyy}-${mm}-${dd}`;
    } catch {
      return '';
    }
  }

  private formatTimeForPicker(isoString: string): string {
    if (!isoString) return '09:00 AM';
    try {
      const d = new Date(isoString);
      let hours = d.getHours();
      const minutes = String(d.getMinutes()).padStart(2, '0');
      const period = hours >= 12 ? 'PM' : 'AM';
      hours = hours % 12;
      hours = hours ? hours : 12;
      return `${String(hours).padStart(2, '0')}:${minutes} ${period}`;
    } catch {
      return '09:00 AM';
    }
  }

  private mergeDateAndTime(dateStr: string, timeStr: string): string {
    const [time, period] = timeStr.split(' ');
    const [parsedHours, minutes] = time.split(':').map(Number);
    let hours = parsedHours;

    if (period === 'PM' && hours < 12) hours += 12;
    if (period === 'AM' && hours === 12) hours = 0;

    const d = new Date(dateStr);
    d.setHours(hours, minutes, 0, 0);
    return d.toISOString();
  }

  onSubmit() {
    if (this.eventForm.invalid) {
      Object.keys(this.eventForm.controls).forEach((key) => {
        this.eventForm.get(key)?.markAsTouched();
      });
      return;
    }

    const formValue = this.eventForm.value;

    // Ensure we handle the date strings correctly for the backend
    const payload = {
      ...formValue,
      start_date: this.mergeDateAndTime(
        formValue.start_date,
        formValue.start_time,
      ),
      end_date: this.mergeDateAndTime(formValue.end_date, formValue.end_time),
    };

    // Remove temp fields before emitting
    delete (payload as any).start_time;
    delete (payload as any).end_time;

    this.save.emit(payload);
    // Reset submitting state will be handled by parent or on success
  }

  onDelete() {
    if (
      this.event?.id &&
      confirm('Are you sure you want to delete this event?')
    ) {
      this.delete.emit(this.event.id);
    }
  }
}
