# Go URL Shortener

En enkel og effektiv URL-forkortingstjeneste skrevet i Go.

## Funksjoner

- Forkorter lange URLer til korte, lette lenker
- Genererer QR-koder for alle forkortede URLer
- Minimalistisk og brukervennlig grensesnitt

## Installasjon

```bash
# Klon repositoriet
git clone https://github.com/bjornhels/go-shortener.git
cd go-shortener

# Last ned avhengigheter
go mod download
```

## Bruk

```bash
# Start serveren
go run main.go
```

Besøk `http://localhost:8009` i nettleseren din for å bruke tjenesten.

## Teknologier

- Go 1.24
- [go-qrcode](https://github.com/skip2/go-qrcode) for QR-kodegenerering
- Standard Go HTTP-bibliotek 