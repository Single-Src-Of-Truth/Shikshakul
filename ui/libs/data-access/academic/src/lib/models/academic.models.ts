export interface AcademicYear {
  id: string;
  name: string;
  start_date: string;
  end_date: string;
  is_current: boolean;
  academic_year_id?: string;
}

export interface ClassGrade {
  id: string;
  name: string;
  sort_order: number;
}

export interface Section {
  id: string;
  name: string;
  capacity: number;
  class_id?: string;
}

export type Status = 'PENDING' | 'ACTIVE' | 'INACTIVE' | 'APPROVED' | 'REJECTED' | 'SUSPENDED';

export type SubjectType = 'THEORY' | 'PRACTICAL' | 'LAB';

export type FeeFrequency = 'MONTHLY' | 'YEARLY' | 'ONE_TIME';

export type PaymentMode = 'CASH' | 'CHEQUE' | 'UPI' | 'ONLINE';

export type FeeStatus = 'PENDING' | 'PARTIAL' | 'PAID' | 'OVERDUE';

export type DayOfWeek = 'MONDAY' | 'TUESDAY' | 'WEDNESDAY' | 'THURSDAY' | 'FRIDAY' | 'SATURDAY' | 'SUNDAY';

export type EvaluationStatus = 'DRAFT' | 'SUBMITTED' | 'PUBLISHED';

export type AttendanceStatus = 'PRESENT' | 'ABSENT' | 'LATE' | 'HALF_DAY' | 'LEAVE';

export type CertificateType = 'ID_CARD' | 'TRANSFER_CERTIFICATE' | 'BONAFIDE' | 'CHARACTER';

export type EventType = 'HOLIDAY' | 'EXAM' | 'EVENT' | 'MEETING';

export interface Subject {
  subject_id?: string;
  name: string;
  code: string;
  type: SubjectType;
  category?: string;
  created_by?: string;
}

export interface StudentOnboardData {
  academic_year_id: string;
  class_id: string;
  section_id?: string;
  first_name: string;
  last_name: string;
  email?: string;
  mobile?: string;
  dob: string;
  gender: string;
  profile_data: any;
  documents?: any;
}

export interface TeacherOnboardData {
  first_name: string;
  last_name: string;
  email: string;
  mobile: string;
  profile_data: {
    dob?: string;
    gender?: string;
    aadhar_number?: string;
    current_address?: any;
    employment?: any;
    banking?: any;
    [key: string]: any;
  };
  documents?: any;
}

export interface FeeHead {
  id?: string;
  name: string;
  type: 'RECURRING' | 'ONE_TIME';
}

export interface FeeStructure {
  id?: string;
  academic_year_id: string;
  class_id: string;
  fee_head_id: string;
  amount: number;
  frequency: FeeFrequency;
  due_date_day?: number;
}

export interface Routine {
  id: string;
  day_of_week: DayOfWeek;
  start_time: string;
  end_time: string;
  subject_name: string;
  teacher_name: string;
  room_number?: string;
}

export interface CreateRoutineRequest {
  academic_year_id: string;
  class_id: string;
  section_id: string;
  subject_id: string;
  teacher_id: string;
  day_of_week: DayOfWeek;
  start_time: string;
  end_time: string;
  room_number?: string;
}

export interface ExamTerm {
  id: string;
  name: string;
  academic_year_id: string;
  start_date: string;
  end_date: string;
  status: Status;
}

export interface CreateExamTermRequest {
  name: string;
  academic_year_id: string;
  start_date: string;
  end_date: string;
}

export interface ExamSchedule {
  id: string;
  exam_term_id: string;
  class_id: string;
  subject_id: string;
  exam_date: string;
  start_time: string;
  end_time: string;
  room_number?: string;
  max_marks: number;
  pass_marks: number;
  subject_name?: string;
  term_name?: string;
}

export interface CreateExamScheduleRequest {
  exam_term_id: string;
  class_id: string;
  subject_id: string;
  exam_date: string;
  start_time: string;
  end_time: string;
  room_number?: string;
  max_marks: number;
  pass_marks: number;
}
