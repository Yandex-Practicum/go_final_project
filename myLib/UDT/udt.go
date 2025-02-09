package udt

// Запрос
type RxFormat struct {
	Date    string `json:"date"`             // дата задачи в формате 20060102;
	Title   string `json:"title"`            // заголовок задачи. Обязательное поле;
	Comment string `json:"comment"`          // комментарий к задаче;
	Repeat  string `json:"repeat,omitempty"` // правило повторения.
}

type RxFormatFull struct {
	ID      string `json:"id"`               // id задачи
	Date    string `json:"date"`             // дата задачи в формате 20060102;
	Title   string `json:"title"`            // заголовок задачи. Обязательное поле;
	Comment string `json:"comment"`          // комментарий к задаче;
	Repeat  string `json:"repeat,omitempty"` // правило повторения.
}

// Для возврата клиенту id добавленной в БД задачи
type ReportAddFormat struct {
	ID string `json:"id"`
}

// Для возврата клиенту сообщения об ошибке
type ErrMsgFormat struct {
	Error string `json:"error"`
}

// Ответ клиенту
type TxFormatEl struct {
	ID      string `json:"id"`      // дата задачи в формате 20060102;
	Date    string `json:"date"`    // заголовок задачи. Обязательное поле;
	Title   string `json:"title"`   // заголовок задачи;
	Comment string `json:"comment"` // комментарий к задаче;
	Repeat  string `json:"repeat"`  // правило повторения.
}

type TxFormat struct {
	Tasks []TxFormatEl `json:"tasks"` // слайс выборки данных БД
}
