export interface ApiResponse<T> {
  data: T;
  message?: string;
  success?: boolean;
}

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
  sections?: Section[];
}

export interface EvaluationStudent {
  student_id: string;
  first_name: string;
  last_name: string;
  roll_number: string;
  marks_obtained?: number;
  is_absent: boolean;
  remarks?: string;
  status?: string; // e.g., 'DRAFT', 'FINALIZED'
}

export interface MarkSheetResponse {
  exam_schedule_id: string;
  subject_name: string;
  term_name?: string;
  class_name?: string;
  max_marks: number;
  pass_marks: number;
  students: EvaluationStudent[];
}

export interface SubmitMarksRequest {
  exam_schedule_id: string;
  finalize: boolean;
  marks: {
    student_id: string;
    marks_obtained: number;
    is_absent: boolean;
    remarks?: string;
  }[];
}

export interface Section {
  id: string;
  name: string;
  capacity: number;
  class_id?: string;
}

export enum StudentStatus {
  PENDING = 'PENDING',
  ACTIVE = 'ACTIVE',
  INACTIVE = 'INACTIVE',
  APPROVED = 'APPROVED',
  REJECTED = 'REJECTED',
  SUSPENDED = 'SUSPENDED',
}

export enum StudentAction {
  APPROVE = 'APPROVE',
  REJECT = 'REJECT',
  PENDING = 'PENDING',
}

export type SubjectType = 'THEORY' | 'PRACTICAL' | 'LAB';

export type FeeFrequency = 'MONTHLY' | 'YEARLY' | 'ONE_TIME';

export type PaymentMode = 'CASH' | 'CHEQUE' | 'UPI' | 'ONLINE';

export type FeeStatus = 'PENDING' | 'PARTIAL' | 'PAID' | 'OVERDUE';

export type DayOfWeek =
  | 'MONDAY'
  | 'TUESDAY'
  | 'WEDNESDAY'
  | 'THURSDAY'
  | 'FRIDAY'
  | 'SATURDAY'
  | 'SUNDAY';

export type EvaluationStatus = 'DRAFT' | 'SUBMITTED' | 'PUBLISHED';

export type AttendanceStatus =
  | 'PRESENT'
  | 'ABSENT'
  | 'LATE'
  | 'HALF_DAY'
  | 'LEAVE';

export type CertificateType =
  | 'ID_CARD'
  | 'TRANSFER_CERTIFICATE'
  | 'BONAFIDE'
  | 'CHARACTER';

export type EventType = 'HOLIDAY' | 'EXAM' | 'EVENT' | 'MEETING';

export interface Subject {
  id: string;
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
  is_published: boolean;
  is_active: boolean;
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
  duration_min: number;
  room_number?: string;
  max_marks: number;
  pass_marks: number;
  subject_name?: string;
  subject?: Subject;
  exam_term?: ExamTerm;
  evaluation_status?: string;
}

export interface CreateExamScheduleRequest {
  exam_term_id: string;
  class_id: string;
  subject_id: string;
  exam_date: string;
  start_time: string;
  end_time: string;
  duration_min: number;
  room_number?: string;
  max_marks: number;
  pass_marks: number;
}

export interface GenerateResultRequest {
  exam_term_id: string;
  class_id: string;
}

export interface PublishResultRequest {
  exam_term_id: string;
  class_id: string;
  publish: boolean;
}

export interface StudentSummary {
  name: string;
  admission_no: string;
  roll_no: string;
  class_name: string;
}

export interface ExamSummary {
  term_name: string;
}

export interface SubjectMark {
  subject_name: string;
  max_marks: number;
  marks_obtained: number;
  grade: string;
  remarks: string;
  is_absent: boolean;
}

export interface ResultSummary {
  total_max: number;
  total_obtained: number;
  percentage: number;
  grade: string;
  status: string;
}

export interface ReportCardResponse {
  student_info: StudentSummary;
  exam_info: ExamSummary;
  subjects: SubjectMark[];
  summary: ResultSummary;
}

export interface AttendanceEntry {
  student_id: string;
  status: AttendanceStatus | string;
  remarks?: string;
}

export interface MarkAttendanceRequest {
  class_id: string;
  section_id: string;
  date: string; // ISO string format for time.Time
  students: AttendanceEntry[];
}

export interface RecordDetail {
  student_id: string;
  student_name: string;
  roll_no: string;
  status: AttendanceStatus | string;
  remarks?: string;
}

export interface AttendanceStat {
  total: number;
  present: number;
  absent: number;
  late: number;
}

export interface AttendanceRegisterResponse {
  date: string;
  section_id: string;
  records: RecordDetail[];
  summary: AttendanceStat;
}

export interface StudentAttendanceHistory {
  date: string;
  status: AttendanceStatus | string;
  remarks?: string;
}

export interface IDCardResponse {
  student_id: string;
  admission_no: string;
  full_name: string;
  class_details: string;
  dob: string;
  blood_group: string;
  father_name: string;
  contact_number: string;
  address: string;
  photo_url: string;
}

export interface IssueTCRequest {
  student_id: string;
  leaving_date: string; // ISO string
  reason: string;
  conduct: string;
  mark_inactive?: boolean;
}

export interface TCResponse {
  certificate_no: string;
  issue_date: string;
  student_name: string;
  father_name: string;
  mother_name: string;
  admission_no: string;
  dob: string;
  nationality: string;
  class_joined: string;
  last_class: string;
  result_status: string;
  reason: string;
  conduct: string;
}

export interface IssueBonafideRequest {
  student_id: string;
  purpose: string;
}

export interface BonafideResponse {
  certificate_no: string;
  issue_date: string;
  student_name: string;
  admission_no: string;
  class_details: string;
  academic_year: string;
  purpose: string;
}

export type CalendarEventType =
  | 'HOLIDAY'
  | 'EXAM'
  | 'PTM'
  | 'CULTURAL'
  | 'OTHER';

export interface CreateEventRequest {
  academic_year_id: string;
  title: string;
  description?: string;
  event_type: CalendarEventType | string;
  start_date: string; // ISO string
  end_date: string; // ISO string
  is_holiday?: boolean;
}

export interface UpdateEventRequest {
  title?: string;
  description?: string;
  event_type?: CalendarEventType | string;
  start_date?: string;
  end_date?: string;
  is_holiday?: boolean;
}

export interface EventResponse {
  id: string;
  title: string;
  description: string;
  event_type: CalendarEventType | string;
  start_date: string;
  end_date: string;
  is_holiday: boolean;
}
