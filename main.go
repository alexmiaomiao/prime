package main

import (
	"log"
	"os"
	"sync"
)

func raw(stop <-chan struct{}) chan int {
	ch := make(chan int)
	go func() {
		for i := 2; ; i++ {
			select {
			case ch <- i:

			case <-stop:
				return
			}
		}
	}()
	return ch
}

func filter(in <-chan int, prime int, stop <-chan struct{}) chan int {
	out := make(chan int)
	go func() {
		for {
			select {
			case i := <-in:
				if i%prime != 0 {
					out <- i
				}
			case <-stop:
				return
			}
		}
	}()
	return out
}

//PrimeGen 按顺序生成n个质数
func PrimeGen(n int) []int {
	primes := make([]int, n)
	stop := make(chan struct{})
	//PrimeGen返回后让生成质数的所有协程退出
	defer close(stop)

	for ch, i := raw(stop), 0; i < n; i++ {
		primes[i] = <-ch
		ch = filter(ch, primes[i], stop)
	}

	return primes
}

func twoSum(primes []int, target int) [2]int {
	m := make(map[int]bool, len(primes))
	for _, prime := range primes {
		m[prime] = true
		if m[target-prime] {
			return [2]int{prime, target - prime}
		}
	}
	return [2]int{0, 0}
}

//Goldbach 验证n个质数之间所有偶数的哥德巴赫猜想
func Goldbach(primes []int) {
	logger := log.New(os.Stdout, "", 0)
	n := len(primes) - 1
	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()

			for j := primes[i] + 1; j < primes[i+1]; j += 2 {
				a := twoSum(primes[:i+1], j)
				if a[0] != 0 {
					//使用logger输出保证线程安全
					logger.Printf("%d + %d = %d\n", a[0], a[1], j)
				} else {
					logger.Fatal(j)
				}
			}
		}(i)
	}

	wg.Wait()
}

func main() {
	Goldbach(PrimeGen(10000))
}
