Oto **20 zadań na Dzień 6: Współbieżność (Concurrency) – Goroutines, Kanały i Selektory**.

Dzisiaj wchodzimy w jeden z najważniejszych i najbardziej charakterystycznych elementów Go. Zamiast jednowątkowej pętli zdarzeń z Node.js (`Event Loop`) i obietnic (`Promises`), poznasz model **CSP (Communicating Sequential Processes)**.

Święta zasada Go brzmi: *"Nie komunikuj się poprzez dzielenie pamięci; dziel pamięć poprzez komunikację"*.

---

## Część 1: Lekkie Wątki (`Goroutines`) i Synchronizacja (Zadania 1–5)

### Zadanie 1: Twoja pierwsza goroutine

Napisz funkcję `sayHello(name string)`, która wypisuje powitanie. Wywołaj ją w `main()` ze słowem kluczowym `go`: `go sayHello("Gopher")`. Zauważ, że program wyłączył się, zanim funkcja cokolwiek wydrukowała. Przeanalizuj, dlaczego tak się stało.

### Zadanie 2: Oczekiwanie na goroutines (`sync.WaitGroup`)

Użyj struktury `sync.WaitGroup`, aby naprawić Zadanie 1:

1. Stwórz `var wg sync.WaitGroup`.
2. Dodaj zadanie do licznika: `wg.Add(1)`.
3. Przekaż wskaźnik do `wg` do goroutine i wywołaj w niej `defer wg.Done()`.
4. Na końcu `main()` wywołaj `wg.Wait()`.

### Zadanie 3: Uruchamianie wielu goroutines w pętli (Pułapka ze zmienną pętli)

Uruchom 5 goroutines w pętli `for i := 0; i < 5; i++`. Przekaż `i` jako argument do funkcji wewnątrz goroutine. Sprawdź, co się stanie, jeśli zamiast przekazać `i` jako parametr, użyjesz go bezpośrednio wewnątrz funkcji anonimowej `go func() { fmt.Println(i) }()`.

### Zadanie 4: Badanie wyścigów pamięci (`Data Race`)

Napisz program, w którym 100 goroutines jednocześnie zwiększa wspólną zmienną `counter++` bez żadnej synchronizacji. Uruchom program w terminalu z flagą detekcji wyścigów:

```bash
go run -race main.go

```

Zaobserwuj raport ze skanera *Data Race*.

### Zadanie 5: Mutex – Ochrona zasobów (`sync.Mutex`)

Napraw wyścig z Zadania 4. Użyj `var mu sync.Mutex` i zabezpiecz operację modyfikacji zmiennej blokadami `mu.Lock()` oraz `mu.Unlock()`. Ponownie uruchom program z flagą `-race` i upewnij się, że wyścig zniknął.

---

## Część 2: Kanały (`Channels`) – Podstawy (Zadania 6–10)

### Zadanie 6: Niebuforowany kanał (Unbuffered Channel)

Stwórz kanał przesyłający napisy: `ch := make(chan string)`. Uruchom goroutine, która wysyła do kanału wiadomość: `ch <- "ping"`. W wątku głównym odebrana treść przypisz do zmiennej: `msg := <-ch` i ją wyświetl.

### Zadanie 7: Blokada przy braku odbiorcy (Deadlock)

Zobacz, co się stanie, gdy spróbujesz wysłać dane do niebuforowanego kanału `ch <- "data"` w tym samym wątku (`main`), bez uruchamiania osobnej goroutine do ich odbioru. Przeanalizuj błąd *fatal error: all goroutines are asleep - deadlock!*.

### Zadanie 8: Buforowany kanał (Buffered Channel)

Stwórz kanał z buforem o rozmiarze 2: `ch := make(chan int, 2)`. Wyślij do niego dwie wartości pod rząd w wątku głównym (`ch <- 1`, `ch <- 2`). Zauważ, że program się nie zablokował. Co się stanie przy próbie wysłania trzeciej wartości?

### Zadanie 9: Zamykanie kanału (`close`) i pętla `range`

Napisz goroutine (producenta), która wysyła w pętli liczby od 1 do 5 do kanału, a po zakończeniu wysyłania **zamyka kanał**: `close(ch)`. W wątku głównym (konsument) użyj pętli `for val := range ch` do odebrania wszystkich danych.

### Zadanie 10: Bezpieczne sprawdzanie czy kanał jest otwarty

Odbierz wartość z zamkniętego kanału używając składni dwuargumentowej: `val, ok := <-ch`. Sprawdź, jaką wartość ma `val` oraz `ok` dla otwartego i zamkniętego kanału.

---

## Część 3: Selektory (`select`) i Zaawansowane Wzorce (Zadania 11–15)

### Zadanie 11: Instrukcja `select`

Stwórz dwa kanały `ch1` i `ch2`. Uruchom dwie goroutines wysyłające dane do tych kanałów z różnymi opóźnieniami (`time.Sleep`). Użyj bloku `select`, aby odebrać wiadomość z tego kanału, który odpowie jako **pierwszy**.

### Zadanie 12: Timeout na operacji za pomocą `select`

Napisz funkcję pobierającą dane z kanału, ale z zabezpieczeniem czasowym. Użyj `select` łącząc Twój kanał z kanałem zwracanym przez `time.After(2 * time.Second)`:

```go
select {
case res := <-dataChan:
    fmt.Println("Received:", res)
case <-time.After(2 * time.Second):
    fmt.Println("Timeout!")
}

```

### Zadanie 13: Kanały jednokierunkowe (Directional Channels)

Dla zwiększenia bezpieczeństwa typów, w funkcjach można ograniczyć prawa do kanału:

* `func produce(ch chan<- int)` – kanał tylko do zapisu (send-only).
* `func consume(ch <-chan int)` – kanał tylko do odczytu (receive-only).
Napisz dwie funkcje z takimi sygnaturami i połącz je kanałem w `main()`.

### Zadanie 14: Anulowanie z użyciem `ctx.Done()` i `select`

Napisz worker pracujący w nieskończonej pętli `for`, który wewnątrz używa `select` do sprawdzania dwóch kanałów: kanału z zadaniami oraz kanału `ctx.Done()` z kontekstu. Po anulowaniu kontekstu worker powinien zakończyć działanie.

### Zadanie 15: Non-blocking Channel Operations (`default` w `select`)

Użyj klauzuli `default` w bloku `select`, aby podjąć próbę odczytu z kanału bez blokowania wątku, jeśli kanał jest akurat pusty.

---

## Część 4: Wzorce Współbieżności w Praktyce (Zadania 16–20)

### Zadanie 16: Worker Pool (Pula Pracowników)

To najpopularniejszy wzorzec produkcyjny w Go!

1. Stwórz kanał `jobs := make(chan int, 100)`.
2. Stwórz kanał `results := make(chan int, 100)`.
3. Uruchom 3 workerów (goroutines), z których każdy w pętli pobiera zadania z `jobs`, wykonuje obliczenie (np. `job * 2`) i wysyła wynik do `results`.
4. Wyślij 10 zadań i zamknij `jobs`. Odbierz 10 wyników.

### Zadanie 17: Fan-Out, Fan-In

* **Fan-Out:** Rozdziel jedno źródło zadań na wiele równolegle pracujących goroutines.
* **Fan-In:** Połącz wyniki z wielu osobnych kanałów w jeden wspólny kanał wyjściowy za pomocą funkcji łączącej (`merge`).

### Zadanie 18: Wyścig zapytania HTTP (First Responder / Hedged Requests)

Napisz funkcję, która wysyła to samo zapytanie (np. pobranie danych z bazy) do 3 różnych serwerów/mirrorów jednocześnie w osobnych goroutines. Użyj buforowanego kanału, aby pobrać **pierwszą odpowiedź**, która nadejdzie, i zignoruj pozostałe.

### Zadanie 19: Ograniczanie przepustowości (Rate Limiting)

Użyj `time.Tick(200 * time.Millisecond)`, aby stworzyć ogranicznik, który pozwala na przetworzenie maksymalnie 5 żądań na sekundę.

### Zadanie 20: Pustokątna struktura w kanałach sygnałowych (`chan struct{}`)

Często kanał służy jedynie do wysłania sygnału ("coś się stało"), a nie przesyłania danych. W Go używa się do tego typu `chan struct{}`, ponieważ `struct{}` zajmuje **0 bajtów** w pamięci! Stwórz kanał `done := make(chan struct{})` i użyj go do powiadomienia `main()` o zakończeniu tła.

---