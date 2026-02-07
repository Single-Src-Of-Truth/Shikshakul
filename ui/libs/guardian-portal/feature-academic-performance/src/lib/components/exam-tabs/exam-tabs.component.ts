import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';

@Component({
  selector: 'shikshakul-exam-tabs',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './exam-tabs.component.html',
  styleUrl: './exam-tabs.component.scss',
})
export class ExamTabsComponent {
  activeTab = 'half';
}
