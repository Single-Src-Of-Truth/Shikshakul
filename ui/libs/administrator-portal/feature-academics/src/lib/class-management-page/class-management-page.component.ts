import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { ClassListItemComponent } from '../components/class-list-item/class-list-item.component';
import { ClassStatsComponent } from '../components/class-stats/class-stats.component';

@Component({
  selector: 'shikshakul-class-management-page',
  standalone: true,
  imports: [CommonModule, ClassListItemComponent, ClassStatsComponent],
  templateUrl: './class-management-page.component.html',
  styleUrl: './class-management-page.component.scss',
})
export class ClassManagementPageComponent {
  classes = [
    {
      name: 'Class X',
      stream: 'General',
      description: 'Standard 10 • Secondary',
      sections: ['A', 'B', 'C', 'D'],
      studentCount: 124,
    },
    {
      name: 'Class XI',
      stream: 'Science',
      description: 'Standard 11 • Senior Secondary',
      sections: ['A', 'B'],
      studentCount: 85,
    },
    {
      name: 'Class XI',
      stream: 'Commerce',
      description: 'Standard 11 • Senior Secondary',
      sections: ['C', 'D'],
      studentCount: 92,
    },
    {
      name: 'Class XII',
      stream: 'Arts',
      description: 'Standard 12 • Senior Secondary',
      sections: ['E'],
      studentCount: 40,
    },
    {
      name: 'Class I',
      stream: 'General',
      description: 'Standard 1 • Primary',
      sections: [],
      studentCount: 0,
    },
  ];
}
