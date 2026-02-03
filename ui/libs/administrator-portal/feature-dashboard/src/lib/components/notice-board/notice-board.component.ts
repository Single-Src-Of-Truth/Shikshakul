import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';

@Component({
  selector: 'shikshakul-notice-board',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './notice-board.component.html',
  styleUrl: './notice-board.component.scss',
})
export class NoticeBoardComponent {
  notices = [
    {
      day: '24',
      month: 'OCT',
      title: 'CBSE Board Exam Date Sheet',
      description:
        'The tentative date sheet for Class X and XII has been released by the board.',
      type: 'info', // Blue badge
    },
    {
      day: '28',
      month: 'OCT',
      title: 'Staff Meeting: Monthly Review',
      description:
        'Mandatory attendance for all department heads in the conference room.',
      type: 'warning', // Orange badge (for meetings/alerts)
    },
    {
      day: '12',
      month: 'NOV',
      title: 'Diwali Holiday Announcement',
      description: 'School will remain closed from Nov 12th to Nov 15th.',
      type: 'urgent', // Red badge (for holidays/closures)
    },
  ];
}
