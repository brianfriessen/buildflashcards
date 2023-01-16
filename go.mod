module example.com/main2

go 1.17

replace example.com/bingImageSearch => ../bingsearch

//replace example.com/forvosearch => ../forvosearch
replace example.com/forvosearch => ../forvosearch

require (
	example.com/bingImageSearch v0.0.0-00010101000000-000000000000
	example.com/forvosearch v0.0.0-00010101000000-000000000000
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
)
