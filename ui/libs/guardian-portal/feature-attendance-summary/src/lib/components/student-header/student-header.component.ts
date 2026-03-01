import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'shikshakul-student-header',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './student-header.component.html',
  styleUrl: './student-header.component.scss',
})
export class StudentHeaderComponent {
  @Input() student: any;
}
