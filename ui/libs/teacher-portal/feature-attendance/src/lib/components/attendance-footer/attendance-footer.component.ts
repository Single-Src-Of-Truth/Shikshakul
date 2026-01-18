import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'shikshakul-attendance-footer',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './attendance-footer.component.html',
  styleUrl: './attendance-footer.component.scss',
})
export class AttendanceFooterComponent {
  @Input() total = 0;
  @Input() present = 0;
  @Input() absent = 0;
  @Input() late = 0;
  @Output() save = new EventEmitter<void>();
  @Output() cancel = new EventEmitter<void>();
}
