// generated by stringer -type=Type; DO NOT EDIT

package chat

import "fmt"

const _Type_name = "InvalidTextTranslateScoreSelector"

var _Type_index = [...]uint8{0, 7, 11, 20, 25, 33}

func (i Type) String() string {
	if i < 0 || i+1 >= Type(len(_Type_index)) {
		return fmt.Sprintf("Type(%d)", i)
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
