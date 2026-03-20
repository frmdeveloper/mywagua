const simple = async(conn, m) => {
    if (!m) return {}

    const msg = {}
    msg.full = m
    if (m.Info) {
        msg.id = m.Info.ID
        msg.from = m.Info.Chat.endsWith("@lid") ? m.Info.SenderAlt?.replace(/:[0-9]+/,"") : m.Info.Chat?.replace(/:[0-9]+/,"")
        msg.fromMe = m.Info.IsFromMe
        msg.isGroup = m.Info.IsGroup
        msg.sender = m.Info.SenderAlt?.replace(/:[0-9]+/,"") || m.Info.Sender?.replace(/:[0-9]+/,"")
        msg.pushname = m.Info.PushName
    }
    if (m.RawMessage) {
        m.message = m.RawMessage
        //if (m?.message?.messageContextInfo) delete m.message.messageContextInfo
        //if (m?.message?.senderKeyDistributionMessage) delete m.message.senderKeyDistributionMessage
        m.message = m.message.viewOnceMessageV2?.message ||
            m.message.documentWithCaptionMessage?.message ||
            m.message.editedMessage?.message?.protocolMessage?.editedMessage ||
            m.message 
        let mtype = Object.keys(m.message)
        msg.type = mtype.find(k => (k === 'conversation' || k.includes('Message')) && k !== 'senderKeyDistributionMessage')
        msg.msg = m.message[msg.type]
        msg.text = m.message.conversation || msg.msg?.text || msg.msg?.caption || msg.msg?.selectedId || ''
        const terpusah = /^(#|\!|\/|\.)( +)/.test(msg.text)
        if (terpusah) msg.text = msg.text.replace(" ", "")
        msg.args = msg.text?.trim().split(/ +/).slice(1)
        msg.prefix = /^[!#%./\\]/.test(msg.text) ? msg.text.match(/^[!#%./\\]/gi) : ''
        msg.command = msg.text?.slice(0).trim().split(/ +/).shift().toLowerCase()
        msg.q = msg.args?.join(" ")
        msg.mentionedJid = msg.msg && msg.msg.contextInfo && msg.msg.contextInfo.mentionedJID && msg.msg.contextInfo.mentionedJID.length && msg.msg.contextInfo.mentionedJID || []
    }
    let quoted = msg?.msg?.contextInfo?.quotedMessage
    msg.quoted = {}
    if (quoted) {
        quoted = quoted.groupMentionedMessage?.message || quoted
        let type = Object.keys(quoted)[0]
        const isi = quoted[type]
        msg.quoted.type = type
        msg.quoted.from = msg.from
        msg.quoted.id = msg.msg.contextInfo.stanzaID
        msg.quoted.sender = msg.msg.contextInfo.participant.replace(/:[0-9]+/,"")
        msg.quoted.fromMe = msg.quoted.sender.replace(/:[0-9]+/,"") == conn.Store().ID
        //msg.quoted.key = {remoteJid: msg.quoted.from, id: msg.quoted.id, fromMe: msg.quoted.fromMe, participant: msg.quoted.sender}
        msg.quoted.text = isi?.caption || isi?.text || isi?.message?.documentMessage?.caption || isi || ""
        msg.quoted.mentionedJid = quoted[type]?.contextInfo?.mentionedJID
        msg.quoted.groupMentions = quoted[type]?.contextInfo?.groupMentions
        msg.quoted.full = quoted
    }
    return msg
}
export default simple