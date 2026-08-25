package main

import (
    "context"
    "encoding/json"
    "fmt"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/store"
    "go.mau.fi/whatsmeow/store/sqlstore"
    "go.mau.fi/whatsmeow/types"
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
    _ "github.com/jackc/pgx/v5/stdlib"
)

var ctx = context.Background()
type J map[string]any
var nextHandle atomic.Uint64
func newHandle() string { return fmt.Sprintf("%d", nextHandle.Add(1)) }
func Throw(env napi.EnvType, err any) any {
    return napi.ThrowError(env, "", fmt.Sprintf("%s", err))
}

var deviceMap sync.Map

type fileLogger struct {
    module   string
    minLevel string
    mu       *sync.Mutex
    w        io.Writer
}

func logLevelRank(s string) int {
    switch strings.ToUpper(s) {
    case "DEBUG": return 0
    case "INFO":  return 1
    case "WARN":  return 2
    case "ERROR": return 3
    default:      return 0
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

func makeLogger(module, level, file string, color bool) waLog.Logger {
    if level == "" { return nil }
    if file != "" {
        f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil { return waLog.Stdout(module, level, color) }
        var mu sync.Mutex
        return &fileLogger{module: module, minLevel: level, mu: &mu, w: f}
    }
    return waLog.Stdout(module, level, color)
}

//go:linkname RegisterNapi sirherobrine23.com.br/Sirherobrine23/napi-go/module.Register
func RegisterNapi(env napi.EnvType, export *napi.Object) {

    containerFn, _ := napi.GoFuncOf(env, func(driver, dsn, logLevel string) any {
        dbLog := makeLogger("Database", logLevel, "", false)

        if driver == "" { driver = "sqlite3" }
        if dsn == "" { dsn = "file:mywagua.db?_foreign_keys=on" }

        container, err := sqlstore.New(ctx, driver, dsn, dbLog)
        if err != nil { return Throw(env, err) }

        result := map[string]any{}

        devInfo := func(dev *store.Device) map[string]any {
            h := newHandle()
            deviceMap.Store(h, dev)
            jid := ""
            if dev.ID != nil { jid = dev.ID.String() }
            return map[string]any{"handle": h, "jid": jid}
        }

        result["GetFirstDevice"] = func() any {
            dev, err := container.GetFirstDevice(ctx)
            if err != nil { return Throw(env, err) }
            return devInfo(dev)
        }

        result["GetAllDevices"] = func() any {
            devs, err := container.GetAllDevices(ctx)
            if err != nil { return Throw(env, err) }
            infos := make([]map[string]any, len(devs))
            for i, dev := range devs { infos[i] = devInfo(dev) }
            return infos
        }

        result["GetDevice"] = func(jid string) any {
            j, err := types.ParseJID(jid)
            if err != nil { return Throw(env, err) }
            dev, err := container.GetDevice(ctx, j)
            if err != nil { return Throw(env, err) }
            if dev == nil { return nil }
            return devInfo(dev)
        }

        result["PutDevice"] = func() any {
            dev := container.NewDevice()
            return devInfo(dev)
        }

        result["DeleteDevice"] = func(handle string) any {
            v, ok := deviceMap.Load(handle)
            if !ok { return Throw(env, fmt.Errorf("device not found: %s", handle)) }
            dev := v.(*store.Device)
            if err := container.DeleteDevice(ctx, dev); err != nil { return Throw(env, err) }
            deviceMap.Delete(handle)
            return nil
        }

        return result
    })

    clientFn, _ := napi.GoFuncOf(env, func(deviceHandle string, cfg any) any {
        v, ok := deviceMap.Load(deviceHandle)
        if !ok { return Throw(env, fmt.Errorf("device not found: %s", deviceHandle)) }
        deviceStore := v.(*store.Device)

        var config struct {
            Logger struct {
                Client string
                Color  bool
                File   string
            }
            OsName string
        }
        if err := json.Unmarshal([]byte(ToJson(cfg)), &config); err != nil { return Throw(env, err) }

        store.DeviceProps.PlatformType = waCompanionReg.DeviceProps_CHROME.Enum()
        if config.OsName == "" { config.OsName = "Chrome" }
        store.DeviceProps.Os = proto.String(config.OsName)

        clientLog := makeLogger("Client", config.Logger.Client, config.Logger.File, config.Logger.Color)
        client := whatsmeow.NewClient(deviceStore, clientLog)
        // pesan pertama dari pengirim yang belum ada sesi/sender key gagal didekripsi;
        // retry receipt ke pengirim sering nggak nyampe, jadi minta plaintext-nya ke hp utama
        client.AutomaticMessageRerequestFromPhone = true

        var (
            queue []string
            mu    sync.Mutex
        )
        client.AddEventHandler(func(evt interface{}) {
            mu.Lock()
            defer mu.Unlock()
            queue = append(queue, ToJson(map[string]interface{}{
                "type": fmt.Sprintf("%T", evt),
                "evt":  evt,
            }))
        })

        vm := goja.New()
        vm.Set("ctx", ctx)
        vm.Set("client", client)
        conn := Sends(env, client)
        p := conn["_params"].(map[string]string)
        conn["run"] = func(value string) any {
            result, err := vm.RunString(value)
            if err != nil { return Throw(env, err) }
            return fmt.Sprintf("%s", result)
        }
        conn["getEvt"] = func() []string {
            mu.Lock()
            defer mu.Unlock()
            if len(queue) == 0 { return nil }
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

    export.Set("Container", containerFn)
    export.Set("Client", clientFn)
}

func main() {}
