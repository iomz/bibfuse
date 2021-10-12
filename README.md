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
    url       = "(OPTIONAL)",
    doi       = "(OPTIONAL)",
    isbn      = "(OPTIONAL)",
    issn      = "(OPTIONAL)",
    journal   = "(TODO)",
    keyword   = "(OPTIONAL)",
    metanote  = "(OPTIONAL)",
    number    = "(TODO)",
    numpages  = "(OPTIONAL)",
    pages     = "(OPTIONAL)",
    publisher = "(TODO)",
    volume    = "(TODO)",
    year      = "(TODO)",
}
```

### Journal articles
```
@article{mizutani2021article 
    title     = {{Title of the Article}},
    author    = {Mizutani, Iori},
    doi       = "(OPTIONAL)",
    isbn      = "(OPTIONAL)",
    issn      = "(OPTIONAL)",
    journal   = "A cool journal",
    keyword   = "(OPTIONAL)",
    metanote  = "(OPTIONAL)",
    number    = 1,
    numpages  = "(OPTIONAL)",
    pages     = "(OPTIONAL)",
    publisher = "A cool publisher",
    url       = "(OPTIONAL)",
    volume    = 1,
    year      = 2021,
}
```

### Books
```
@book{mizutani2021book,
    title     = {{Title of the Book}},
    author    = "Mizutani, Iori"
    url       = "(OPTIONAL)",
    address   = "London, United Kingdom",
    doi       = "(OPTIONAL)",
    edition   = "(TODO)",
    isbn      = "(OPTIONAL)",
    issn      = "(OPTIONAL)",
    metanote  = "(OPTIONAL)",
    publisher = "Pearson Education",
    year      = 2021,
}
```

### Conference papers, lecture notes, extended abstract, etc.
```
@inproceedings{mizutani2012inproceedings,
    title     = {{Title of the Conference Paper}},
    author    = "Mizutani, Iori"
    url       = "(OPTIONAL)",
    booktitle = "Proceedings of the Cool Conference 2021",
    doi       = "(OPTIONAL)",
    isbn      = "(OPTIONAL)",
    issn      = "(OPTIONAL)",
    keyword   = "(OPTIONAL)",
    location  = "(OPTIONAL)",
    metanote  = "(OPTIONAL)",
    numpages  = "(OPTIONAL)",
    pages     = "(OPTIONAL)",
    publisher = "(TODO)",
    series    = "(OPTIONAL)",
    year      = 2021,
}
```

### Online resources, artifacts, etc.
```
@misc{mizutani2021misc,
    title       = "Title of the Resource",
    author      = "Mizutani, Iori"
    url         = "(TODO)",
    institution = "(TODO)",
    metanote    = "(OPTIONAL)",
    month       = "(OPTIONAL)",
    note        = "(TODO)",
    year        = 2021,
}
```

### Standards, specifications, white papers, etc.
```
@techreport{mizutani2021techreport,
    title       = {{Title of the Technical Document}},
    author      = {{Mizutani, Iori}},
    url         = "(TODO)",
    institution = {{TODO}},
    metanote    = "(OPTIONAL)",
    month       = "(OPTIONAL)",
    series      = "(OPTIONAL)",
    version     = "(OPTIONAL)",
    year        = 2021,
}
```
