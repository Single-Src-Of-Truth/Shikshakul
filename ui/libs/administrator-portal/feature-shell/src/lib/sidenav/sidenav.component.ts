import { Component, EventEmitter, Input, Output, inject, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { AuthStore } from '@shikshakul/auth';

interface MenuItem {
  label: string;
  icon: string;
  route: string;
  /** If set, the item is only visible when user has this IAM permission */
  requiredPermission?: string;
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

  readonly authStore = inject(AuthStore);

  /** User's initials for the avatar */
  readonly initials = computed(() => {
    const p = this.authStore.profile();
    if (!p) return '?';
    return `${p.first_name[0] ?? ''}${p.last_name?.[0] ?? ''}`.toUpperCase();
  });

  /** Full display name */
  readonly displayName = computed(() => {
    const p = this.authStore.profile();
    if (!p) return 'Loading...';
    return `${p.first_name} ${p.last_name ?? ''}`.trim();
  });

  /** Primary role label */
  readonly primaryRole = computed(() => {
    const roles = this.authStore.roles();
    return roles[0] ?? '';
  });

  private readonly allMenuGroups: MenuGroup[] = [
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
        { label: 'Roles & Permissions', icon: 'manage_accounts', route: '/settings/roles' },
      ],
    },
    {
      // Super Admin only section — hidden for normal school users
      header: 'SUPER ADMIN',
      items: [
        {
          label: 'Manage Schools',
          icon: 'corporate_fare',
          route: '/super-admin/tenants',
          requiredPermission: 'iam:tenants:read',
        },
      ],
    },
  ];

  /**
   * Menu groups filtered by the current user's permissions.
   * Groups with no visible items are removed entirely.
   */
  readonly menuGroups = computed(() => {
    return this.allMenuGroups
      .map((group) => ({
        ...group,
        items: group.items.filter((item) => {
          if (!item.requiredPermission) return true;
          return this.authStore.hasPermission(item.requiredPermission);
        }),
      }))
      .filter((group) => group.items.length > 0);
  });

  triggerSidebarToggle(): void {
    this.toggleCollapse.emit();
  }

  logout(): void {
    this.authStore.logout();
  }
}
