import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'shikshakul-role-list',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './role-list.component.html',
  styleUrl: './role-list.component.scss',
})
export class RoleListComponent {
  @Input() roles: any[] = [];
  @Input() selectedRoleId: string | null = null;
  @Output() selectRole = new EventEmitter<string>();
  @Output() createRole = new EventEmitter<void>();
}
