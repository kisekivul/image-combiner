/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"errors"
	"log"
)

// import (
// 	"image-combiner/cmd"
// )

// func main() {
// 	cmd.Execute()
// }

func main() {
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	err = errors.New("abc")
}
