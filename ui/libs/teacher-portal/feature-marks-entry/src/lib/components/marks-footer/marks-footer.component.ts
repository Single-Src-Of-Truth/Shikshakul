import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Output } from '@angular/core';

@Component({
  selector: 'shikshakul-marks-footer',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './marks-footer.component.html',
  styleUrl: './marks-footer.component.scss',
})
export class MarksFooterComponent {
  @Output() saveDraft = new EventEmitter<void>();
  @Output() submit = new EventEmitter<void>();
}
