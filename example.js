import { makeWASocket } from "./index.js"
import util from "util"

const conn = makeWASocket({
    DbPath: "anu.db"
    //Logger: { Database: "DEBUG", Client: "DEBUG", Color: true }
})
console.log(conn)
if (!conn.Store().ID) {
    conn.Connect()
    console.log(conn.PairPhone("xxx"))
    setInterval(() => {
        if (!conn.IsLoggedIn()) {
            conn.Disconnect()
        }
    }, 60000)
} else {
    conn.Connect()
    const a = conn.SetGroupLocked("xxx@g.us", true)
    console.log(a)
}
conn.Event((a) => {
    //console.log(a)
})