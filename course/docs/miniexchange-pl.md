Oto **20 zadań projektu: Mini-giełda / Order Matching Engine (w Go)**.

Znasz już goroutines, kanały, `select`, `WaitGroup` i mutexy z Dnia 7 oraz projekty typu parser / gohooks / minivm. Dziś zbudujesz **uproszczony silnik kojarzenia zleceń** — abstrakcję tego, co stoi za giełdą (matching engine):

* traderzy produkujący zlecenia równolegle,
* jedna goroutine właścicielem książki zleceń (order book),
* kojarzenie bid/ask w transakcje (`Trade`),
* portfele / salda chronione mutexem,
* fan-out notowań do subskrybentów przez kanały,
* czyste zamknięcie sesji (`close` + `WaitGroup` + brak race).

Cel nie jest produkcyjny (to nie jest FIX, nie ma persystencji ani prawdziwych pieniędzy). Chodzi o to, żeby **zobaczyć, gdzie w prawdziwym systemie idą kanały, a gdzie mutex** — i dlaczego order book zwykle ma jednego właściciela.

---

```text
miniexchange/
├── main.go
└── exchange/
    ├── order.go         # Side, Order, Trade
    ├── book.go          # order book + Match
    ├── portfolio.go     # salda pod mutexem
    ├── engine.go        # pętla silnika na kanałach
    ├── trader.go        # równolegli traderzy
    ├── marketdata.go    # fan-out notowań
    └── exchange_test.go
```

Nazwy plików możesz lekko zmienić, ale trzymaj podział: **model / book / portfolio / engine / traderzy / market data**.

---

## Model domeny (zanim zaczniesz kod)

Przyjmij konwencję i trzymaj się jej w całym projekcie:

```text
                    ┌─────────────┐
   trader A ──┐     │   orders    │      ┌──────────────┐
   trader B ──┼────►│   chan      │─────►│    Engine    │
   trader C ──┘     └─────────────┘      │  (1 goroutine)│
                                         │  owns Book   │
                                         └──────┬───────┘
                                                │ trades chan
                        ┌───────────────────────┼───────────────────────┐
                        ▼                       ▼                       ▼
                 Portfolio.Apply            MarketData              (testy / log)
                 (mutex na mapie)           fan-out do N subów
```

* **Jeden symbol** na start (np. `"GOOG"`) — nie buduj multi-asset routera.
* **Ceny i ilości:** `int` (np. cena w groszach / centach). Unikaj `float64` w matching’u.
* **Order book:** bids posortowane **malejąco** po cenie, asks **rosnąco**. Przy remisie ceny — FIFO (kolejność przyjścia).
* **Matching:** najlepszy bid ≥ najlepszy ask → powstaje `Trade` po cenie **spoczywającego** zlecenia (maker) — albo zawsze po cenie ask; **wybierz jedną regułę i udokumentuj**.
* **Owner booka:** tylko `Engine` mutuje książkę. Traderzy **nie** sięgają do `Book` bezpośrednio.
* **Mutex:** chroni portfele (mapa trader → cash/pozycja), nie książkę. Book jest single-writer przez design.

Nie mieszaj odpowiedzialności: kanał = przepływ zdarzeń; mutex = dzielony stan odczytywany/zapisywany z wielu goroutines.

---

## Faza 1: Model zleceń i transakcji (Kroki 1–4)

### Krok 1: Katalog i pakiet

Stwórz katalog `miniexchange` (obok `minivm` / `gohooks`) w tym samym module:

```bash
mkdir -p miniexchange/exchange
```

Pakiet implementacji: `exchange`. `main.go` tylko składa demo sesji.

### Krok 2: `Side`, `Order`, `Trade`

W `exchange/order.go`:

```go
type Side int

const (
    SideBuy Side = iota
    SideSell
)

type Order struct {
    ID       int64
    TraderID string
    Side     Side
    Symbol   string
    Price    int // np. centy
    Qty      int // pozostała ilość (zmniejsza się przy partial fill)
}

type Trade struct {
    Symbol   string
    Price    int
    Qty      int
    BuyID    int64  // ID zlecenia kupna
    SellID   int64  // ID zlecenia sprzedaży
    BuyTrader  string
    SellTrader string
}
```

Dodaj `String() string` dla `Side` (wygodne w logach).

### Krok 3: Walidacja zlecenia

Napisz:

```go
func ValidateOrder(o Order) error
```

Odrzuć:

* puste `TraderID` / `Symbol`,
* `Price <= 0`, `Qty <= 0`,
* nieznany `Side`.

Zwracaj sentinel errors (`var ErrInvalidOrder = errors.New(...)`) i testuj przez `errors.Is`.

### Krok 4: Testy modelu

W `exchange_test.go`:

1. Poprawne zlecenie przechodzi walidację.
2. `Qty == 0` → błąd.
3. `String()` dla Buy/Sell jest czytelny.

```bash
go test ./miniexchange/...
```

---

## Faza 2: Order book i matching (Kroki 5–9)

### Krok 5: Struktura `Book`

W `exchange/book.go`:

```go
type Book struct {
    bids []Order // malejąco po Price, potem FIFO
    asks []Order // rosnąco po Price, potem FIFO
}
```

Na start wystarczy posortowany slice + wstawianie w dobre miejsce (binary search opcjonalnie — linear insert jest OK edukacyjnie).

Nie eksportuj mutacji booka poza pakietem bez potrzeby — metody na `*Book` wystarczą.

### Krok 6: `Add` bez natychmiastowego matcha

Zaimplementuj wewnętrzne wstawianie zlecenia na właściwą stronę książki z zachowaniem porządku cenowego + FIFO przy remisie.

Na tym etapie **nie** kojarz jeszcze — najpierw przetestuj samą kolejność.

### Krok 7: `Match` / `Accept`

Napisz metodę, która przyjmuje nowe zlecenie i:

1. Próbuje skojarzyć je z przeciwną stroną książki, dopóki się da.
2. Generuje zero lub więcej `Trade` (partial fills dozwolone).
3. Zmniejsza `Qty` na zleceniach; usuwa w pełni wypełnione z booka.
4. Resztę agresywnego zlecenia (jeśli `Qty > 0`) dokłada do booka jako maker.

Przykład:

* w booku leży Sell `@100 qty=5`
* przychodzi Buy `@100 qty=3` → Trade `qty=3`, ask maleje do `2`
* przychodzi Buy `@100 qty=10` → Trade `qty=2`, reszta Buy `qty=8` ląduje jako bid

### Krok 8: Odczyt top-of-book

Dodaj:

```go
func (b *Book) BestBid() (Order, bool)
func (b *Book) BestAsk() (Order, bool)
func (b *Book) Depth() (bids, asks int)
```

Przyda się w testach i w market data.

### Krok 9: Testy matchingu (table-driven)

Scenariusze obowiązkowe:

1. Brak przeciwieństwa → zlecenie ląduje w booku, brak trade’ów.
2. Pełne skojarzenie 1:1.
3. Partial fill (agresywne większe niż spoczywające).
4. Jedno agresywne zjada kilka poziomów cenowych.
5. Zlecenie po gorszej cenie **nie** matchuje (Buy `@90` vs Ask `@100`).
6. FIFO przy tej samej cenie: starsze zlecenie wypełniane pierwsze.

To serce projektu — bez zielonych testów nie idź dalej.

---

## Faza 3: Silnik na kanałach (Kroki 10–13)

### Krok 10: `Engine` i kanały

W `exchange/engine.go`:

```go
type Engine struct {
    orders chan Order
    trades chan Trade
    // book, portfolio, nextID, wg / done — wg Twojego designu
}
```

Konstruktor np.:

```go
func NewEngine(orderBuf, tradeBuf int) *Engine
```

* `orders` — wejście (traderzy piszą),
* `trades` — wyjście (każdy skojarzony trade),
* bufory > 0, żeby demo nie deadlockowało przy burstach; w testach możesz użyć małych buforów świadomie.

### Krok 11: Pętla `Run`

```go
func (e *Engine) Run()
```

Uruchamiana **w jednej goroutine**:

```go
for o := range e.orders {
    // nadaj ID jeśli trzeba
    // ValidateOrder
    // trades := book.Accept(o)
    // wyślij każdy trade na e.trades
    // zaktualizuj portfolio
}
close(e.trades) // gdy orders się zamknie i pętla wyjdzie
```

Zasada: **tylko ta goroutine** woła metody mutujące `Book`.

### Krok 12: Submit i stop

Dodaj API:

```go
func (e *Engine) Submit(o Order) error  // non-blocking lub zwykły send — wybierz i opisz
func (e *Engine) Close()                // close(orders); poczekaj aż Run skończy
```

`Close` musi być idempotentne albo jasno udokumentowane jako one-shot. Po `Close` kolejne `Submit` → błąd (np. `ErrSessionClosed`).

Wzorzec stopu: `close(orders)` → `Run` domyka matching → `close(trades)` → konsumenci market data kończą `range`.

### Krok 13: Test silnika (bez race)

Test integracyjny:

1. `go e.Run()`
2. Wyślij 2 komplementarne zlecenia z testu.
3. Odbierz dokładnie 1 `Trade` z oczekiwaną ceną/ilością.
4. `Close()`, upewnij się że `trades` się zamyka (pętla `range` kończy się).

```bash
go test -race ./miniexchange/...
```

---

## Faza 4: Portfele (mutex) i traderzy (parallelism) (Kroki 14–16)

### Krok 14: `Portfolio` pod `sync.Mutex`

W `exchange/portfolio.go`:

```go
type Portfolio struct {
    mu       sync.Mutex
    cash     map[string]int // TraderID → gotówka
    position map[string]int // TraderID → szt. instrumentu
}
```

API:

* `Seed(traderID string, cash, position int)`
* `Apply(trade Trade) error` — buyer traci `price*qty` cash, zyskuje pozycję; seller odwrotnie
* `Snapshot(traderID string) (cash, position int)` — odczyt pod `RLock` jeśli użyjesz `RWMutex`

Odrzuć trade, który zrobiłby ujemne cash / short (albo pozwól na short — **zadecyduj i przetestuj**). Rekomendacja edukacyjna: **odrzucaj brak środków na etapie Accept/Apply** (prostszy invariant).

Uwaga projektowa: jeśli odrzucasz po zmatchowaniu, book i portfolio mogą się rozjechać. Czyściej: sprawdzaj salda *przed* finalizacją trade’u albo traktuj seed jako wystarczająco duży w demo i w `Apply` tylko panikuj/loguj inconsistency w testach. Opisz wybraną strategię w komentarzu nad `Apply`.

### Krok 15: Równolegli traderzy

W `exchange/trader.go`:

```go
func RunTrader(e *Engine, id string, orders []Order, wg *sync.WaitGroup)
```

Każdy trader:

1. `defer wg.Done()`
2. Wysyła swoją listę zleceń przez `e.Submit` (z małym `time.Sleep` / jitterem opcjonalnie).
3. Nie dotyka booka ani map portfela bezpośrednio.

W teście / demo odpal 3–5 traderów naraz (`go RunTrader(...)`) i poczekaj `wg.Wait()`, potem `e.Close()`.

### Krok 16: Test wyścigu portfeli

Napisz test, w którym wiele trade’ów aktualizuje te same portfele równolegle z perspektywy konsumentów — albo engine woła `Apply` z jednego wątku (OK), a **równoległe `Snapshot`** lecą z innych goroutines podczas sesji.

Odpal:

```bash
go test -race ./miniexchange/...
```

Cel: zero raportów race. Jeśli `-race` krzyczy na `Book` — złamałeś regułę single-owner.

---

## Faza 5: Market data (fan-out) i domknięcie (Kroki 17–20)

### Krok 17: Broadcaster notowań

W `exchange/marketdata.go` zaimplementuj fan-out:

* wejście: `<-chan Trade` (z engine),
* rejestracja subskrybentów: każdy dostaje własny `chan Trade` (buforowany),
* goroutine czyta trade i rozsyła do wszystkich subów,
* po zamknięciu wejścia — zamknij kanały subskrybentów.

```go
type MarketData struct { /* ... */ }

func NewMarketData(in <-chan Trade) *MarketData
func (m *MarketData) Subscribe(buf int) <-chan Trade
func (m *MarketData) Run()
```

Uwaga: wolny subskrybent nie powinien blokować całej giełdy w nieskończoność. W wersji bazowej: buforowany kanał + drop / skip gdy pełny (`select` + `default`) — **udokumentuj** wybraną politykę.

### Krok 18: Sygnał końca sesji (`chan struct{}`)

Dodaj kanał sygnałowy (0 bajtów payloadu):

```go
done := make(chan struct{})
```

Po `Close()` / po wyjściu `Run` zamknij `done` (albo wyślij `struct{}{}`), żeby `main` i subskrybenci mogli selectem zareagować na koniec sesji. To ćwiczenie z Day 7 zadanie 20 — w prawdziwym kontekście.

### Krok 19: Demo w `main.go`

Złóż czytelną sesję (nie spaghetti):

1. Seed portfeli 3 traderów.
2. Start `Engine.Run` + `MarketData.Run` w goroutines.
3. Subskrybent notowań wypisuje trade’y.
4. Odpal traderów z zestawem zleceń, które **na pewno** wygenerują kilka transakcji + coś zostawią w booku.
5. `wg.Wait` → `Engine.Close` → poczekaj na `done` / domknięcie market data.
6. Wypisz snapshoty portfeli + `Depth()` booka (możesz dodać `Engine.BookSnapshot()` tylko do debug — uważaj na race: albo wołaj po stopie, albo kopiuj pod kontrolą engine).

```bash
go run ./miniexchange
go test -race ./miniexchange/...
```

### Krok 20: Kryteria jakości i pokrycie

Dopisz / dogryź testy:

1. Matching table-driven (Faza 2).
2. Engine: submit → trade → close.
3. Portfolio: Apply zmienia cash/position atomowo względem Snapshot.
4. MarketData: 2 subskrybentów dostaje ten sam trade.
5. Całość z `-race`.

```bash
go test -race -cover ./miniexchange/...
```

---

## Kryteria „gotowe”

Masz działający projekt, gdy:

* [ ] `Order` / `Trade` / walidacja działają i mają testy
* [ ] `Book` sortuje bid/ask i matchuje z partial fill + FIFO
* [ ] `Engine` czyta z `orders`, pisze na `trades`, jest **jedynym** mutującym book
* [ ] `Close` domyka sesję bez deadlocka (`close` + koniec `range`)
* [ ] `Portfolio` chronione mutexem; salda spójne po trade’ach
* [ ] Wielu traderów wysyła zlecenia równolegle (`WaitGroup`)
* [ ] Market data robi fan-out trade’ów do ≥2 subskrybentów
* [ ] `go test -race ./miniexchange/...` przechodzi czysto
* [ ] Demo CLI pokazuje sesję: zlecenia → trade’y → snapshoty

---

## Wskazówki (bez gotowca)

* **Book = single writer.** Jeśli kuszi Cię `mutex` na booku „na wszelki wypadek”, najpierw sprawdź czy nie wystarczy jedna goroutine — to główna lekcja architektury.
* Mutex trzymaj blisko mapy (`Portfolio`), nie „na całym Engine”.
* Zamyka kanał **tylko nadawca**. Traderzy nie zamykają `orders` — robi to `Engine.Close` / `main` po `WaitGroup`.
* Nie wysyłaj na zamknięty kanał — stąd `ErrSessionClosed` albo `sync.Once` wokół `close`.
* `-race` jest częścią Definition of Done, nie opcją.
* Nie buduj UI, REST, WebSocket, persystencji, autentykacji, short-sell margin call — świadomie pomijamy.
* Ceny jako `int` > `float64`. Matching na floatach to osobny koszmar.

---

## Opcjonalne rozszerzenia (po domknięciu 1–20)

* anulowanie zlecenia (`Cancel` jako osobny typ event’u na tym samym kanale — `interface{}` / sum type)
* wiele symboli (mapa `symbol → *Book` nadal w jednej goroutine engine)
* `context.Context` do cancel sesji zamiast tylko `Close`
* limit zleceń / rate limit per trader (`time.Ticker`)
* stop-loss / order types (`Limit` vs `Market`)
* snapshot order book okresowo na osobny kanał (market depth)
* journal trade’ów do pliku JSONL
* hedged „best price” między dwoma bookami (ćwiczenie first-responder z Day 7)

Nie rób ich, zanim nie przejdziesz krytycznej ścieżki 1–20.
