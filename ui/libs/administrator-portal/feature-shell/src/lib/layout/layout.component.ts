import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { SidenavComponent } from '../sidenav/sidenav.component';

import { AcademicYearService, ApiResponse, AcademicYear } from '@shikshakul/data-access/academic';
import { Observable } from 'rxjs';

import { HostListener, OnInit } from '@angular/core';

@Component({
  selector: 'shikshakul-layout',
  standalone: true,
  imports: [CommonModule, RouterOutlet, SidenavComponent],
  templateUrl: './layout.component.html',
  styleUrl: './layout.component.scss',
})
export class LayoutComponent implements OnInit {
  isSidebarCollapsed = false;
  isMobile = false;
  showMobileMenu = false;

  private academicYearService = inject(AcademicYearService);
  currentYear$: Observable<ApiResponse<AcademicYear>> = this.academicYearService.getCurrentAcademicYear();

  ngOnInit() {
    this.checkScreenSize();
  }

  @HostListener('window:resize')
  onResize() {
    this.checkScreenSize();
  }

  private checkScreenSize() {
    const width = window.innerWidth;
    this.isMobile = width <= 768;

    // Auto-collapse on tablet, hide on mobile
    if (width <= 1024 && width > 768) {
      this.isSidebarCollapsed = true;
    } else if (width > 1024) {
      this.isSidebarCollapsed = false;
    }

    if (!this.isMobile) {
      this.showMobileMenu = false;
    }
  }

  toggleSidebar() {
    if (this.isMobile) {
      this.showMobileMenu = !this.showMobileMenu;
    } else {
      this.isSidebarCollapsed = !this.isSidebarCollapsed;
    }
  }

  closeMobileMenu() {
    this.showMobileMenu = false;
  }
}
