package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/OblivionTime/simple-logger/zap"
	"github.com/OblivionTime/simple-logger/zap/zapcore"

	"gopkg.in/natefinch/lumberjack.v2"
)

type myLog struct {
	ZapLog *zap.Logger
	debug  bool
}

var Log *myLog

func (log *myLog) Info(v ...interface{}) {
	if !log.debug {
		return
	}
	message := fmt.Sprintln(v...)
	fields := []zap.Field{}
	if ce := log.ZapLog.CheckCe(zap.InfoLevel, message); ce != nil {
		ce.Write(fields...)
	}
}
func (log *myLog) Infof(format string, v ...interface{}) {
	if !log.debug {
		return
	}
	message := fmt.Sprintf(format, v...)
	fields := []zap.Field{}
	if ce := log.ZapLog.CheckCe(zap.InfoLevel, message); ce != nil {
		ce.Write(fields...)
	}
}
func (log *myLog) Error(v ...interface{}) {
	if !log.debug {
		return
	}
	message := fmt.Sprintln(v...)
	fields := []zap.Field{}
	if ce := log.ZapLog.CheckCe(zap.WarnLevel, message); ce != nil {
		ce.Write(fields...)
	}
}
func (log *myLog) Errorf(format string, v ...interface{}) {
	if !log.debug {
		return
	}
	message := fmt.Sprintf(format, v...)
	fields := []zap.Field{}
	if ce := log.ZapLog.CheckCe(zap.WarnLevel, message); ce != nil {
		ce.Write(fields...)
	}
}

func (log *myLog) Debug(v ...interface{}) {
	if !log.debug {
		return
	}
	message := fmt.Sprintln(v...)
	fields := []zap.Field{}
	if ce := log.ZapLog.CheckCe(zap.DebugLevel, message); ce != nil {
		ce.Write(fields...)
	}
}
func (log *myLog) Debugf(format string, v ...interface{}) {
	if !log.debug {
		return
	}
	message := fmt.Sprintf(format, v...)
	fields := []zap.Field{}
	if ce := log.ZapLog.CheckCe(zap.DebugLevel, message); ce != nil {
		ce.Write(fields...)
	}
}
func (log *myLog) Fatal(v ...interface{}) {
	if !log.debug {
		return
	}
	message := fmt.Sprintln(v...)
	fields := []zap.Field{}
	if ce := log.ZapLog.CheckCe(zap.FatalLevel, message); ce != nil {
		ce.Write(fields...)
	}
	os.Exit(1)
}
func (log *myLog) Fatalf(format string, v ...interface{}) {
	if !log.debug {
		return
	}
	message := fmt.Sprintf(format, v...)
	fields := []zap.Field{}
	if ce := log.ZapLog.CheckCe(zap.FatalLevel, message); ce != nil {
		ce.Write(fields...)
	}
	os.Exit(1)
}
func (log *myLog) Panic(v ...interface{}) {
	if !log.debug {
		return
	}
	message := fmt.Sprintln(v...)
	fields := []zap.Field{}
	if ce := log.ZapLog.CheckCe(zap.PanicLevel, message); ce != nil {
		ce.Write(fields...)
	}
	panic(message)
}
func (log *myLog) Patalf(format string, v ...interface{}) {
	if !log.debug {
		return
	}
	message := fmt.Sprintf(format, v...)
	fields := []zap.Field{}
	if ce := log.ZapLog.CheckCe(zap.PanicLevel, message); ce != nil {
		ce.Write(fields...)
	}
	panic(message)
}

var levelMap = map[string]zapcore.Level{

	"debug": zapcore.DebugLevel,

	"info": zapcore.InfoLevel,

	"warn": zapcore.WarnLevel,

	"error": zapcore.ErrorLevel,

	"dpanic": zapcore.DPanicLevel,

	"panic": zapcore.PanicLevel,

	"fatal": zapcore.FatalLevel,
}

//??????????????????
func setLoggerFile(lev_name string, encoder zapcore.Encoder, LogPath string) zapcore.Core {
	//????????????
	priority := zap.LevelEnablerFunc(func(lev2 zapcore.Level) bool { //info???debug??????,debug??????????????????
		return levelMap[lev_name] == lev2
	})
	if LogPath == "" {
		LogPath = "./log"
	}
	// filename := fmt.Sprintf(`%s/%s-%s-%s/%s.log`, conf.Config.Log.Path, time.Now().Format("2006"), time.Now().Format("01"), time.Now().Format("02"), lev_name)
	filename := fmt.Sprintf(`%s/%s-%s-%s/%s.log`, LogPath, time.Now().Format("2006"), time.Now().Format("01"), time.Now().Format("02"), lev_name)
	//info??????writeSyncer
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filename, //??????????????????????????????????????????????????????????????????
		MaxSize:    10,       //??????????????????,??????MB
		MaxBackups: 100,      //??????????????????????????????
		MaxAge:     30,       //????????????????????????
		Compress:   false,    //??????????????????
	})
	return zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), priority)

}

func InitLog(debug bool, LogPath string) {
	var coreArr []zapcore.Core
	Log = &myLog{
		ZapLog: nil,
		debug:  debug,
	}
	//???????????????
	encoderConfig := zap.NewProductionEncoderConfig()            //NewJSONEncoder()??????json?????????NewConsoleEncoder()????????????????????????
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        //??????????????????
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder //???????????????????????????????????????????????????zapcore.CapitalLevelEncoder????????????
	//encoderConfig.EncodeCaller = zapcore.FullCallerEncoder        //????????????????????????
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	for key := range levelMap {
		fileCore := setLoggerFile(key, encoder, LogPath)
		coreArr = append(coreArr, fileCore)
	}
	Log.ZapLog = zap.New(zapcore.NewTee(coreArr...), zap.AddCaller())
}
