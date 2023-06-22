package faltservice

import "github.com/kylemcc/parse"

type ProcessingFileParse struct {
	parse.Base
	Status    int16
	FpsFileID int
}
