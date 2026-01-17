import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { SubjectStatsComponent } from '../components/subject-stats/subject-stats.component';

@Component({
  selector: 'shikshakul-subject-setup-page',
  standalone: true,
  imports: [CommonModule, SubjectStatsComponent],
  templateUrl: './subject-setup-page.component.html',
  styleUrl: './subject-setup-page.component.scss',
})
export class SubjectSetupPageComponent {
  subjects = [
    {
      name: 'Mathematics',
      code: 'MATH-041',
      type: 'Theory',
      category: 'Compulsory (Core)',
      createdBy: 'John Doe',
      icon: 'M',
      color: 'blue',
    },
    {
      name: 'Science (Physics)',
      code: 'SCI-PHY-086',
      type: 'Practical',
      category: 'Compulsory (Core)',
      createdBy: 'John Doe',
      icon: 'S',
      color: 'green',
    },
    {
      name: 'English Core',
      code: 'ENG-301',
      type: 'Theory',
      category: 'Language',
      createdBy: 'Anita S.',
      icon: 'E',
      color: 'red',
    },
    {
      name: 'Art Education',
      code: 'ART-502',
      type: 'Co-Scholastic',
      category: 'Elective',
      createdBy: 'Anita S.',
      icon: 'A',
      color: 'teal',
    },
  ];
}
