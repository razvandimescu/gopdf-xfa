# gopdf-xfa

Translate the labels in an XFA-dynamic PDF form by rewriting its template stream.

XFA forms carry their entire layout — fields, captions, scripts — as an XML
packet inside the PDF rather than as page content. Translating one therefore
means editing that XML and putting it back without disturbing anything else the
document references. This tool does that in two passes, so the translation
itself stays a plain text file you can edit, diff, and keep under version
control.

```sh
gopdf-xfa dump form.pdf labels.tsv     # extract translatable labels
$EDITOR labels.tsv                     # fill in the en column
gopdf-xfa apply form.pdf labels.tsv out.pdf
```

`dump` is safe to re-run: translations already filled in are carried over, and
only labels new to the form are appended. If the existing file cannot be read it
refuses rather than overwriting, since the output path is the same file.

## Installation

```sh
go install github.com/razvandimescu/gopdf-xfa@latest
```

## How it works

The template packet is located through the catalog's `AcroForm/XFA` array, and
the entry is verified to carry the XFA template namespace before it is trusted —
some generators list an empty `template` while the real content sits in
`xdp:xdp`. When the AcroForm linkage is missing or unusable, every indirect
object is scanned for the namespace marker instead, in ascending object number
so the result is reproducible.

Writing back uses [`gopdf`](https://github.com/razvandimescu/gopdf)'s
`Reader.Rewrite`, which clones the document and substitutes named stream objects
in place. Everything outside the template — fonts, images, the structure tree,
signatures — is copied unchanged.

## Scope

The label extraction and the built-in correction rules were developed against
Romanian ANAF fiscal forms. The mechanism is generic to any XFA template; the
rule table in `apply.go` is not, and is the first thing to adjust for another
form family.

## License

MIT
