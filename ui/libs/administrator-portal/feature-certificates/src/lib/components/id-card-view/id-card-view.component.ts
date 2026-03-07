import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { IDCardResponse } from '@shikshakul/data-access/academic';

@Component({
  selector: 'shikshakul-id-card-view',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './id-card-view.component.html',
  styleUrl: './id-card-view.component.scss',
})
export class IdCardViewComponent {
  @Input() idCards: IDCardResponse[] = [];
  @Output() close = new EventEmitter<void>();

  printPage() {
    window.print();
  }
}
