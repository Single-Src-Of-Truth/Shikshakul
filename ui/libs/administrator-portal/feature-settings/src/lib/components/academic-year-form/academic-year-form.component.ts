import { Component, EventEmitter, Input, Output, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { AcademicYearService } from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';
import { DatepickerComponent } from '@shikshakul/shared/ui/datepicker';

@Component({
    selector: 'shikshakul-academic-year-form',
    standalone: true,
    imports: [CommonModule, ReactiveFormsModule, DatepickerComponent],
    templateUrl: './academic-year-form.component.html',
    styleUrl: './academic-year-form.component.scss',
})
export class AcademicYearFormComponent implements OnInit {
    @Input() isOpen = false;
    @Output() closeDrawer = new EventEmitter<void>();
    @Output() yearCreated = new EventEmitter<void>();

    private fb = inject(FormBuilder);
    private acadYearService = inject(AcademicYearService);
    private snackbar = inject(SnackbarService);

    yearForm!: FormGroup;
    submitting = false;
    dynamicYearExample = '';

    ngOnInit() {
        const currentYear = new Date().getFullYear();
        this.dynamicYearExample = `${currentYear}-${currentYear + 1}`;
        this.initForm();
    }

    initForm() {
        this.yearForm = this.fb.group({
            name: ['', [Validators.required, Validators.pattern(/^\d{4}-\d{4}$/)]],
            start_date: ['', Validators.required],
            end_date: ['', Validators.required],
        });
    }

    close() {
        this.closeDrawer.emit();
    }

    resetForm() {
        if (confirm('Are you sure you want to clear the form?')) {
            this.yearForm.reset();
        }
    }

    onSubmit() {
        // Mark all as touched to show validation errors inline if invalid
        if (this.yearForm.invalid) {
            Object.keys(this.yearForm.controls).forEach(key => {
                this.yearForm.get(key)?.markAsTouched();
            });
            return;
        }

        this.submitting = true;
        const formValue = this.yearForm.value;

        // Convert dates to ISO-8601 format
        const payload = {
            ...formValue,
            start_date: new Date(formValue.start_date).toISOString(),
            end_date: new Date(formValue.end_date).toISOString()
        };


        this.acadYearService.createAcademicYear(payload).subscribe({
            next: () => {
                this.snackbar.success('Success', 'Academic year created successfully.');
                this.submitting = false;
                this.yearCreated.emit();
                this.yearForm.reset();
                this.close();
            },
            error: (err) => {
                this.submitting = false;
                this.snackbar.error('Error', err.error?.message || 'Failed to create academic year.');
            }
        });
    }
}
