package binding

type trustedString string

func Bless(s string) trustedString {
	return trustedString(s)
}
