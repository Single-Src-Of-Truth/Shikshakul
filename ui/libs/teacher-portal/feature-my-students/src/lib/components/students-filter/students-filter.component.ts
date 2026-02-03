import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';

@Component({
  selector: 'shikshakul-students-filter',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './students-filter.component.html',
  styleUrl: './students-filter.component.scss',
})
export class StudentsFilterComponent {}
