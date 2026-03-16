import { simple } from "./simple.js"
import { createRequire } from "module"
const require = createRequire(import.meta.url)
const addon = require("./main.node")

export const makeWASocket = (config) => {
    const go = addon.create(config)
    const sock = mappingSock(go)
    return sock
}
function mappingSock(go) {
const types = go.run(`JSON.stringify(Object.keys(client).map(a => ({name:a,type:typeof client[a]})))`)
const mapped = { ctx:"golangContextBackground()" }
JSON.parse(types).forEach(i => {
    if (i.type == "function") mapped[i.name] = null
})
const sock = {
    ...mapped, ...go, simple,
    Event(callback) {
        setInterval(() => go.getEvt().forEach(i=>callback(JSON.parse(i))),100)
    },
    Call(name, ...arg) {
        const command = arg.map(a => {
            return a == this.ctx ? "ctx" :
            /boolean|number/.test(typeof a) ? a :
            JSON.stringify(a)
        }).join(", ")
        console.log(`client.${name}(${command})`)
        return go.run(`client.${name}(${command})`)
    },
    Store() { return JSON.parse(go.GetStore()) }
}
return sock
}