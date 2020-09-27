/**
 * 针对timer.go做单元测试，主要测试定时器相关接口 依赖模块delayFunc.go
 */
package ztimer

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// 定义一个超时函数
func myFunc(v ...interface{}) {
	fmt.Printf("No.%d function calld. delay %d second(s)\n", v[0].(int), v[1].(int))
}

func TestTimer(t *testing.T) {
	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(12*time.Second))
	defer cancel()

	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func(i int) {
			NewTimerAfter(NewDelayFunc(myFunc, []interface{}{i, 2 * i}), time.Duration(2*i)*time.Second).Run(ctx, &wg)
		}(i)
	}

	wg.Wait()
}
