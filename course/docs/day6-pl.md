Oto **20 zadań na Dzień 6: Pakiety, Struktura Projektu i Standardowa Biblioteka**.

Dzisiaj zrobisz krok od pisania skryptów w jednym pliku do tworzenia **prawdziwych, ustrukturyzowanych aplikacji w Go**. Opanujesz podział na pakiety, hermetyzację (widoczność pól) oraz kluczowe pakiety biblioteki standardowej (`net/http`, `context`, `slog`, `time`).

---

## Część 1: Pakiety, Widoczność i Struktura Projektu (Zadania 1–5)

### Zadanie 1: Twój pierwszy pod-pakiet

Stwórz folder `config/` a w nim plik `config.go` z deklaracją `package config`. Stwórz w nim strukturę `AppConfig` z polem `Port int`. Zaimportuj ten pakiet w głównym pliku `main.go` i użyj tej struktury.

### Zadanie 2: Publiczne vs Prywatne (Wielkość liter)

W pakiecie `config` stwórz dwie funkcje: `Load()` (publiczna) oraz `parseEnv()` (prywatna). Zobacz, co się stanie, gdy spróbujesz wywołać `config.parseEnv()` w pliku `main.go`.

### Zadanie 3: Enkapsulacja i Gettery/Settery

W nowym pakiecie `user/` zdefiniuj strukturę `User` z prywatnym polem `email string`. Udostępnij publiczne metody `SetEmail(e string) error` (z walidacją `@`) oraz `Email() string` (getter).

> **Wskazówka idiomu Go:** W Go nie używa się przedrostka `Get` dla getterów. Zamiast `GetEmail()` stosuje się po prostu `Email()`.

### Zadanie 4: Aliasy importów i unikanie konfliktów

Wyobraź sobie, że importujesz dwa pakiety o takiej samej nazwie końcowej (np. `math/rand` oraz `crypto/rand`). Użyj aliasu importu w `main.go`, aby używać obu jednocześnie:

```go
import (
    crand "crypto/rand"
    mrand "math/rand"
)

```

### Zadanie 5: Anonimowe importy (Side-effects `_`)

Zobacz, jak działa import tylko dla efektów ubocznych (np. rejestracja sterownika bazy danych): `import _ "[github.com/lib/pq](https://github.com/lib/pq)"`. Dowiedz się, do czego służy specjalna funkcja `init()` w pakietach.

---

## Część 2: Kontekst (`context.Context`) – Kluczowy Koncept w Go (Zadania 6–10)

### Zadanie 6: Tworzenie bazowego kontekstu

W Go `context.Context` służy do przekazywania sygnałów anulowania, deadline'ów oraz metadanych żądania. Stwórz bazowy kontekst za pomocą `ctx := context.Background()`.

### Zadanie 7: Przekazywanie wartości w kontekście (`context.WithValue`)

Stwórz funkcję `processRequest(ctx context.Context)`. Do kontekstu dodaj id żądania: `ctx = context.WithValue(ctx, "request_id", "abc-123")`. Wewnątrz funkcji wyciągnij tę wartość i sprawdź jej typ za pomocą asercji typu.

### Zadanie 8: Anulowanie operacji (`context.WithCancel`)

Stwórz kontekst z funkcją anulującą: `ctx, cancel := context.WithCancel(context.Background())`. Uruchom symulowaną długa operację w pętli `select` nasłuchującej na `<-ctx.Done()`. Wywołaj `cancel()` i zaobserwuj, jak operacja natychmiast przerwie działanie.

### Zadanie 9: Timeouty (`context.WithTimeout`)

Stwórz kontekst, który automatycznie anuluje się po 100 milisekundach: `ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)` (zawsze pamiętaj o `defer cancel()!`). Symuluj operację trwającą 500 ms i obsłuż przerywanie z powodu przekroczenia czasu.

### Zadanie 10: Przekazywanie kontekstu jako PIERWSZY argument

Zgodnie z konwencją Go, jeśli funkcja przyjmuje kontekst, **musi być on pierwszym argumentem**: `func FetchData(ctx context.Context, id string) error`. Przenosząc ten nawyk z TS (gdzie czasy wykonania czy opcje wrzuca się na koniec), dostosuj swoje funkcje do standardu Go.

---

## Część 3: Tworzenie HTTP Serwera (`net/http`) (Zadania 11–15)

### Zadanie 11: Najprostszy serwer HTTP

Stwórz serwer HTTP bez używania zewnętrznych frameworków (jak Express w Node.js). Użyj `http.HandleFunc` z nowym routingiem wprowadzanym od Go 1.22:

```go
http.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("OK"))
})
http.ListenAndServe(":8080", nil)

```

### Zadanie 12: Serwowanie danych JSON

Napisz handler `GET /api/user`, który tworzy strukturę `User`, ustawia nagłówek `w.Header().Set("Content-Type", "application/json")` i serializuje dane bezpośrednio do strumienia odpowiedzi za pomocą `json.NewEncoder(w).Encode(user)`.

### Zadanie 13: Odczytywanie JSON z Body (`POST`)

Napisz handler `POST /api/user`, który dekoduje body żądania do struktury za pomocą `json.NewDecoder(r.Body).Decode(&user)`. Zadbaj o obsługę błędów przy niepoprawnym JSON-ie.

### Zadanie 14: Parametry ze ścieżki (Path Values w Go 1.22+)

Napisz handler `GET /users/{id}`, który pobiera zmienną z adresu URL używając wbudowanej metody `id := r.PathValue("id")`.

### Zadanie 15: Prosty Middleware HTTP

Middleware w Go to funkcja przyjmująca `http.Handler` i zwracająca `http.Handler`. Napisz middleware `LoggingMiddleware`, który mierzy czas wykonania każdego żądania HTTP (użyj `time.Now()` i `time.Since()`) i wypisuje metodę oraz ścieżkę.

---

## Część 4: Nowoczesne Logowanie (`slog`) i Czas (`time`) (Zadania 16–20)

### Zadanie 16: Strukturyzowane logi z pakietem `log/slog`

Od wersji Go 1.21 w bibliotece standardowej znajduje się pakiet `slog`. Zamiast zwykłego `fmt.Println`, użyj `slog.Info("user logged in", "user_id", 42, "role", "admin")`.

### Zadanie 17: Logowanie w formacie JSON

Skonfiguruj `slog`, aby wypluwał logi w formacie JSON (idealne dla środowisk produkcyjnych i Datadog/Grafana):

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
slog.SetDefault(logger)

```

### Zadanie 18: Bezpieczna praca z czasem (`time.Time`)

W JS operacje na datach potrafią być koszmarem. W Go masz potężny pakiet `time`. Stwórz dwie daty, odejmij je od siebie (`diff := t2.Sub(t1)`), uzyskując typ `time.Duration`, i sprawdź ile to sekund lub milisekund.

### Zadanie 19: Formatowanie i parsowanie dat

W Go daty formatuje się za pomocą **konkretnego punktu w czasie**: `Mon Jan 2 15:04:05 MST 2006` (zapamiętaj układ cyfr: 1 2 3 4 5 6 7). Sformatuj aktualny czas do formatu `YYYY-MM-DD` używając patternu `"2006-01-02"`.

### Zadanie 20: Ticker i Timer

Użyj `time.NewTicker(1 * time.Second)`, aby stworzyć pętlę wykonującą akcję co sekundę (np. wypisanie statusu w konsoli). Pamiętaj o zatrzymaniu tickera za pomocą `defer ticker.Stop()`.

---
