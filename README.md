[![Build](https://github.com/Bookshelf-Writer/TreeOfLanguages/actions/workflows/pylint.yml/badge.svg?branch=main)](https://github.com/Bookshelf-Writer/TreeOfLanguages/actions/workflows/pylint.yml)

![GitHub repo file or directory count](https://img.shields.io/github/directory-file-count/Bookshelf-Writer/TreeOfLanguages?color=orange)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/Bookshelf-Writer/TreeOfLanguages?color=green)
![GitHub repo size](https://img.shields.io/github/repo-size/Bookshelf-Writer/TreeOfLanguages)

## Script for working with language trees

```bash
go build -o treelang.bin
```

```bash
./treelang.bin
```

```bash
./treelang.bin --help
```

```bash
./treelang.bin info
```

---

```bash
./treelang.bin diff -h
```

```bash
./treelang.bin diff --from treeLang/
```

```bash
./treelang.bin diff --from treeLang/ --mode all
```

```bash
./treelang.bin diff --from treeLang/ --master treeLang/en.json
```

```bash
./treelang.bin diff --from treeLang/ --master treeLang/en.json --full
```

```bash
./treelang.bin diff --master treeLang/en.json --slave treeLang/pl.json
```

```bash
./treelang.bin diff --master treeLang/en.json --slave treeLang/pl.json --full
```

---

```bash
./treelang.bin generate -h
```

```bash
./treelang.bin generate --from treeLang --to TEMP --json --master TEMP/treelang_en.gen.yml
```

```bash
./treelang.bin generate --from treeLang/en.json --to TEMP --json
```

```bash
./treelang.bin generate --from treeLang --to TEMP --json --map
```

```bash
./treelang.bin generate --from treeLang/en.json --to TEMP --yml
```

```bash
./treelang.bin generate --from treeLang --to TEMP --yml --map
```

```bash
./treelang.bin generate --from treeLang/en.json --to TEMP --go-package main
```

```bash
./treelang.bin generate --from treeLang/en.json --to TEMP --go-package main --func-png
```

```bash
./treelang.bin generate --from treeLang --to TEMP --go-package main --map
```

---

---

### Mirrors

- https://git.bookshelf-writer.fun/Bookshelf-Writer/TreeOfLanguages
 