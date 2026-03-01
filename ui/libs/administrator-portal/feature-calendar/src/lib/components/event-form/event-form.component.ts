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
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import {
  EventResponse,
  CalendarEventType,
} from '@shikshakul/data-access/academic';

@Component({
  selector: 'shikshakul-event-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './event-form.component.html',
  styleUrl: './event-form.component.scss',
})
export class EventFormComponent implements OnInit, OnChanges {
  @Input() mode: 'view' | 'add' | 'edit' = 'add';
  @Input() event: EventResponse | null = null;

  @Output() save = new EventEmitter<any>();
  @Output() delete = new EventEmitter<string>();
  @Output() cancel = new EventEmitter<void>();

  private fb = inject(FormBuilder);

  eventForm = this.fb.group({
    title: ['', Validators.required],
    description: [''],
    event_type: ['OTHER' as CalendarEventType, Validators.required],
    start_date: ['', Validators.required],
    end_date: ['', Validators.required],
    is_holiday: [false],
  });

  ngOnInit() {
    this.initForm();
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes['event'] || changes['mode']) {
      this.initForm();
    }
  }

  private initForm() {
    if (this.mode === 'edit' && this.event) {
      this.eventForm.patchValue({
        title: this.event.title,
        description: this.event.description,
        event_type: this.event.event_type as CalendarEventType,
        start_date: this.formatDateForInput(this.event.start_date),
        end_date: this.formatDateForInput(this.event.end_date),
        is_holiday: this.event.is_holiday,
      });
    } else {
      this.eventForm.reset({
        event_type: 'OTHER',
        is_holiday: false,
      });
    }
  }

  private formatDateForInput(isoString: string): string {
    if (!isoString) return '';
    try {
      const d = new Date(isoString);
      // slice(0, 16) gets 'YYYY-MM-DDTHH:mm' for datetime-local input
      return d.toISOString().slice(0, 16);
    } catch {
      return '';
    }
  }

  onSubmit() {
    if (this.eventForm.valid) {
      const formValue = this.eventForm.value;
      const payload = {
        ...formValue,
        start_date: new Date(formValue.start_date!).toISOString(),
        end_date: new Date(formValue.end_date!).toISOString(),
      };

      this.save.emit(payload);
    } else {
      Object.keys(this.eventForm.controls).forEach((key) => {
        this.eventForm.get(key)?.markAsTouched();
      });
    }
  }

  onDelete() {
    if (this.event?.id) {
      this.delete.emit(this.event.id);
    }
  }
}
