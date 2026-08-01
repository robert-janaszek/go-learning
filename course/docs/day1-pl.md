Oto **zadania na Dzień 1: Narzędzia Go (CLI)**.

Zanim wejdziesz w składnię i współbieżność, oswoisz się z podstawowym toolingiem. W Go większość pracy zaczyna się od kilku komend w terminalu — odpowiedników tego, co w JS/TS znasz jako `node`, bundler czy Prettier.

Cel dnia: uruchomić projekt, zbudować binarkę, sformatować kod i odpalić prosty linter.

---

## Część 1: Uruchamianie i budowanie (Zadania 1–4)

### Zadanie 1: Uruchomienie w pamięci

W katalogu projektu (tam, gdzie jest `go.mod` i `main.go`) uruchom:

```bash
go run .
```

albo:

```bash
go run main.go
```

Zobacz, co wypisuje program. To kompilacja „w locie” + od razu start — bez trwałego pliku binarnego (wygodne przy nauce).

### Zadanie 2: Kompilacja do pojedynczej binarki

Zbuduj program:

```bash
go build
```

albo z nazwą wyjścia:

```bash
go build -o go-learning
```

Uruchom powstały plik (np. `./go-learning`). Porównaj z `go run`: `build` zostawia artefakt na dysku.

### Zadanie 3: Formatowanie (`go fmt`)

Zmień celowo wcięcia lub spacje w jakimś pliku `.go`, potem:

```bash
go fmt ./...
```

To odpowiednik Prettiera / formatera w ekosystemie Go — jeden oficjalny styl dla całego ekosystemu.

### Zadanie 4: Podstawowy linter (`go vet`)

Uruchom:

```bash
go vet ./...
```

`vet` szuka typowych pomyłek (np. podejrzane `Printf`, nieskuteczne `Lock`). Nie zastępuje pełnego lintera (jak `staticcheck`), ale to dobry pierwszy krok w CI.

---

## Część 2: Moduł i testy (Zadania 5–8)

### Zadanie 5: Podejrzyj `go.mod`

Otwórz `go.mod`. Zanotuj nazwę modułu i wersję języka (`go 1.xx`). To „manifest” zależności projektu.

### Zadanie 6: Lista pakietów

```bash
go list ./...
```

Zobacz, które pakiety Go widzi w repozytorium (`course`, `json-parser`, …).

### Zadanie 7: Testy

```bash
go test ./...
```

Nawet jeśli dziś nie piszesz testów, warto wiedzieć, jak odpalić cały zestaw.

### Zadanie 8: Pomoc CLI

```bash
go help
go help build
go help test
```

Przejrzyj krótkie opisy — dokumentacja jest wbudowana w narzędzie.

---

## Część 3: Workflow nauki (Zadania 9–10)

### Zadanie 9: Scratch `main.go`

W `main.go` odkomentuj / podepnij wywołanie ćwiczenia z `course` (od Dnia 2 wzwyż). Uruchom ponownie przez `go run .`.

### Zadanie 10: Checklista dnia

Upewnij się, że potrafisz bez ściągi:

1. uruchomić kod (`go run`),
2. zbudować binarkę (`go build`),
3. sformatować (`go fmt ./...`),
4. sprawdzić `go vet ./...`.

Od Dnia 2 wchodzisz w język: zmienne, typy i wskaźniki.
