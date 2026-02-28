import { CommonModule } from '@angular/common';
import {
  Component,
  EventEmitter,
  inject,
  Input,
  OnInit,
  Output,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import {
  ClassManagementService,
  Subject,
  SubjectService,
} from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';

interface SelectableSubject extends Subject {
  isSelected: boolean;
  weeklyLectures: number;
  isOptional: boolean;
}

@Component({
  selector: 'shikshakul-subject-allocation-dialog',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './subject-allocation-dialog.component.html',
  styleUrl: './subject-allocation-dialog.component.scss',
})
export class SubjectAllocationDialogComponent implements OnInit {
  private subjectService = inject(SubjectService);
  private classService = inject(ClassManagementService);
  private snackbar = inject(SnackbarService);

  @Input() classId!: string;
  @Input() className!: string;
  @Output() close = new EventEmitter<void>();
  @Output() save = new EventEmitter<void>();

  subjects: SelectableSubject[] = [];
  loading = true;
  isSubmitting = false;

  ngOnInit(): void {
    this.loadSubjects();
  }

  loadSubjects() {
    this.loading = true;
    this.subjectService.getSubjects().subscribe({
      next: (data) => {
        this.subjects = data.map((s) => ({
          ...s,
          isSelected: false,
          weeklyLectures: 5,
          isOptional: false,
        }));
        this.loading = false;
      },
      error: () => {
        this.snackbar.error('Error', 'Failed to load subjects list.');
        this.loading = false;
      },
    });
  }

  toggleSelection(subject: SelectableSubject) {
    subject.isSelected = !subject.isSelected;
  }

  onSubmit() {
    const selected = this.subjects.filter((s) => s.isSelected);

    if (selected.length === 0) {
      this.snackbar.info('Info', 'Please select at least one subject.');
      return;
    }

    this.isSubmitting = true;

    const payload = {
      subjects: selected.map((s) => ({
        subject_id: s.subject_id!,
        is_optional: s.isOptional,
        weekly_lectures: s.weeklyLectures,
      })),
    };

    this.classService.assignSubjectToClass(this.classId, payload).subscribe({
      next: () => {
        this.snackbar.success(
          'Success',
          `Subjects assigned to ${this.className}`,
        );
        this.save.emit();
        this.close.emit();
      },
      error: (err) => {
        console.error(err);
        this.snackbar.error('Error', 'Failed to assign subjects.');
        this.isSubmitting = false;
      },
    });
  }

  onCancel() {
    this.close.emit();
  }
}
