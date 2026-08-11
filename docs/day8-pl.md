Oto **20 zadań na Dzień 8: Podsumowanie Tygodnia 1 – Projekt Praktyczny CLI Task Manager**.

Gratulacje z okazji ukończenia pierwszego tygodnia! Znasz już składnię, wskaźniki, struktury, interfejsy, obsługę błędów, pakiety oraz podstawy współbieżności.

Dzisiaj nie będziemy robić pojedynczych ćwiczeń. Zamiast tego **zbudujesz od zera kompletną, produkcyjną aplikację CLI (Task Manager)**. Podzieliłem ten proces na 20 kroków architektonicznych – każdy krok dodaje jeden konkretny element systemu.

---

```text
taskmanager/
├── go.mod
├── main.go
├── task/
│   ├── task.go        # Struktury i logika biznesowa
│   └── task_test.go   # Testy jednostkowe
└── storage/
    ├── storage.go     # Interfejs dla magazynu danych
    └── json.go        # Implementacja oparta o pliki JSON

```

---

## Faza 1: Podstawowy Model Danych i Repozytorium (Kroki 1–5)

### Krok 1: Inicjalizacja modułu i struktury plików

Stwórz katalog `taskmanager` i zainicjalizuj moduł (`go mod init taskmanager`). Stwórz foldery `task/` oraz `storage/`.

### Krok 2: Model zadania (`task/task.go`)

W pakiecie `task` zdefiniuj strukturę `Task` z polami:

* `ID int`
* `Title string`
* `Done bool`
* `CreatedAt time.Time`
Dodaj tagi JSON do wszystkich pól (np. `json:"id"`).

### Krok 3: Walidacja i Konstruktor

Napisz funkcję `NewTask(id int, title string) (*Task, error)`. Jeśli `title` jest puste, zwróć błąd `errors.New("title cannot be empty")`. W przeciwnym razie ustaw `CreatedAt` na `time.Now()`.

### Krok 4: Interfejs Magazynu Danych (`storage/storage.go`)

W pakiecie `storage` zdefiniuj interfejs `Storage`:

```go
type Storage interface {
    Save(tasks []task.Task) error
    Load() ([]task.Task, error)
}

```

### Krok 5: Implementacja JSON (`storage/json.go`)

Stwórz strukturę `JSONStorage` posiadającą pole `filename string`. Zaimplementuj metody `Save` oraz `Load` przy użyciu standardowych pakietów `os` oraz `encoding/json` (`json.MarshalIndent` / `json.Unmarshal`).

---

## Faza 2: Logika Menedżera Zadań (Kroki 6–10)

### Krok 6: Struktura `TaskManager`

W pakiecie `task` zdefiniuj strukturę `TaskManager`, która posiada:

* `tasks []Task`
* `storage storage.Storage` (użyj interfejsu zamiast konkretnej struktury!)

### Krok 7: Konstruktor `NewManager`

Napisz funkcję `NewManager(s storage.Storage) (*TaskManager, error)`, która przy starcie wywołuje `s.Load()` i zasila wewnętrzną listę `tasks`.

### Krok 8: Metoda `Add`

Dodaj do `TaskManager` metodę `Add(title string) error`. Wygeneruj unikalne ID (np. `len(tasks) + 1`), stwórz nowe zadanie, dodaj je do wycinka i wywołaj `t.storage.Save(t.tasks)`.

### Krok 9: Metoda `MarkDone`

Dodaj metodę `MarkDone(id int) error`. Przeiteruj po zadaniach, znajdź zadanie o podanym ID i zmień `Done` na `true`. Jeśli zadanie nie istnieje, zwróć własny błąd (sentinel error) `ErrTaskNotFound`. Zapisz stan do pliku.

### Krok 10: Metoda `List`

Dodaj metodę `List(showAll bool) []Task`. Jeśli `showAll` wynosi `false`, zwróć tylko zadania niewykonane (`Done == false`).

---

## Faza 3: Parser CLI i Interfejs Użytkownika (Kroki 11–15)

### Krok 11: Zastosowanie `flag` (Standardowa biblioteka)

W pliku `main.go` użyj standardowego pakietu `flag` do parsowania argumentów wiersza poleceń:

* `-add <title>` (dodawanie zadania)
* `-done <id>` (oznaczanie jako wykonane)
* `-list` (wyświetlanie zadań)
* `-all` (flaga modyfikująca `-list`)

### Krok 12: Inicjalizacja zależności w `main.go`

Połącz wszystkie elementy w `main()`:

1. Utwórz instancję `JSONStorage` na pliku `tasks.json`.
2. Utwórz `TaskManager` przekazując do niego storage.
3. Obsłuż ewentualne błędy startowe za pomocą `log.Fatalf`.

### Krok 13: Obsługa komend w `main.go`

Napisz logikę sterującą na podstawie przekazanych flag (użyj `switch` lub `if`). Przykład:

```bash
go run main.go -add "Kupić mleko"
go run main.go -list
go run main.go -done 1

```

### Krok 14: Ładne formatowanie w konsoli (`tabwriter`)

Do wyświetlania listy zadań użyj standardowego pakietu `text/tabwriter`, aby kolumny (ID, Status, Tytuł, Data) wyrównały się automatycznie w konsoli.

### Krok 15: Graceful Shutdown i zapis przy wyjściu

Użyj `os/signal` oraz `syscall.SIGINT` / `syscall.SIGTERM`, aby w przypadku przerwania programu skrótem `Ctrl+C` bezpiecznie dokończyć ewentualny zapis pliku.

---

## Faza 4: Testy Jednostkowe i Współbieżność (Kroki 16–20)

### Krok 16: Twój pierwszy test w Go (`task/task_test.go`)

Stwórz plik `task/task_test.go`. Napisz test `TestNewTask(t *testing.T)`, który sprawdza czy puste pole `title` zwraca błąd, a poprawne pole tworzy strukturę z niezerową datą. Uruchom test za pomocą:

```bash
go test ./...

```

### Krok 17: Table-Driven Tests (Idiomatyczny wzorzec testów w Go)

Przepisz `TestNewTask` na wzorzec **Table-Driven Tests**:

```go
tests := []struct {
    name    string
    title   string
    wantErr bool
}{
    {"valid title", "Kup mleko", false},
    {"empty title", "", true},
}

```

### Krok 18: Mockowanie Storage w testach

Napisz w pliku testowym prostą strukturę `MockStorage` implementującą interfejs `storage.Storage` w pamięci (bez zapisu na dysku). Użyj jej do przetestowania metody `TaskManager.Add()`.

### Krok 19: Asynchroniczne powiadomienie (Goroutine + Channel)

Dodaj do `TaskManager` funkcję, która przy oznaczeniu zadania jako wykonane wysyła asynchronicznie (w goroutine) informację na kanał `chan string` (np. `"Task #1 completed!"`). Odebraną wiadomość wypisz w konsoli.

### Krok 20: Budowanie binarce (`go build`) i testy pokrycia

1. Uruchom testy z flagą pokrycia kodu: `go test -cover ./...`.
2. Zbuduj produkcyjny plik wykonywalny:
```bash
go build -o taskmgr main.go

```


3. Przetestuj gotowy plik wykonywalny w terminalu: `./taskmgr -list`.

---