import { CommonModule } from '@angular/common';
import { Component, inject, OnInit, ChangeDetectorRef } from '@angular/core';
import { SubjectStatsComponent } from '../components/subject-stats/subject-stats.component';
import { SubjectDialogComponent } from '../components/subject-dialog/subject-dialog.component';
import { Subject, SubjectService } from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';

interface UISubject extends Subject {
  icon: string;
  color: string;
  createdBy: string;
}

@Component({
  selector: 'shikshakul-subject-setup-page',
  standalone: true,
  imports: [CommonModule, SubjectStatsComponent, SubjectDialogComponent],
  templateUrl: './subject-setup-page.component.html',
  styleUrl: './subject-setup-page.component.scss',
})
export class SubjectSetupPageComponent implements OnInit {
  private subjectService = inject(SubjectService);
  private snackbar = inject(SnackbarService);
  private cdr = inject(ChangeDetectorRef);

  subjects: UISubject[] = [];
  loading = true;

  stats = { total: 0, theory: 0, practical: 0, other: 0 };

  showDialog = false;
  selectedSubject: Partial<Subject> = { name: '', code: '', type: 'THEORY' };
  isEditMode = false;

  ngOnInit(): void {
    this.loadSubjects();
  }

  loadSubjects() {
    this.loading = true;
    this.cdr.detectChanges();
    this.subjectService.getSubjects().subscribe({
      next: (res: any) => {
        const data = Array.isArray(res) ? res : (res?.data || []);
        this.subjects = data.map((sub: any) => ({
          ...sub,
          icon: sub.name?.charAt(0)?.toUpperCase() || 'S',
          color: this.getColorForType(sub.type),
          createdBy: 'Admin User',
          category: this.getCategoryForType(sub.type),
        }));

        this.calculateStats();
        this.loading = false;
        this.cdr.detectChanges();
      },
      error: (err) => {
        console.error(err);
        this.snackbar.error('Error', 'Failed to load subjects');
        this.loading = false;
        this.cdr.detectChanges();
      },
    });
  }

  calculateStats() {
    this.stats.total = this.subjects.length;
    this.stats.theory = this.subjects.filter((s) => s.type === 'THEORY').length;
    this.stats.practical = this.subjects.filter(
      (s) => s.type === 'LAB' || s.type === 'PRACTICAL',
    ).length;
    this.stats.other = 0;
  }

  getColorForType(type: string): string {
    switch (type) {
      case 'THEORY':
        return 'blue';
      case 'LAB':
      case 'PRACTICAL':
        return 'green';
      default:
        return 'blue';
    }
  }

  getCategoryForType(type: string): string {
    switch (type) {
      case 'THEORY':
        return 'Compulsory (Core)';
      case 'LAB':
      case 'PRACTICAL':
        return 'Lab / Practical';
      default:
        return 'General';
    }
  }

  openAddDialog() {
    this.selectedSubject = { name: '', code: '', type: 'THEORY' };
    this.isEditMode = false;
    this.showDialog = true;
  }

  closeDialog() {
    this.showDialog = false;
  }

  onSave(subject: Partial<Subject>) {
    if (!subject.name || !subject.code) return;

    this.subjectService.createSubject(subject as Subject).subscribe({
      next: () => {
        this.snackbar.success(
          'Subject Created',
          `${subject.name} added successfully.`,
        );
        this.loadSubjects();
        this.closeDialog();
      },
      error: () => this.snackbar.error('Error', 'Failed to create subject.'),
    });
  }
}
