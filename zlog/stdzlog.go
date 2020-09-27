package zlog

/*
 * 全局默认提供一个Log对外句柄，可以直接使用API系列调用
 * 全局日志对象 StdZinxLog
 */

import "os"

// StdZinxLog 全局日志对象
var StdZinxLog = NewZinxLog(os.Stderr, "", BitDefault)

func init() {
	// 因为StdZinxLog对象 对所有输出方法做了一层包裹，所以在打印调用函数的时候，
	// 比正常的logger对象多一层调用
	// 一般的zinxLogger对象 calldDepth=2, StdZinxLog的calldDepth=3
	StdZinxLog.calldDepth = 3
}

// Flags 获取StdZinxLog 标记位
func Flags() int {
	return StdZinxLog.Flags()
}

// ResetFlags 设置StdZinxLog标记位
func ResetFlags(flag int) {
	StdZinxLog.ResetFlags(flag)
}

// AddFlag 添加flag标记
func AddFlag(flag int) {
	StdZinxLog.AddFlag(flag)
}

// SetPrefix 设置StdZinxLog 日志头前缀
func SetPrefix(prefix string) {
	StdZinxLog.SetPrefix(prefix)
}

// SetLogFile 设置StdZinxLog绑定的日志文件
func SetLogFile(fileDir string, fileName string) {
	StdZinxLog.SetLogFile(fileDir, fileName)
}

// CloseDebug 设置关闭debug
func CloseDebug() {
	StdZinxLog.CloseDebug()
}

// OpenDebug 设置打开debug
func OpenDebug() {
	StdZinxLog.OpenDebug()
}

// SetLevel 设置日志级别
func SetLevel(level int) {
	StdZinxLog.SetLevel(level)
}

// Debugf 同 Debug
func Debugf(format string, v ...interface{}) {
	StdZinxLog.Debugf(format, v...)
}

// Debug 指出细粒度信息事件对调试应用程序是非常有帮助的，主要用于开发过程中打印一些运行信息
func Debug(v ...interface{}) {
	StdZinxLog.Debug(v...)
}

// Infof 同 Info
func Infof(format string, v ...interface{}) {
	StdZinxLog.Infof(format, v...)
}

// Info 消息在粗粒度级别上突出强调应用程序的运行过程。
// 打印一些你感兴趣的或者重要的信息，这个可以用于生产环境中输出程序运行的一些重要信息，
// 但是不能滥用，避免打印过多的日志
func Info(v ...interface{}) {
	StdZinxLog.Info(v...)
}

// Warnf 同 Warn
func Warnf(format string, v ...interface{}) {
	StdZinxLog.Warnf(format, v...)
}

// Warn 表明会出现潜在错误的情形，有些信息不是错误信息，但是也要给程序员的一些提示
func Warn(v ...interface{}) {
	StdZinxLog.Warn(v...)
}

// Errorf 同 Error
func Errorf(format string, v ...interface{}) {
	StdZinxLog.Errorf(format, v...)
}

// Error 指出虽然发生错误事件，但仍然不影响系统的继续运行。
// 打印错误和异常信息，如果不想输出太多的日志，可以使用这个级别
func Error(v ...interface{}) {
	StdZinxLog.Error(v...)
}

// Fatalf 同 Fatal
func Fatalf(format string, v ...interface{}) {
	StdZinxLog.Fatalf(format, v...)
}

// Fatal 指出每个严重的错误事件将会导致应用程序的退出。这个级别比较高了。
// 重大错误，这种级别需要终止程序
func Fatal(v ...interface{}) {
	StdZinxLog.Fatal(v...)
}

// Panicf 同 Panic
func Panicf(format string, v ...interface{}) {
	StdZinxLog.Panicf(format, v...)
}

// Panic 这种级别需要终止程序
func Panic(v ...interface{}) {
	StdZinxLog.Panic(v...)
}

// Stack 打印 堆栈日志
func Stack(v ...interface{}) {
	StdZinxLog.Stack(v...)
}
