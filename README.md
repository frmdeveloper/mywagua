<div align="center">
<img src="assets/bukankahinimy.jpg" width="30%">

# 🫪 My WA Gua 🫪

### *Whatsmeow tapi Node.js, njir parah*
---
</div>

## Apaan sih ini wkwk

Njir jadi gini ceritanya wkwk — whatsmeow (Go) tapi bisa dipake dari Node.js lewat N-API. Hadeh ribet amat? Kagak kok, malah lebih kenceng karena backend-nya Go. Njir.

## Yang harus ada dulu

- Node.js 20+
- Go 1.21+
- ffmpeg (opsional, buat konvert video ke audio pas voice call)

## Install

```bash
npm install mywagua@github:frmdeveloper/mywagua
```

---

## Flow Konek

```js
import { Container, makeClient } from "mywagua"

const db     = Container("sqlite3", "file:bot.db?_foreign_keys=on", "INFO")
const device = db.GetFirstDevice()
const conn   = makeClient(device, { Logger: { Client: "INFO" } })

conn.Connect()

conn.Event(async ({ type, evt }) => {
    if (type === "*events.Connected") {
        console.log("Connected njir!")
    }

    if (type === "*events.Message") {
        const m = await conn.simple(evt)
        console.log(m.sender, m.text)
    }
})
```

---

## API

### Container

Basically `sqlstore.New`-nya whatsmeow tapi versi JS njir.

```js
const db = Container(driver, dsn, logLevel)
```

```js
Container("sqlite3", "file:bot.db?_foreign_keys=on", "INFO")

Container("pgx", "postgres://user:pass@localhost:5432/dbname", "DEBUG")

Container()
```

Log level: `"DEBUG"` `"INFO"` `"WARN"` `"ERROR"` — kosongkan kalo mau silent njir wkwk.

#### Method-method db

```js
db.GetFirstDevice()
db.GetAllDevices()
db.GetDevice("6281234567890@s.whatsapp.net")
db.PutDevice()
db.DeleteDevice(handle)
```

Semua device method return object `{handle, dbPath, jid}`:

```js
const device = db.GetFirstDevice()
// { handle: '1', dbPath: 'bot.db', jid: '6281234567890:4@s.whatsapp.net' }
```

---

### makeClient(device, config)

```js
const conn = makeClient(device, {
    Logger: { Client: "INFO" },
    OsName: "Chrome"
})
```

---

### Banyak Bot njir kocak parah wkwk

```js
const db = Container("pgx", "postgres://user:pass@localhost/db")

const devices = db.GetAllDevices()
// [{ handle: '1', jid: '628xx...' }, { handle: '2', jid: '628yy...' }]

const bot1 = makeClient(devices[0], { OsName: "Chrome" })
const bot2 = makeClient(devices[1], { OsName: "Chrome" })

const botA = makeClient(devices.find(d => d.jid.startsWith("62851")), config)
```

---

## Koneksi

```js
conn.Connect()
conn.Disconnect()
conn.IsConnected()
conn.IsLoggedIn()
conn.Logout()
conn.ResetConnection()
conn.WaitForConnection(seconds)
conn.PairPhone("6281234567890")
conn.Store()
```

---

## Event

```js
const stop = conn.Event(async ({ type, evt }) => {
    // ...
})

stop()
```

| Type | Kapan |
|---|---|
| `*events.Message` | Ada pesan masuk njir wkwk |
| `*events.Receipt` | Read/delivery receipt |
| `*events.Presence` | Online/offline |
| `*events.ChatPresence` | Lagi ngetik/rekam hadeh |
| `*events.Connected` | Nyambung njir finally |
| `*events.Disconnected` | Putus njir hadeh |
| `*events.LoggedOut` | Kena banned/logout |
| `*events.CallOffer` | Ada telpon masuk |
| `meowcaller.IncomingCall` | Telpon masuk (meowcaller) |
| `meowcaller.CallReady` | Media call aktif |
| `meowcaller.CallEnd` | Telpon selesai |
| `meowcaller.AudioFrame` | Frame audio live |

---

## Kirim Pesan

```js
await conn.sendMessage(jid, { text: "Halo njir!" })

await conn.sendMessage(jid, {
    text: "Mention nih @628xxx",
    mentions: ["628xxx@s.whatsapp.net"]
})

await conn.sendMessage(jid, { image: { url: "https://..." }, caption: "Foto njir" })

await conn.sendMessage(jid, { text: "Balas" }, { quoted: evt })

conn.MarkRead([msgId], Date.now(), jid, sender)

conn.SendChatPresence(jid, "composing", "")
conn.SendChatPresence(jid, "paused", "")

conn.RevokeMessage(chat, sender, id)

conn.BuildReaction(chat, sender, id, "👍")

conn.BuildPollCreation("Pilih dong wkwk", ["A", "B", "C"], 1)

conn.BuildEdit(chat, id, { conversation: "Diedit hadeh" })

conn.GenerateMessageID()
```

---

## Media

```js
conn.Upload(args)
conn.DownloadAny(message)
```

---

## Grup

```js
conn.GetGroupInfo(jid)
conn.GetGroupInfoFromLink(link)
conn.GetGroupInviteLink(jid, false)
conn.CreateGroup("Nama Grup", ["628xxx@s.whatsapp.net"])
conn.LeaveGroup(jid)
conn.JoinGroupWithLink(link)
conn.JoinGroupWithInvite(inviter, jid, code, expiration)
conn.UpdateGroupParticipants(jid, ["628xxx@s.whatsapp.net"], "add")
conn.UpdateGroupParticipants(jid, ["628xxx@s.whatsapp.net"], "remove")
conn.UpdateGroupParticipants(jid, ["628xxx@s.whatsapp.net"], "promote")
conn.UpdateGroupParticipants(jid, ["628xxx@s.whatsapp.net"], "demote")
conn.SetGroupName(jid, "Nama Baru")
conn.SetGroupDescription(jid, "Deskripsi")
conn.SetGroupTopic(jid, prevID, newID, "Topic")
conn.SetGroupPhoto(jid, args)
conn.SetGroupAnnounce(jid, true)
conn.SetGroupLocked(jid, true)
conn.SetDisappearingTimer(jid, 86400)
conn.GetJoinedGroups()
conn.GetGroupRequestParticipants(jid)
conn.UpdateGroupRequestParticipants(jid, participants, "approve")
conn.GetSubGroups(jid)
conn.LinkGroup(parent, child)
conn.UnlinkGroup(parent, child)
```

---

## Profil & Privasi

```js
conn.GetUserInfo(["628xxx@s.whatsapp.net"])
conn.GetProfilePictureInfo(jid, options)
conn.GetBusinessProfile(jid)
conn.SetStatusMessage("Lagi sibuk njir")
conn.SubscribePresence(jid)
conn.GetPrivacySettings()
conn.SetPrivacySetting("last_seen", "all")
conn.GetStatusPrivacy()
conn.UpdateBlocklist(jid, "block")
conn.UpdateBlocklist(jid, "unblock")
conn.GetBlocklist()
conn.IsOnWhatsApp(["628xxx"])
```

---

## Newsletter

```js
conn.GetNewsletterInfo(jid)
conn.GetSubscribedNewsletters()
conn.FollowNewsletter(jid)
conn.UnfollowNewsletter(jid)
conn.NewsletterToggleMute(jid, true)
conn.NewsletterSendReaction(jid, serverID, "👍", messageID)
conn.NewsletterMarkViewed(jid, [serverID])
conn.CreateNewsletter("Nama", "Deskripsi", pictureArgs)
conn.GetNewsletterMessages(jid, count, before)
conn.GetNewsletterMessageUpdates(jid, count)
conn.NewsletterSubscribeLiveUpdates(jid)
```

---

## Voice Call

Pake [meowcaller](https://github.com/purpshell/meowcaller) di baliknya njir nama library-nya aja udah kocak parah wkwk.

```js
conn.Event(async ({ type, evt }) => {
    if (type === "meowcaller.IncomingCall") {
        console.log("Ada telpon njir dari", evt.peer, evt.isVideo ? "(video)" : "(audio)")
        conn.answerCall(evt.callId)
    }

    if (type === "meowcaller.CallReady") {
        conn.playAudio(evt.callId, "audio.mp3")
        conn.receiveAudio(evt.callId, "rekaman.wav")
    }

    if (type === "meowcaller.CallEnd") {
        console.log("Telpon kelar:", evt.reason)
    }
})

const callId = await conn.placeCall("+6281234567890")
conn.hangupCall(callId)
conn.rejectCall(callId, callerJID)
conn.playAudio(callId, "audio.mp3")
conn.playAudio(callId, "video.mp4")
conn.receiveAudio(callId, "rekaman.wav")
conn.receivePCM(callId)
```

Format audio: `.mp3` `.ogg` `.opus` `.wav` — file video (`.mp4` `.mkv` `.avi` dll) langsung dikonvert otomatis ke mp3 via ffmpeg njir gak perlu ribet konvert manual wkwk serius dah.

> ⚠️ Video call (`playVideo`) belum divalidasi meowcaller hadeh jangan ngarep bisa jalan mulus.

---

## Utility

```js
conn.simple(evt)
conn.decodeJid("628xxx:4@s.whatsapp.net")
conn.getDevice(messageID)
conn.ParseMention("@628xxx teks")
conn.generateWAMessageFromContent(jid, content, options)
conn.sendMessage(jid, content, options)
conn.getContentType(content)
```

---

<div align="center">

**Made with 💩**

*tai lu bang*

<sub>readme by <img src="https://www.anthropic.com/favicon.ico" width="12" height="12"> claude</sub>

</div>
