package dto

type GenerateResultRequest struct {
	ExamTermID string `json:"exam_term_id" binding:"required"`
	ClassID    string `json:"class_id" binding:"required"`
}

type PublishResultRequest struct {
	ExamTermID string `json:"exam_term_id" binding:"required"`
	ClassID    string `json:"class_id" binding:"required"`
	Publish    bool   `json:"publish"`
}

type ReportCardResponse struct {
	StudentInfo StudentSummary `json:"student_info"`
	ExamInfo    ExamSummary    `json:"exam_info"`
	Subjects    []SubjectMark  `json:"subjects"`
	Summary     ResultSummary  `json:"summary"`
}

type StudentSummary struct {
	Name        string `json:"name"`
	AdmissionNo string `json:"admission_no"`
	RollNo      string `json:"roll_no"`
	ClassName   string `json:"class_name"`
}

type ExamSummary struct {
	TermName string `json:"term_name"`
}

type SubjectMark struct {
	SubjectName   string  `json:"subject_name"`
	MaxMarks      float64 `json:"max_marks"`
	MarksObtained float64 `json:"marks_obtained"`
	Grade         string  `json:"grade"`
	Remarks       string  `json:"remarks"`
	IsAbsent      bool    `json:"is_absent"`
}

type ResultSummary struct {
	TotalMax      float64 `json:"total_max"`
	TotalObtained float64 `json:"total_obtained"`
	Percentage    float64 `json:"percentage"`
	Grade         string  `json:"grade"`
	Status        string  `json:"status"`
}
