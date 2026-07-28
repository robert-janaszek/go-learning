Oto **20 zadań na Dzień 4: Interfejsy (`Interfaces`), Duck Typing i Polimorfizm**.

Twój cel na dziś: zrozumieć, dlaczego interfejsy w Go definiuje się po stronie *konsumenta* danych, jak działają interfejsy niejawne (implicit) oraz jak pisać czysty, luźno powiązany kod bez klas i dziedziczenia.

---

## Część 1: Podstawy Interfejsów i Niejawna Implementacja (Zadania 1–5)

### Zadanie 1: Twój pierwszy interfejs

Zdefiniuj interfejs `Stringer` posiadający jedną metodę: `String() string`. Stwórz strukturę `User` (`Name string`, `Age int`) i zaimplementuj dla niej metodę `String() string`. Zauważ, że **nie używasz żądnego słowa kluczowego typu `implements**`.

### Zadanie 2: Polimorfizm w funkcji

Napisz funkcję `PrintInfo(s Stringer)`, która przyjmuje cokolwiek, co spełnia interfejs `Stringer`, i wywołuje na tym metodę `String()`. Wywołaj ją w `main()` z instancją `User`.

### Zadanie 3: Wiele struktur spełniających ten sam interfejs

Stwórz drugą strukturę `Book` (`Title string`, `Author string`) i również zaimplementuj dla niej metodę `String() string`. Przekaż instancję `Book` do tej samej funkcji `PrintInfo`.

### Zadanie 4: Wartość domyślna interfejsu (`nil interface`)

Zadeklaruj zmienną typu interfejsowego: `var s Stringer` bez przypisywania do niej żadnej struktury. Sprawdź warunkiem `if s == nil`, czy jest pusta. Zobacz, co się stanie, jeśli spróbujesz wywołać `s.String()` na pustym interfejsie.

### Zadanie 5: Wytwórnia/Lista interfejsów

Stwórz wycinek (slice) elementów typu interfejsowego: `items := []Stringer{user, book}`. Iteruj po nim pętlą `for _, item := range items` i wywołaj `item.String()` dla każdego z nich.

---

## Część 2: Pointer Receiver vs Value Receiver w Interfejsach (Zadania 6–10)

### Zadanie 6: Interfejs z Pointer Receiverem – Pułapka rzutowania

Stwórz interfejs `Saver` z metodą `Save() error`. Stwórz strukturę `Document` i zaimplementuj metodę `Save()` używając **pointer receivera** `(d *Document)`.

### Zadanie 7: Próba przypisania wartości (Sprawdź błąd kompilatora!)

Spróbuj przypisać do interfejsu zwykłą wartość: `var s Saver = Document{}`. Zobacz błąd kompilatora! Dlaczego musisz przekazać wskaźnik `var s Saver = &Document{}`? *(To kluczowy niuans w Go!)*.

### Zadanie 8: Interfejs z Value Receiverem

Stwórz drugą strukturę `Note` i zaimplementuj `Save()` z **value receiverem** `(n Note)`. Sprawdź, czy do `Saver` możesz przypisać zarówno `Note{}` (wartość), jak i `&Note{}` (wskaźnik).

### Zadanie 9: Standardowy interfejs `io.Reader` i `io.Writer`

Go ma genialne wbudowane interfejsy. Zapoznaj się z `io.Writer` (`Write(p []byte) (n int, err error)`). Napisz funkcję `WriteHello(w io.Writer)`, która zapisuje bajty `"Hello Go"` do cokolwiek, co przyjmuje ten interfejs (np. `os.Stdout` lub `bytes.Buffer`).

### Zadanie 10: Zapis do pliku i do konsoli tym samym kodem

Użyj funkcji `WriteHello` z Zadania 9 dwa razy: raz przekazując `os.Stdout` (wypisanie w konsoli), a drugi raz przekazując plik utworzony przez `os.Create("test.txt")`. Zauważ siłę polimorfizmu bez tworzenia skomplikowanych hierarchii klas!

---

## Część 3: Type Assertion, Type Switch i Empty Interface (Zadania 11–15)

### Zadanie 11: Pusty interfejs (`any` / `interface{}`)

W Go 1.18 wprowadzono alias `any` (odpowiednik `interface{}`). Odpowiada on typowi `unknown` w TypeScript. Stwórz funkcję `Describe(i any)`, która przyjmuje dowolny typ i go drukuje.

### Zadanie 12: Type Assertion (Asercja Typu)

Mając zmienną `var val any = "hello Go"`, wyciągnij z niej oryginalny typ `string` za pomocą rzutowania: `s := val.(string)`. Wydrukuj długość tego stringa (`len(s)`).

### Zadanie 13: Bezpieczna Asercja Typu (Wzorzec `val, ok`)

Co się stanie, jeśli zrobisz `num := val.(int)` na zmiennej trzymającej stringa? (Aplikacja się wyłoży!). Napisz bezpieczną wersję używając składni `n, ok := val.(int)` i obsłuż przypadek, gdy `ok == false`.

### Zadanie 14: Type Switch (Odpowiednik `match` / `switch typeof`)

Napisz funkcję `ProcessInput(v any)`, która używa składni `switch v.(type)` do sprawdzenia typu przekazanego parametru (`int`, `string`, `bool`, `Player`). Dla każdego typu wypisz odpowiedni komunikat.

### Zadanie 15: Wyciąganie wartości ze struktury ukrytej za interfejsem

Stwórz interfejs `Payer`. Stwórz strukturę `CreditCard` posiadającą unikalne pole `CardNumber string`. Przypisz `CreditCard` do zmiennej typu `Payer`. Użyj asercji typu, aby wydobyć `CardNumber`.

---

## Część 4: Kompozycja Interfejsów i Dobre Praktyki (Zadania 16–20)

### Zadanie 16: Łączenie interfejsów (Interface Embedding)

Zdefiniuj dwa małe interfejsy:

1. `Reader` z metodą `Read() string`
2. `Writer` z metodą `Write(data string)`
Stwórz trzeci interfejs `ReadWriter`, który **osadza w sobie oba te interfejsy**.

### Zadanie 17: Złota zasada Go: "Kompaktowe interfejsy"

W TypeScript często tworzy się gigantyczne interfejsy z kilkunastoma metodami. W Go idealny interfejs ma **1 lub 2 metody**. Stwórz strukturę `FileHandler`, która spełnia interfejs `ReadWriter` z Zadania 16.

### Zadanie 18: Interfejs definiowany po stronie odbiorcy (Consumer-defined Interface)

To najważniejszy koncept w Go. Stwórz pod-pakiet `store/` ze strukturą `PostgresStore` posiadającą metodę `GetUsers() []string`. W pakiecie `main` zdefiniuj interfejs `UserGetter` i użyj go w usłudze `UserService`. *(Zauważ: `PostgresStore` w ogóle nie wie o istnieniu interfejsu w `main`!)*.

### Zadanie 19: Łatwe mockowanie w testach (preview)

Dzięki podejściu z Zadania 18, stwórz w `main` strukturę `MockStore` implementującą `GetUsers() []string` ze sztucznymi danymi. Podmień `PostgresStore` na `MockStore` w `UserService`.

### Zadanie 20: Domyślne interfejsy ze standardowej biblioteki (`error`)

Czy wiesz, że wbudowany typ `error` w Go to po prostu interfejs z jedną metodą?

```go
type error interface {
    Error() string
}

```

Stwórz własną strukturę `CustomError` (`Code int`, `Message string`), zaimplementuj dla niej metodę `Error() string` i zwróć ją jako zwykły `error` z funkcji.

---