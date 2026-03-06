import { Component, OnInit, signal, computed, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CalendarHeaderComponent } from '../components/calendar-header/calendar-header.component';
import { CalendarGridComponent } from '../components/calendar-grid/calendar-grid.component';
import { CalendarStatsComponent } from '../components/calendar-stats/calendar-stats.component';
import { CalendarUpcomingComponent } from '../components/calendar-upcoming/calendar-upcoming.component';
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
    CalendarStatsComponent,
    CalendarUpcomingComponent,
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
  isUpcomingDrawerOpen = signal<boolean>(false);

  viewMode = signal<'view' | 'add' | 'edit'>('view');
  selectedEvent = signal<EventResponse | null>(null);
  selectedDate = signal<Date | null>(null);
  loading = signal<boolean>(false);
  isSubmitting = signal<boolean>(false);

  // Derived signals
  currentMonth = computed(() => this.currentDate().getMonth() + 1);
  currentYear = computed(() => this.currentDate().getFullYear());

  ngOnInit() {
    this.loadEvents();
  }

  loadEvents() {
    this.loading.set(true);
    this.calendarService
      .getEvents(this.currentMonth(), this.currentYear())
      .subscribe({
        next: (res: any) => {
          this.events.set(res.data || []);
          this.loading.set(false);
        },
        error: () => this.loading.set(false),
      });
  }

  onMonthChange(newDate: Date) {
    this.currentDate.set(newDate);
    this.loadEvents();
  }

  onAddEvent(date?: Date) {
    this.selectedEvent.set(null);
    this.selectedDate.set(date || null);
    this.viewMode.set('add');
  }

  onEditEvent(event: EventResponse) {
    this.selectedEvent.set(event);
    this.viewMode.set('edit');
  }

  onSaveEvent(data: any) {
    this.isSubmitting.set(true);
    this.academicYearService.getCurrentAcademicYear().subscribe({
      next: (res: any) => {
        const yearData = res?.data || res;
        const actualYear = Array.isArray(yearData) ? yearData[0] : yearData;
        const currentAcademicYearId =
          actualYear?.id || actualYear?.academic_year_id || res?.data?.id;

        if (!currentAcademicYearId) {
          this.snackbar.error('Error', 'No active academic year found.');
          this.isSubmitting.set(false);
          return;
        }

        if (this.viewMode() === 'add') {
          const payload: CreateEventRequest = {
            academic_year_id: currentAcademicYearId,
            ...data,
          };
          this.calendarService.createEvent(payload).subscribe({
            next: () => {
              this.snackbar.success('Success', 'Event created successfully');
              this.isSubmitting.set(false);
              this.viewMode.set('view');
              this.loadEvents();
            },
            error: (err) => {
              this.snackbar.error(
                'Error',
                err.error?.message || 'Failed to create event',
              );
              this.isSubmitting.set(false);
            },
          });
        } else if (this.viewMode() === 'edit' && this.selectedEvent()) {
          const payload: UpdateEventRequest = { ...data };
          this.calendarService
            .updateEvent(this.selectedEvent()!.id, payload)
            .subscribe({
              next: () => {
                this.snackbar.success('Success', 'Event updated successfully');
                this.isSubmitting.set(false);
                this.viewMode.set('view');
                this.loadEvents();
              },
              error: (err) => {
                this.snackbar.error(
                  'Error',
                  err.error?.message || 'Failed to update event',
                );
                this.isSubmitting.set(false);
              },
            });
        }
      },
      error: () => {
        this.snackbar.error('Error', 'Cannot fetch active academic year.');
        this.isSubmitting.set(false);
      },
    });
  }

  onDeleteEvent(id: string) {
    this.isSubmitting.set(true);
    this.calendarService.deleteEvent(id).subscribe({
      next: () => {
        this.snackbar.success('Success', 'Event deleted successfully');
        this.isSubmitting.set(false);
        this.viewMode.set('view');
        this.loadEvents();
      },
      error: (err) => {
        this.snackbar.error(
          'Error',
          err.error?.message || 'Failed to delete event',
        );
        this.isSubmitting.set(false);
      },
    });
  }
}
