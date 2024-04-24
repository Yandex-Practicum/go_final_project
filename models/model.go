package models

type Task struct {
	//gorm.Model
	ID      uint   `json:"id,string" gorm:"unique;primaryKey;autoIncrement"` //	 автоинкрементный идентификатор
	Date    string `json:"date" gorm:"index"`                                // дата задачи, которая будет хранится в формате YYYYMMDD или в Go-представлении 20060102;
	Title   string `json:"title"`                                            // заголовок задачи;
	Comment string `json:"comment"`                                          // комментарий к задаче;
	Repeat  string `json:"repeat" gorm:"size:128"`                           // строковое поле не более 128 символов
}

type Tabler interface {
	TableName() string
}

// TableName overrides the table name used by Task to `scheduler` in gorm automigrate
func (Task) TableName() string {
	return "scheduler"
}

type ResponseError struct {
	Error string `json:"error"`
}

type ResponseTaskId struct {
	Id uint `json:"id"`
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}
