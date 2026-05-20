import { Component, EventEmitter, Input, Output } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ClassGrade, Section } from '@shikshakul/data-access/academic';

@Component({
  selector: 'shikshakul-attendance-filters',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './attendance-filters.component.html',
  styleUrl: './attendance-filters.component.scss',
})
export class AttendanceFiltersComponent {
  @Input() classes: ClassGrade[] = [];
  @Input() availableSections: Section[] = [];

  @Input() selectedClassId = '';
  @Input() selectedSectionId = '';
  @Input() selectedDate = '';

  @Output() filterChange = new EventEmitter<{
    classId: string;
    sectionId: string;
    date: string;
  }>();

  onClassChange(classId: string) {
    this.selectedClassId = classId;
    this.selectedSectionId = ''; // Reset section
    this.emitChange();
  }

  onSectionChange(sectionId: string) {
    this.selectedSectionId = sectionId;
    this.emitChange();
  }

  onDateChange(event: any) {
    this.selectedDate = event.target.value;
    this.emitChange();
  }

  private emitChange() {
    this.filterChange.emit({
      classId: this.selectedClassId,
      sectionId: this.selectedSectionId,
      date: this.selectedDate,
    });
  }
}

