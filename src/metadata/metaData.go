package metadata

// todo RelPos is the only actual interesting pos for after the lexing phase?

type MetaData struct {
	Source            string
	Pos, RelPos, Line int
}
