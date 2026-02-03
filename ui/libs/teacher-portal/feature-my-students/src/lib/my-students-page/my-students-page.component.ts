import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { StudentsFilterComponent } from '../components/students-filter/students-filter.component';
import { StudentsTableComponent } from '../components/students-table/students-table.component';

@Component({
  selector: 'shikshakul-my-students-page',
  standalone: true,
  imports: [CommonModule, StudentsFilterComponent, StudentsTableComponent],
  templateUrl: './my-students-page.component.html',
  styleUrl: './my-students-page.component.scss',
})
export class MyStudentsPageComponent {
  students = [
    {
      id: '12',
      rollNo: '01',
      name: 'Aarav Patel',
      initials: 'AP',
      avatarClass: 'blue',
      gender: 'Male',
      status: 'Active',
      statusClass: 'active',
    },
    {
      id: '13',
      rollNo: '02',
      name: 'Aditi Sharma',
      initials: 'AS',
      avatarClass: 'purple',
      gender: 'Female',
      status: 'Active',
      statusClass: 'active',
    },
    {
      id: '14',
      rollNo: '03',
      name: 'Arjun Singh',
      initials: '',
      avatarClass: 'avatar-img',
      gender: 'Male',
      status: 'Active',
      statusClass: 'active',
      avatarUrl: 'assets/student-arjun.jpg',
    },
    {
      id: '15',
      rollNo: '04',
      name: 'Diya Kapoor',
      initials: 'DK',
      avatarClass: 'orange',
      gender: 'Female',
      status: 'Absent',
      statusSub: 'Today',
      statusClass: 'absent',
    },
    {
      id: '16',
      rollNo: '05',
      name: 'Ishaan Verma',
      initials: 'IV',
      avatarClass: 'green',
      gender: 'Male',
      status: 'Active',
      statusClass: 'active',
    },
    {
      id: '17',
      rollNo: '06',
      name: 'Meera Reddy',
      initials: 'MR',
      avatarClass: 'pink',
      gender: 'Female',
      status: 'Active',
      statusClass: 'active',
    },
  ];
}
