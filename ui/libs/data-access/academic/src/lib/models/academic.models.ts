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
