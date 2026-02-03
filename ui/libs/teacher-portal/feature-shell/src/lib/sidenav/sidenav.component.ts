import { Component, EventEmitter, Input, Output } from '@angular/core';
import { CommonModule } from '@angular/common';
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
  selector: 'shikshakul-teacher-sidenav',
  standalone: true,
  imports: [CommonModule, RouterLink, RouterLinkActive],
  templateUrl: './sidenav.component.html',
  styleUrls: ['./sidenav.component.scss'],
})
export class SidenavComponent {
  @Input() collapsed = false;
  @Output() toggleCollapse = new EventEmitter<void>();

  menuGroups: MenuGroup[] = [
    {
      items: [{ label: 'Dashboard', icon: 'dashboard', route: '/dashboard' }],
    },
    {
      header: 'PEOPLE',
      items: [{ label: 'My Students', icon: 'school', route: '/students' }],
    },
    {
      header: 'ACADEMICS',
      items: [
        { label: 'My Classes', icon: 'class', route: '/classes' },
        // { label: 'Timetable', icon: 'calendar_month', route: '/timetable' },
        { label: 'Attendance', icon: 'fact_check', route: '/attendance' },
      ],
    },
    {
      header: 'EXAMINATION',
      items: [
        { label: 'Marks Entry', icon: 'edit_note', route: '/marks' },
        // { label: 'Exam Reports', icon: 'bar_chart', route: '/reports' },
      ],
    },
  ];

  triggerSidebarToggle() {
    this.toggleCollapse.emit();
  }
}
