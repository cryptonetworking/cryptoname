package handy

//NoCopy warns in go vet to struct be copied
type NoCopy Empty

//Lock does nothing
//go:noinline
func (*NoCopy) Lock() {}

//Unlock does nothing
//go:noinline
func (*NoCopy) Unlock() {}
