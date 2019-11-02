package main

import (
    "fmt"
    "github.com/timmyyuan/go-gitdiff/gitdiff"
    "log"
    "os"
)

func main() {
    patch, err := os.Open("changes.patch")
    if err != nil {
        log.Fatal(err)
    }

    _, ps, err := gitdiff.ParsePatch(patch)
    if err != nil {
        log.Fatal(err)
    }
    for k, v := range ps.ShaToFiles {
        log.Printf("%s : \n", k)
        for _, f := range v {
            for _, frag := range f.TextFragments {
                funcnames := frag.FuncNames()
                for _, name := range funcnames {
                    fmt.Printf("%s, ", name)
                }
            }
        }
    }
}