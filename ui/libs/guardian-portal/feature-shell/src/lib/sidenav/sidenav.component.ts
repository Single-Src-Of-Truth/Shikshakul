import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';

interface MenuItem {
  label: string;
  icon: string;
  route: string;
}

interface MenuGroup {
  header?: string;
  items: MenuItem[];
}

@Component({
  selector: 'shikshakul-guardian-sidenav',
  standalone: true,
  imports: [CommonModule, RouterLink, RouterLinkActive],
  templateUrl: './sidenav.component.html',
  styleUrl: './sidenav.component.scss',
})
export class SidenavComponent {
  @Input() collapsed = false;
  @Output() toggleCollapse = new EventEmitter<void>();

  menuGroups: MenuGroup[] = [
    {
      items: [{ label: 'Dashboard', icon: 'dashboard', route: '/dashboard' }],
    },
    {
      header: 'ACADEMICS',
      items: [
        { label: 'Attendance', icon: 'date_range', route: '/attendance' },
        { label: 'Report Cards', icon: 'assignment', route: '/reports' },
        { label: 'Exam Schedule', icon: 'event_note', route: '/exams' },
      ],
    },
    {
      header: 'COMMUNICATION',
      items: [{ label: 'Class Teacher', icon: 'school', route: '/teacher' }],
    },
  ];

  triggerSidebarToggle() {
    this.toggleCollapse.emit();
  }
}
