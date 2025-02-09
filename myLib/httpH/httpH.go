package httph

import (
	"database/sql"
	"encoding/json"
	"errors"
	udt "go_final_project/myLib/UDT"
	dbA "go_final_project/myLib/dataBase"
	genDate "go_final_project/myLib/dateGen"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Обработчик правил повторения
func NextDateH(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("Получен запрос")

	// Чтение значений Query параметров
	//
	nowV := r.URL.Query().Get("now")
	dateV := r.URL.Query().Get("date")
	repeatV := r.URL.Query().Get("repeat")

	// Проверка корректности содержимого nowV
	//
	nowD, err := time.Parse("20060102", nowV)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Генерация даты по данным из запроса
	//
	repStr, err := genDate.NextDate(nowD, dateV, repeatV)
	if err != nil {
		http.Error(w, `"Bad request"`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(repStr))

}

// Обработчик добавления задач в БД
func AddTasksH(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Подготовка аварийного сообщения в случае ошибки обработчика
	//
	dataSendError, err := generationArrMsg("Не указан заголовок задачи")
	if err != nil {
		http.Error(w, `"Internal Server Error"`, http.StatusInternalServerError)
		return
	}

	// Чтение тела
	//
	bodyReq, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `"Bad Request"`, http.StatusBadRequest)
		return
	}
	defer func() { _ = r.Body.Close() }()

	var msgRx udt.RxFormat

	err = json.Unmarshal(bodyReq, &msgRx)
	if err != nil {
		w.Write(dataSendError)
		return
	}

	tn := time.Now()

	// Проверка содержимого поля title на отсутствие данных
	//
	if msgRx.Title == "" {
		w.Write(dataSendError)
		return
	}

	var tr time.Time

	// Проверка принятых данных в поле date
	//
	if msgRx.Date == "" {
		msgRx.Date = tn.Format("20060102")
	} else { // Проверка корректности данных в поле date
		tr, err = time.Parse("20060102", msgRx.Date)
		if err != nil {
			w.Write(dataSendError)
			return
		}
	}

	if tr.Before(tn) && tn.Format("20060102") != msgRx.Date {
		// если правило повторения не указано или равно пустой строке
		if msgRx.Repeat == "" {
			msgRx.Date = tn.Format("20060102")
		} else { // при указанном правиле повторения
			msgRx.Date, err = genDate.NextDate(tn, msgRx.Date, msgRx.Repeat)
			if err != nil {
				w.Write(dataSendError)
				return
			}
		}
	}

	// Подключение к БД
	//
	sqlDB, dbCl, err := dbA.CheckCreateDB()
	if err != nil {
		if dbCl != nil {
			dbCl()
		}
		log.Fatal(err)
		os.Exit(1)
	}
	defer func() {
		err = sqlDB.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Добавление записи в БД
	//
	res, err := sqlDB.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", msgRx.Date),
		sql.Named("title", msgRx.Title),
		sql.Named("comment", msgRx.Comment),
		sql.Named("repeat", msgRx.Repeat))

	if err != nil {
		http.Error(w, `"Internal Server Error"`, http.StatusInternalServerError)
		return
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		http.Error(w, `"Internal Server Error"`, http.StatusInternalServerError)
		return
	}

	var sID udt.ReportAddFormat

	sID.ID = strconv.Itoa(int(lastID))

	dataSend, err := json.Marshal(sID)
	if err != nil {
		http.Error(w, `"Internal Server Error"`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(dataSend)
}

// Обработчик чтения задач из БД
func ReadTasksH(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Подготовка аварийного сообщения в случае ошибки обработчика
	//
	dataSendError, err := generationArrMsg("шибка запроса данных")
	if err != nil {
		http.Error(w, `"Internal Server Error"`, http.StatusInternalServerError)
		return
	}

	// Подключение к БД
	//
	sqlDB, dbCl, err := dbA.CheckCreateDB()
	if err != nil {
		if dbCl != nil {
			dbCl()
		}
		log.Fatal(err)
		os.Exit(1)
	}
	defer func() {
		err = sqlDB.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	tn := time.Now()
	tnS := tn.Format("20060102")

	// Формирование запроса
	//
	rows, err := sqlDB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date >= :date ORDER BY date LIMIT 15 ",
		sql.Named("date", tnS))

	if err != nil {
		http.Error(w, `"Internal Server Error"`, http.StatusInternalServerError)
		w.Write(dataSendError)
		os.Exit(1)
	}
	defer func() {
		_ = rows.Close()
	}()

	// Выворд содержмимого запроса
	//
	var t udt.TxFormat = udt.TxFormat{}

	for rows.Next() {

		task := udt.TxFormatEl{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		if err != nil {
			http.Error(w, `"Internal Server Error"`, http.StatusInternalServerError)
			w.Write(dataSendError)
			os.Exit(1)
		}

		t.Tasks = append(t.Tasks, task)
	}

	if t.Tasks == nil {
		t.Tasks = make([]udt.TxFormatEl, 0)
	}

	b, err := json.Marshal(t)
	if err != nil {
		http.Error(w, `"Internal Server Error"`, http.StatusInternalServerError)
		w.Write(dataSendError)
		os.Exit(1)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// Возврат всех параметров задачи по его ID
func GetTasksH(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Подготовка аварийного сообщения если id не найдет
	//
	idNotFound, err := generationArrMsg("Задача не найдена")
	if err != nil {
		http.Error(w, `"Internal Server Error"`, http.StatusInternalServerError)
		return
	}

	// Подготовка аварийного сообщения если id не указан
	//
	idNot, err := generationArrMsg("Не указан идентификатор")
	if err != nil {
		http.Error(w, `"Internal Server Error"`, http.StatusInternalServerError)
		return
	}

	// Подключение к БД
	//
	sqlDB, dbCl, err := dbA.CheckCreateDB()
	if err != nil {
		if dbCl != nil {
			dbCl()
		}
		log.Fatal(err)
		os.Exit(1)
	}
	defer func() {
		err = sqlDB.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Чтение параметра запроса
	//
	id := r.URL.Query().Get("id")

	idN, err := qualId(id)
	if err != nil {
		w.Write(idNot)
		return
	}

	// Формирование запроса
	//
	task := udt.TxFormatEl{}

	row := sqlDB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id ",
		sql.Named("id", idN))

	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		w.Write(idNotFound)
		return
	}

	// Подготовка данных для передачи
	//
	b, err := json.Marshal(task)
	if err != nil {
		http.Error(w, `"Internal Server Error"`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// Сохранение данных задачи по его ID
func SaveTasksH(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Подготовка аварийного сообщения в случае ошибки обработчика
	//
	dataSendError, err := generationArrMsg("Задача не найдена")
	if err != nil {
		http.Error(w, `"Internal Server Error"`, http.StatusInternalServerError)
		return
	}

	// Чтение тела
	//
	bodyReq, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `"Bad Request"`, http.StatusBadRequest)
		return
	}
	defer func() { _ = r.Body.Close() }()

	var msgRx udt.RxFormatFull

	err = json.Unmarshal(bodyReq, &msgRx)
	if err != nil {
		w.Write(dataSendError)
		return
	}

	// Проверка поля id
	//
	id, err := qualId(msgRx.ID)
	if err != nil {
		w.Write(dataSendError)
		return
	}

	// Проверка принятых данных в поле date
	//
	msgRx.Date, err = qualDate(msgRx.Date, msgRx.Repeat)
	if err != nil {
		w.Write(dataSendError)
		return
	}

	// Проверка поля title
	//
	if msgRx.Title == "" {
		w.Write(dataSendError)
		return
	}

	// Подключение к БД
	//
	sqlDB, dbCl, err := dbA.CheckCreateDB()
	if err != nil {
		if dbCl != nil {
			dbCl()
		}
		log.Fatal(err)
		os.Exit(1)
	}
	defer func() {
		err = sqlDB.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Добавление данных задачи в БД
	//
	_, err = sqlDB.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", msgRx.Date),
		sql.Named("title", msgRx.Title),
		sql.Named("comment", msgRx.Comment),
		sql.Named("repeat", msgRx.Repeat),
		sql.Named("id", id))

	if err != nil {
		w.Write(dataSendError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))

}

// Завершение задачи
func DoneTasksH(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Подготовка аварийного сообщения в случае ошибки обработчика
	//
	dataSendError, err := generationArrMsg("Задача не найдена")
	if err != nil {
		http.Error(w, `"Internal Server Error"`, http.StatusInternalServerError)
		return
	}

	// Чтение параметра запроса
	//
	idQ := r.URL.Query().Get("id")

	id, err := qualId(idQ)
	if err != nil {
		w.Write(dataSendError)
		return
	}

	// Подключение к БД
	//
	sqlDB, dbCl, err := dbA.CheckCreateDB()
	if err != nil {
		if dbCl != nil {
			dbCl()
		}
		log.Fatal(err)
		os.Exit(1)
	}
	defer func() {
		err = sqlDB.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Получение данных задачи по его ID
	//
	var msg udt.RxFormatFull

	row := sqlDB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id",
		sql.Named("id", id))

	err = row.Scan(&msg.ID, &msg.Date, &msg.Title, &msg.Comment, &msg.Repeat)

	if err != nil {
		w.Write(dataSendError)
		return
	}

	// Анализ содержимого в поле date после проверки
	// Если поле пустое - происходит удаление задачи по принятому параметру id
	// Если поле не пустое - сохраняются новые данные по принятому параметру id
	//
	if msg.Repeat == "" { // Удаление записи

		_, err = sqlDB.Exec("DELETE FROM scheduler WHERE id = :id",
			sql.Named("id", id))

		if err != nil {
			w.Write(dataSendError)
			return
		}

	} else { // Обновление записи

		// Проверка поля repeat
		//
		tn := time.Now()

		if err != nil {
			w.Write(dataSendError)
			return
		}

		msg.Date, err = genDate.NextDate(tn, msg.Date, msg.Repeat)
		if err != nil {
			w.Write(dataSendError)
			return
		}

		_, err = sqlDB.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
			sql.Named("date", msg.Date),
			sql.Named("title", msg.Title),
			sql.Named("comment", msg.Comment),
			sql.Named("repeat", msg.Repeat),
			sql.Named("id", id))

		if err != nil {
			w.Write(dataSendError)
			return
		}

	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))

}

// Удаление задачи
func DelTasksH(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Подготовка аварийного сообщения в случае ошибки обработчика
	//
	dataSendError, err := generationArrMsg("Задача не найдена")
	if err != nil {
		http.Error(w, `"Internal Server Error"`, http.StatusInternalServerError)
		return
	}

	// Чтение параметра запроса
	//
	idQ := r.URL.Query().Get("id")

	id, err := qualId(idQ)
	if err != nil {
		w.Write(dataSendError)
		return
	}

	// Подключение к БД
	//
	sqlDB, dbCl, err := dbA.CheckCreateDB()
	if err != nil {
		if dbCl != nil {
			dbCl()
		}
		log.Fatal(err)
		os.Exit(1)
	}
	defer func() {
		err = sqlDB.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Удаление задачи из БД
	//
	_, err = sqlDB.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))

	if err != nil {
		w.Write(dataSendError)
		return
	}

	// Успешное завершение
	//
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))

}

// Функция возвращает массив байт и ошибку,
//
// Параметры:
//
// str - строка по данным которой формируется массив байт
func generationArrMsg(str string) (b []byte, err error) {

	errMsg := udt.ErrMsgFormat{
		Error: str,
	}
	b, err = json.Marshal(errMsg)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Функция производит проверку содержимого id. Возвращает признак id формате int и ошибку
//
// Параметры:
//
// id - строка с номером
func qualId(idQ string) (int, error) {

	if idQ == "" {
		return 0, errors.New("Ошибка")
	}

	id, err := strconv.Atoi(idQ)
	if err != nil {
		return 0, errors.New("Ошибка")
	}

	return id, nil
}

// Функция производит проверку содержимого date. Возвращает обработанное date и ошибку
//
// Параметры:
//
// date - дата подлежащая обработке
// repeat - условие повторения
func qualDate(date, repeat string) (string, error) {

	tn := time.Now()

	if date == "" {
		return "", errors.New("Ошибка")
	}

	_, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("Ошибка")
	}

	dateNew, err := genDate.NextDate(tn, date, repeat)
	if err != nil {
		return "", errors.New("Ошибка")

	}

	return dateNew, errors.New("Ошибка")
}
