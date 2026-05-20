import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { TeacherProfileCardComponent } from '../components/teacher-profile-card/teacher-profile-card.component';
import { AboutTeacherComponent } from '../components/about-teacher/about-teacher.component';
import { SubjectsTaughtComponent } from '../components/subjects-taught/subjects-taught.component';
import { ContactInfoComponent } from '../components/contact-info/contact-info.component';
import { VisitingHoursComponent } from '../components/visiting-hours/visiting-hours.component';

@Component({
  selector: 'shikshakul-class-teacher-page',
  standalone: true,
  imports: [
    CommonModule,
    TeacherProfileCardComponent,
    AboutTeacherComponent,
    SubjectsTaughtComponent,
    ContactInfoComponent,
    VisitingHoursComponent,
  ],
  templateUrl: './class-teacher-page.component.html',
  styleUrl: './class-teacher-page.component.scss',
})
export class ClassTeacherPageComponent {}
