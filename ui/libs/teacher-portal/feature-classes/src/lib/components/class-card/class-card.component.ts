import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';

@Component({
  selector: 'shikshakul-class-card',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './class-card.component.html',
  styleUrl: './class-card.component.scss',
})
export class ClassCardComponent {
  @Input() classData: any;
}
