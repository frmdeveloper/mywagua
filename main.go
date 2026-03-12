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
    "sync"
    "sync/atomic"
    waLog "go.mau.fi/whatsmeow/util/log"
    _ "unsafe"
    _ "sirherobrine23.com.br/Sirherobrine23/napi-go/module"
    _ "github.com/ncruces/go-sqlite3/driver"
    _ "github.com/ncruces/go-sqlite3/embed"
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
    }
    DbPath string
    OsName string
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
    
    dbLog := waLog.Stdout("Database", config.Logger.Database, config.Logger.Color)
    if config.Logger.Database == "" { dbLog = nil }
    
    if config.DbPath == "" { config.DbPath = "mywagua.db" }
    container, err := sqlstore.New(ctx, "sqlite3", "file:"+config.DbPath+"?_foreign_keys=on", dbLog)
    if err != nil { return Throw(env,err) }

    deviceStore, err := container.GetFirstDevice(ctx)
    if err != nil { return Throw(env,err) }

    clientLog := waLog.Stdout("Client", config.Logger.Client, config.Logger.Color)
    if config.Logger.Client == "" { clientLog = nil }
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
    return conn
})
export.Set("create", sock)
}

func main() {}