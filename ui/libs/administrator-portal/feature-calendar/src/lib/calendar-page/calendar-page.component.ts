import { Component, OnInit, signal, computed, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CalendarHeaderComponent } from '../components/calendar-header/calendar-header.component';
import { CalendarGridComponent } from '../components/calendar-grid/calendar-grid.component';
import { CalendarSidebarComponent } from '../components/calendar-sidebar/calendar-sidebar.component';
import { EventFormComponent } from '../components/event-form/event-form.component';
import {
  CalendarService,
  EventResponse,
  CreateEventRequest,
  UpdateEventRequest,
  AcademicYearService,
} from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';

@Component({
  selector: 'shikshakul-calendar-page',
  standalone: true,
  imports: [
    CommonModule,
    CalendarHeaderComponent,
    CalendarGridComponent,
    CalendarSidebarComponent,
    EventFormComponent,
  ],
  templateUrl: './calendar-page.component.html',
  styleUrl: './calendar-page.component.scss',
})
export class CalendarPageComponent implements OnInit {
  private calendarService = inject(CalendarService);
  private academicYearService = inject(AcademicYearService);
  private snackbar = inject(SnackbarService);

  currentDate = signal<Date>(new Date());
  events = signal<EventResponse[]>([]);

  viewMode = signal<'view' | 'add' | 'edit'>('view');
  selectedEvent = signal<EventResponse | null>(null);

  // Derived signals
  currentMonth = computed(() => this.currentDate().getMonth() + 1);
  currentYear = computed(() => this.currentDate().getFullYear());

  ngOnInit() {
    this.loadEvents();
  }

  loadEvents() {
    this.calendarService
      .getEvents(this.currentMonth(), this.currentYear())
      .subscribe((res: any) => {
        this.events.set(res.data || []);
      });
  }

  onMonthChange(newDate: Date) {
    this.currentDate.set(newDate);
    this.loadEvents();
  }

  onAddEvent() {
    this.selectedEvent.set(null);
    this.viewMode.set('add');
  }

  onEditEvent(event: EventResponse) {
    this.selectedEvent.set(event);
    this.viewMode.set('edit');
  }

  onSaveEvent(data: any) {
    this.academicYearService.getCurrentAcademicYear().subscribe({
      next: (res: any) => {
        console.log('Full Academic Year Response:', JSON.stringify(res));

        // Try all possible ways to find the ID
        const yearData = res?.data || res;
        const actualYear = Array.isArray(yearData) ? yearData[0] : yearData;

        const currentAcademicYearId =
          actualYear?.id ||
          actualYear?.academic_year_id ||
          actualYear?._id ||
          res?.data?.id;

        console.log('Extracted Year ID:', currentAcademicYearId);

        if (!currentAcademicYearId) {
          this.snackbar.error(
            'Error',
            'No active academic year ID found in response.',
          );
          return;
        }

        if (this.viewMode() === 'add') {
          const payload: CreateEventRequest = {
            academic_year_id: currentAcademicYearId,
            ...data,
          };
          this.calendarService.createEvent(payload).subscribe(() => {
            this.viewMode.set('view');
            this.loadEvents();
          });
        } else if (this.viewMode() === 'edit' && this.selectedEvent()) {
          const payload: UpdateEventRequest = { ...data };
          this.calendarService
            .updateEvent(this.selectedEvent()!.id, payload)
            .subscribe(() => {
              this.snackbar.success('Success', 'Event updated successfully');
              this.viewMode.set('view');
              this.loadEvents();
            });
        }
      },
      error: (err: any) => {
        console.error('Calendar Page - Error fetching Academic Year:', err);
        this.snackbar.error('Error', 'Cannot fetch active academic year.');
      },
    });
  }

  onDeleteEvent(id: string) {
    if (confirm('Are you sure you want to delete this event?')) {
      this.calendarService.deleteEvent(id).subscribe(() => {
        this.snackbar.success('Success', 'Event deleted successfully');
        this.viewMode.set('view');
        this.loadEvents();
      });
    }
  }
}
