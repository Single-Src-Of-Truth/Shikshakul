import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { SidenavComponent } from '../sidenav/sidenav.component';

import { AcademicYearService, ApiResponse, AcademicYear } from '@shikshakul/data-access/academic';
import { Observable } from 'rxjs';

@Component({
  selector: 'shikshakul-layout',
  imports: [CommonModule, RouterOutlet, SidenavComponent],
  templateUrl: './layout.component.html',
  styleUrl: './layout.component.scss',
})
export class LayoutComponent {
  isSidebarCollapsed = false;
  private academicYearService = inject(AcademicYearService);

  currentYear$: Observable<ApiResponse<AcademicYear>> = this.academicYearService.getCurrentAcademicYear();

  toggleSidebar() {
    this.isSidebarCollapsed = !this.isSidebarCollapsed;
  }
}
