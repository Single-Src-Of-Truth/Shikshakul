import { CommonModule } from '@angular/common';
import { Component, EventEmitter, inject, OnInit, Output, ChangeDetectorRef } from '@angular/core';
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
  private cdr = inject(ChangeDetectorRef);

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
    this.cdr.detectChanges();
    this.classService.getClasses().subscribe({
      next: (apiData) => {
        const dataArr = Array.isArray(apiData) ? apiData : ((apiData as any)?.data || []);
        this.classes = dataArr
          .sort((a: any, b: any) => (a.sort_order || 0) - (b.sort_order || 0))
          .map((cls: any) => ({
            ...cls,
            sections: [],
            studentCount: 0,
          }));
        this.loading = false;
        this.cdr.detectChanges();
      },
      error: () => {
        this.snackbar.error('Error', 'Failed to load classes.');
        this.loading = false;
        this.cdr.detectChanges();
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

    if (this.isEditMode && data.id) {
      this.classService
        .updateClass(data.id, {
          name: data.name,
          sort_order: Number(data.sort_order),
        })
        .subscribe({
          next: () => {
            this.loadClasses();
            this.closeClassDialog();
            this.snackbar.success('Class Updated', `Class ${data.name} updated.`);
          },
          error: (err) => {
            console.error(err);
            this.snackbar.error('Error', 'Failed to update class.');
          }
        });
    } else {
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
  }

  onEditClass(cls: UIClassItem) {
    this.selectedClassData = {
      id: cls.id,
      name: cls.name,
      sort_order: cls.sort_order,
    };
    this.isEditMode = true;
    this.showClassDialog = true;
  }

  onDeleteClass(cls: UIClassItem) {
    if (confirm(`Are you sure you want to delete class ${cls.name}?`)) {
      this.classService.deleteClass(cls.id).subscribe({
        next: () => {
          this.snackbar.success('Class Deleted', `Class ${cls.name} was removed.`);
          this.loadClasses();
        },
        error: (err) => {
          this.snackbar.error('Error', 'Failed to delete class.');
          console.error(err);
        }
      });
    }
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

  onRemoveSection(cls: UIClassItem, sectionName: string) {
    if (confirm(`Are you sure you want to remove section ${sectionName} from ${cls.name}?`)) {
      this.classService.getSectionsByClass(cls.id).subscribe({
        next: (sections) => {
          const sectionToDelete = sections.find(s => s.name === sectionName);
          if (sectionToDelete) {
            this.classService.deleteSection(sectionToDelete.id).subscribe({
              next: () => {
                this.snackbar.success('Section Removed', `Section ${sectionName} deleted.`);
                this.loadClasses(); // Reload classes to refresh sections
              },
              error: (err) => {
                this.snackbar.error('Error', 'Failed to remove section.');
                console.error(err);
              }
            });
          } else {
            this.snackbar.error('Error', 'Section not found.');
          }
        },
        error: (err) => {
          this.snackbar.error('Error', 'Could not fetch sections to delete.');
          console.error(err);
        }
      });
    }
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
