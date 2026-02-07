import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { ProfileBannerComponent } from '../components/profile-banner/profile-banner.component';
import { PersonalDetailsComponent } from '../components/personal-details/personal-details.component';
import { FamilyInfoComponent } from '../components/family-info/family-info.component';
import { ContactDetailsComponent } from '../components/contact-details/contact-details.component';
import { AcademicInfoComponent } from '../components/academic-info/academic-info.component';
import { DocumentsListComponent } from '../components/documents-list/documents-list.component';

@Component({
  selector: 'shikshakul-student-profile-page',
  standalone: true,
  imports: [
    CommonModule,
    ProfileBannerComponent,
    PersonalDetailsComponent,
    FamilyInfoComponent,
    ContactDetailsComponent,
    AcademicInfoComponent,
    DocumentsListComponent,
  ],
  templateUrl: './student-profile-page.component.html',
  styleUrl: './student-profile-page.component.scss',
})
export class StudentProfilePageComponent {}
