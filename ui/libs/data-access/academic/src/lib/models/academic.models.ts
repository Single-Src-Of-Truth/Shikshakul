export interface AcademicYear {
  id: string;
  name: string;
  start_date: string;
  end_date: string;
  is_current: boolean;
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

export interface Subject {
  id?: string;
  name: string;
  code: string;
  type: 'THEORY' | 'PRACTICAL' | 'CO_SCHOLASTIC';
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
