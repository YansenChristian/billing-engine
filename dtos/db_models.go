package dtos

type LoanModel struct {
}

func (r *LoanModel) GetAll() []interface{} {
	return []interface{}{}
}

func (r *LoanModel) GetTableName() string {
	return "loans_tab"
}
