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
    if ("image" in content) {
        let mediacontent
        if ("url" in content.image) {
            if (isUrl(content.image.url)) mediacontent = { Url: content.image.url }
            if (existsSync(content.image.url)) mediacontent = { File: content.image.url }
        }
        if ("base64" in content.image) {
            mediacontent = { Base64: content.image.base64 }
        }
        const upload = this.Upload(mediacontent, "WhatsApp Image Keys")
        message = { imageMessage: upload }
        if ("caption" in content) message.imageMessage.caption = content.caption
    }
    const key = getContentType(message)
    if ("contextInfo" in content) {
        message[key].contextInfo = content.contextInfo
    }
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