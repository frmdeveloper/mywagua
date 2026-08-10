import {swmeow} from "./simple.js"
import * as message from "./message.js"
import { createRequire } from "module"
const require = createRequire(import.meta.url)
const addon = require("./main.node")

export function Container(driver = "", dsn = "", logLevel = "") {
    const container = addon.Container(driver, dsn, logLevel)

    const dbPath = (driver === "sqlite3" || driver === "")
        ? (dsn || "file:mywagua.db").replace(/^file:/, "").split("?")[0]
        : null

    const wrap = (item) => {
        const handle = typeof item === "string" ? item : item.handle
        const jid    = typeof item === "string" ? undefined : item.jid
        return { handle, dbPath, jid }
    }

    const origDelete = container.DeleteDevice.bind(container)
    container.DeleteDevice = (device) => {
        const handle = typeof device === "string" ? device : device.handle
        return origDelete(handle)
    }

    const orig = {
        GetFirstDevice:  container.GetFirstDevice.bind(container),
        GetAllDevices:   container.GetAllDevices.bind(container),
        GetDevice:       container.GetDevice.bind(container),
        PutDevice:       container.PutDevice.bind(container),
    }

    container.GetFirstDevice = ()      => wrap(orig.GetFirstDevice())
    container.GetAllDevices  = ()      => (orig.GetAllDevices() || []).map(wrap)
    container.GetDevice      = (jid)   => { const h = orig.GetDevice(jid); return h ? wrap(h) : null }
    container.PutDevice      = ()      => wrap(orig.PutDevice())

    return container
}

export function makeClient(device, config = {}) {
    const handle = typeof device === "string" ? device : device.handle
    const go = addon.Client(handle, config)
    return mappingSock(go)
}

function mappingSock(go) {
const types = go.run(`JSON.stringify(Object.keys(client).map(a => ({name:a,type:typeof client[a]})))`)
const mapped = { ctx:"golangContextBackground()" }
JSON.parse(types).forEach(i => {
    if (i.type == "function") mapped[i.name] = null
})

const params = go._params || {}
const wrappedGo = {}
for (const [key, val] of Object.entries(go)) {
    if (key === '_params') continue
    if (typeof val === 'function') {
        const p = params[key] ?? ''
        const wrapper = { [key](...args) { return val(...args) } }[key]
        Object.defineProperty(wrapper, 'toString', {
            value: () => `function ${key}(${p}) { [native code] }`,
            configurable: true
        })
        wrappedGo[key] = wrapper
    } else {
        wrappedGo[key] = val
    }
}

const sock = {
    whatsmeow: true,
    ...mapped, ...wrappedGo, simple:swmeow,
    Event(callback) {
        return setInterval(() => {
            const evts = go.getEvt()
            if (evts) evts.forEach(i => callback(JSON.parse(i)))
        }, 100)
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
    Store() { return JSON.parse(go.GetStore()) },
    getDevice(id) {
        return /^3A.{18}$/.test(id) ? 'ios'
            : /^3E.{20}$/.test(id) ? 'web'
            : /^(.{21}|.{32})$/.test(id) ? 'android'
            : /^(3F|.{18}$)/.test(id) ? 'desktop'
            : 'unknown'
    },
    decodeJid(jid) { return jid.replace(/:[0-9]+/, "") }
}

return {...sock, ...binder(message, go)}
}

function binder(target, fill) {
    const binded = {}
    for (const key in target) {
        if (typeof target[key] === "function") {
            binded[key] = target[key].bind(fill)
        }
    }
    return binded
}

