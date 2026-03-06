import { CommonModule } from '@angular/common';
import { Component, inject, OnInit, ChangeDetectorRef } from '@angular/core';
import { forkJoin, of } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { ClassStatsComponent } from '../components/class-stats/class-stats.component';
import { ClassFormComponent } from '../components/class-form/class-form.component';
import {
  ClassGrade,
  ClassManagementService,
} from '@shikshakul/data-access/academic';
import { SectionFormComponent } from '../components/section-form/section-form.component';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';
import { SubjectAllocationDialogComponent } from '../components/subject-allocation-dialog/subject-allocation-dialog.component';
import {
  StudentService,
  AcademicYearService,
} from '@shikshakul/data-access/academic';
import { ClassListItemComponent } from '../components/class-list-item/class-list-item.component';

interface UISectionItem {
  id: string;
  name: string;
  capacity: number;
  studentCount: number;
}

export interface UIClassItem extends ClassGrade {
  stream: string;
  description: string;
  sections: UISectionItem[];
  sectionNames: string[];
  studentCount: number;
}

@Component({
  selector: 'shikshakul-class-management-page',
  standalone: true,
  imports: [
    CommonModule,
    ClassStatsComponent,
    ClassFormComponent,
    SectionFormComponent,
    SubjectAllocationDialogComponent,
    ClassListItemComponent,
  ],
  templateUrl: './class-management-page.component.html',
  styleUrl: './class-management-page.component.scss',
})
export class ClassManagementPageComponent implements OnInit {
  private classService = inject(ClassManagementService);
  private studentService = inject(StudentService);
  private acadYearService = inject(AcademicYearService);
  private snackbar = inject(SnackbarService);
  private cdr = inject(ChangeDetectorRef);

  currentSessionName = '...';

  classes: UIClassItem[] = [];
  loading = true;

  totalClasses = 0;
  totalSections = 0;
  totalStudents = 0;

  showDrawer = false;
  showSectionDialog = false;
  showSubjectDialog = false;

  selectedClassData: Partial<ClassGrade> = { name: '', sort_order: 0 };
  selectedClassForSection: UIClassItem | null = null;
  selectedClassForSubjects: UIClassItem | null = null;

  isEditMode = false;

  ngOnInit(): void {
    this.fetchCurrentYear();
    this.loadClasses();
  }

  fetchCurrentYear() {
    this.acadYearService.getCurrentAcademicYear().subscribe({
      next: (res) => {
        if (res.data) {
          this.currentSessionName = res.data.name;
          this.cdr.detectChanges();
        }
      },
    });
  }

  loadClasses() {
    this.loading = true;
    this.cdr.detectChanges();
    this.classService.getClasses().subscribe({
      next: (apiData) => {
        let dataArr = Array.isArray(apiData)
          ? apiData
          : (apiData as any)?.data || [];

        dataArr = dataArr
          .sort((a: any, b: any) => (a.sort_order || 0) - (b.sort_order || 0))
          .map((cls: any) => ({
            ...cls,
            id: cls.id || cls.class_id || cls._id,
            name: cls.name || cls.class_name,
            sectionNames: [],
            studentCount: 0,
          }));

        if (dataArr.length === 0) {
          this.classes = [];
          this.loading = false;
          this.cdr.detectChanges();
          return;
        }

        const requests = dataArr.map((cls: any) => {
          const sectionsReq = this.classService
            .getSectionsByClass(cls.id)
            .pipe(catchError(() => of({ data: [] })));
          const studentsReq = this.studentService
            .getStudents({ class_id: cls.id })
            .pipe(catchError(() => of({ data: [] })));

          return forkJoin({
            sections: sectionsReq,
            students: studentsReq,
          }).pipe(
            map(({ sections, students }) => {
              const sectionData = Array.isArray(sections)
                ? sections
                : sections?.data || [];
              const studentData = Array.isArray(students)
                ? students
                : students?.data || [];

              cls.sections = sectionData.map((s: any) => ({
                id: s.id,
                name: s.name,
                capacity: s.capacity || 40,
                studentCount: studentData.filter(
                  (stu: any) => stu.section_id === s.id,
                ).length,
              }));

              cls.sectionNames = cls.sections.map((s: any) => s.name);
              cls.studentCount = studentData.length;
              return cls;
            }),
          );
        });

        forkJoin(requests).subscribe({
          next: (enrichedClasses: any) => {
            this.classes = enrichedClasses;
            this.calculateTotals();

            this.loading = false;
            this.cdr.detectChanges();
          },
          error: () => {
            this.classes = dataArr; // fallback
            this.loading = false;
            this.cdr.detectChanges();
          },
        });
      },
      error: () => {
        this.snackbar.error('Error', 'Failed to load classes.');
        this.loading = false;
        this.cdr.detectChanges();
      },
    });
  }

  openDrawer() {
    this.selectedClassData = {
      name: '',
      sort_order: this.classes.length + 1,
    };
    this.isEditMode = false;
    this.showDrawer = true;
  }

  closeClassDialog() {
    this.showDrawer = false;
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
            this.snackbar.success(
              'Class Updated',
              `Class ${data.name} updated.`,
            );
          },
          error: (err) => {
            console.error(err);
            this.snackbar.error('Error', 'Failed to update class.');
          },
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

  editClass(cls: UIClassItem) {
    this.selectedClassData = {
      id: cls.id,
      name: cls.name,
      sort_order: cls.sort_order,
    };
    this.isEditMode = true;
    this.showDrawer = true;
  }

  deleteClass(classId: string) {
    const cls = this.classes.find((c) => c.id === classId);
    if (!cls) return;

    if (confirm(`Are you sure you want to delete class ${cls.name}?`)) {
      this.classService.deleteClass(cls.id).subscribe({
        next: () => {
          this.snackbar.success(
            'Class Deleted',
            `Class ${cls.name} was removed.`,
          );
          this.loadClasses();
        },
        error: (err) => {
          this.snackbar.error('Error', 'Failed to delete class.');
          console.error(err);
        },
      });
    }
  }

  openSectionDialog(cls: UIClassItem) {
    this.selectedClassForSection = cls;
    this.showSectionDialog = true;
  }

  closeSectionDialog() {
    this.showSectionDialog = false;
    // Retain selectedClassForSection so formData is not lost when reopening for the same class
  }

  onSectionSaved() {
    this.loadClasses();
    this.snackbar.success(
      'Section Added',
      `New section added to ${this.selectedClassForSection?.name}.`,
    );
  }

  removeSection(classId: string, sectionName: string) {
    const cls = this.classes.find((c) => c.id === classId);
    if (!cls) return;

    if (
      confirm(
        `Are you sure you want to remove section ${sectionName} from ${cls.name}?`,
      )
    ) {
      this.classService.getSectionsByClass(cls.id).subscribe({
        next: (sections) => {
          const sectionToDelete = sections.data?.find(
            (s) => s.name === sectionName,
          );
          if (sectionToDelete) {
            this.classService.deleteSection(sectionToDelete.id).subscribe({
              next: () => {
                this.snackbar.success(
                  'Section Removed',
                  `Section ${sectionName} deleted.`,
                );
                this.loadClasses();
              },
              error: (err) => {
                this.snackbar.error('Error', 'Failed to remove section.');
                console.error(err);
              },
            });
          } else {
            this.snackbar.error('Error', 'Section not found.');
          }
        },
        error: (err) => {
          this.snackbar.error('Error', 'Could not fetch sections to delete.');
          console.error(err);
        },
      });
    }
  }

  openSubjectDialog(cls: UIClassItem) {
    this.selectedClassForSubjects = cls;
    this.showSubjectDialog = true;
  }

  closeSubjectDialog() {
    this.showSubjectDialog = false;
    // Retain selectedClassForSubjects so selections are not lost when reopening for the same class
  }

  private calculateTotals() {
    this.totalClasses = this.classes.length;
    this.totalSections = this.classes.reduce(
      (acc, cls) => acc + (cls.sections?.length || 0),
      0,
    );
    this.totalStudents = this.classes.reduce(
      (acc, cls) => acc + (cls.studentCount || 0),
      0,
    );
  }
}
