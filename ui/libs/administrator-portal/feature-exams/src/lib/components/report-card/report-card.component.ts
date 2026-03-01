import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { ReportCardResponse } from '@shikshakul/data-access/academic';

@Component({
    selector: 'lib-report-card',
    standalone: true,
    imports: [CommonModule],
    templateUrl: './report-card.component.html',
    styleUrl: './report-card.component.scss',
})
export class ReportCardComponent {
    @Input({ required: true }) data!: ReportCardResponse;

    print() {
        window.print();
    }
}
