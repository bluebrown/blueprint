package types

import "strings"

type StringSlice []string

func (i *StringSlice) String() string {
	return strings.Join([]string(*i), ",")
}

func (i *StringSlice) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *StringSlice) Type() string {
	return "list"
}
