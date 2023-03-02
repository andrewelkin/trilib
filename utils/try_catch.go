package utils

import "fmt"

type TryBlock struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

type Exception interface{}

func Throw(up Exception) {
	panic(up)
}

func Throwf(format string, a ...interface{}) {
	panic(fmt.Sprintf(format, a...))
}

func (tcf TryBlock) Do() {
	if tcf.Finally != nil {

		defer tcf.Finally()
	}
	if tcf.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				tcf.Catch(r)
			}
		}()
	}
	tcf.Try()
}

/*
How to:

	utils.TryBlock{
		Try: func() {
			utils.Throw("Oops...")
		},
		Catch: func(e utils.Exception) {
			fmt.Printf("Caught %v\n", e)
		},
		Finally: func() {
			fmt.Println("Finally...")
		},
	}.Do()


*/
