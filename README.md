[![Build](https://github.com/Bookshelf-Writer/treelang/actions/workflows/pylint.yml/badge.svg?branch=main)](https://github.com/Bookshelf-Writer/treelang/actions/workflows/pylint.yml)

![GitHub repo file or directory count](https://img.shields.io/github/directory-file-count/Bookshelf-Writer/treelang?color=orange)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/Bookshelf-Writer/treelang?color=green)
![GitHub repo size](https://img.shields.io/github/repo-size/Bookshelf-Writer/treelang)

## Script for working with language trees

```bash
 go generate ./...
```

```bash
go build -ldflags="-s -w" -o treelang.linux.bin  ./cli/...
```

```bash
./treelang.linux.bin
```

```bash
./treelang.linux.bin --help
```

```bash
./treelang.linux.bin info
```

---

```bash
./treelang.linux.bin diff -h
```

```bash
./treelang.linux.bin diff --from treeLang/
```

```bash
./treelang.linux.bin diff --from treeLang/ --mode all
```

```bash
./treelang.linux.bin diff --from treeLang/ --master treeLang/en.json
```

```bash
./treelang.linux.bin diff --from treeLang/ --master treeLang/en.json --full
```

```bash
./treelang.linux.bin diff --master treeLang/en.json --slave treeLang/pl.json
```

```bash
./treelang.linux.bin diff --master treeLang/en.json --slave treeLang/pl.json --full
```

---

```bash
./treelang.linux.bin generate -h
```

```bash
./treelang.linux.bin generate --from treeLang --to TEMP --json --master TEMP/treelang_en.gen.yml
```

```bash
./treelang.linux.bin generate --from treeLang/en.json --to TEMP --json
```

```bash
./treelang.linux.bin generate --from treeLang --to TEMP --json --map
```

```bash
./treelang.linux.bin generate --from treeLang/en.json --to TEMP --yml
```

```bash
./treelang.linux.bin generate --from treeLang --to TEMP --yml --map
```

```bash
./treelang.linux.bin generate --from treeLang/en.json --to TEMP --go-package main
```

```bash
./treelang.linux.bin generate --from treeLang/en.json --to TEMP --go-package main --func-png
```

```bash
./treelang.linux.bin generate --from treeLang --to TEMP --go-package main --map
```

---

```bash
./treelang.linux.bin generate --from tmp/localization/iBregus_Secret_App/uk.json --to tmp/localization --go-package localization
```

```bash
./treelang.linux.bin generate --from tmp/localization/iBregus_Secret_App/uk.json --to tmp/localization --go-package localization --func-png
```

```bash
./treelang.linux.bin generate --from tmp/localization/iBregus_Secret_App/uk.json --to tmp/localization --go-package localization --map
```

---

### Mirrors

- https://git.bookshelf-writer.fun/Bookshelf-Writer/treelang
 
