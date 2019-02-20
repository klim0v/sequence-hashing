package main

import (
	"fmt"
	"github.com/klim0v/sequence-hashing/pkg/entity"
	"github.com/klim0v/sequence-hashing/pkg/redis"
	"golang.org/x/sync/errgroup"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("Не верное количество агрументов: 1.Исходное число (минимум 6 знаков). 2.Число последовательностей для генерации следующих чисел.")
		os.Exit(1)
	}

	number, err := strconv.Atoi(args[0])
	if err != nil || number < 99999 {
		fmt.Println("Исходное число (минимум 6 знаков)")
		os.Exit(1)
	}

	count, err := strconv.Atoi(args[1])
	if err != nil || count < 0 || count > entity.MaxCount {
		fmt.Println("Число последовательностей для генерации следующих чисел.")
		os.Exit(1)
	}

	var g errgroup.Group
	for i := 0; i <= count; i++ {
		g.Go(func() error {
			return redis.NewClient().Push(entity.NewResult(number))
		})
	}
	if err := g.Wait(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
