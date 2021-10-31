bibfuse
=======
[![Test](https://github.com/iomz/bibfuse/actions/workflows/test.yml/badge.svg)](https://github.com/iomz/bibfuse/actions/workflows/test.yml)
[![Docker](https://github.com/iomz/bibfuse/actions/workflows/docker.yml/badge.svg)](https://github.com/iomz/bibfuse/actions/workflows/docker.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/iomz/bibfuse)](https://goreportcard.com/report/github.com/iomz/bibfuse)
[![codecov](https://codecov.io/gh/iomz/bibfuse/branch/main/graph/badge.svg?token=fN1tyc6ssX)](https://codecov.io/gh/iomz/bibfuse)
[![License](https://img.shields.io/github/license/iomz/bibfuse.svg)](https://github.com/iomz/bibfuse/blob/main/LICENSE)
[![GoDoc](https://godoc.org/github.com/iomz/bibfuse?status.svg)](https://godoc.org/github.com/iomz/bibfuse)

A CLI tool to manage bibtex entries using [nickng/bibtex](https://github.com/nickng/bibtex).

bibfuse creates an SQLite database file (`--db`) from given BibTex files (`*.bib`), and generates a single *clean* `.bib` file (`--out`).

The filtering formats can be defined in the config file (`--config`). bibfuse takes the `bibfuse.toml` in this package by default.

If no `.bib` files are given, it just reads the database and updates the BibTex file.

# Table of Contents

* [Synopsis](#synopsis)
  * [Usage](#usage)
  * [Usage with Docker](#docker)
* [bibfuse filters for BibTex format](#filters)
  * [`todos` and `optionals` filters](#todo-optional)
  * [`oneof_` filters with `-smart`](#oneof)
  * [Citation Types](#cite-type)
    * [@article](#article)
    * [@book](#book)
    * [@incollection](#incollection)
    * [@inproceedings](#inproceedings)
    * [@mastersthesis](#mastersthesis)
    * [@misc](#misc)
    * [@phdthesis](#phdthesis)
    * [@techreport](#techreport)
    * [@unpublished](#unpublished)
* [Contribution](#contribution)
* [License](#license)
* [Author](#author)

# Synopsis <a name="synopsis"/>

## Install <a name="install"/>
```console
% go get -u github.com/iomz/bibfuse/...
```

## Usage <a name="usage"/>

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
        Do not hide empty fields in the resulting bibtex.
  -smart
        Use oneof selectively filters when importing bibtex.
  -verbose
        Print verbose messages.
  -version
        Print version.
```

### Example
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
    metanote    = "(OPTIONAL)",
    number      = "(OPTIONAL)",
    numpages    = "(OPTIONAL)",
    pages       = "(OPTIONAL)",
    publisher   = "(OPTIONAL)",
    volume      = "(OPTIONAL)",
    year        = "(TODO)",
}
```

## Usage with Docker <a name="docker"/>
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

# bibfuse filters for BibTex format <a name="filters"/>
bibfuse reflects rather subjective opinion to filter and flag the required fields depending on the type, aiming for the compatibility with most of the research publication requirements.

## `todos` and `optionals` filters <a name="todo-optional"/>

bibfuse filters fields for each entry depending on the type: `@article`, `@book`, `@incollection`, `@inproceedings`, `@mastersthesis`, `@misc`, `@phdthesis`, `@techreport`, and `@unpublished` as defined in `bibfuse.toml`.

Mandatory fields are filled with `(TODO)` while optional fileds are filled with `(OPTIONAL)`.

## `oneof_` filters with `-smart` <a name="oneof"/>

In addition, you can define additional filters whose name starting with `oneof_` to selectively _discard_ some fields in presence of a specific fields with the `-smart` option. For example, by default for the `@article` type the default `bibfuse.toml` has the following `oneof_` filter:

```toml
oneof_doi_page = [
    "doi",
    "pages",
    "numpages"
]
```

This checks each field in the order "doi", "pages", and "numpages" then after finding the first one of neither empty, "(OPTIONAL)", nor "(TODO)", bibfuse discards the remaining fields and does _NOT_ store them in the database. Note that if you want to retain those fields, you should not use this `-smart` option.

This feature enables rather concise bibliography in your manuscript while maintaining the accessibility to the cited documents through more efficient identities (e.g., DOI).

## Citation Types <a name="cite-type"/>

### Journal articles <a name="article"/>
```latex
@article{mizutani2021article
    title     = {{Title of the Article}},
    author    = "(TODO)",
    journal   = "(TODO)",
    year      = "(TODO)",
    doi       = "(OPTIONAL)",
    isbn      = "(OPTIONAL)", % removed if doi exists
    issn      = "(OPTIONAL)", % removed if doi exists
    metanote  = "(OPTIONAL)",
    number    = "(OPTIONAL)", % removed if doi exists
    numpages  = "(OPTIONAL)", % removed if doi exists
    pages     = "(OPTIONAL)", % removed if doi exists
    publisher = "(OPTIONAL)", % removed if doi exists
    url       = "(OPTIONAL)", % removed if doi exists
    volume    = "(OPTIONAL)", % removed if doi exists
}
```

### Books <a name="book"/>
```latex
@book{mizutani2021book,
    title     = {{Title of the Book}},
    author    = "(TODO)"
    publisher = "(TODO)",     % removed if doi exists
    year      = "(TODO)",
    doi       = "(OPTIONAL)",
    edition   = "(OPTIONAL)", % removed if doi exists
    isbn      = "(OPTIONAL)", % removed if doi exists
    issn      = "(OPTIONAL)", % removed if doi exists
    metanote  = "(OPTIONAL)",
    url       = "(OPTIONAL)", % removed if doi exists
}
```

### Chapters or articles in a book <a name="incollection"/>
```latex
@incollection{mizutani2012incollection,
    title     = {{Title of the Book Chapter}},
    author    = "(TODO)"
    booktitle = "(TODO)",
    publisher = "(TODO)",     % removed if doi exists
    year      = "(TODO)",
    doi       = "(OPTIONAL)",
    isbn      = "(OPTIONAL)", % removed if doi exists
    issn      = "(OPTIONAL)", % removed if doi exists
    metanote  = "(OPTIONAL)",
    numpages  = "(OPTIONAL)", % removed if doi exists
    pages     = "(OPTIONAL)", % removed if doi exists
    series    = "(OPTIONAL)",
    url       = "(OPTIONAL)", % removed if doi exists
}
```

### Conference papers, lecture notes, extended abstract, etc. <a name="inproceedings"/>
```latex
@inproceedings{mizutani2012inproceedings,
    title     = {{Title of the Conference Paper}},
    author    = "(TODO)"
    booktitle = "(TODO)",
    year      = "(TODO)",
    doi       = "(OPTIONAL)",
    isbn      = "(OPTIONAL)", % removed if doi exists
    issn      = "(OPTIONAL)", % removed if doi exists
    metanote  = "(OPTIONAL)",
    numpages  = "(OPTIONAL)", % removed if doi exists
    pages     = "(OPTIONAL)", % removed if doi exists
    publisher = "(OPTIONAL)", % removed if doi exists
    series    = "(OPTIONAL)",
    url       = "(OPTIONAL)", % removed if doi exists
}
```

### Master's theses <a name="mastersthesis"/>
```latex
@mastersthesis{mizutani2021mastersthesis,
    title       = {{Title of the Master's Thesis}},
    author      = "(TODO)",
    url         = "(OPTIONAL)",
    metanote    = "(OPTIONAL)",
    school      = "(TODO)",
    year        = "(TODO)",
}

```

### Online resources, artifacts, etc. <a name="misc"/>
```latex
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

### Ph.D. theses / dissertations <a name="phdthesis"/>
```latex
@phdthesis{mizutani2021phdthesis,
    title       = {{Title of the Ph.D. Thesis}},
    author      = "(TODO)",
    url         = "(OPTIONAL)",
    metanote    = "(OPTIONAL)",
    school      = "(TODO)",
    year        = "(TODO)",
}
```

### Standards, specifications, white papers, etc. <a name="techreport"/>
```latex
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

### Documents not formally published. <a name="unpublished"/>
```latex
@unpublished{mizutani2021unpublished,
    title       = {{Title of the Unpublished Work}},
    author      = "(TODO)",
    url         = "(TODO)",
    metanote    = "(OPTIONAL)",
    note        = "(TODO)",
}
```

# Contribution <a name="contribution"/>
See `CONTRIBUTING.md`.

# License <a name="license"/>
See `LICENSE`.

# Author <a name="author"/>
Iori Mizutani ([@iomz](https://github.com/iomz))
