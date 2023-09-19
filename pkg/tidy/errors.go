package tidy

import "fmt"

type SortingError struct {
	Filename string
	AbsPath  string
	Sort     bool // indicates whether the error occured in a sort or unsort operation
	Err      error
}

func (se *SortingError) Error() string {
	return fmt.Sprintf("Sorting Error: Could not move file to desired destination.\n\tFile:\t[%s]\n\tDest:\t[%s]\n\tOperation:\t[%s]\n\n\tError:\t%s",
		se.Filename, se.AbsPath, se.sortOperation(), se.Err.Error())
}

func (se *SortingError) sortOperation() string {
	if se.Sort {
		return "Sort"
	}
	return "Unsort"
}
