import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'shikshakul-class-list-item',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './class-list-item.component.html',
  styleUrl: './class-list-item.component.scss',
})
export class ClassListItemComponent {
  @Input() classData: any; // We'll define the interface later
  @Output() edit = new EventEmitter<void>();
  @Output() delete = new EventEmitter<void>();
  @Output() addSection = new EventEmitter<void>();
  @Output() removeSection = new EventEmitter<string>();
}
