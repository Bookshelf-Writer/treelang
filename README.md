[![Build](https://github.com/Bookshelf-Writer/TreeOfLanguages/actions/workflows/pylint.yml/badge.svg?branch=main)](https://github.com/Bookshelf-Writer/TreeOfLanguages/actions/workflows/pylint.yml)

![GitHub repo file or directory count](https://img.shields.io/github/directory-file-count/Bookshelf-Writer/TreeOfLanguages?color=orange)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/Bookshelf-Writer/TreeOfLanguages?color=green)
![GitHub repo size](https://img.shields.io/github/repo-size/Bookshelf-Writer/TreeOfLanguages)

## Скрипт для работы с языковыми деревьями.

Особенности:
- Жестко заданная базовая структура, с минимальным количеством обязательных полей
- Поддержка картинок (флаги)
- Привязка к ISO 639-2 (указывается только код)
- Автогенерация недостающих полей
- Строгая обработка структуры (должны быть ТОЛЬКО строки)

Функционал:
- Проверка языковых файлов на корректность
- Перезапись языкового файла, удаляя некорректное, сортируя по ключам и форматируя
- Сравнение файлов между собой
- Применение "изменений" к другим файлам относительно эталона (добавление новых полей, или удаление отсутствующих)
- Генерация файла-карты

---

### Скомпилированы е файлы
- [Windows](/cmd/TreeOfLanguages.windows.exe)
- [Linux](/cmd/TreeOfLanguages.linux.bin)
- [MacOS](/cmd/TreeOfLanguages.macos.bin)

[Папка с исполняемыми файлами](/cmd)

_Если есть какие-то проблемы в работе, пробуйте указывать абсолютные пути_

---

### Внешние библиотеки
- ` pip install argparse `  - Работа с входными параметрами. Аннотации и проверка наличия.
- ` pip install colorama `  - Поддержка цветов для разных платформ, а не только Linux


---

---

### Примеры использования

#### --check _(Проверка файлов)_

- `--file` `-f` - Путь к языковому файлу
- `--dir` `-d`  - Путь к директории в которой лежат языковые файлы. Недопустимо: [ `/` | ` `  ]
- `--safe` `-s` - Безопасный режим. При перезаписи будет создаваться копия старой версии.

##### Проверка конкретного файла.
```bash
python3 script.py --check --file 'somedir/langTreeFile.en.json' 
```

##### Проверка всех файлов в директории.
```bash
python3 script.py --check --dir 'somedir' 
```

---

#### --compiling _(Компиляция и перекомпиляция файлов)_
- `--file` `-f` - Путь к языковому файлу
- `--dir` `-d`  - Путь к директории в которой лежат языковые файлы
- `--tree`      - Режим работы с языковыми деревьями
- `--routs`     - Режим работы с файлом карты
- `--reactTS`  - Режим компиляции файла для интеграции в React-проект с TypeScript

##### Перекомпиляция конкретного файла. В случае нахождения недопустимых полей они будут удалены.
```bash
python3 script.py --compiling --tree --file 'somedir/langTreeFile.en.json' 
```

##### Перекомпиляция всех файлов в директории.
```bash
python3 script.py --compiling --tree --dir 'somedir' 
```

##### Создание файла карты по конкретному файлу.
```bash
python3 script.py --compiling --routs --file 'somedir/langTreeFile.en.json' 
```

##### Создание файла карты используя языковые файлы в директории
```bash
python3 script.py --compiling --routs --dir 'somedir' 
```

---

#### --compare _(Сравнение файлов)_

- `--file` `-f` - Путь к языковому файлу
- `--dir` `-d`  - Путь к директории в которой лежат языковые файлы
- `--pattern`   - Путь к "эталону"

##### Проверка различий между двумя файлами
```bash
python3 script.py --compare --pattern 'somedir/langTreeFile.ru.json' --file 'somedir/langTreeFile.en.json' 
```

##### Сравнение эталонного файла со всеми файлами в директории.
```bash
python3 script.py --compare --pattern 'somedir/langTreeFile.ru.json' --dir 'somedir' 
```

---

#### --clone _(Клонирование изменений)_

- `--file` `-f` - Путь к языковому файлу
- `--dir` `-d`  - Путь к директории в которой лежат языковые файлы
- `--pattern`   - Путь к "эталону"

##### Перенос изменений из эталонного файла в другой, конкретный, файл
```bash
python3 script.py --clone --pattern 'somedir/langTreeFile.ru.json' --file 'somedir/langTreeFile.en.json' 
```

##### Перенос изменений из эталонного файла во все файлы в директории.
```bash
python3 script.py --clone --pattern 'somedir/langTreeFile.ru.json' --dir 'somedir' 
```

---

#### --iso _(Работа со вшитыми языками)_

- `--code`      - Поиск по вшитым языкам. По коду.
- `--name`      - Поиск по вшитым языкам По названию.
- `--codeOnly`  - Вывод при поиске только коды.
- `--fileJSON`  - Генерация сжатого base64 по ссылке на json-файл.

##### Поиск по вшитым языкам по имени. Вернет массив с объектами если есть совпадения.
```bash
python3 script.py --iso --name 'бело' 
```

##### Поиск по вшитым языкам по коду. Вернет массив строк с кодами если есть совпадения.
```bash
python3 script.py --iso --code 'ru' --codeOnly
```

##### Генерация сжатого представления json-файла. Вернет bash64-строку.
```bash
python3 script.py --iso --fileJSON './someFile.json' 
```

#### Использования в реальных задачах

##### Обновление интеграции с React-TS
```bash
python3 script.py --clone --pattern './treeLang/langTreeFile.ru.json' --dir './treeLang' 
python3 script.py --compiling --reactTS --dir './treeLang' 
```
- Все языковые файлы в `dir` приводятся к единой структуре, недостающие поля копируются из `pattern`
- Создание (или обновление) файла `languages.ts` с встраиванием всех языковых файлов в `dir`

##### Обновление интеграции с JS
```bash
python3 script.py --clone --pattern './treeLang/langTreeFile.ru.json' --dir './treeLang' 
python3 script.py --compiling --JS --dir './treeLang' 
```
- Все языковые файлы в `dir` приводятся к единой структуре, недостающие поля копируются из `pattern`
- Создание (или обновление) файла `languages.js` с встраиванием всех языковых файлов в `dir`


---

---

### Структура

Пример файла языкового дерева `example.ru.json`
```json
{
  "info": {
    "iso": "ru",
    "date": "2024-05-18",
    "hash": "d7f8dcd7be4a924c63826de09af356bc",
    "img": "data:image/svg+xml;base64,PD94...."
  },
  "data": {
    "string_value_1": "test text",
    "stringValue2": "new text",
    "groupString1": {
      "a1": "",
      "a2": "",
      "groupString2": {
        "b1": "",
        "b2": ""
      }
    },
    "arrayString1": [
        "c1",
        "c2"
    ],
    "arrayString2": [
        {
            "h1": "",
            "h2": ""
        }
    ]
  }
}
```

- Блоки **info** и **data** обязательные.
- Поле **iso** в блоке **info** является обязательным (код языка).
- Поля **date** и **hash** в блоке **info** автоматически генерируются при перекомпиляции.
- Поле **img** в блоке **info** принимает только строку `data:image` в формате `base64`. Поле может быть пустым или отсутствовать.
- В блоке **data** ключ состоит из латинских букв, цифр и нижнего подчеркивания. Ключ не может начинаться с цифры.
- В блоке **data** допустимы только строки и не пустые массивы\объекты. Массив\объект может быть заполнен любым разрешенным типом.
- При пере компиляции все объекты в **date** сортируются по ключам. Гаданая очередность сохранятся только в массивах.

---

Пример файла карты `languages.json`
```json
{
  "ru": {
    "file": "ru.json",
    "img": "data:image/svg+xml;.....",
    "name": {
      "short": "RU",
      "full": "Russian",
      "native": "Русский"
    }
  },
  "en": {
    "file": "en.json",
    "img": "",
    "name": {
      "short": "EN",
      "full": "English",
      "native": "English"
    }
  }
}
```

- Файл карты полностью генерируемый.
- Первичный ключ является **iso** из файла языкового дерева.
- Поле **img** копируется из языкового дерева. Если такого поля там нет, то оставляет пустую строку
- Поле **name** автоматически генерируется. **short** это первичный ключ в верхнем регистре. Поля **full** и **native** берутся из вшитого в скрипт словаря.
- Порядок блоков соответствует тому, как скрипт прочитал файлы из директории. Очередность можно задать именованием файлов.

---

---

### Mirrors

- https://git.bookshelf-writer.fun/Bookshelf-Writer/TreeOfLanguages
 