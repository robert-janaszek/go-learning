Oto **18 zadań projektu: Mini-runtime w stylu React Hooks (w Go)**.

W React masz `useState`, `useEffect` i komponenty funkcyjne. Dziś zbudujesz **uproszczony odpowiednik tego modelu w Go** — nie po to, by zastąpić Reacta, tylko żeby zrozumieć:

* jak działa kolejność wywołań hooków,
* gdzie żyje stan między „renderami”,
* jak efekt zależy od dependency list,
* jak zamknięcie (closure) łączy się z mutowalnym stanem runtime’u.

To ćwiczenie architektoniczne. Implementujesz mały silnik + jeden przykładowy komponent (np. licznik).

---

```text
gohooks/
├── go.mod
├── main.go
├── hook/
│   ├── runtime.go     # kontekst renderu, kolejka hooków, Render()
│   ├── state.go       # UseState
│   ├── effect.go      # UseEffect
│   └── runtime_test.go
└── components/
    └── counter.go     # przykładowy komponent używający hooków
```

Nazwy możesz lekko zmienić, ale trzymaj podział: **runtime / state / effect / komponent**.

---

## Faza 1: Model renderu i runtime (Kroki 1–5)

### Krok 1: Moduł i struktura

Stwórz katalog `gohooks` (obok `taskmanager`) i zainicjalizuj moduł:

```bash
go mod init gohooks
```

Dodaj foldery `hook/` oraz `components/`.

### Krok 2: Typ komponentu

W pakiecie `hook` zdefiniuj typ funkcji komponentu, np.:

```go
type Component func()
```

Na razie komponent nic nie zwraca — „render” to efekt uboczny (`fmt.Println` UI w konsoli). Później możesz zmienić na `func() string`, jeśli wolisz zwracać tekst do wypisania.

### Krok 3: Stan jednego hooka w runtime

Zdefiniuj wewnętrzną strukturę slotu hooka, np. `hookState`, która potrafi trzymać:

* wartość stanu (`any` / `interface{}`),
* opcjonalnie dane efektu (deps, cleanup) — możesz rozbić to na osobne typy w kolejnych krokach.

Na start wystarczy pole na wartość (`any`).

### Krok 4: Runtime renderu

Stwórz strukturę `Runtime` (nazwa dowolna), która trzyma:

* listę / wycinek slotów hooków (`[]hookState` lub podobnie),
* indeks bieżącego hooka podczas jednego renderu (`hookIndex int`),
* referencję do aktualnie renderowanego `Component`,
* flagę „trzeba przerenderować” (np. `dirty bool`) albo kolejkę update’ów.

Zasada jak w React: **podczas jednego wywołania komponentu hooki są czytane/tworzone w stałej kolejności**, a `hookIndex` rośnie przy każdym `UseState` / `UseEffect`.

### Krok 5: Funkcja `Render`

Napisz `func (r *Runtime) Render(c Component)` (lub `Mount`), która:

1. Zeruje `hookIndex` na `0` (kolejność od początku).
2. Ustawia bieżący komponent.
3. Wywołuje `c()`.
4. Po renderze uruchamia efekty (na razie możesz zostawić pusty krok — dojdziesz w Fazie 3).

Uruchom z `main` pusty komponent `func() { fmt.Println("hello") }` i upewnij się, że się wypisuje.

---

## Faza 2: `UseState` (Kroki 6–10)

### Krok 6: Dostęp do bieżącego runtime

Hooki muszą wiedzieć, *w którym* runtime działają. Wybierz jedno podejście (prostsze → trudniejsze):

1. **Globalny / pakietowy wskaźnik** na aktywny `*Runtime` (ustawiany w `Render`) — najprostsze na start, jak wczesny React.
2. **Context przekazany inaczej** — później; nie komplikuj na początku.

Zaimplementuj wariant 1: na wejściu do `Render` ustaw `activeRuntime = r`, na wyjściu wyczyść (np. `defer`).

### Krok 7: `UseState` — pierwsze wywołanie vs kolejne

Zaimplementuj:

```go
func UseState[T any](initial T) (T, func(T))
```

(albo bez generics: `UseState(initial any) (any, func(any))`, jeśli wolisz na razie unikać typów parametrycznych).

Zachowanie:

* Przy **pierwszym** renderze, gdy nie ma jeszcze slotu pod `hookIndex`: utwórz slot, zapisz `initial`.
* Przy **kolejnych** renderach: odczytaj wartość ze slotu o indeksie `hookIndex` (nie resetuj do `initial`).
* Zawsze zwiększ `hookIndex` o 1 przed wyjściem z funkcji.

### Krok 8: Setter i oznaczenie do re-renderu

Setter zwrócony z `UseState` powinien:

1. Zapisać nową wartość w odpowiednim slocie (uwaga: setter żyje dłużej niż jeden render — musi pamiętać **indeks slotu**, nie polegać na bieżącym `hookIndex`).
2. Oznaczyć runtime jako „brudny” / zaplanować kolejny render.

Na tym etapie setter **nie musi** od razu synchronicznie wołać `Render` — wystarczy flaga. W następnym kroku dodasz pętlę.

### Krok 9: Pętla re-renderów

Dodaj metodę w stylu `func (r *Runtime) Run(c Component)` która:

1. Woła `Render(c)`.
2. Dopóki flaga „dirty” jest ustawiona: czyści flagę i znowu woła `Render(c)`.
3. Uważaj na nieskończoną pętlę (np. setter wołany bezwarunkowo w ciele komponentu) — możesz dodać limit renderów (np. max 50) i `panic` / błąd przy przekroczeniu.

### Krok 10: Komponent `Counter`

W `components/counter.go` napisz komponent, który:

* woła `UseState(0)` dla licznika,
* wypisuje aktualną wartość,
* w jakiś sposób zwiększa stan (np. od razu w pierwszym renderze tylko raz — albo przez prostą „symulację zdarzenia” z `main`).

Na razie wystarczy scenariusz testowy z `main`:

* zamontuj `Counter`,
* z zewnątrz wywołaj setter (jeśli go wyeksponujesz) **albo** wewnątrz komponentu po wypisaniu zaplanuj `+1` tylko gdy wartość `< N`.

Cel: zobaczyć w konsoli kolejne wartości `0, 1, 2, …` pochodzące z re-renderów, nie z ręcznej pętli `for` w komponencie omijającej hooki.

---

## Faza 3: `UseEffect` (Kroki 11–15)

### Krok 11: Slot efektu

Rozszerz model slotu (lub dodaj osobny typ) o dane efektu:

* funkcja efektu: `func() func()` (efekt zwraca opcjonalny `cleanup`; cleanup może być `nil`),
* poprzednie zależności: np. `[]any`,
* flaga „czy uruchomić po tym renderze”.

Podpis docelowy:

```go
func UseEffect(effect func() func(), deps ...any)
```

Jeśli wolisz bez variadic: `UseEffect(effect func() func(), deps []any)`.

### Krok 12: Porównanie dependency list

Podczas `UseEffect`:

1. Weź / utwórz slot pod bieżącym `hookIndex`, zwiększ indeks.
2. Porównaj nowe `deps` z zapisanymi poprzednimi.
3. Jeśli to pierwszy raz **albo** deps się zmieniły → oznacz efekt do uruchomienia po renderze.
4. Jeśli deps się **nie** zmieniły → nie uruchamiaj efektu (jak w React).

Porównanie: na start wystarczy proste `==` elementów po rzucie / `reflect.DeepEqual`. Dla nauki `DeepEqual` jest OK.

**Uwaga Reactowa:** `deps == nil` vs `deps == []` — możesz uprościć: zawsze wymagaj listy deps (nawet pustej). Pusty wycinek = „odpal raz po mount”.

### Krok 13: Uruchamianie efektów po renderze

Po wywołaniu komponentu w `Render`:

1. Przejdź po slotach oznaczonych do odpalenia.
2. Jeśli slot miał poprzedni `cleanup` — wywołaj go **zanim** odpalisz nowy efekt.
3. Wywołaj efekt, zapisz zwrócony cleanup.
4. Zapisz nowe deps w slocie.

### Krok 14: Efekt w `Counter`

Rozszerz komponent:

* `count, setCount := UseState(0)`
* `UseEffect(func() func() { fmt.Println("effect: count =", count); return nil }, count)`

Albo z cleanup:

```text
effect start for count=X
… re-render …
cleanup for count=X
effect start for count=X+1
```

Upewnij się, że przy braku zmiany `count` efekt **nie** odpala się ponownie (zrób dwa rendery z tą samą wartością — np. setter z tą samą liczbą — i zaobserwuj brak loga).

### Krok 15: Unmount / cleanup na koniec

Dodaj `func (r *Runtime) Unmount()`:

* wywołaj wszystkie zapisane `cleanup` z efektów,
* wyczyść sloty / oznacz runtime jako niemontowany.

Z `main`: `Run` → trochę update’ów → `Unmount` → upewnij się, że ostatnie cleanup’y poszły.

---

## Faza 4: Zasady hooków, testy, dopieszczenie (Kroki 16–18)

### Krok 16: Zasada kolejności (Rules of Hooks)

Celowo zepsuj komponent: wywołaj `UseState` w `if count > 0 { ... }`. Uruchom kilka renderów i zaobserwuj, co się dzieje (złe wartości, panic, desync indeksów).

Potem dodaj w runtime prostą ochronę, np.:

* zapamiętaj liczbę hooków z pierwszego renderu,
* przy kolejnych, jeśli końcowy `hookIndex` się różni — `panic("hooks order mismatch")`.

To jeden z najważniejszych wniosków całego projektu.

### Krok 17: Testy jednostkowe

W `hook/runtime_test.go` napisz testy (table-driven tam, gdzie ma sens):

1. `UseState` zachowuje wartość między renderami.
2. Setter powoduje kolejny render z nową wartością.
3. `UseEffect` odpala się po pierwszym renderze.
4. `UseEffect` nie odpala się, gdy deps się nie zmieniły.
5. Cleanup woła się przed kolejnym efektem i przy `Unmount`.

Użyj liczników (`int` w closurach / polach testowych), nie tylko `fmt.Println`.

```bash
go test ./...
go test -cover ./...
```

### Krok 18: `main` jako mini demo + build

W `main.go` złóż czytelne demo:

1. Stwórz `Runtime`.
2. Odpal komponent `Counter` (stan + efekt z logami).
3. Symuluj kilka „kliknięć” (settery / zaplanowane update’y).
4. `Unmount`.

Zbuduj binarkę:

```bash
go build -o gohooks main.go
./gohooks
```

---

## Kryteria „gotowe”

Masz działający projekt, gdy:

* [ ] Komponent funkcyjny renderuje się przez runtime
* [ ] `UseState` trzyma stan między renderami i setter planuje update
* [ ] Jest pętla re-renderów (z bezpiecznym limitem)
* [ ] `UseEffect` honoruje deps + cleanup
* [ ] `Unmount` czyści efekty
* [ ] Naruszenie kolejności hooków jest wykrywane
* [ ] Testy pokrywają state i effect
* [ ] Demo w CLI pokazuje zachowanie zbliżone do intuicji Reacta

---

## Wskazówki (bez gotowca)

* Hook musi zamykać **indeks slotu** w setterze/cleanupie — nie używaj globalnego `hookIndex` w momencie późniejszego wywołania settera.
* Generics (`UseState[T]`) są wygodne, ale najpierw możesz zrobić wersję na `any`, potem przepisać.
* Nie buduj Virtual DOM — wystarczy „render” = wypisanie stanu.
* Nie celuj w pełną zgodność z React 18 (batching, concurrent) — to świadomy model edukacyjny.
* Jeśli utkniesz na cyklu importów: `components` → `hook`, nigdy odwrotnie; `main` skleja oba.

---

## Opcjonalne rozszerzenia (po domknięciu 1–18)

* `UseRef` — mutowalna komórka bez re-renderu
* `UseMemo` — cache wartości przy niezmienionych deps
* batching setterów (wiele `setState` = jeden re-render)
* `UseState` z updaterem `func(prev T) T`

Nie rób ich, zanim nie przejdziesz krytycznej ścieżki 1–18.
