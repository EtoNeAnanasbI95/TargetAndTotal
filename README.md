# TargetAndTotal CLI

* Это утилита написанная на Go с использованием llm — ollama, для студентов обречённых писать тонны отчётов к различного рода практическим работам

![Пример](./preview/preview.gif)

## Зависимости
* Go
* Ollama llm
## Установка
> [!WARNING]
> Это важное предупреждение, которое нужно учитывать.
> 
> Утилита имеет абсолютную зависимость от запущенной ollama, если у вас её нет, то вам необходимо поставить [её](https://hub.docker.com/r/ollama/ollama) и запустить соответствующим образом

* Выполнить установку можно с помощью утилиты go install
```
go install github.com/EtoNeAnanasbI95/TargetAndTotal/cmd/TargetAndTotal
```