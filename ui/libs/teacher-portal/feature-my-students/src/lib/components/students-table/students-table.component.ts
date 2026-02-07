import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'shikshakul-students-table',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './students-table.component.html',
  styleUrl: './students-table.component.scss',
})
export class StudentsTableComponent {
  @Input() students: any[] = [];
}
