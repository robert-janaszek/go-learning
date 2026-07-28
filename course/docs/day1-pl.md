Oto **20 krótkich, wysoce praktycznych zadań**, które przeprowadzą Cię od podstaw zmiennych i typów, przez manipulację pamięcią, aż po niuanse przekazywania danych do funkcji.

Podzieliłem je na sekcje, abyś stopniowo budował intuicję wokół wskaźników.

---

## Część 1: Zmienne, Typy i Deklaracje (Zadania 1–5)

### Zadanie 1: Różne sposoby deklaracji

Zadeklaruj zmienną `age` typu `int` na 3 sposoby: z użyciem `var` z podaniem typu, z użyciem `var` z inferencją typu oraz z użyciem krótkiej deklaracji `:=`. Wydrukuj ich typy za pomocą `fmt.Printf("%T\n", age)`.

### Zadanie 2: Konwersja typów (Brak implicit casting)

W JS możesz dodać `5 + "5"`. W Go nie ma automatycznej rzutowania typów. Stwórz zmienną `a int = 42` oraz `b float64 = 3.14`. Przekonwertuj `a` na `float64`, dodaj do `b` i wynik przypisz do nowej zmiennej.

### Zadanie 3: Zero Values (Wartości Domyślne)

W JS niezainicjalizowana zmienna to `undefined`. Zadeklaruj za pomocą `var` bez przypisywania wartości: `int`, `float64`, `string`, `bool` oraz wskaźnik `*int`. Wydrukuj je i sprawdź, jakie mają wartości domyślne (*zero values*).

### Zadanie 4: Stałe (`const`) i `iota`

Stwórz blok stałych reprezentujących dni tygodnia (od `Monday` do `Sunday`) z użyciem generatora `iota`. Wyświetl ich wartości liczbowe w konsoli.

### Zadanie 5: Shadowing (Przysłanianie zmiennych)

Stwórz zmienną `x := 10` w zewnętrznym bloku. Otwórz nowy blok kodu `{ ... }`, zadeklaruj w nim `x := 20` i wyświetl `x`. Poza blokiem wyświetl ponownie `x`. Przeanalizuj, co się stało.

---

## Część 2: Podstawy Wskaźników – Adresy i Dereferencja (Zadania 6–10)

### Zadanie 6: Pobieranie adresu (`&`)

Stwórz zmienną `score := 100`. Stwórz zmienną `ptr`, która przechowuje adres zmiennej `score`. Wydrukuj wartość `score`, adres `score` oraz typ zmiennej `ptr`.

### Zadanie 7: Dereferencja (`*`)

Mając wskaźnik `ptr` z Zadania 6, zmień wartość `score` na `200` **używając tylko wskaźnika `ptr**` (operator dereferencji `*ptr = ...`). Wydrukuj `score`.

### Zadanie 8: Dwa wskaźniki do jednej zmiennej

Stwórz `x := 50`. Stwórz dwa wskaźniki `p1` oraz `p2`, oba wskazujące na `x`. Zmień wartość przez `p1` na `100`, a następnie wydrukuj wartość pobraną przez `*p2`.

### Zadanie 9: Wskaźnik typu Nil (`nil pointer`)

Zadeklaruj wskaźnik `var p *int` (bez przypisywania adresu). Sprawdź warunkiem `if p == nil`, czy wskaźnik jest pusty. Zobacz, co się stanie (i jaki błąd zgłosi runtime), jeśli spróbujesz zrobić `*p = 10` bez inicjalizacji (tzw. *nil pointer dereference panic*).

### Zadanie 10: Podwójny wskaźnik (`**int`)

Stwórz zmienną `val := 42`. Stwórz wskaźnik `p` wskazujący na `val`. Stwórz wskaźnik do wskaźnika `pp` (typ `**int`) wskazujący na `p`. Odczytaj wartość `42` używając tylko `pp`.

---

## Część 3: Wskaźniki w Funkcjach (Zadania 11–15)

### Zadanie 11: Swap (Zamiana wartości)

Napisz funkcję `swap(a, b *int)`, która zamienia wartości dwóch zmiennych miejscami. Przetestuj ją w `main()` na dwóch zmiennych `x := 1` i `y := 2`.

### Zadanie 12: Modyfikacja napisu (`string`)

Napisz funkcję `uppercase(s *string)`, która przyjmuje wskaźnik do `stringa` i zmienia jego treść na wielkie litery (użyj `strings.ToUpper`). Sprawdź wynik w `main()`.

### Zadanie 13: Bezpieczne dzielenie z opcjonalnym wynikiem

Napisz funkcję `safeDivide(a, b float64, result *float64) bool`. Jeśli `b == 0`, funkcja zwraca `false`. Jeśli nie, zapisuje wynik pod adres `result` i zwraca `true`.

### Zadanie 14: Licznik (Incementator)

Stwórz strukturę/zmienną reprezentującą licznik. Napisz funkcję `increment(val *int)`, która zwiększa wartość o 1 przy każdym wywołaniu. Wywołaj ją 3 razy w pętli.

### Zadanie 15: Funkcja zwracająca wskaźnik

Napisz funkcję `createInt(val int) *int`, która tworzy lokalną zmienną wewnątrz funkcji i zwraca jej adres `&localVal`.
*(Kontekst dla programisty C/C++: W Go to jest w pełni bezpieczne! Kompilator zrobi Escape Analysis i zaalokuje tę zmienną na stercie zamiast na stosie).*

---

## Część 4: Wskaźniki i Struktury (Zadania 16–20)

### Zadanie 16: Struktura Gracz i modyfikacja stanu

Zdefiniuj strukturę `Player` z polami `Name string` i `Health int`. Napisz funkcję `takeDamage(p *Player, amount int)`, która zmniejsza punkty życia gracza.

### Zadanie 17: Receiver wartości vs Receiver wskaźnika

Do struktury `Player` z dodaj dwie metody:

1. `HealValue(amount int)` z *value receiverem* `(p Player)`
2. `HealPointer(amount int)` z *pointer receiverem* `(p *Player)`
Wywołaj obie w `main()` i zaobserwuj, która metoda faktycznie leczy gracza.

### Zadanie 18: Automatyczna dereferencja w strukturach

Stwórz wskaźnik do struktury `p := &Player{Name: "Gopher", Health: 100}`. Zmień jego pole `Health` na `90` wpisując po prostu `p.Health = 90`. Zauważ, że w Go nie musisz pisać `(*p).Health = 90` – język robi to automatycznie!

### Zadanie 19: Opcjonalne pola w strukturach (Koncept z JS/TS)

W TypeScript masz opcjonalne pola `age?: number` (które mogą być `undefined`). W Go robi się to wskaźnikami!
Stwórz strukturę `User` z polami `Name string` oraz `Age *int`. Stwórz jednego użytkownika bez podanego wieku (`Age = nil`) i jednego z wiekiem.

### Zadanie 20: Wyciek pamięci / Pętla na strukturach

Stwórz wycinek (slice) struktur `[]Player`. Prziteruj po nim pętlą `for _, player := range players`. Spróbuj zmienić `player.Health = 0` wewnątrz pętli. Sprawdź, dlaczego to **nie działa** (zmienna w pętli `range` jest tylko kopią) i jak to naprawić używając indeksów `players[i].Health = 0`.

---