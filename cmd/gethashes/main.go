package main

import (
	"fmt"
	"github.com/klim0v/sequence-hashing/pkg/entity"
	"github.com/klim0v/sequence-hashing/pkg/redis"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"strconv"
)

var (
	argCountErrorMessage = "не верное количество агрументов"
	arg1ErrorMessage     = "исходное число (минимум 6 знаков), %s"
	arg2ErrorMessage     = "число последовательностей для генерации следующих чисел, %s"
)

func main() {
	number, count, err := Parse(os.Args[1:])
	if err != nil {
		log.Println(err)
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

func Parse(args []string) (number uint64, count int, err error) {
	if len(args) != 2 {
		err = fmt.Errorf(argCountErrorMessage)
		return
	}
	number, err = strconv.ParseUint(args[0], 10, 64)
	if err != nil || number < 99999 {
		err = fmt.Errorf(arg1ErrorMessage, args[0])
		return
	}
	count, err = strconv.Atoi(args[1])
	if err != nil || count < 0 || count > entity.MaxCount {
		err = fmt.Errorf(arg2ErrorMessage, args[1])
		return
	}
	return
}
