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
        { label: 'Exams', icon: 'assignment', route: '/exams/setup' },
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
