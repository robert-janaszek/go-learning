Oto **20 zadań na Dzień 3: Struktury (`Structs`), Metody i Brak Klas**.

Twój cel na dziś: opanować kompozycję zamiast dziedziczenia oraz nauczyć się pisać metody przyczepiane do struktur z odpowiednim odbiorcą (*value* vs *pointer receiver*).

---

## Część 1: Definiowanie i Inicjalizacja Struktur (Zadania 1–5)

### Zadanie 1: Podstawowa struktura

Zdefiniuj strukturę `Book` z polami: `Title` (string), `Author` (string), `Pages` (int), `IsRead` (bool). Utwórz instancję tej struktury w `main()` podając nazwy pól (*struct literal*) i wyświetl ją.

### Zadanie 2: Różne sposoby inicjalizacji

Utwórz 3 instancje struktury `Book`:

1. Używając nazw pól (`Book{Title: "...", ...}`).
2. Bez nazw pól (uważaj na kolejność!).
3. Pustą instancję (`b := Book{}`) i uzupełnij pola w osobnych liniach (`b.Title = "..."`).

### Zadanie 3: Wskaźnik do struktury i funkcja fabrykująca (Constructor Pattern)

W Go nie ma słowa `new` w kontekście klas. Tworzy się funkcje typu `NewBook`. Napisz funkcję `NewBook(title, author string, pages int) *Book`, która zwraca **wskaźnik** do nowo utworzonej struktury.

### Zadanie 4: Anonimowa struktura (Ad-hoc Struct)

Często w Go (np. w testach lub przy wyciąganiu JSON-a) tworzy się struktury jednorazowe. Stwórz anonimową strukturę zawierającą pola `ConfigName` oraz `Port`, zainicjalizuj ją od razu wartościami i wyświetl.

### Zadanie 5: Porównywanie struktur

Stwórz dwie instancje `Book` o identycznych wartościach. Sprawdź warunkiem `if b1 == b2`, czy Go potrafi je porównać. Następnie dodaj do struktury pole typu slice `Tags []string` i sprawdź, dlaczego kod przestał się kompilować (struktury z typami referencyjnymi nie są porównywalne operatorem `==`).

---

## Część 2: Metody – Pointer vs Value Receiver (Zadania 6–10)

### Zadanie 6: Pierwsza metoda z Value Receiverem

Dodaj do struktury `Book` metodę `Summary() string`. Metoda ma zwracać napis w formacie `"Title" by Author (X pages)`. Zastosuj *value receiver* `(b Book)`.

### Zadanie 7: Metoda z Pointer Receiverem (Modyfikacja stanu)

Dodaj do `Book` metodę `MarkAsRead()`. Przemyśl: czy odbierak metody powinien być wskaźnikiem `(b *Book)` czy wartością `(b Book)`? Przetestuj w `main()`, wywołując tę metodę na książce, która miała `IsRead = false`.

### Zadanie 8: Wywołanie metody z pointer receiverem na wartości

Stwórz zmienną `b := Book{Title: "Go in Action"}` (zwykła wartość, nie wskaźnik). Wywołaj na niej metodę `MarkAsRead()` z Zadania 7. Zauważ, że Go **automatycznie pobiera adres** (`(&b).MarkAsRead()`) – nie musisz zamieniać zmiennej na wskaźnik.

### Zadanie 9: Metody na typach własnych (podstawowych)

W Go metody można przyczepiać nie tylko do struktur! Zdefiniuj własny typ: `type Celsius float64`. Dodaj do niego metodę `ToFahrenheit() float64`. Przetestuj w `main()`.

### Zadanie 10: Metoda zmieniająca typ własny

Do typu `Celsius` dodaj metodę `Add(degrees float64)`. Wybierz odpowiedni receiver, aby metoda faktycznie modyfikowała wartość temperatury, na której została wywołana.

---

## Część 3: Kompozycja i Osadzanie (Embedding) Zamiast Dziedziczenia (Zadania 11–15)

### Zadanie 11: Zwykłe zagnieżdżenie struktur

Stwórz strukturę `Address` (`City string`, `ZipCode string`). Stwórz strukturę `User` z polami `Name string` oraz `Addr Address`. Zainicjalizuj `User` i wyświetl miasto użytkownika (`u.Addr.City`).

### Zadanie 12: Anonymous Struct Embedding (Promowane pola)

Zmodyfikuj `User` tak, aby pole `Address` było **anonimowe** (tzw. osadzone/embedded struct):

```go
type User struct {
    Name string
    Address // Brak nazwy pola, tylko typ!
}

```

Sprawdź w `main()`, jak działa tzw. *field promotion* – uzyskaj dostęp do miasta pisząc po prostu `u.City`.

### Zadanie 13: Nadpisywanie pól i metod (Shadowing w kompozycji)

Dodaj do `Address` metodę `FullAddress() string`. Następnie dodaj do `User` własną metodę `FullAddress() string`, która zwraca imię i adres. Wywołaj obie w `main()` i sprawdź, jak Go rozwiązuje konflikty nazw.

### Zadanie 14: Osadzanie wskaźnika do struktury

Stwórz strukturę `Engine` (`HorsePower int`). Stwórz strukturę `Car` z osadzonym wskaźnikiem `*Engine`. Sprawdź, co się stanie, gdy wywołasz metodę na `Car`, jeśli `Engine` jest równy `nil`.

### Zadanie 15: Kompozycja z wielu struktur

Stwórz dwie małe struktury: `Logger` (metoda `Log(msg string)`) oraz `Database` (metoda `Connect()`). Stwórz strukturę `Server`, która osadza **obie** te struktury. Wywołaj `server.Log()` oraz `server.Connect()`.

---

## Część 4: Praktyczne wzorce, JSON i Tagi (Zadania 16–20)

### Zadanie 16: Tagi struktur (`Struct Tags`)

Zdefiniuj strukturę `Product` z polami `ID int`, `Name string`, `Price float64`. Dodaj tagi JSON, np. `json:"product_id"`. Użyj `json.Marshal(p)` z pakietu `encoding/json`, aby zmienić strukturę na bajty JSON i wyświetl wynik w konsoli (`string(bytes)`).

### Zadanie 17: Ukrywanie pól w JSON (`json:"-"` i `omitempty`)

Dodaj do struktury `Product` pola:

* `InternalCode string` – ma być ignorowane przez JSON (`json:"-"`).
* `Discount float64` – ma być pomijane w JSON-ie, jeśli jest równe 0 (`json:"discount,omitempty"`).
Sprawdź działanie `json.Marshal`.

### Zadanie 18: Unmarshaling (JSON -> Struct)

Stwórz zmienną ze stringiem reprezentującym JSON: `jsonData := []byte('{"name":"Laptop", "price": 2500}')`. Użyj `json.Unmarshal(jsonData, &p)` do wczytania danych do struktury `Product`. Uważaj: dlaczego musisz przekazać `&p` (wskaźnik), a nie samo `p`?

### Zadanie 19: Koszyk sklepowy z metodą wyliczającą sumę

Zrób rozgrzewkę przed projektem: Stwórz `CartItem` (`Product Product`, `Quantity int`). Stwórz `Cart` z polem `Items []CartItem`. Dodaj metody: `AddItem(p Product, qty int)` oraz `Total() float64`.

### Zadanie 20: Enkapsulacja i prywatne pola

Stwórz pod-pakiet w folderze `bank/`. Zdefiniuj w nim strukturę `Account` z **prywatnym** polem `balance float64`. Udostępnij publiczne metody `Deposit(amount float64)`, `Withdraw(amount float64) error` oraz `Balance() float64`. Sprawdź w `main.go`, że nie możesz zmodyfikować pola `balance` bezpośrednio.

---