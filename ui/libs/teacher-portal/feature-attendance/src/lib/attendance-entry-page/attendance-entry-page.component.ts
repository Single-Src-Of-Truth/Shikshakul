import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { AttendanceFiltersComponent } from '../components/attendance-filters/attendance-filters.component';
import { StudentListComponent } from '../components/student-list/student-list.component';
import { AttendanceFooterComponent } from '../components/attendance-footer/attendance-footer.component';

@Component({
  selector: 'shikshakul-attendance-entry-page',
  standalone: true,
  imports: [
    CommonModule,
    AttendanceFiltersComponent,
    StudentListComponent,
    AttendanceFooterComponent,
  ],
  templateUrl: './attendance-entry-page.component.html',
  styleUrl: './attendance-entry-page.component.scss',
})
export class AttendanceEntryPageComponent {
  students = [
    {
      rollNo: '01',
      name: 'Aarav Sharma',
      initials: 'AS',
      status: 'P',
      remark: '',
    },
    {
      rollNo: '02',
      name: 'Aditi Rao',
      initials: 'AR',
      status: 'P',
      remark: '',
    },
    {
      rollNo: '03',
      name: 'Arjun Singh',
      initials: 'AS',
      status: 'A',
      remark: 'Sick Leave',
    },
    {
      rollNo: '04',
      name: 'Diya Patel',
      initials: 'DP',
      status: 'P',
      remark: '',
    },
    {
      rollNo: '05',
      name: 'Ishaan Gupta',
      initials: 'IG',
      status: 'L',
      remark: 'Bus Delayed',
    },
    {
      rollNo: '06',
      name: 'Kavya Verma',
      initials: 'KV',
      status: 'P',
      remark: '',
    },
  ];

  get total() {
    return this.students.length;
  }
  get present() {
    return this.students.filter((s) => s.status === 'P').length;
  }
  get absent() {
    return this.students.filter((s) => s.status === 'A').length;
  }
  get late() {
    return this.students.filter((s) => s.status === 'L').length;
  }

  updateStats() {
    // Triggers change detection for getters automatically
  }

  onSave() {
    console.log('Saving attendance:', this.students);
  }
}
