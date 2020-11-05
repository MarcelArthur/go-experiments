package main

import (
	"fmt"
	"sync"
)

// 加锁实现单例模式
var(
	lock *sync.Mutex = &sync.Mutex{}
	instance *Singletons
)

type Singletons struct{}


func GetInstance() *Singletons{
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil{
			instance = &Singletons{}
			fmt.Println("instance")
		}
	}
	return instance
}

func main() {
	var s *Singletons
	s = GetInstance()
	s = GetInstance()
	fmt.Println(s)
}

