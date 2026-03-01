import { Component, EventEmitter, Input, Output } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TCResponse, BonafideResponse } from '@shikshakul/data-access/academic';

@Component({
  selector: 'shikshakul-certificate-preview',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './certificate-preview.component.html',
  styleUrl: './certificate-preview.component.scss',
})
export class CertificatePreviewComponent {
  @Input() data: TCResponse | BonafideResponse | null = null;
  @Input() type: 'TC' | 'BONAFIDE' = 'TC';
  @Output() close = new EventEmitter<void>();

  printCertificate() {
    window.print();
  }

  isTC(data: any): data is TCResponse {
    return this.type === 'TC';
  }

  isBonafide(data: any): data is BonafideResponse {
    return this.type === 'BONAFIDE';
  }
}
