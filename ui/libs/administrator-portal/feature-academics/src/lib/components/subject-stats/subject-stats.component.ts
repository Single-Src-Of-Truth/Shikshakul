import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';

@Component({
  selector: 'shikshakul-subject-stats',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './subject-stats.component.html',
  styleUrl: './subject-stats.component.scss',
})
export class SubjectStatsComponent {
  @Input() total = 42;
  @Input() theory = 28;
  @Input() practical = 14;
  @Input() other = 6;
}
