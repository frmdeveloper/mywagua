import { writeFileSync, existsSync, rmSync } from "fs"
import {randomBytes} from "crypto"
import {generateThumbnail} from "./thumbnail.js"

export async function sendMessage(jid, content = {}, options = {}) {
    const message = await generateWAMessageFromContent.bind(this)(jid, content, options)
    return this.relayMessage(jid, message, options)
}
export async function generateWAMessageFromContent(jid, content = {}, options = {}) {
    let message = {}
    if ("text" in content) {
        message = { extendedTextMessage: { text: content.text } }
    }
    const mediatype = Object.keys(content)[0]
                    ?.match(/(audio|image|sticker|video|document)/i)?.[0]
    if (mediatype) {
        const isSticker = mediatype == "sticker"
        let mediacontent
        let rand
        if (Buffer.isBuffer(content[mediatype])) {
            if (content[mediatype].length <= 5242880) {
                const b64 = content[mediatype].toString("base64")
                content[mediatype] = {base64:b64}
            } else {
                rand = randomBytes(14).toString('hex')
                writeFileSync("assets/"+rand, content[mediatype])
                content[mediatype] = {url:"assets/"+rand}
            }
        }
        if ("url" in content[mediatype]) {
            if (isUrl(content[mediatype].url)) mediacontent = { Url: content[mediatype].url }
            if (existsSync(content[mediatype].url)) mediacontent = { File: content[mediatype].url }
        }
        if ("base64" in content[mediatype]) mediacontent = { Base64: content[mediatype].base64 }
        const types = isSticker ? "Image" : mediatype.replace(/^./, ma => ma.toUpperCase())
        if (!mediacontent) throw new Error("Media not found")
        const thumb = /^(image|video)$/.test(mediatype) ? await generateThumbnail(mediatype, mediacontent) : null
        if (thumb?.file) mediacontent = { File: thumb.file }
        else if (thumb?.bytes && mediacontent.Url && thumb.bytes.length <= 5242880) {
            mediacontent = { Base64: thumb.bytes.toString("base64") }
        }
        let upload
        try {
            upload = this.Upload(mediacontent, "WhatsApp "+types+" Keys")
        } finally {
            if (rand) rmSync("assets/"+rand)
            if (thumb?.file) rmSync(thumb.file, { force: true })
        }
        message[mediatype+"Message"] = upload
        if (thumb) {
            message[mediatype+"Message"].JPEGThumbnail = thumb.JPEGThumbnail
            message[mediatype+"Message"].width = thumb.width
            message[mediatype+"Message"].height = thumb.height
        }
        if (mediatype == "audio") message[mediatype+"Message"].mimetype = "audio/mpeg"
        if (mediatype == "image") message[mediatype+"Message"].mimetype = "image/jpeg"
        if (mediatype == "video") message[mediatype+"Message"].mimetype = "video/mp4"
        if (mediatype == "sticker") message[mediatype+"Message"].mimetype = "image/webp"
    } 
    const key = getContentType(message)
    if ("caption" in content) message[key].caption = content.caption
    if ("mimetype" in content) message[key].mimetype = content.mimetype
    if ("fileName" in content) message[key].fileName = content.fileName
    if ("contextInfo" in content) message[key].contextInfo = content.contextInfo
    if ("mentions" in content) {
        message[key].contextInfo = message[key].contextInfo || {}
        message[key].contextInfo.mentionedJID = content.mentions
    }
    if (content.parseMention) {
        const teks = [content.text, content.caption].filter(a => typeof a === "string").join(" ")
        message[key].contextInfo = message[key].contextInfo || {}
        message[key].contextInfo.mentionedJID = [...new Set([
            ...(message[key].contextInfo.mentionedJID || []),
            ...this.ParseMention(teks),
        ])]
    }
    if (options.quoted) {
        message[key].contextInfo = message[key].contextInfo || {}
        message[key].contextInfo.stanzaID = options.quoted.Info.ID
        message[key].contextInfo.participant = options.quoted.Info.Sender
        message[key].contextInfo.quotedMessage = options.quoted.RawMessage
    }
    if (("text" in content) && !message[key].contextInfo) {
        message = { conversation:content.text }
    }
    if ("delete" in content) {
        message = { protocolMessage:{ key:content.delete.key ?? content.delete, type:0 } }
    }
    if ("edit" in content) {
        const editedMessage = content.text
            ? { conversation: content.text }
            : message
        message = { protocolMessage:{ key:content.edit.key ?? content.edit, type:14, editedMessage } }
    }
    return message
}
export async function editMessage(jid, key, newContent = {}, options = {}) {
    return sendMessage.call(this, jid, { edit: key, ...newContent }, options)
}
export function getContentType(content) {
    if (content) {
        const keys = Object.keys(content)
        return keys.find(k => (k === 'conversation' || k.includes('Message')) && k !== 'senderKeyDistributionMessage')
    }
}
const isUrl = (url) => {
    return url.match(new RegExp(/https?:\/\/(www\.)?[-a-zA-Z0-9@:%._+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_+.~#?&/=]*)/, 'gi'))
}