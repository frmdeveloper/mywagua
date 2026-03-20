export async function generateWAMessageFromContent(jid, content = {}, options = {}) {
    let message = {}
    if ("text" in content) {
        message = { extendedTextMessage: { text: content.text } }
    }
    const key = getContentType(message)
    if ("contextInfo" in content) {
        message[key].contextInfo = content.contextInfo
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
export async function sendMessage(jid, content = {}, options = {}) {
    const message = await generateWAMessageFromContent(jid, content, options)
    return this.relayMessage(jid, message, options)
}
export function getContentType(content) {
    if (content) {
        const keys = Object.keys(content)
        return keys.find(k => (k === 'conversation' || k.includes('Message')) && k !== 'senderKeyDistributionMessage')
    }
}