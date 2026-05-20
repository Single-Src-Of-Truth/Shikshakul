import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { FormGroup, ReactiveFormsModule } from '@angular/forms';
import { ClassGrade, Section } from '@shikshakul/data-access/academic';

@Component({
  selector: 'shikshakul-academic-details',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './academic-details.component.html',
  styleUrl: './academic-details.component.scss',
})
export class AcademicDetailsComponent {
  @Input() formGroup!: FormGroup;
  @Input() classes: ClassGrade[] = [];
  @Input() sections: Section[] = [];
}
