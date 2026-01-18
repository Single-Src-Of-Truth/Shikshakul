import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { ClassStatsComponent } from '../components/class-stats/class-stats.component';
import { ClassCardComponent } from '../components/class-card/class-card.component';

@Component({
  selector: 'shikshakul-my-classes-page',
  standalone: true,
  imports: [CommonModule, ClassStatsComponent, ClassCardComponent],
  templateUrl: './my-classes-page.component.html',
  styleUrl: './my-classes-page.component.scss',
})
export class MyClassesPageComponent {
  classes = [
    {
      name: 'Class X - A',
      subject: 'Mathematics',
      students: 42,
      nextTime: 'Next: Today, 10:00 AM',
      room: 'Room 101, Main Block',
      color: 'blue',
      isToday: true,
    },
    {
      name: 'Class IX - B',
      subject: 'Science (Physics)',
      students: 38,
      nextTime: 'Next: Tomorrow, 09:00 AM',
      room: 'Lab 2, Science Block',
      color: 'purple',
      isToday: false,
    },
    {
      name: 'Class XII - Science',
      subject: 'Lab Assistant',
      students: 20,
      nextTime: 'Next: Today, 02:00 PM',
      room: 'Physics Lab',
      color: 'green',
      isToday: true,
    },
    {
      name: 'Class VIII - C',
      subject: 'Hindi',
      students: 40,
      nextTime: 'Next: Mon, 11:30 AM',
      room: 'Room 204',
      color: 'red',
      isToday: false,
    },
    {
      name: 'Class XI - B',
      subject: 'Computer Science',
      students: 28,
      nextTime: 'Next: Today, 01:00 PM',
      room: 'Computer Lab 1',
      color: 'teal',
      isToday: true,
    },
  ];
}
