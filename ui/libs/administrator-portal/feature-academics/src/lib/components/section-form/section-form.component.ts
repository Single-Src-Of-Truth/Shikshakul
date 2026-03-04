import { Component, EventEmitter, Input, Output, OnInit, OnChanges, SimpleChanges, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { ClassManagementService } from '@shikshakul/data-access/academic';

@Component({
    selector: 'shikshakul-section-form',
    standalone: true,
    imports: [CommonModule, ReactiveFormsModule],
    templateUrl: './section-form.component.html',
    styleUrl: './section-form.component.scss',
})
export class SectionFormComponent implements OnInit, OnChanges {
    @Input() classId!: string;
    @Input() className = '';
    @Input() isOpen = false;

    @Output() closeDrawer = new EventEmitter<void>();
    @Output() sectionSaved = new EventEmitter<void>();

    private fb = inject(FormBuilder);
    private classService = inject(ClassManagementService);

    sectionForm!: FormGroup;
    submitting = false;
    private lastClassId: string | null = null;

    ngOnInit() {
        this.initForm();
    }

    ngOnChanges(changes: SimpleChanges) {
        if (changes['classId'] && this.classId !== this.lastClassId) {
            this.lastClassId = this.classId;
            if (this.sectionForm) {
                this.sectionForm.reset({
                    name: '',
                    capacity: 40
                });
            }
        }
    }

    initForm() {
        this.sectionForm = this.fb.group({
            name: ['', Validators.required],
            capacity: [40, [Validators.required, Validators.min(1)]]
        });
    }

    onSubmit() {
        if (this.sectionForm.invalid || !this.classId) {
            Object.keys(this.sectionForm.controls).forEach(key => {
                this.sectionForm.get(key)?.markAsTouched();
            });
            return;
        }

        this.submitting = true;
        this.classService
            .createSection(this.classId, {
                name: this.sectionForm.value.name,
                capacity: this.sectionForm.value.capacity,
            })
            .subscribe({
                next: () => {
                    this.submitting = false;
                    this.sectionSaved.emit();
                    this.closeDrawer.emit();
                    this.sectionForm.reset({ name: '', capacity: 40 });
                },
                error: () => {
                    this.submitting = false;
                    alert('Failed to create section.');
                },
            });
    }

    onCancel() {
        this.closeDrawer.emit();
    }

    resetForm() {
        if (confirm('Are you sure you want to clear the form?')) {
            this.sectionForm.reset({
                name: '',
                capacity: 40
            });
        }
    }
}
