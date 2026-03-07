import { Component, EventEmitter, Input, Output, OnInit, OnChanges, SimpleChanges, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { ClassGrade } from '@shikshakul/data-access/academic';

@Component({
  selector: 'shikshakul-class-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './class-form.component.html',
  styleUrl: './class-form.component.scss',
})
export class ClassFormComponent implements OnInit, OnChanges {
  @Input() classData: Partial<ClassGrade> = { name: '', sort_order: 0 };
  @Input() isEditMode = false;
  @Input() isOpen = false;

  @Output() closeDrawer = new EventEmitter<void>();
  @Output() save = new EventEmitter<Partial<ClassGrade>>();

  private fb = inject(FormBuilder);
  classForm!: FormGroup;
  submitting = false;

  ngOnInit() {
    this.initForm();
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes['classData'] && this.classForm) {
      if (this.isEditMode) {
        this.classForm.patchValue({
          name: this.classData.name || '',
          sort_order: this.classData.sort_order || 0
        });
      } else if (changes['classData'].firstChange || !this.classForm.get('name')?.value) {
        this.classForm.patchValue({
          name: this.classData.name || '',
          sort_order: this.classData.sort_order || 0
        });
      }
    }
  }

  initForm() {
    this.classForm = this.fb.group({
      name: [this.classData.name || '', Validators.required],
      sort_order: [this.classData.sort_order || 0, [Validators.required, Validators.min(1)]]
    });
  }

  onSubmit() {
    if (this.classForm.invalid) {
      Object.keys(this.classForm.controls).forEach(key => {
        this.classForm.get(key)?.markAsTouched();
      });
      return;
    }
    this.save.emit({
      ...this.classData,
      name: this.classForm.value.name,
      sort_order: Number(this.classForm.value.sort_order)
    });
  }

  onCancel() {
    this.closeDrawer.emit();
  }

  resetForm() {
    if (confirm('Are you sure you want to clear the form?')) {
      this.classForm.reset({
        name: '',
        sort_order: this.classData.sort_order || 1
      });
    }
  }
}
