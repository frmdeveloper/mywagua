import { readFileSync, writeFileSync, rmSync, createWriteStream } from "fs"
import { createRequire } from "module"
import { pathToFileURL } from "url"
import { execFile } from "child_process"
import { promisify } from "util"
import { randomBytes } from "crypto"
import { Readable } from "stream"
import { pipeline } from "stream/promises"
import { tmpdir } from "os"
import { join } from "path"

const run = promisify(execFile)

const MAX_DIM = 200
const QUALITY = 50
const HEADER_TIMEOUT = 60000
const REUSE_MAX = 20971520
const PROBE_BYTES = 524288
const UA = "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Mobile Safari/537.36"

let jimp = null
let searched = false

async function loadJimp() {
    if (searched) return jimp
    searched = true
    const attempts = [
        () => import("jimp"),
        () => import(pathToFileURL(createRequire(process.cwd() + "/").resolve("jimp")).href),
    ]
    for (const attempt of attempts) {
        try {
            const mod = (await attempt()).default
            if (typeof mod?.read === "function") { jimp = mod; break }
        } catch {}
    }
    return jimp
}

let ffmpeg = null
async function hasFfmpeg() {
    if (ffmpeg === null) {
        try { await run("ffmpeg", ["-version"], { timeout: 10000 }); ffmpeg = true }
        catch { ffmpeg = false }
    }
    return ffmpeg
}

const tempName = () => join(tmpdir(), "mywagua-" + randomBytes(8).toString("hex"))

async function get(url) {
    const ctrl = new AbortController()
    const timer = setTimeout(() => ctrl.abort(), HEADER_TIMEOUT)
    try {
        const res = await fetch(url, { headers: { "User-Agent": UA }, signal: ctrl.signal })
        if (!res.ok) throw new Error("HTTP " + res.status)
        return res
    } finally {
        clearTimeout(timer)
    }
}

async function* head(body, limit) {
    let taken = 0
    for await (const chunk of body) {
        yield chunk
        taken += chunk.length
        if (taken >= limit) return
    }
}

async function fetchFile(url) {
    const res = await get(url)
    const length = Number(res.headers.get("content-length")) || 0
    const whole = length > 0 && length <= REUSE_MAX
    const path = tempName()
    const body = Readable.fromWeb(res.body)
    await pipeline(whole ? body : head(body, PROBE_BYTES), createWriteStream(path))
    return { path, whole }
}

async function fetchBytes(url) {
    const res = await get(url)
    const length = Number(res.headers.get("content-length")) || 0
    if (length > REUSE_MAX) {
        await res.body.cancel()
        throw new Error("larger than " + REUSE_MAX + " bytes")
    }
    return Buffer.from(await res.arrayBuffer())
}

async function firstFrame(path) {
    try {
        const { stdout } = await run("ffmpeg", ["-v", "error", "-y", "-i", path, "-frames:v", "1", "-an", "-f", "image2", "-c:v", "mjpeg", "pipe:1"], {
            encoding: "buffer",
            maxBuffer: 33554432,
            timeout: 60000,
        })
        return stdout
    } catch (e) {
        const detail = e.stderr?.toString().trim().split("\n").pop() || e.message
        throw new Error("ffmpeg: " + detail.slice(0, 120))
    }
}

export async function generateThumbnail(type, media) {
    const Jimp = await loadJimp()
    if (!Jimp) return null
    let temp
    let whole = false
    let keep = false
    try {
        let raw
        if (type == "video") {
            if (!await hasFfmpeg()) return null
            let path = media.File
            if (!path && media.Url) {
                const got = await fetchFile(media.Url)
                temp = path = got.path
                whole = got.whole
            }
            if (!path && media.Base64) {
                temp = path = tempName()
                writeFileSync(path, Buffer.from(media.Base64, "base64"))
            }
            if (!path) return null
            raw = await firstFrame(path)
        } else if (media.Base64) raw = Buffer.from(media.Base64, "base64")
        else if (media.File) raw = readFileSync(media.File)
        else if (media.Url) {
            raw = await fetchBytes(media.Url)
            whole = true
        }
        if (!raw?.length) return null
        const img = await Jimp.read(raw)
        const { width, height } = img.bitmap
        const small = await img.scaleToFit(MAX_DIM, MAX_DIM).quality(QUALITY).getBufferAsync(Jimp.MIME_JPEG)
        const thumb = { JPEGThumbnail: small.toString("base64"), width, height }
        if (!media.Url || !whole) return thumb
        keep = type == "video"
        return keep ? { ...thumb, file: temp } : { ...thumb, bytes: raw }
    } catch (e) {
        console.log("thumbnail skipped:", e.message)
        return null
    } finally {
        if (temp && !keep) rmSync(temp, { force: true })
    }
}
