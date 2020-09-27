package ztimer

import (
	"context"
	"sync"
	"time"

	"github.com/YungMonk/zinx/zlog"
)

const (
	// HourName 小时级时间轮名称
	HourName = "HOUR"
	// HourInterval 小时级刻度之间的duration时间间隔
	HourInterval = 60 * 60 * 1e3 // ms为精度（即 1小时）
	// HourScales 小时级时间轮的轮盘一共多少个刻度(如我们正常的时钟就是12个刻度)
	HourScales = 12

	// MinuteName 分钟级时间轮名称
	MinuteName = "MINUTE"
	// MinuteInterval 分钟级刻度之间的duration时间间隔
	MinuteInterval = 60 * 1e3 // ms为精度 （即 1分钟）
	// MinuteScales 分钟级时间轮的轮盘一共多少个刻度(如我们正常的时钟就是12个刻度)
	MinuteScales = 60

	// SecondName 秒级时间轮名称
	SecondName = "SECOND"
	// SecondInterval 秒级刻度之间的duration时间间隔
	SecondInterval = 1e3 // ms为精度 （即 1秒）
	// SecondScales 秒级时间轮的轮盘一共多少个刻度(如我们正常的时钟就是12个刻度)
	SecondScales = 60

	// TimersMaxCap 每个时间轮刻度挂载定时器的最大个数
	TimersMaxCap = 2048
)

/**
 * 注意：
 *  有关时间的几个换算
 *  time.Second(秒) = time.Millisecond * 1e3
 *	time.Millisecond(毫秒) = time.Microsecond * 1e3
 *	time.Microsecond(微秒) = time.Nanosecond * 1e3
 *
 *	time.Now().UnixNano() ==> time.Nanosecond (纳秒)
 */

// Timer 定时器实现
type Timer struct {
	// 延迟调用函数
	delayFunc *DelayFunc
	// 调用时间(unix 时间， 单位ms)
	unixts int64
}

// NewTimerAt 创建一个定时器，在指定的时间触发 定时器方法
// 类似 Redis 中的 expireat
//     df: DelayFunc类型的延迟调用函数类型    unixNano: unix计算机从1970-1-1至今经历的纳秒数
func NewTimerAt(df *DelayFunc, unixNano int64) *Timer {
	return &Timer{
		// 要延迟的函数
		delayFunc: df,
		// 将纳秒转换成对应的毫秒 ms ，定时器以ms为最小精度
		unixts: unixNano / 1e6,
	}
}

// NewTimerAfter 创建一个定时器，在当前时间延迟 duration 之后触发 定时器方法
// 类似 Redis 中的 expire
//     df: DelayFunc类型的延迟调用函数类型    duration: 多少纳秒后触发
func NewTimerAfter(df *DelayFunc, duration time.Duration) *Timer {
	return NewTimerAt(df, time.Now().UnixNano()+int64(duration))
}

// UnixMilli 返回1970-1-1至今经历的毫秒数
func UnixMilli() int64 {
	return time.Now().UnixNano() / 1e6
}

// Run 启动定时器，用一个go承载
func (t *Timer) Run(ctx context.Context, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()

		now := UnixMilli()
		// 设置的定时器是否在当前时间之后
		if t.unixts > now {
			select {
			// 是否取消当前过程
			case <-ctx.Done():
				zlog.Info("timer is cancel")
			// 睡眠，直至时间超时，已微秒为单位进行睡眠
			case <-time.After(time.Duration(t.unixts-now) * time.Millisecond):
				// 调用事先注册好的超时延迟方法
				t.delayFunc.Call()
			}
		} else {
			// 定时器在当前时间之前，立即执行
			t.delayFunc.Call()
		}
	}()
}
