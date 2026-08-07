# GoHooks — rozszerzenia (po bazowych krokach 1–18)

Ten dokument jest **kontynuacją** [`gohooks-pl.md`](./gohooks-pl.md). Zakładamy, że masz działający runtime: `UseState`, `UseEffect` (deps + cleanup), `Run` / `Unmount`, ochronę kolejności hooków i podstawowe testy.

Cel: przejść od „pętli dirty do limitu 50” do **modelu zdarzeniowego** (kanały), mądrzejszego planowania re-renderu oraz **drzewa komponentów**.

Nie kopiuj tu gotowca z internetu pod React Fiber — to nadal świadomy, edukacyjny silnik.

---

## Co zmieniamy względem bazy

| Było (1–18) | Cel rozszerzeń |
|-------------|----------------|
| `Run` kręci się w pętli, max ~50 renderów | Pętla „żyje” dłużej; budzi się, gdy jest powód |
| Globalny limit renderów | Limit **per hook**: ten sam slot nie może odpalać update’u N razy z rzędu |
| Flaga `dirty` + busy-loop w jednym wywołaniu `Run` | Sygnał na kanale → re-render tylko po realnej zmianie |
| `Component func()` bez dzieci | Komponent może zwracać listę dzieci do renderu |

---

## Faza A: Re-render na kanale (zamiast busy-loop)

### Krok A1: Sygnał „trzeba przerenderować”

Zamiast (albo obok) samej flagi `dirty`, runtime ma kanał, np. `updates chan struct{}` (albo `chan update` z indeksem slotu — patrz A3/A4).

Setter `UseState` po zapisaniu nowej wartości **nie** zakłada, że ktoś zaraz zwinie pętlę `for` w tym samym stosie wywołań. Wysyła sygnał (nieblokująco albo z domyślnym buforem 1 — decyzja projektowa).

Do przemyślenia:

* kanał buforowany `1` vs niebuforowany,
* co, gdy setter wołany z goroutine poza `Render` (na start możesz wymagać: tylko z „wątku” pętli runtime / po `Render`),
* czy wielokrotny `set` przed kolejnym renderem składa się w **jeden** sygnał (zalążek batchingu).

### Krok A2: Pętla zdarzeniowa `Run`

`Run` wygląda mniej więcej tak (szkielet myślowy, nie gotowiec):

1. Pierwszy `Render` (mount).
2. Potem: `for { select { case <-updates: Render; case <-ctx.Done()/quit: return } }`.
3. Koniec życia: zamknięcie kanału / `Unmount` / cancel kontekstu — wybierz jedną czytelną strategię i opisz ją w komentarzu API.

Ważne: **nie** spinuj pustej pętli sprawdzającej `dirty` w hot path. Brak sygnału = goroutine czeka na kanale.

### Krok A3: Re-render tylko gdy stan się faktycznie zmienił

W setterze:

* jeśli nowa wartość jest równa poprzedniej (`==` dla typów porównywalnych albo `DeepEqual` dla nauki) → **nie** wysyłaj sygnału,
* w przeciwnym razie zapisz i wyślij sygnał.

Dodaj test: dwa razy `set(taSamaWartosc)` → jeden render (albo zero dodatkowych), nie dwa.

### Krok A4: `UseSelect` — selektor bez zbędnych renderów

Dodaj hook w stylu:

```go
func UseSelect[S, T any](get func() S, selectFn func(S) T) T
```

albo prostszy wariant na Twoim modelu stanu (np. select z wartości już trzymanej w slocie / ze „store”).

Zachowanie:

* wylicza wartość wybraną (`selected`),
* trzyma poprzednią `selected` w slocie,
* jeśli `selected` się **nie** zmieniło → nie planuj re-renderu (nawet gdy „surowy” stan źródłowy mógł się zmienić w inny sposób — zależnie od tego, jak podepniesz źródło),
* jeśli się zmieniło → zapisz i zaplanuj update (kanał).

To ćwiczenie uczy: **subskrypcja wycinka stanu**, nie całego świata — intuicja jak `useSelector` w ekosystemie React, bez Reduxa.

Demo: komponent pokazuje `count/2` albo `count > 10`; inkrement, który nie zmienia wyniku selecta, nie powinien produkować kolejnego „UI print”.

---

## Faza B: Limit „10 razy pod rząd” (per hook), nie globalne 50

### Krok B1: Licznik kolejnych update’ów slotu

Globalny `max 50 renderów` wywal albo zostaw tylko jako awaryjny hard-cap (opcjonalnie).

Główna ochrona:

* każdy slot, który potrafi planować re-render (`UseState`, `UseSelect`, …), ma licznik `consecutiveSchedules` (nazwa dowolna),
* gdy slot **skutecznie** planuje kolejny render, zwiększ licznik tego slotu,
* gdy **inny** slot planuje update albo gdy render przebiegł „spokojnie” bez stormu — zresetuj liczniki (zdefiniuj regułę i trzymaj się jej),
* po **10** skutecznych zaplanowaniach **tego samego** slotu z rzędu → `panic` / błąd z komunikatem zawierającym indeks hooka.

Cel dydaktyczny: złapać `set(x+1)` bezwarunkowo w ciele komponentu, bez sztucznego „50 renderów i już”.

### Krok B2: Test stormu

Napisz test, który:

1. celowo robi nieskończoną aktualizację z jednego settera,
2. oczekuje panic / błędu po ~10,
3. osobno: naprzemienne update’y dwóch slotów **nie** powinny odpalić limitu jednego z nich (albo powinny — zależnie od Twojej reguły resetu; udokumentuj wybór).

### Krok B3 (opcjonalnie): rozróżnij „schedule” od „render”

Czasem jeden sygnał na kanale = jeden render, ale wiele setterów zbatchowanych. Licznik „pod rząd” liczy sensowniej **zaplanowane flush’e z tego slotu**, niekoniecznie każdą linię `set` w batchu. Zdecyduj i opisz w README / komentarzu.

---

## Faza C: Komponent zwraca listę dzieci

### Krok C1: Nowy typ komponentu

Przejdź z `func()` na coś w stylu:

```go
type Component func() []Component
```

albo:

```go
type Component func() Result
// Result trzyma dzieci + opcjonalny „output” (string do konsoli)
```

Pusty wynik / `nil` = liść (jak obecny Counter, tylko z jawnym return).

### Krok C2: Render drzewa

`Render` po wywołaniu komponentu dostaje wycinek dzieci i **rekurencyjnie** (albo kolejką DFS/BFS) renderuje każde dziecko.

Ustal reguły hooków dla drzewa — to najtrudniejsza część fazy:

* **Wariant prostszy (zalecany na start):** jeden aktywny „węzeł” naraz; przy wejściu w dziecko push kontekstu hooków (własna lista slotów / własny `hookIndex`), przy wyjściu pop. Każdy komponent-instancja ma swój magazyn hooków.
* **Wariant trudniejszy:** jedna globalna kolejka jak w React (kolejność = kolejność wywołań w DFS) — łatwiej o bugi, bliżej prawdziwego modelu.

W docs bazowych kolejność hooków była per jeden komponent — tu musisz mieć **tożsamość instancji**.

### Krok C3: Tożsamość dziecka (key)

Gdy rodzic zwraca listę:

```text
[ChildA, ChildB] → [ChildB, ChildA]
```

bez kluczy zepsujesz stan (jak w React). Dodaj prosty mechanizm:

```go
type Element struct {
    Key       string
    Component Component
}
```

Rodzic zwraca `[]Element`. Runtime mapuje `key → hook storage`. Brak klucza = indeks w liście (z dokumentowanym ograniczeniem).

### Krok C4: Mount / update / unmount dzieci

Przy kolejnym renderze rodzica:

* nowe key → mount (puste sloty, potem efekty),
* zniknięte key → `Unmount` gałęzi (cleanup efektów),
* te same key → re-render z zachowanym stanem.

Demo: lista 3 mini-liczników; usuń środkowy; upewnij się, że cleanup środkowego poszedł, a stan sąsiadów został.

### Krok C5: Testy drzewa

1. Rodzic z dwojgiem dzieci — każde ma własny `UseState`.
2. Reorder bez key vs z key — zaobserwuj różnicę (może być testem behawioralnym na wartościach stanu).
3. Unmount rodzica unmountuje potomków.

---

## Faza D: Dopieszczenie API i demo

### Krok D1: `main` jako mini aplikacja zdarzeniowa

* start `Run` w goroutine albo `Run(ctx)`,
* z `main` / tickera / `stdin` wysyłaj „kliki” (wołania setterów wystawionych bezpiecznie albo przez kolejkę eventów),
* na sygnał quit → `Unmount` + wyjście.

### Krok D2: Batching setterów

Wiele `set` synchronicznie w jednym evencie = **jeden** re-render. Naturalne przy kanale z buforem 1 i flagą „już zaplanowane”.

### Krok D3: Kryteria „gotowe” (rozszerzenia)

* [ ] Pętla runtime czeka na kanale; nie busy-loopuje
* [ ] Setter bez zmiany wartości nie budzi renderu
* [ ] `UseSelect` (lub równoważny) pomija zbędne rendery
* [ ] Limit ~10 kolejnych schedule’y **tego samego** hooka; brak sztywnego „50 i koniec” jako jedynej ochrony
* [ ] Komponent może zwrócić listę dzieci; dzieci mają osobny stan
* [ ] Key (lub świadomy fallback indeksowy) + unmount gałęzi
* [ ] Testy pokrywają kanał, select, storm limit, drzewo

---

## Co jeszcze ma sens dodać? (kolejność sugerowana)

Poniższe nie są wymagane do domknięcia A–D, ale dobrze łączą się z tym modelem:

1. **`UseRef`** — komórka mutowalna bez sygnału na kanał (idealny kontrast do `UseState`).
2. **`UseMemo` / `UseCallback`** — cache przy deps; w Go `UseCallback` mniej boli niż w JS, ale uczy tego samego slotu co memo.
3. **`UseState(updater func(prev T) T)`** — bezpieczniejsze przy batchingu i eventach z kanału.
4. **Prosty Context** (`UseProvide` / `UseContext`) — unikasz prop-drilling w drzewie z fazy C; implementacja: stos na runtime przy DFS.
5. **Jednowątkowa reguła runtime** — dokumentuj: wszystkie hooki i rendery na jednej goroutine; inne goroutine tylko przez `r.Enqueue(fn)` na kanał. To zapobiega race bez mutexów wszędzie.
6. **`UseReducer`** — `dispatch(action)` zamiast wielu setterów; ładnie siedzi na kanale eventów.
7. **Dev-only: log drzewa** (indent DFS) — pomaga debugować mount/unmount dzieci.
8. **Świadomie pominąć:** Virtual DOM diffing stringów HTML, concurrent rendering, Suspense — koszt duży, nauka z A–D już gęsta.

Jeśli wybierasz tylko jedną rzecz po A–D: **Context + reguła jednej goroutine** — domyka drzewo komponentów w praktyczny sposób.

---

## Wskazówki (bez gotowca)

* Setter zamyka `*Runtime` + indeks slotu (jak w bazie); dodatkowo zamyka / używa kanału update’ów należącego do tego runtime.
* `select` + `ctx.Done()` w `Run` uczy cancelacji lepiej niż sam `close(chan)`.
* Przy dzieciach najpierw zrób wariant „osobny storage per instancja + key”, dopiero potem ewentualnie jedną globalną kolejkę hooków.
* Limit „10 pod rząd” testuj panic-em przez `defer recover` w teście — wzorzec idiomatyczny w Go.
* Nie mieszaj w jednym PR: refaktor na kanały **oraz** pełne drzewo dzieci — rozbij na dwa kamienie milowe (Faza A+B, potem C).

---

## Mapowanie na pliki (propozycja)

```text
gohooks/
├── hook/
│   ├── runtime.go      # Run na select/chan, Unmount
│   ├── state.go        # UseState + sygnał + equality
│   ├── select.go       # UseSelect
│   ├── effect.go
│   ├── tree.go         # Element, key, mount/unmount dzieci
│   └── *_test.go
└── components/
    ├── counter.go
    └── list_demo.go    # rodzic zwracający []Element
```

Nazwy dowolne — trzymaj granicę: **sygnały / select / drzewo** czytelnie oddzielone od efektów.
