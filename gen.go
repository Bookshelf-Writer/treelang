package treelang

//go:generate bash -c "rm -rf target/*"
//go:generate bash -c "rm -rf tmp/*"

//go:generate bash "./_run/scripts/creator_const_Go.sh"
//go:generate go run ./_generate/dependencies
//go:generate go run ./_generate/methods

const ()
