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
  @Input() mobileActive = false;
  @Output() toggleCollapse = new EventEmitter<void>();

  menuGroups: MenuGroup[] = [
    {
      items: [{ label: 'Dashboard', icon: 'dashboard', route: '/dashboard' }],
    },
    {
      header: 'SETUP',
      items: [
        { label: 'Academic Years', icon: 'calendar_today', route: '/settings/academic-years' },
        { label: 'Class Setup', icon: 'class', route: '/academics/classes' },
        { label: 'Subjects', icon: 'library_books', route: '/academics/subjects' },
      ],
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
        { label: 'Attendance', icon: 'fact_check', route: '/academics/attendance' },
        { label: 'Timetable', icon: 'calendar_view_week', route: '/academics/timetable' },
      ],
    },
    {
      header: 'EXAMS',
      items: [
        { label: 'Exam Setup', icon: 'assignment', route: '/exams/setup' },
        { label: 'Results', icon: 'auto_stories', route: '/exams/results' },
      ],
    },
    {
      header: 'ADMINISTRATION',
      items: [
        { label: 'Fee Management', icon: 'payments', route: '/fees' },
        { label: 'School Calendar', icon: 'event', route: '/academics/calendar' },
      ],
    },
  ];

  triggerSidebarToggle() {
    this.toggleCollapse.emit();
  }
}
