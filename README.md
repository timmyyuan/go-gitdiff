# go-gitdiff

[![GoDoc](https://godoc.org/github.com/bluekeyes/go-gitdiff/gitdiff?status.svg)](http://godoc.org/github.com/bluekeyes/go-gitdiff/gitdiff) [![Go Report Card](https://goreportcard.com/badge/github.com/bluekeyes/go-gitdiff)](https://goreportcard.com/report/github.com/bluekeyes/go-gitdiff)

A Go library for parsing and applying patches generated by `git diff`, `git
show`, and `git format-patch`. It can also parse and apply unified diffs
generated by the standard `diff` tool.

It supports both standard line-oriented text patches and Git binary patches.

```golang
patch, err := os.Open("changes.patch")
if err != nil {
    log.Fatalf(err)
}

files, preamble, err := gitdiff.Parse(patch)
if err != nil {
    log.Fatalf(err)
}

// files is a slice of *gitdiff.File describing the files changed in the patch
// preamble is a string of the content of the patch before the first file
```

## Status

In development, expect API changes. Patch parsing works, but has not been
tested extensively against real-world patches. Patch application has not been
implemented yet.

## Why another git/unified diff parser?

[Several][sourcegraph] [packages][sergi] with [similar][waigani]
[functionality][seletskiy] exist, so why did I write another?

1. No other packages I found support binary diffs, as generated with the
   `--binary` flag. This is the main reason for writing a new packages, as the
   format is pretty different from line-oriented diffs and is unique to Git.

2. Most other packages only parse patches, so you need another package to apply
   them (and if they do support applies, it is only for text files.)

3. This package aims to accept anything that `git apply` accepts, and closely
   follows the logic in [`apply.c`][apply.c].

4. It seemed like a fun thing to write and a good way to learn more about Git.

[sourcegraph]: https://github.com/sourcegraph/go-diff
[sergi]: https://github.com/sergi/go-diff
[waigani]: https://github.com/waigani/diffparser
[seletskiy]: https://github.com/seletskiy/godiff

[apply.c]: https://github.com/git/git/blob/master/apply.c

## Differences From Git

1. Certain types of invalid input that are accepted by `git apply` generate
   errors. These include:

   - Numbers immediately followed by non-numeric characters
   - Trailing characters on a line after valid or expected content

2. Errors for invalid input are generally more verbose and specific than those
   from `git apply`.

3. The translation from C to Go may have introduced inconsistencies in the way
   Unicode file names are handled; these are bugs, so please report any issues
   of this type.

4. When reading headers, there is no validation that OIDs present on an `index`
   line are shorter than or equal to the maximum hash length, as this requires
   knowing if the repository used SHA1 or SHA256 hashes.

5. When reading "traditional" patches (those not produced by `git`), prefixes
   are not stripped from file names; `git apply` attempts to remove prefixes
   that match the current repository directory/prefix.
