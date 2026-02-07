import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { ProfileHeaderComponent } from '../components/profile-header/profile-header.component';
import { PerformanceStatsComponent } from '../components/performance-stats/performance-stats.component';
import { MarksHistoryComponent } from '../components/marks-history/marks-history.component';
import { GuardianCardComponent } from '../components/guardian-card/guardian-card.component';

@Component({
  selector: 'shikshakul-student-details-page',
  standalone: true,
  imports: [
    CommonModule,
    ProfileHeaderComponent,
    PerformanceStatsComponent,
    MarksHistoryComponent,
    GuardianCardComponent,
  ],
  templateUrl: './student-details-page.component.html',
  styleUrl: './student-details-page.component.scss',
})
export class StudentDetailsPageComponent {
  student = {
    name: 'Aarav Patel',
    class: 'Class 10-B',
    rollNo: '12',
    house: 'Red',
    admissionNo: '2023-098',
    dob: '14 Aug 2008',
    boardRegNo: 'R/2024/5567',
    category: 'General',
    avatarUrl: 'assets/student-aarav.jpg',
  };

  exams = [
    {
      name: 'Unit Test 1',
      date: '15 Jul 2023',
      max: 50,
      obtained: 42,
      grade: 'A2',
      remarks: 'Good concepts',
    },
    {
      name: 'Half Yearly',
      date: '22 Sep 2023',
      max: 100,
      obtained: 88,
      grade: 'A1',
      remarks: 'Excellent improvement',
    },
    {
      name: 'Unit Test 2',
      date: '10 Dec 2023',
      max: 50,
      obtained: 38,
      grade: 'B1',
      remarks: 'Needs focus on Algebra',
    },
    {
      name: 'Project Work',
      date: '05 Jan 2024',
      max: 20,
      obtained: 19,
      grade: 'A1',
      remarks: 'Submitted on time',
    },
  ];

  observation =
    'Aarav is a bright student who participates actively in class discussions. While he has shown excellent grasp of Geometry, he tends to make calculation errors in Algebra. I would recommend practicing more problems from Chapter 5 and 7.';
}
