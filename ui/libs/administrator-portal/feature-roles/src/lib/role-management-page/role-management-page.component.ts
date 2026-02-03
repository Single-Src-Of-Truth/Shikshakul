import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RoleListComponent } from '../components/role-list/role-list.component';
import { PermissionDetailComponent } from '../components/permission-detail/permission-detail.component';

@Component({
  selector: 'shikshakul-role-management-page',
  standalone: true,
  imports: [CommonModule, RoleListComponent, PermissionDetailComponent],
  templateUrl: './role-management-page.component.html',
  styleUrl: './role-management-page.component.scss',
})
export class RoleManagementPageComponent {
  roles = [
    {
      id: '1',
      name: 'Super Admin',
      description: 'System Owner',
      icon: 'admin_panel_settings',
      colorClass: '',
    },
    {
      id: '2',
      name: 'Principal',
      description: 'Full Access (minus settings)',
      icon: 'person',
      colorClass: 'purple',
    },
    {
      id: '3',
      name: 'Class Teacher',
      description: 'Academic & Student Mgmt',
      icon: 'school',
      colorClass: 'blue',
    },
    {
      id: '4',
      name: 'Accountant',
      description: 'Finance Module Only',
      icon: 'account_balance_wallet',
      colorClass: 'green',
    },
    {
      id: '5',
      name: 'Transport In-charge',
      description: 'Routes & Vehicles',
      icon: 'directions_bus',
      colorClass: 'orange',
    },
  ];

  selectedRoleId = '3';

  // Current Permissions for selected role
  permissionGroups = [
    {
      name: 'Student Information System',
      icon: 'school',
      permissions: [
        {
          label: 'View Student Profile',
          description: 'Basic details, parents info',
          enabled: true,
        },
        {
          label: 'Edit Student Profile',
          description: 'Update address, phone',
          enabled: false,
        },
        {
          label: 'Manage Attendance',
          description: 'Mark daily attendance',
          enabled: true,
        },
        {
          label: 'Issue TC',
          description: 'Generate Transfer Cert',
          enabled: false,
        },
      ],
    },
    {
      name: 'Examinations & Results',
      icon: 'assignment',
      permissions: [
        {
          label: 'Enter Subject Marks',
          description: 'Internal/Term assessments',
          enabled: true,
        },
        {
          label: 'Generate Report Cards',
          description: 'View/Print student report',
          enabled: true,
        },
        {
          label: 'Publish Results',
          description: 'Make visible to parents',
          enabled: false,
        },
      ],
    },
    {
      name: 'Fee Management',
      icon: 'payments',
      permissions: [
        {
          label: 'View Fee Status',
          description: 'Check pending dues',
          enabled: true,
        },
        {
          label: 'Collect Fees',
          description: 'Process payments',
          enabled: false,
        },
      ],
    },
  ];

  get selectedRole() {
    return this.roles.find((r) => r.id === this.selectedRoleId);
  }

  onRoleSelect(id: string) {
    this.selectedRoleId = id;
    // In a real app, you would fetch permissions for this ID here
  }

  onSave() {
    console.log(
      'Saving permissions for role:',
      this.selectedRoleId,
      this.permissionGroups,
    );
  }
}
