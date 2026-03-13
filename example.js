import { makeWASocket } from "./index.js"
import util from "util"

const conn = makeWASocket({
    DbPath: "anu.db"
    Logger: { Database: "DEBUG", Client: "DEBUG", Color: true }
})
console.log(conn)
if (!conn.Store().ID) {
    conn.Connect()
    console.log(conn.PairPhone("62xxx"))
    setInterval(() => {
        if (!conn.IsLoggedIn()) {
            conn.Disconnect()
        }
    }, 60000)
} else {
    conn.Connect()
    conn.SendPresence("available")
    conn.SendPresence("unavailable")
}
conn.Event((a) => {
    //console.log(a)
})