package utils

func Assert(condition bool) {
	if !condition {

		panic("Assert Failed")
	}
}
