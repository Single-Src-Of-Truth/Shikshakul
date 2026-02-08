import { CommonModule } from '@angular/common';
import { Component, EventEmitter, inject, OnInit, Output } from '@angular/core';
import { ClassListItemComponent } from '../components/class-list-item/class-list-item.component';
import { ClassStatsComponent } from '../components/class-stats/class-stats.component';
import { ClassDialogComponent } from '../components/class-dialog/class-dialog.component';
import {
  ClassGrade,
  ClassManagementService,
} from '@shikshakul/data-access/academic';
import { SectionDialogComponent } from '../components/section-dialog/section-dialog.component';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';
import { SubjectAllocationDialogComponent } from '../components/subject-allocation-dialog/subject-allocation-dialog.component';

interface UIClassItem extends ClassGrade {
  stream: string;
  description: string;
  sections: string[];
  studentCount: number;
}

@Component({
  selector: 'shikshakul-class-management-page',
  standalone: true,
  imports: [
    CommonModule,
    ClassListItemComponent,
    ClassStatsComponent,
    ClassDialogComponent,
    SectionDialogComponent,
    SubjectAllocationDialogComponent,
  ],
  templateUrl: './class-management-page.component.html',
  styleUrl: './class-management-page.component.scss',
})
export class ClassManagementPageComponent implements OnInit {
  private classService = inject(ClassManagementService);
  private snackbar = inject(SnackbarService);

  classes: UIClassItem[] = [];
  loading = true;

  showClassDialog = false;
  showSectionDialog = false;
  showSubjectDialog = false;

  selectedClassData: Partial<ClassGrade> = { name: '', sort_order: 0 };
  selectedClassForSection: UIClassItem | null = null;
  selectedClassForSubjects: UIClassItem | null = null;

  isEditMode = false;

  ngOnInit(): void {
    this.loadClasses();
  }

  loadClasses() {
    this.loading = true;
    this.classService.getClasses().subscribe({
      next: (apiData) => {
        this.classes = apiData
          .sort((a, b) => a.sort_order - b.sort_order)
          .map((cls) => ({
            ...cls,
            stream: 'General',
            description: `Standard ${cls.sort_order} • Academic`,
            sections: [],
            studentCount: 0,
          }));
        this.loading = false;
      },
      error: () => {
        this.snackbar.error('Error', 'Failed to load classes.');
        this.loading = false;
      },
    });
  }

  openCreateClassDialog() {
    this.selectedClassData = {
      name: '',
      sort_order: (this.classes.length + 1) * 10,
    };
    this.isEditMode = false;
    this.showClassDialog = true;
  }

  closeClassDialog() {
    this.showClassDialog = false;
  }

  saveClass(data: Partial<ClassGrade>) {
    if (!data.name || !data.sort_order) return;

    this.classService
      .createClass({
        name: data.name,
        sort_order: Number(data.sort_order),
      })
      .subscribe({
        next: () => {
          this.loadClasses();
          this.closeClassDialog();
          this.snackbar.success(
            'Class Created',
            `Class ${data.name} has been successfully added.`,
          );
        },
        error: (err) => {
          console.error(err);
          this.snackbar.error(
            'Error',
            'Failed to create the class. Please try again.',
          );
        },
      });
  }

  // TODO
  onEditClass(cls: UIClassItem) {
    // Implement Edit Logic
    this.snackbar.info(
      'Coming Soon',
      `Edit functionality for ${cls.name} is under development.`,
    );
  }

  // TODO
  onDeleteClass(cls: UIClassItem) {
    // Call delete API here
    this.snackbar.info(
      'Coming Soon',
      `Delete functionality is under development.`,
    );
  }

  openSectionDialog(cls: UIClassItem) {
    this.selectedClassForSection = cls;
    this.showSectionDialog = true;
  }

  closeSectionDialog() {
    this.showSectionDialog = false;
    this.selectedClassForSection = null;
  }

  onSectionSaved() {
    if (!this.selectedClassForSection?.id) return;
    const classId = this.selectedClassForSection.id;

    this.classService.getSectionsByClass(classId).subscribe({
      next: (sections) => {
        const classIdx = this.classes.findIndex((c) => c.id === classId);
        if (classIdx !== -1) {
          this.classes[classIdx].sections = sections.map((s) => s.name);
        }

        this.snackbar.success(
          'Section Added',
          `New section added to ${this.selectedClassForSection?.name}.`,
        );
      },
      error: () => {
        this.snackbar.error(
          'Warning',
          'Section saved, but failed to refresh list.',
        );
      },
    });
  }

  // TODO
  onRemoveSection(cls: UIClassItem, sectionName: string) {
    this.snackbar.info(
      'Coming Soon',
      `Removing sections is under development.`,
    );
  }

  openSubjectDialog(cls: UIClassItem) {
    this.selectedClassForSubjects = cls;
    this.showSubjectDialog = true;
  }

  closeSubjectDialog() {
    this.showSubjectDialog = false;
    this.selectedClassForSubjects = null;
  }
}
