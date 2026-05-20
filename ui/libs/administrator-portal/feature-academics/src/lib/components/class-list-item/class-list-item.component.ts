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
  @Input() classData: any;

  @Output() addSectionClick = new EventEmitter<void>();
  @Output() removeSection = new EventEmitter<string>();
  @Output() manageSubjects = new EventEmitter<void>();
  @Output() edit = new EventEmitter<void>();
  @Output() delete = new EventEmitter<void>();

  isExpanded = false;

  onAddSection(event: Event) {
    event.stopPropagation();
    this.addSectionClick.emit();
  }

  toggleExpand() {
    this.isExpanded = !this.isExpanded;
  }
}
