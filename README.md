bibfuse
=======
[![Test](https://github.com/iomz/bibfuse/actions/workflows/test.yml/badge.svg)](https://github.com/iomz/bibfuse/actions/workflows/test.yml)
[![Docker](https://github.com/iomz/bibfuse/actions/workflows/docker.yml/badge.svg)](https://github.com/iomz/bibfuse/actions/workflows/docker.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/iomz/bibfuse)](https://goreportcard.com/report/github.com/iomz/bibfuse)
[![codecov](https://codecov.io/gh/iomz/bibfuse/branch/main/graph/badge.svg?token=fN1tyc6ssX)](https://codecov.io/gh/iomz/bibfuse)

A CLI tool to manage bibtex entries using [nickng/bibtex](https://github.com/nickng/bibtex).

bibfuse creates an SQLite database file (`--db`) from given BibTex files (`*.bib`), and generates a single *clean* `.bib` file (`--out`).

The filtering formats can be defined in the config file (`--config`). bibfuse takes the `bibfuse.toml` in this package by default.

If no `.bib` files are given, it just reads the database and update the BibTex file.

```console
% bibfuse -h
Usage of bibfuse: [options] [.bib ... .bib]
  -config string
        The bibfuse.[toml|yml] defining the filters. (default "bibfuse.toml")
  -db string
        The SQLite file to read/write. (default "bib.db")
  -no-optional
        Suppress "OPTIONAL" fields in the resulting bibtex.
  -no-todo
        Suppress "TODO" fields in the resulting bibtex.
  -out string
        The resulting bibtex to write (it overrides if exists). (default "out.bib")
  -show-empty
        Suppress empty fields in the resulting bibtex.
  -verbose
        Print verbose messages.
  -version
        Print version.
```

# Synopsis
This tool takes `.bib` files and filter fields for each entry depending on the type: article, book, inproceedings, misc, and techreport.
The mandatory fields are filled with `(TODO)` and optional fileds are filled with `(OPTIONAL)` by default.

## Install

```console
% go get -u github.com/iomz/bibfuse
```

## Usage

```console
% cat ref.bib
@article{someone2021a,
    title     = {{A Journal Article}},
}
% bibfuse ref.bib
2021/10/17 15:47:32 parsing ref.bib
2021/10/17 15:47:32 +1 new entries
2021/10/17 15:47:32 bib.db contains 1 entries
2021/10/17 15:47:32 1 entries written to out.bib
% cat out.bib
@article{someone2021a,
    title       = {{A Journal Article}},
    author      = "(TODO)",
    url         = "(OPTIONAL)",
    doi         = "(OPTIONAL)",
    isbn        = "(OPTIONAL)",
    issn        = "(OPTIONAL)",
    journal     = "(TODO)",
    keyword     = "(OPTIONAL)",
    metanote    = "(OPTIONAL)",
    number      = "(OPTIONAL)",
    numpages    = "(OPTIONAL)",
    pages       = "(OPTIONAL)",
    publisher   = "(OPTIONAL)",
    volume      = "(OPTIONAL)",
    year        = "(TODO)",
}
```

## Usage with Docker

```console
% cat ref.bib
@article{someone2021a,
    title     = {{A Journal Article}},
}
% docker run -v $(pwd):$(pwd) -w $(pwd) --rm iomz/bibfuse ref.bib
2021/10/17 13:53:33 parsing ref.bib
2021/10/17 13:53:33 +0 new entries
2021/10/17 13:53:33 bib.db contains 1 entries
2021/10/17 13:53:33 1 entries written to out.bib
% sqlite3 bib.db "SELECT * FROM entries;"
1|someone2021a|article|(TODO)|{A Journal Article}||(OPTIONAL)||(OPTIONAL)|(OPTIONAL)||(TODO)|(OPTIONAL)||(OPTIONAL)||(OPTIONAL)|(OPTIONAL)|(OPTIONAL)|(OPTIONAL)||||(OPTIONAL)||(OPTIONAL)|(TODO)
```

# BibTex entry format

bibfuse reflects rather subjective opinion to filter and flag the required fields depending on the type.

Aiming for the compatibility with most of the research publication requirements.

## Journal articles
```
@article{mizutani2021article
    title     = {{Title of the Article}},
    author    = "(TODO)",
    journal   = "(TODO)",
    year      = "(TODO)",
    doi       = "(OPTIONAL)",
    isbn      = "(OPTIONAL)",
    issn      = "(OPTIONAL)",
    keyword   = "(OPTIONAL)",
    metanote  = "(OPTIONAL)",
    number    = "(OPTIONAL)",
    numpages  = "(OPTIONAL)",
    pages     = "(OPTIONAL)",
    publisher = "(OPTIONAL)",
    url       = "(OPTIONAL)",
    volume    = "(OPTIONAL)",
}
```

## Books
```
@book{mizutani2021book,
    title     = {{Title of the Book}},
    author    = "(TODO)"
    publisher = "(TODO)",
    year      = "(TODO)",
    doi       = "(OPTIONAL)",
    edition   = "(OPTIONAL)",
    isbn      = "(OPTIONAL)",
    issn      = "(OPTIONAL)",
    metanote  = "(OPTIONAL)",
    url       = "(OPTIONAL)",
}
```

## Chapters or articles in a book
```
@incollection{mizutani2012incollection,
    title     = {{Title of the Book Chapter}},
    author    = "(TODO)"
    booktitle = "(TODO)",
    publisher = "(TODO)",
    year      = "(TODO)",
    url       = "(OPTIONAL)",
    doi       = "(OPTIONAL)",
    isbn      = "(OPTIONAL)",
    issn      = "(OPTIONAL)",
    keyword   = "(OPTIONAL)",
    location  = "(OPTIONAL)",
    metanote  = "(OPTIONAL)",
    numpages  = "(OPTIONAL)",
    pages     = "(OPTIONAL)",
    series    = "(OPTIONAL)",
}
```

## Conference papers, lecture notes, extended abstract, etc.
```
@inproceedings{mizutani2012inproceedings,
    title     = {{Title of the Conference Paper}},
    author    = "(TODO)"
    booktitle = "(TODO)",
    year      = "(TODO)",
    doi       = "(OPTIONAL)",
    isbn      = "(OPTIONAL)",
    issn      = "(OPTIONAL)",
    keyword   = "(OPTIONAL)",
    location  = "(OPTIONAL)",
    metanote  = "(OPTIONAL)",
    numpages  = "(OPTIONAL)",
    pages     = "(OPTIONAL)",
    publisher = "(OPTIONAL)",
    series    = "(OPTIONAL)",
    url       = "(OPTIONAL)",
}
```

## Master's theses
```
@mastersthesis{mizutani2021mastersthesis,
    title       = {{Title of the Master's Thesis}},
    author      = "(TODO)",
    url         = "(OPTIONAL)",
    metanote    = "(OPTIONAL)",
    school      = "(TODO)",
    year        = "(TODO)",
}

```

## Online resources, artifacts, etc.
```
@misc{mizutani2021misc,
    title       = "Title of the Resource",
    author      = "(TODO)"
    note        = "(TODO)",
    url         = "(TODO)",
    year        = "(TODO)",
    institution = "(OPTIONAL)",
    metanote    = "(OPTIONAL)",
}
```

## Ph.D. theses / dissertations
```
@phdthesis{mizutani2021phdthesis,
    title       = {{Title of the Ph.D. Thesis}},
    author      = "(TODO)",
    url         = "(OPTIONAL)",
    metanote    = "(OPTIONAL)",
    school      = "(TODO)",
    year        = "(TODO)",
}
```

## Standards, specifications, white papers, etc.
```
@techreport{mizutani2021techreport,
    title       = {{Title of the Technical Document}},
    author      = "(TODO)",
    institution = "(TODO)",
    year        = "(TODO)",
    metanote    = "(OPTIONAL)",
    series      = "(OPTIONAL)",
    url         = "(OPTIONAL)",
    version     = "(OPTIONAL)",
}
```

## Documents not formally published.
```
@unpublished{mizutani2021unpublished,
    title       = {{Title of the Unpublished Work}},
    author      = "(TODO)",
    url         = "(TODO)",
    metanote    = "(OPTIONAL)",
    note        = "(TODO)",
}
```

# Contribution

See `CONTRIBUTING.md`.

# License

See `LICENSE`.

# Author

Iori Mizutani ([@iomz](https://github.com/iomz))
