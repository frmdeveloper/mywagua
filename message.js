import { existsSync } from "fs"

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
                    ?.match(/(image|sticker|video|document)/i)?.[0]
    if (mediatype) {
        const isSticker = mediatype == "sticker"
        let mediacontent
        if ("url" in content[mediatype]) {
            if (isUrl(content[mediatype].url)) mediacontent = { Url: content[mediatype].url }
            if (existsSync(content[mediatype].url)) mediacontent = { File: content[mediatype].url }
        }
        if ("base64" in content[mediatype]) mediacontent = { Base64: content[mediatype].base64 }
        const types = isSticker ? "Image" : mediatype.replace(/^./, ma => ma.toUpperCase())
        const upload = this.Upload(mediacontent, "WhatsApp "+types+" Keys")
        message[mediatype+"Message"] = upload
        if (isSticker) message[mediatype+"Message"].mimetype = "image/webp"
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
    if (options.quoted) {
        message[key].contextInfo = message[key].contextInfo || {}
        message[key].contextInfo.stanzaID = options.quoted.Info.ID
        message[key].contextInfo.participant = options.quoted.Info.Sender
        message[key].contextInfo.quotedMessage = options.quoted.RawMessage
    }
    if (("text" in content) && !message[key].contextInfo) {
        message = { conversation:content.text }
    }
    return message
}
export function getContentType(content) {
    if (content) {
        const keys = Object.keys(content)
        return keys.find(k => (k === 'conversation' || k.includes('Message')) && k !== 'senderKeyDistributionMessage')
    }
}
const isUrl = (str) => {
  return /^(https?:\/\/)[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!$&'()*+,;.=]+$/.test(str)
}