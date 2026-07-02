package main

import (
    "context"
    "encoding/json"
    "fmt"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/store"
    "go.mau.fi/whatsmeow/store/sqlstore"
    "go.mau.fi/whatsmeow/proto/waCompanionReg"
    "google.golang.org/protobuf/proto"
    "github.com/dop251/goja"
    "sirherobrine23.com.br/Sirherobrine23/napi-go"
    "io"
    "os"
    "strings"
    "sync"
    "sync/atomic"
    "time"
    waLog "go.mau.fi/whatsmeow/util/log"
    meowcaller "github.com/purpshell/meowcaller"
    _ "unsafe"
    _ "sirherobrine23.com.br/Sirherobrine23/napi-go/module"
    _ "github.com/mattn/go-sqlite3"
)

var ctx = context.Background()
type J map[string]any
var nextHandle atomic.Uint64
func newHandle() string { return string(nextHandle.Add(1)) }
func Throw(env napi.EnvType, err any) any {
    return napi.ThrowError(env, "", fmt.Sprintf("%s",err))
}

type Config struct {
    Logger struct {
        Database string
        Client string
        Color bool
        File string
    }
    DbPath string
    OsName string
}

type fileLogger struct {
    module   string
    minLevel string
    mu       *sync.Mutex
    w        io.Writer
}

func logLevelRank(s string) int {
    switch strings.ToUpper(s) {
    case "DEBUG": return 0
    case "INFO": return 1
    case "WARN": return 2
    case "ERROR": return 3
    default: return 0
    }
}

func (f *fileLogger) write(level, msg string, args ...interface{}) {
    if logLevelRank(level) < logLevelRank(f.minLevel) { return }
    f.mu.Lock()
    defer f.mu.Unlock()
    fmt.Fprintf(f.w, "%s [%s] [%s] %s\n", time.Now().Format("15:04:05.000"), level, f.module, fmt.Sprintf(msg, args...))
}
func (f *fileLogger) Debugf(msg string, args ...interface{}) { f.write("DEBUG", msg, args...) }
func (f *fileLogger) Infof(msg string, args ...interface{})  { f.write("INFO", msg, args...) }
func (f *fileLogger) Warnf(msg string, args ...interface{})  { f.write("WARN", msg, args...) }
func (f *fileLogger) Errorf(msg string, args ...interface{}) { f.write("ERROR", msg, args...) }
func (f *fileLogger) Sub(module string) waLog.Logger {
    return &fileLogger{module: f.module + "/" + module, minLevel: f.minLevel, mu: f.mu, w: f.w}
}

//go:linkname RegisterNapi sirherobrine23.com.br/Sirherobrine23/napi-go/module.Register
func RegisterNapi(env napi.EnvType, export *napi.Object) {
sock, _ := napi.GoFuncOf(env, func(cfg any) any {
    var config Config
    err := json.Unmarshal([]byte(ToJson(cfg)), &config)
    if err != nil { return Throw(env,err) }
    
    latestVer, err := whatsmeow.GetLatestVersion(ctx,nil)
    if err != nil { return Throw(env,err) }
    store.SetWAVersion(*latestVer)
    store.DeviceProps.PlatformType = waCompanionReg.DeviceProps_DESKTOP.Enum()
    if config.OsName == "" {
        config.OsName = "My WA Gua"
    }
    store.DeviceProps.Os = proto.String(config.OsName)
    
    var logFile *os.File
    if config.Logger.File != "" {
        logFile, err = os.OpenFile(config.Logger.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil { return Throw(env, err) }
    }
    var logMu sync.Mutex
    makeLogger := func(module, level string) waLog.Logger {
        if level == "" { return nil }
        if logFile != nil {
            return &fileLogger{module: module, minLevel: level, mu: &logMu, w: logFile}
        }
        return waLog.Stdout(module, level, config.Logger.Color)
    }

    dbLog := makeLogger("Database", config.Logger.Database)
    
    if config.DbPath == "" { config.DbPath = "mywagua.db" }
    container, err := sqlstore.New(ctx, "sqlite3", "file:"+config.DbPath+"?_foreign_keys=on", dbLog)
    if err != nil { return Throw(env,err) }

    deviceStore, err := container.GetFirstDevice(ctx)
    if err != nil { return Throw(env,err) }

    clientLog := makeLogger("Client", config.Logger.Client)
    client := whatsmeow.NewClient(deviceStore, clientLog)

    var (
        queue []string
        mu sync.Mutex
    )
    client.AddEventHandler(func(evt interface {}) {
        mu.Lock()
        defer mu.Unlock()
        au := map[string]interface{}{
            "type": fmt.Sprintf("%T", evt),
            "evt": evt,
        }
        queue = append(queue, ToJson(au))
    })

    vm := goja.New()
    vm.Set("ctx", ctx)
    vm.Set("client", client)
    conn := Sends(env, client)
    p := conn["_params"].(map[string]string)
    conn["run"] = func(value string) any {
        result, err :=  vm.RunString(value)
        if err != nil { return Throw(env,err) }
        return fmt.Sprintf("%s",result)
    }
    conn["getEvt"] = func()[]string{
        mu.Lock()
        defer mu.Unlock()
        if len(queue) == 0 {
            return nil
        }
        result := make([]string, len(queue))
        copy(result, queue)
        queue = queue[:0]
        return result
    }
    push := func(v any) {
        mu.Lock()
        defer mu.Unlock()
        queue = append(queue, ToJson(v))
    }
    caller := meowcaller.NewClient(client)
    SetupCaller(env, caller, push, conn, p)
    return conn
})
export.Set("create", sock)
}

func main() {}