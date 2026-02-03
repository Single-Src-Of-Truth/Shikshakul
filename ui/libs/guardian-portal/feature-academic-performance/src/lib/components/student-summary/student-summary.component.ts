import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';

@Component({
  selector: 'shikshakul-student-summary',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './student-summary.component.html',
  styleUrl: './student-summary.component.scss',
})
export class StudentSummaryComponent {}
