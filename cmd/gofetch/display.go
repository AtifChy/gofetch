package main

type Display struct {
	Width       int32
	Height      int32
	RefreshRate int32
	IsPrimary   bool
}

type rect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}
