package safe

func Bless(level TrustLevel, s string) String {
	switch level {
	case Untrusted:
		return UntrustedString(s)
	case HTMLSafe:
		return HTML{s}
	case TextSafe:
		return Text{s}
	case URLSafe:
		return URL{s}
	case AttributeSafe:
		return Attribute{s}
	case FullyTrusted:
		return constantString(s)
	default:
		panic("invalid TrustLevel")
	}
}
