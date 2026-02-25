import { CommonModule } from '@angular/common';
import { Component, inject, OnInit, ChangeDetectorRef } from '@angular/core';
import { FormArray, FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { StudentService } from '@shikshakul/data-access/academic';
import { SnackbarService } from '@shikshakul/shared/ui/snackbar';

@Component({
  selector: 'lib-student-schema-config',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './student-schema-config.component.html',
  styleUrl: './student-schema-config.component.scss',
})
export class StudentSchemaConfigComponent implements OnInit {
  private fb = inject(FormBuilder);
  private studentService = inject(StudentService);
  private snackbar = inject(SnackbarService);
  private cdr = inject(ChangeDetectorRef);

  schemaForm!: FormGroup;
  loading = false;
  saving = false;

  ngOnInit(): void {
    this.initForm();
    this.loadSchema();
  }

  initForm() {
    this.schemaForm = this.fb.group({
      schema_name: ['Student Admission Schema', Validators.required],
      fields: this.fb.array([]),
    });
  }

  get fields() {
    return this.schemaForm.get('fields') as FormArray;
  }

  addField(existingData?: any) {
    const fieldForm = this.fb.group({
      field_name: [existingData?.field_name || '', Validators.required],
      field_type: [existingData?.field_type || 'text', Validators.required],
      is_required: [existingData?.is_required || false],
      order_index: [existingData?.order_index || this.fields.length],
    });
    this.fields.push(fieldForm);
  }

  removeField(index: number) {
    this.fields.removeAt(index);
  }

  loadSchema() {
    this.loading = true;
    this.studentService.getActiveSchema().subscribe({
      next: (res) => {
        const schema = res?.data || res;
        if (schema && schema.fields) {
          this.schemaForm.patchValue({ schema_name: schema.schema_name });
          schema.fields.forEach((f: any) => this.addField(f));
        } else {
          // Add a default field if no schema exists
          this.addField({ field_name: 'Custom Info', field_type: 'text' });
        }
        this.loading = false;
        this.cdr.detectChanges();
      },
      error: () => {
        this.loading = false;
        // Start fresh on error (assuming 404 No Schema)
        this.addField({ field_name: 'Custom Info', field_type: 'text' });
        this.cdr.detectChanges();
      }
    });
  }

  saveSchema() {
    if (this.schemaForm.invalid) {
      this.snackbar.error('Error', 'Please fill all required fields');
      return;
    }
    this.saving = true;
    this.studentService.createSchema(this.schemaForm.value).subscribe({
      next: () => {
        this.snackbar.success('Success', 'Schema updated successfully');
        this.saving = false;
        this.cdr.detectChanges();
      },
      error: (err) => {
        this.snackbar.error('Error', 'Failed to update schema');
        this.saving = false;
        console.error(err);
        this.cdr.detectChanges();
      }
    });
  }

  deleteSchema() {
    if (confirm('Are you sure you want to delete the active schema? All dynamic student fields will be lost.')) {
      this.saving = true;
      this.studentService.deleteSchema().subscribe({
        next: () => {
          this.snackbar.success('Success', 'Schema deleted');
          this.fields.clear();
          this.addField(); // Reset to empty state
          this.saving = false;
          this.cdr.detectChanges();
        },
        error: () => {
          this.snackbar.error('Error', 'Failed to delete schema');
          this.saving = false;
          this.cdr.detectChanges();
        }
      });
    }
  }
}
