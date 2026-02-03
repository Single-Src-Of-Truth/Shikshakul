import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'shikshakul-permission-detail',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './permission-detail.component.html',
  styleUrl: './permission-detail.component.scss',
})
export class PermissionDetailComponent {
  @Input() role: any;
  @Input() permissionGroups: any[] = [];
  @Output() save = new EventEmitter<void>();

  togglePermission(groupIndex: number, permIndex: number) {
    const perm = this.permissionGroups[groupIndex].permissions[permIndex];
    perm.enabled = !perm.enabled;
  }

  toggleGroup(groupIndex: number, event: any) {
    const isChecked = event.target.checked;
    this.permissionGroups[groupIndex].permissions.forEach(
      (p: any) => (p.enabled = isChecked),
    );
  }
}
