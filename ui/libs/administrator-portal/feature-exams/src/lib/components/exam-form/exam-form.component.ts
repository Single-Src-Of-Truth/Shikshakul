import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import {
  FormArray,
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';

@Component({
  selector: 'shikshakul-exam-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './exam-form.component.html',
  styleUrl: './exam-form.component.scss',
})
export class ExamFormComponent {
  examForm: FormGroup;

  constructor(private fb: FormBuilder) {
    this.examForm = this.fb.group({
      examName: ['', Validators.required],
      academicYear: ['2024-2025', Validators.required],
      classGrade: ['', Validators.required],
      gradingSystem: ['Marks Based (CBSE)', Validators.required],
      subjects: this.fb.array([
        this.createSubjectRow('Mathematics', 80, 26),
        this.createSubjectRow('Science', 80, 26),
      ]),
    });
  }

  // Helper to create a new row
  createSubjectRow(subject = '', max = 0, pass = 0): FormGroup {
    return this.fb.group({
      subjectName: [subject, Validators.required],
      maxMarks: [max, Validators.required],
      passMarks: [pass, Validators.required],
    });
  }

  // Getter for easy access in HTML
  get subjectsArray(): FormArray {
    return this.examForm.get('subjects') as FormArray;
  }

  addSubject() {
    this.subjectsArray.push(this.createSubjectRow('', 100, 33));
  }

  removeSubject(index: number) {
    this.subjectsArray.removeAt(index);
  }

  onSubmit() {
    if (this.examForm.valid) {
      console.log(this.examForm.value);
    }
  }
}
