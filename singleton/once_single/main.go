package main


// sync.Doc实现单例

import (
	"fmt"
	"sync"
)


type Singleton struct{}

var singleton *Singleton

var once sync.Once



func GetSingletonObj() *Singleton{
	once.Do(func() {
		fmt.Println("create obj")
		singleton = new(Singleton)
	})

	return singleton
}

func main() {
	var s *Singleton
	s = GetSingletonObj()
	s = GetSingletonObj()
	fmt.Println(s)
}