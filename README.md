## Инструкция

---
Клонируй проект:
```
git clone https://github.com/adevvvv/note_app
```
Открой проект в терминале и запусти команду:
```
docker-compose up --build
```
Открой ссылку в браузере для просмотра возможностей и тестирования проекта:
```
http://localhost:8000/swagger/index.html#/
```
---
## Реализовано

---
### Авторизация пользователя:
- [x]  Пользователь должен авторизоваться, отправив логин и пароль в приложение.
- [x]  Приложение проверяет корректность полученных данных и возвращает авторизационный токен в случае успеха.
---
### Регистрация пользователей:
- [x]  Регистрация осуществляется посредством отправки в приложение логина и пароля.
- [x]  Реализованы разумные ограничения на формат логина и пароля.
1. **Логин:**
    - Длина логина должна быть от 4 до 20 символов.
    - Логин может содержать только буквы (латинские), цифры и символ подчеркивания.
2. **Пароль:**
    - Длина пароля должна быть от 6 до 20 символов.
    - Пароль может содержать буквы (латинские), цифры и следующие специальные символы: **`!@#$%^&*()-+=`**
---
### Размещение заметки:
- [x]  Размещение заметки происходит посредством отправки данных в формате JSON: заголовок, текст.
- [x]  Заметки могут размещать только авторизованные пользователи.
- [x]  Реализованы разумные ограничения на длину заголовка и текста.
- [x]  В успешном ответе возвращаются данные добавленной заметки
---
### Редактирование заметки:
- [x]  Происходит как создание заметки.
- [x]  Пользователи могут редактировать только свои заметки если срок размещения не больше 1 дня
---
### Отображение списка заметок:
- [x]  Лента представляет собой список заметок, отсортированных по дате добавления
- [x]  Реализована постраничная навигация и возможность фильтрации по определенным датам или диапазонам добавления,
  пользователю.
- [x]  Для каждой заметки возвращается: заголовок, текст, логин автора.
- [x]  Для авторизованных пользователей возвращается признак принадлежности заметки текущему пользователю.
---
### Дополнительно:
- [x]  Добавлена возможность удаления заметок.
- [x]  Обработка и хранение паролей с использованием хэширования.
- [x]  Документация к АРI c помощью Swagger и комментарии к коду.
- [x]  Реализована обработка ошибок.
- [x]  Упаковка приложения и БД (Postgres) в Docker с инструкцией развертывания.
