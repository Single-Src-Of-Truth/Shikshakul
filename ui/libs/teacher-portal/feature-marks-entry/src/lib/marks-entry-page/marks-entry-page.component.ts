import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { MarksFiltersComponent } from '../components/marks-filters/marks-filters.component';
import { MarksTableComponent } from '../components/marks-table/marks-table.component';
import { MarksFooterComponent } from '../components/marks-footer/marks-footer.component';

@Component({
  selector: 'shikshakul-marks-entry-page',
  standalone: true,
  imports: [
    CommonModule,
    MarksFiltersComponent,
    MarksTableComponent,
    MarksFooterComponent,
  ],
  templateUrl: './marks-entry-page.component.html',
  styleUrl: './marks-entry-page.component.scss',
})
export class MarksEntryPageComponent {
  students = [
    {
      rollNo: '01',
      name: 'Aarav Patel',
      theory: 72,
      practical: 18,
      total: 90,
      grade: 'A1',
      remark: '',
      isAbsent: false,
    },
    {
      rollNo: '02',
      name: 'Aditi Sharma',
      theory: 65,
      practical: 19,
      total: 84,
      grade: 'A2',
      remark: '',
      isAbsent: false,
    },
    {
      rollNo: '03',
      name: 'Arjun Singh',
      theory: 0,
      practical: 0,
      total: 0,
      grade: 'AB',
      remark: 'Medical',
      isAbsent: true,
    },
    {
      rollNo: '04',
      name: 'Diya Gupta',
      theory: 78,
      practical: 20,
      total: 98,
      grade: 'A1',
      remark: 'Excellent Perform',
      isAbsent: false,
    },
    {
      rollNo: '05',
      name: 'Ishaan Kumar',
      theory: 85,
      practical: 15,
      total: 0,
      grade: 'Error',
      remark: '',
      isAbsent: false,
    }, // Intentional error for demo
    {
      rollNo: '06',
      name: 'Kavya Reddy',
      theory: null,
      practical: null,
      total: 0,
      grade: '-',
      remark: '',
      isAbsent: false,
    },
  ];
}
