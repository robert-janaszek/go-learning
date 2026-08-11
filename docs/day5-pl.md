Oto **20 zadań na Dzień 5: Obsługa Błędów (`if err != nil`), Pakiet `errors` i Brak Wyjątków**.

Twój cel na dziś: oduczyć się myślenia w kategoriach `try/catch/throw`, opanować idiomy obsługi błędów jako wartości (*errors are values*) oraz poznać nowoczesny wzorzec zawijania błędów (error wrapping) wprowadzony w nowszych wersjach Go.

---

## Część 1: Podstawy i Zwracanie Błędów (Zadania 1–5)

### Zadanie 1: Twój pierwszy błąd (`errors.New`)

Napisz funkcję `ValidateAge(age int) error`. Jeśli wiek jest mniejszy niż 0 lub większy niż 120, zwróć błąd utworzony za pomocą `errors.New("invalid age")`. W przeciwnym razie zwróć `nil`. Obsłuż wynik w `main()`.

### Zadanie 2: Formatowanie błędu z kontekstem (`fmt.Errorf`)

Zmodyfikuj funkcję z Zadania 1 tak, aby w treści błędu znajdowała się przekazana wartość, np. `fmt.Errorf("age %d is out of range [0-120]", age)`.

### Zadanie 3: Happy Path on the Left (Czysta struktura kodu)

Napisz funkcję procesującą zamówienie `ProcessOrder(id int, amount float64) error`. Wykonaj 3 walidacje: ID > 0, amount > 0, amount < 10000. Napisz kod tak, aby w przypadku błędu od razu go zwracać (`if err != nil { return err }`), unikając zagnieżdżonych bloków `else`.

### Zadanie 4: Sentinel Errors (Stałe/zmienne błędów)

Zadeklaruj na poziomie pakietu zmienne błędów (konwencja nazewnicza: `Err...`):

* `var ErrNotFound = errors.New("item not found")`
* `var ErrPermissionDenied = errors.New("permission denied")`
Napisz funkcję `FindUser(id int) (*User, error)`, która zwraca `ErrNotFound`, gdy ID wynosi 0. Sprawdź ten błąd w `main()` zwykłym porównaniem `if err == ErrNotFound`.

### Zadanie 5: Ignorowanie błędu (Co robi `_`)

Wywołaj funkcję ze standardowej biblioteki (np. `strconv.Atoi("123")`), która zwraca `(int, error)`. Zignoruj błąd operatorem `_`. Przeanalizuj, dlaczego lintery (np. `golangci-lint`) uważają to za antywzorzec w kodzie produkcyjnym.

---

## Część 2: Własne Typy Błędów i Custom Structs (Zadania 6–10)

### Zadanie 6: Struktura jako błąd

Stwórz własną strukturę błędu `ValidationError`:

```go
type ValidationError struct {
    Field   string
    Message string
}

```

Zaimplementuj dla niej interfejs `error` (metoda `Error() string`, która ładnie formatuje pole i treść).

### Zadanie 7: Zwracanie własnej struktury jako `error`

Napisz funkcję `RegisterUser(email, password string) error`. Jeśli email nie zawiera `@`, zwróć `&ValidationError{Field: "email", Message: "missing @"}`. Zauważ, że funkcja zwraca typ interfejsowy `error`!

### Zadanie 8: Wyciąganie własnego błędu (Asercja Typu)

W `main()` wywołaj `RegisterUser` z błędnym mailem. Odbierz błąd typu `error` i użyj asercji typu (`ve, ok := err.(*ValidationError)`), aby uzyskać dostęp do konkretnych pól `Field` i `Message`.

### Zadanie 9: Pułapka `nil` z typami strukturalnymi (Kluczowy niuans!)

Napisz funkcję `badValidate() error`, w której zadeklarujesz wskaźnik do swojej struktury `var customErr *ValidationError = nil`, a na końcu zwrócisz `customErr`. Sprawdź w `main()`, dlaczego warunek `if err != nil` wyewoluuje do `true`! *(Podpowiedź: interfejs z typem `*ValidationError` i wartością `nil` sam w sobie NIE JEST równe `nil`)*.

### Zadanie 10: Poprawne zwracanie `nil` z własnych struktur

Napraw błąd z Zadania 9. Przeanalizuj, dlaczego zawsze należy zwracać jawnie literale `nil` (`return nil`), zamiast zmiennych wskaźnikowych reprezentujących błąd.

---

## Część 3: Error Wrapping i Nowoczesne Pakiet `errors` (Zadania 11–15)

### Zadanie 11: Owijanie błędów konfiksorem `%w`

Napisz funkcję `ReadConfig(path string) error`. Wywołaj w niej `os.ReadFile(path)`. Jeśli dostaniesz błąd, zwróć go owiniętego o własny kontekst:
`return fmt.Errorf("failed to read config file %s: %w", path, err)`.

### Zadanie 12: Sprawdzanie błędów za pomocą `errors.Is`

Napisz funkcję, która wywołuje `ReadConfig` z nieistniejącą ścieżką. W `main()` użyj nowej metody `errors.Is(err, os.ErrNotExist)` zamiast zwykłego `==`, aby sprawdzić, czy pod spodem (wewnątrz łańcucha owinięć) kryje się błąd braku pliku.

### Zadanie 13: Rozwijanie błędów za pomocą `errors.As`

Użyj funkcji `errors.As` do bezpiecznego wyciągnięcia struktury `*os.PathError` z owiniętego błędu zwróconego przez operacje na plikach. Wyświetl pole `Path` tego błędu.

### Zadanie 14: Łączenie wielu błędów (`errors.Join` - Go 1.20+)

Napisz funkcję walidującą formularz, która zbiera wszystkie błędy walidacji do wycinka `[]error` i na końcu zwraca połączony błąd za pomocą `errors.Join(err1, err2, err3)`. Sprawdź, jak wygląda wydrukowany błąd.

### Zadanie 15: Przeplatane owijanie błędów w architekturze warstwowej

Stwórz 3-warstwowy łańcuch wywołań: `Repository` $\rightarrow$ `Service` $\rightarrow$ `Handler`.

1. `Repository` zwraca `ErrNotFound`.
2. `Service` owija go: `fmt.Errorf("user service: %w", err)`.
3. `Handler` sprawdza `errors.Is(err, ErrNotFound)` i zwraca kod HTTP 404.

---

## Część 4: Panic, Recover i Sytuacje Nadzwyczajne (Zadania 16–20)

### Zadanie 16: Kiedy używać `panic`?

W Go `panic` używa się niezmiernie rzadko (odpowiednik krytycznego wypróżnienia aplikacji, np. brak pliku konfiguracyjnego przy starcie). Napisz funkcję `MustParseURL(rawURL string)`, która robi `panic("invalid URL")`, jeśli URL jest pusty.

### Zadanie 17: Przechwytywanie paniki (`defer` + `recover`)

Napisz funkcję `SafeExecute(fn func())`, która wywołuje przekazaną funkcję `fn`. Użyj w niej bloku `defer` oraz `recover()`, aby przechwycić ewentualną panikę i zapobiec crashowi całej aplikacji.

### Zadanie 18: Konwersja `panic` na `error`

Zmodyfikuj `SafeExecute` tak, aby zwracała `error`. Jeśli wewnątrz nastąpiła panika, przekształć odzyskaną wartość (`recover()`) w zwykły `error` i go zwróć.

### Zadanie 19: Porządkowanie zasobów z `defer` przy błędach

Otwórz plik za pomocą `file, err := os.Open(...)`. Od razu po sprawdzeniu błędu dopisz `defer file.Close()`. Wywołaj pomocniczą funkcję, która zwraca błąd i upewnij się, że plik został zamknięty przed wyjściem z funkcji głównej.

### Zadanie 20: Tworzenie middleware logującego błędy

Napisz prostą funkcję-opakowanie (middleware): `ExecuteWithLogging(fn func() error)`. Funkcja ta wywołuje `fn()`, sprawdza czy zwrócono błąd, i jeśli tak – loguje go w konsoli ze stemplem czasowym z użyciem standardowego pakietu `log`.

---