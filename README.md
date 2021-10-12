# bibfuse
A CLI tool to manage bibtex entries using [nickng/bibtex](https://github.com/nickng/bibtex).

## Synopsis
This tool takes `.bib` files and filter fields for each entry depending on the type: article, book, inproceedings, misc, and techreport.
The mandatory fields are filled with `(TODO)` and optional fileds are filled with `(OPTIONAL)` by default.

```console
% go get -u github.com/iomz/bibfuse
% cat ref.bib
@article{someone2021a,
    title     = {{A Journal Article}},
}
% bibfuse -in ref.bib
@article{someone2021a,
    title     = {{A Journal Article}},
    author    = "(TODO)",
    journal   = "(TODO)",
    year      = "(TODO)",
    url       = "(OPTIONAL)",
    doi       = "(OPTIONAL)",
    isbn      = "(OPTIONAL)",
    issn      = "(OPTIONAL)",
    keyword   = "(OPTIONAL)",
    metanote  = "(OPTIONAL)",
    number    = "(OPTIONAL)",
    numpages  = "(OPTIONAL)",
    pages     = "(OPTIONAL)",
    publisher = "(OPTIONAL)",
    volume    = "(OPTIONAL)",
}
```

### Journal articles
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

### Books
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

### Chapters or articles in a book
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

### Conference papers, lecture notes, extended abstract, etc.
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

### Online resources, artifacts, etc.
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

### Standards, specifications, white papers, etc.
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
