Oto **20 zadań projektu: Mini-VM z wirtualnym RAM-em (w Go)**.

Znasz już wskaźniki z Dnia 2, escape analysis („stos vs sterta w Go”) oraz projekty typu parser / gohooks. Dziś zbudujesz **edukacyjną maszynę wirtualną z własną pamięcią**:

* jednolity blok bajtów (wirtualny RAM),
* wskaźniki jako *offsety* w tej pamięci (nie natywne `*T` z Go),
* stos rosnący w dół + stack pointer (SP),
* stertę z `Alloc` / `Free`,
* ramki wywołań (`CALL` / `RET`),
* mały zestaw instrukcji i pętlę wykonania.

Cel nie jest produkcyjny (to nie jest WASM ani prawdziwy GC Go). Chodzi o to, żeby **zobaczyć layout pamięci na własne oczy** — bo w zwykłym Go kompilator chowa go przed Tobą.

---

```text
minivm/
├── main.go
└── vm/
    ├── memory.go      # RAM, Addr, Load/Store słów
    ├── stack.go       # SP, Push/Pop
    ├── heap.go        # Alloc / Free
    ├── frame.go       # ramki wywołań (opcjonalnie osobny plik)
    ├── opcode.go      # kody instrukcji
    ├── machine.go     # VM: rejestry, Execute
    └── vm_test.go
```

Nazwy plików możesz lekko zmienić, ale trzymaj podział: **pamięć / stos / sterta / instrukcje / maszyna**.

---

## Model pamięci (zanim zaczniesz kod)

Przyjmij konwencję i trzymaj się jej w całym projekcie:

```text
niski adres                                              wysoki adres
0 ─────────────────────────────────────────────────────► Size-1
[  HEAP rośnie w górę → ........ wolne ........ ← STACK ]
                          ↑                      ↑
                       HeapBrk                  SP
```

* **Słowo maszyny:** `uint32` (4 bajty) albo `uint64` (8 bajtów) — wybierz jedno i używaj konsekwentnie. Poniżej zakładamy **`uint32` / 4 bajty** (prościej w dumpach).
* **Endianness:** little-endian przy zapisie/odczycie słów z `[]byte`.
* **`Addr`:** typ `type Addr uint32` — to Twój „wskaźnik”. `0` = null (niedozwolony do dereferencji).
* **Stos:** rośnie w dół (SP maleje przy `Push`).
* **Sterta:** rośnie w górę od niskich adresów (po nagłówku / regionie kodu — decyzja w Fazie 1).
* **Kolizja:** jeśli `HeapBrk` dogoni `SP` → błąd (stack overflow / out of memory).

Nie mieszaj natywnych wskaźników Go (`*int`) z adresami VM — w API VM operujesz na `Addr` i `[]byte`.

---

## Faza 1: Wirtualny RAM i wskaźniki (Kroki 1–5)

### Krok 1: Moduł i struktura

Stwórz katalog `minivm` (obok `gohooks` / `taskmanager`) w tym samym module albo jako podpakiet:

```bash
mkdir -p minivm/vm
```

Pakiet implementacji: `vm`. `main.go` tylko składa demo.

### Krok 2: Typ `Addr` i pamięć

W `vm/memory.go`:

```go
type Addr uint32

const WordSize = 4 // bajty; dostosuj jeśli wybierzesz uint64

type Memory struct {
    data []byte
}
```

Napisz:

* `NewMemory(size int) *Memory` — alokuje `make([]byte, size)`, zeruje (zero value slice i tak jest zerami).
* `Size() int`
* walidację wyrównania: adres musi być podzielny przez `WordSize` przy operacjach na słowach.

### Krok 3: `Load` / `Store` słowa

Zaimplementuj:

```go
func (m *Memory) Load(addr Addr) (uint32, error)
func (m *Memory) Store(addr Addr, value uint32) error
```

Zasady:

* odrzuć `addr == 0` (null dereference),
* odrzuć adres poza zakresem (`addr+WordSize > size`),
* odrzuć adres niewyrównany,
* użyj `encoding/binary.LittleEndian` (albo ręcznego składania bajtów).

Dodaj też warianty na pojedynczy bajt (`Load8` / `Store8`), jeśli przydadzą się później — nie są obowiązkowe na start.

### Krok 4: Dump pamięci (debug)

Napisz `func (m *Memory) Dump(start, end Addr) string` (albo `WriteDump(w io.Writer, ...)`), który wypisuje hex w czytelnych liniach (np. 16 bajtów na linię + adres).

Przyda się przy każdym kolejnym kroku — bez dumpu będziesz zgadywać.

### Krok 5: Testy pamięci

W `vm/vm_test.go` (lub `memory_test.go`):

1. `Store` + `Load` zwraca tę samą wartość.
2. `Load(0)` → błąd.
3. Adres poza zakresem → błąd.
4. Adres niewyrównany (np. `1`) → błąd.
5. Little-endian: po `Store(addr, 0x01020304)` bajty w pamięci to `04 03 02 01`.

```bash
go test ./minivm/...
```

---

## Faza 2: Stos i stack pointer (Kroki 6–9)

### Krok 6: Maszyna z SP

W `vm/machine.go` (lub `stack.go`) zdefiniuj:

```go
type VM struct {
    mem  *Memory
    sp   Addr   // stack pointer
    // później: hp (heap break), ip (instruction pointer), fp (frame pointer)...
}
```

Konwencja startowa:

* `mem` ma rozmiar np. `4096` bajtów,
* początkowe `sp = Addr(size)` — „tuż za końcem” / na końcu regionu (ustal jednoznacznie: czy SP wskazuje wolne miejsce, czy ostatni zajęty slot — **udokumentuj w komentarzu** i trzymaj konsekwentnie).

Rekomendacja (klasyczna): **SP wskazuje na wierzchołek (ostatnią wypchniętą wartość)**; pusty stos → `sp == Addr(size)` (poza regionem / sentinel). Alternatywnie: SP wskazuje pierwsze wolne słowo. Wybierz jeden model.

### Krok 7: `Push` / `Pop`

```go
func (v *VM) Push(value uint32) error
func (v *VM) Pop() (uint32, error)
```

`Push`:

1. Oblicz nowy adres wierzchołka (`sp - WordSize`).
2. Sprawdź kolizję ze stertą / dolnym limitem (na razie dolny limit może być `0` lub stały `heapStart`).
3. `Store` wartości.
4. Zaktualizuj `sp`.

`Pop`: odwrotnie; błąd przy pustym stosie.

### Krok 8: Peek i głębokość

Dodaj:

* `Peek() (uint32, error)` — czyta wierzchołek bez zdejmowania,
* `StackDepth() int` — ile słów leży na stosie.

### Krok 9: Testy stosu

1. Push 1, 2, 3 → Pop → 3, 2, 1 (LIFO).
2. Pop na pustym stosie → błąd.
3. Wypchnij tyle, żeby prawie zapełnić pamięć → kolejny Push → błąd (overflow).
4. Dump pokazuje wartości w wysokich adresach.

---

## Faza 3: Wskaźniki VM — pośrednie Load/Store (Kroki 10–11)

### Krok 10: Instrukcyjny podgląd bez pełnego ISA

Zanim zrobisz pełny dekoder, dodaj na `VM` operacje „jak z asemblera”:

```go
func (v *VM) LoadIndirect() error  // Pop addr; Push mem[addr]
func (v *VM) StoreIndirect() error // Pop value; Pop addr; mem[addr] = value
```

(albo z jawnymi argumentami — ważne, żeby **adres był wartością na stosie**, nie natywnym wskaźnikiem Go).

Przykładowy scenariusz ręczny w teście:

1. Wybierz adres w regionie danych, np. `0x100`.
2. `Store(0x100, 42)` bezpośrednio.
3. `Push(0x100)` → `LoadIndirect` → na stosie `42`.

### Krok 11: Wskaźnik do wskaźnika (double indirection)

Zademonstruj w teście / demo:

1. Pod adresem `A` leży wartość `7`.
2. Pod adresem `B` leży wartość `A` (wskaźnik).
3. Sekwencja: Push `B` → LoadIndirect → LoadIndirect daje `7`.

To odpowiednik dnia 2 (`**int`), ale w Twoim RAM-ie.

---

## Faza 4: Sterta — `Alloc` / `Free` (Kroki 12–15)

### Krok 12: Bump allocator (najprostsza sterta)

Dodaj do `VM` (lub osobnej struktury `Heap`):

* `heapStart Addr` — stały początek sterty (np. `0x40`, zostaw miejsce na „statyczne” dane),
* `heapBrk Addr` — aktualny koniec zajętej sterty (break).

```go
func (v *VM) Alloc(nbytes uint32) (Addr, error)
```

Zachowanie bump:

1. Wyrównaj `nbytes` w górę do wielokrotności `WordSize`.
2. Opcjonalnie zapisz nagłówek bloku (rozmiar) — **zalecane**, bo ułatwi `Free` później.
3. Zwróć adres **payloadu** (za nagłówkiem).
4. Przesuń `heapBrk`.
5. Jeśli `heapBrk` zbliży się niebezpiecznie do `sp` → błąd `ErrOutOfMemory`.

Na tym etapie `Free` może być no-op albo jeszcze nie istnieć.

### Krok 13: Nagłówek bloku

Ustal layout, np.:

```text
[ size:u32 ][ payload... ]
            ^
            zwracany Addr
```

Albo z flagą zajętości:

```text
[ size:u32 ][ free:u32 ][ payload... ]
```

Opisz layout w komentarzu nad `Alloc`. Test: po dwóch `Alloc(8)` adresy różnią się o wyrównany rozmiar + nagłówek.

### Krok 14: Free-list (prawdziwy `Free`)

Rozszerz alokator:

```go
func (v *VM) Free(ptr Addr) error
```

Wymagania minimalne:

* `Free` oznacza blok jako wolny i wstawia go na **listę wolnych** (free-list) — wskaźnik `next` może żyć w payloadzie / w polu nagłówka,
* kolejny `Alloc` najpierw szuka bloku na free-list (first-fit wystarczy), dopiero potem bump,
* `Free(0)` → błąd,
* podwójny `Free` tego samego bloku → błąd (detect double-free — choćby flagą `free`),
* `Free` wskaźnika, który nie wskazuje na payload bloku → błąd.

Nie musisz robić coalescingu (łączenia sąsiadów) w wersji bazowej — to opcjonalne rozszerzenie.

### Krok 15: Testy sterty

1. `Alloc` + `Store`/`Load` przez zwrócony adres działa.
2. `Free` + ponowny `Alloc` tego samego rozmiaru **reusing** adresu (albo z free-list — asercja że nie tylko bump rośnie w nieskończoność przy alloc/free w pętli).
3. Alloc aż do kolizji ze stosem → OOM.
4. Double-free → błąd.
5. Stos i sterta: wypchnij dużo na stos, potem `Alloc` musi respektować SP.

---

## Faza 5: Ramki wywołań i ISA (Kroki 16–18)

### Krok 16: Frame pointer i `CALL` / `RET`

Dodaj rejestr `fp Addr` (frame pointer) oraz `ip int` (indeks w slice instrukcji — patrz Krok 17).

Konwencja ramki (uproszczenie — możesz uprościć jeszcze bardziej):

```text
wysokie adresy
  ... poprzednia ramka ...
  [ arg N ]
  [ arg 1 ]
  [ return addr ]  ← po CALL
  [ saved FP    ]  ← FP wskazuje tu (albo na saved FP)
  [ local 0     ]
  [ local 1     ]
niskie adresy (kierunek wzrostu stosu)
```

Zaimplementuj metody (na razie wołane z Go, nie z opcode):

* `Call(returnAddr uint32, nargs int) error` — zapisuje return addr + stary FP, ustawia nowy FP, opcjonalnie rezerwuje sloty na lokale,
* `Ret() (returnAddr uint32, error)` — odtwarza FP i SP, zwraca adres powrotu.

Przetestuj zagnieżdżone „wywołanie” 2–3 poziomów i sprawdź, że lokale nie nadpisują się nawzajem.

### Krok 17: Opcodes i program

W `opcode.go`:

```go
type Op byte

const (
    OpHalt Op = iota
    OpPush      // operand: u32 literal
    OpPop
    OpAdd       // pop a, pop b, push b+a
    OpLoad      // pop addr, push mem[addr]
    OpStore     // pop val, pop addr, mem[addr]=val
    OpAlloc     // pop nbytes, push ptr
    OpFree      // pop ptr
    OpCall      // operand: target ip; nargs na stosie wg Twojej konwencji
    OpRet
    OpDup       // duplikuj wierzchołek (wygodne)
    OpPrint     // pop i loguj — tylko do demo
)
```

```go
type Instr struct {
    Op   Op
    Arg  uint32 // używane gdy instrukcja ma operand
}
```

`VM` trzyma `code []Instr` oraz `ip int`.

### Krok 18: Pętla `Execute`

```go
func (v *VM) Execute(code []Instr) error
```

Pętla:

1. `for { instr := code[ip]; ip++; switch instr.Op ... }`
2. `OpHalt` → return nil.
3. Błędy z `Push`/`Pop`/`Alloc` propaguj w górę (z numerem `ip` w komunikacie — bardzo pomaga).
4. Limit instrukcji (np. 100_000) przeciw nieskończonej pętli — nawet jeśli jeszcze nie masz skoków warunkowych.

Zaimplementuj najpierw: `Push`, `Pop`, `Add`, `Load`, `Store`, `Halt`. Potem doklej `Alloc`/`Free`/`Call`/`Ret`.

---

## Faza 6: Demo, bezpieczeństwo, domknięcie (Kroki 19–20)

### Krok 19: Program demonstracyjny w `main`

Złóż w `main.go` czytelne demo (nie jeden spaghetti dump):

**Program A — arytmetyka na stosie**

* Push 40, Push 2, Add, Print → `42`.

**Program B — sterta**

* Alloc 4 bajtów, Store wartości `7` pod zwrócony wskaźnik, Load z powrotem, Print, Free.

**Program C — „funkcja” przez CALL/RET**

* Prosta procedura: bierze 1 argument ze stosu, dodaje `1`, zwraca przez Push wyniku + `Ret`.
* Main: Push arg, Call, Print.

Po każdym programie wypisz krótki dump (SP, FP, heapBrk + wycinek pamięci).

```bash
go run ./minivm
```

### Krok 20: Testy integracyjne + kryteria jakości

Dopisz testy table-driven na `Execute`:

1. Program `40 + 2` → wierzchołek / Print path = 42 (możesz zbierać output przez callback zamiast `fmt`).
2. Alloc/Store/Load/Free w bytecode.
3. Call/Ret zagnieżdżone.
4. Stack overflow przy Push w pętli.
5. OOM przy Alloc.
6. Null load (`Push 0`, `OpLoad`) → błąd.

Użyj `t.Helper()` i czytelnych komunikatów z `ip` / `sp`.

```bash
go test ./minivm/...
go test -cover ./minivm/...
```

---

## Kryteria „gotowe”

Masz działający projekt, gdy:

* [ ] Jest jednolity RAM (`[]byte`) z `Load`/`Store` słów i typem `Addr`
* [ ] Null, OOB i niewyrównany dostęp zwracają błędy
* [ ] Stos LIFO z SP rosnącym w dół + ochrona overflow
* [ ] `Alloc` (bump) + `Free` (free-list first-fit) na tej samej pamięci co stos
* [ ] Kolizja heap↔stack jest wykrywana
* [ ] Double-free / free(null) są odrzucane
* [ ] `CALL`/`RET` utrzymują ramki (FP) bez kasowania lokalnych danych wyżej
* [ ] `Execute` odpala bytecode z `Halt`, arytmetyką, load/store, alloc/free
* [ ] Testy jednostkowe + integracyjne przechodzą
* [ ] Demo w CLI pokazuje A/B/C i dump pamięci

---

## Wskazówki (bez gotowca)

* **Jedna prawda o SP:** napisz 5 linii komentarza przy polu `sp` i nie zmieniaj semantyki w połowie projektu.
* Trzymaj `error` wartościami pakietowymi: `var ErrNullDeref = errors.New(...)`, `ErrStackOverflow`, `ErrOutOfMemory`, `ErrDoubleFree` — testuj przez `errors.Is`.
* Nie używaj `panic` na ścieżce „zły program użytkownika” — panic zostaw na bugi VM (np. niespójny stan wewnętrzny).
* Dump hex > `fmt.Printf` ad hoc przy debugowaniu alokatora.
* Nie implementuj GC, stronicowania, MMU ani wątków — świadomie pomijamy.
* Natywne `*T` w Go mogą trzymać *strukturę* VM; **dane programu** żyją tylko w `Memory.data`.

---

## Opcjonalne rozszerzenia (po domknięciu 1–20)

* coalescing wolnych bloków (łącz sąsiadów przy `Free`)
* `OpJmp` / `OpJmpIfZero` + prosta pętla w bytecode
* osobny region **code** vs **data** (teraz code jest w Go `[]Instr` — możesz przenieść immediates do RAM)
* 2–3 rejestry ogólnego użytku obok stosu (akumulator)
* serializacja dumpu pamięci do testów golden-file
* „use after free” detector: stub wypełniający freed payload wzorcem `0xDE` i flaga generacji

Nie rób ich, zanim nie przejdziesz krytycznej ścieżki 1–20.
