import { CommonModule } from '@angular/common';
import {
  Component,
  inject,
  OnInit,
  signal,
  ChangeDetectorRef,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import {
  ExamService,
  MarkSheetResponse,
  EvaluationStudent,
  SubmitMarksRequest,
  StudentService,
} from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';

@Component({
  selector: 'shikshakul-evaluation-mark-sheet',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './evaluation-mark-sheet.component.html',
  styleUrl: './evaluation-mark-sheet.component.scss',
})
export class EvaluationMarkSheetComponent implements OnInit {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private examService = inject(ExamService);
  private studentService = inject(StudentService);
  private snackbar = inject(SnackbarService);
  private cdr = inject(ChangeDetectorRef);

  scheduleId = '';
  markSheet = signal<MarkSheetResponse | null>(null);
  loading = signal(false);
  submitting = signal(false);

  ngOnInit(): void {
    this.scheduleId = this.route.snapshot.params['scheduleId'];
    if (this.scheduleId) {
      this.loadMarkSheet();
    }
  }

  loadMarkSheet() {
    this.loading.set(true);

    this.examService.getMarksSheet(this.scheduleId).subscribe({
      next: (res: any) => {
        const rawData = res?.data || res;
        let students: EvaluationStudent[] = [];
        let metadata: any = {};

        if (Array.isArray(rawData)) {
          students = rawData;
        } else {
          students = rawData.students || rawData.marks || [];
          metadata = rawData;
        }

        const needsNames = students.length > 0 && !students[0].first_name;

        if (needsNames && metadata.class_id) {
          this.studentService
            .getStudents({ class_id: metadata.class_id })
            .subscribe({
              next: (studentListRes: any) => {
                const allStudents = studentListRes?.data || studentListRes;
                const enrichedStudents = students.map((marksEntry) => {
                  const profile = allStudents.find(
                    (p: any) => p.id === marksEntry.student_id,
                  );
                  return {
                    ...marksEntry,
                    first_name: profile?.first_name || 'N/A',
                    last_name: profile?.last_name || '',
                    roll_number: profile?.roll_number || 'N/A',
                  };
                });

                this.setMarkSheet(metadata, enrichedStudents);
              },
              error: () => this.setMarkSheet(metadata, students),
            });
        } else {
          this.setMarkSheet(metadata, students);
        }
      },
      error: () => {
        this.snackbar.error('Error', 'Failed to load mark sheet.');
        this.loading.set(false);
        this.cdr.detectChanges();
      },
    });
  }

  private setMarkSheet(metadata: any, students: EvaluationStudent[]) {
    this.markSheet.set({
      exam_schedule_id: this.scheduleId,
      subject_name: metadata.subject_name || metadata.subject?.name,
      term_name: metadata?.exam_term?.name || metadata?.term_name,
      class_name: metadata?.class?.name || metadata?.class_name,
      max_marks: metadata.max_marks,
      pass_marks: metadata.pass_marks,
      students: students,
    });
    this.loading.set(false);
    this.cdr.detectChanges();
  }

  onMarksChange(student: EvaluationStudent, value: string) {
    const marks = parseFloat(value);
    if (isNaN(marks)) {
      student.marks_obtained = undefined;
      return;
    }

    const maxMarks = this.markSheet()?.max_marks || 100;
    if (marks > maxMarks) {
      this.snackbar.warning(
        'Warning',
        `Marks cannot exceed maximum (${maxMarks})`,
      );
      student.marks_obtained = maxMarks;
    } else if (marks < 0) {
      student.marks_obtained = 0;
    } else {
      student.marks_obtained = marks;
    }
  }

  toggleAbsent(student: EvaluationStudent) {
    student.is_absent = !student.is_absent;
    if (student.is_absent) {
      student.marks_obtained = 0;
    }
  }

  save(finalize = false) {
    const sheet = this.markSheet();
    if (!sheet) return;

    if (
      finalize &&
      !confirm(
        'Are you sure you want to finalize? This will lock the marks and they cannot be edited further.',
      )
    ) {
      return;
    }

    this.submitting.set(true);
    const payload: SubmitMarksRequest = {
      exam_schedule_id: this.scheduleId,
      finalize,
      marks: sheet.students.map((s) => ({
        student_id: s.student_id,
        marks_obtained: s.marks_obtained || 0,
        is_absent: s.is_absent,
        remarks: s.remarks,
      })),
    };

    this.examService.submitMarks(payload).subscribe({
      next: () => {
        this.snackbar.success(
          'Success',
          finalize
            ? 'Marks finalized successfully.'
            : 'Draft saved successfully.',
        );
        this.submitting.set(false);
        if (finalize) {
          this.loadMarkSheet(); // Reload to get read-only state
        }
      },
      error: (err) => {
        this.submitting.set(false);
        this.snackbar.error(
          'Error',
          err.error?.message || 'Failed to submit marks.',
        );
      },
    });
  }

  goBack() {
    this.router.navigate(['/exams/evaluation']);
  }

  isFinalized(): boolean {
    return this.markSheet()?.students[0]?.status === 'FINALIZED';
  }

  // UI Helpers
  getInitials(s: EvaluationStudent): string {
    const f = s.first_name ? s.first_name[0] : '';
    const l = s.last_name ? s.last_name[0] : '';
    return (f + l).toUpperCase() || '?';
  }

  getAvatarColor(name: string): string {
    const colors = [
      'bg-blue',
      'bg-purple',
      'bg-teal',
      'bg-orange',
      'bg-pink',
      'bg-green',
    ];
    const idx = name ? name.charCodeAt(0) % colors.length : 0;
    return colors[idx];
  }
}
