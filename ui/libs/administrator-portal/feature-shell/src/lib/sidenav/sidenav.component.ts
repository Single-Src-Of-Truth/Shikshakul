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
  selector: 'shikshakul-sidenav',
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
      items: [
        { label: 'Students', icon: 'school', route: '/academics/students' },
        { label: 'Staff', icon: 'badge', route: '/staff' },
      ],
    },
    {
      header: 'ACADEMICS',
      items: [
        { label: 'Class Setup', icon: 'class', route: '/academics/classes' },
        {
          label: 'Subjects',
          icon: 'library_books',
          route: '/academics/subjects',
        },
        {
          label: 'Attendance',
          icon: 'fact_check',
          route: '/academics/attendance',
        },
        {
          label: 'Timetable',
          icon: 'calendar_view_week',
          route: '/academics/timetable',
        },
        {
          label: 'Certificates',
          icon: 'badge',
          route: '/academics/certificates',
        },
        {
          label: 'Calendar',
          icon: 'event',
          route: '/academics/calendar',
        },
        { label: 'Exams', icon: 'assignment', route: '/exams/setup' },
        { label: 'Evaluation', icon: 'edit_note', route: '/exams/evaluation' },
        { label: 'Results', icon: 'auto_stories', route: '/exams/results' },
      ],
    },
    {
      header: 'ADMINISTRATION',
      items: [
        {
          label: 'Fee Management',
          icon: 'payments',
          route: '/fees',
        },
        // { label: 'Transport', icon: 'directions_bus', route: '/transport' }, // Future
      ],
    },
    {
      header: 'CONFIGURATION',
      items: [
        {
          label: 'Roles & Permissions',
          icon: 'admin_panel_settings',
          route: '/settings/roles',
        },
        {
          label: 'Academic Years',
          icon: 'calendar_today',
          route: '/settings/academic-years',
        },
      ],
    },
  ];

  triggerSidebarToggle() {
    this.toggleCollapse.emit();
  }
}
