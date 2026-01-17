import { Component, EventEmitter, Input, Output } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink, RouterLinkActive } from '@angular/router';

interface MenuItem {
  label: string;
  icon: string;
  route: string;
}

interface MenuGroup {
  header?: string; // Optional header like "ACADEMICS"
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

  // Organized into Visual Groups (No clicking required)
  menuGroups: MenuGroup[] = [
    {
      // Main Dashboard (No header needed)
      items: [{ label: 'Dashboard', icon: 'dashboard', route: '/dashboard' }],
    },
    {
      header: 'PEOPLE',
      items: [
        { label: 'Students', icon: 'school', route: '/students/create' },
        { label: 'Staff', icon: 'badge', route: '/staff/create' },
      ],
    },
    {
      header: 'ACADEMICS',
      items: [
        { label: 'Class Setup', icon: 'class', route: '/academics' },
        {
          label: 'Subjects',
          icon: 'library_books',
          route: '/academics/subjects',
        },
        { label: 'Exams', icon: 'assignment', route: '/exams/setup' },
      ],
    },
    {
      header: 'ADMINISTRATION',
      items: [
        { label: 'Fees & Finance', icon: 'payments', route: '/fees' },
        { label: 'Transport', icon: 'directions_bus', route: '/transport' },
      ],
    },
  ];

  triggerSidebarToggle() {
    this.toggleCollapse.emit();
  }
}
