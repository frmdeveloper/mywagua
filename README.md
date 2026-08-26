<div align="center">
<img src="assets/bukankahinimy.jpg" width="30%">

# 🫪 My WA Gua 🫪

### *Whatsmeow tapi Node.js, njir parah*

<sub>`whatsmeow` (Go) → N-API → JavaScript · sqlite / postgres · voice call · goja escape hatch</sub>

</div>

---

## Apaan sih ini wkwk

Jadi gini ceritanya njir — [whatsmeow](https://github.com/tulir/whatsmeow) itu library WhatsApp
paling waras yang ada, tapi dia Go. Sementara ekosistem bot WA isinya Node.js semua. Hadeh.

Jadi daripada rewrite protokolnya (kagak sanggup wkwk), library ini ngebungkus whatsmeow jadi
**native Node addon** lewat N-API. Hasilnya: nulis JavaScript, tapi yang ngurus enkripsi Signal,
websocket, retry, media upload semuanya Go.

```
┌──────────────────────────────────────────────────┐
│  bot lu (JavaScript)                             │
│    index.js  · message.js  · simple.js           │
├──────────────────────────────────────────────────┤
│  main.node   ← N-API bridge (napi-go)            │
├──────────────────────────────────────────────────┤
│  whatsmeow · meowcaller · goja  (Go)             │
├──────────────────────────────────────────────────┤
│  sqlite3  /  postgres (pgx)                      │
└──────────────────────────────────────────────────┘
```

### Kenapa mau pake ini

| | |
|---|---|
| ⚡ **Kenceng** | Crypto & socket dikerjain Go, bukan JS. Event queue-nya lock-protected. |
| 🔀 **Multi-akun** | Satu proses, satu database, banyak bot. Tinggal `GetAllDevices()`. |
| 📞 **Voice call** | Beneran bisa terima/telpon + play audio, via [meowcaller](https://github.com/purpshell/meowcaller). |
| 🖼️ **Thumbnail auto** | Image & video di-generate thumbnail-nya sendiri (jimp + ffmpeg). |
| 🚪 **Escape hatch** | Method whatsmeow yang belum dibungkus? Panggil langsung lewat `conn.run()`. |
| 🗄️ **SQLite / Postgres** | Sesi disimpen di DB, bukan folder JSON yang gampang korup. |

---

## Yang harus ada dulu

| | Wajib? | Buat apa |
|---|---|---|
| **Node.js 20+** | wajib | ESM, `Readable.fromWeb`, native fetch |
| **Go 1.26+** | wajib | `main.node` dibuild pas install |
| **ffmpeg** | opsional | thumbnail video + konvert video→audio pas voice call |
| **jimp** | opsional | generate JPEG thumbnail image/video |

> Tanpa ffmpeg/jimp tetep jalan njir, cuma pesan medianya gak ada thumbnail-nya. Gak fatal.

---

## Install

```bash
npm install mywagua@github:frmdeveloper/mywagua
```

`postinstall` bakal jalanin `go get -u && go build -buildmode=c-shared -o main.node .` — jadi
install pertama agak lama njir sabar, dia compile whatsmeow dari nol.

Kalo mau build ulang manual:

```bash
npm run build
```

<details>
<summary><b>Kena error pas install?</b></summary>

- `go: command not found` → Go belum keinstall / gak ada di `PATH`.
- `gcc failed` → butuh C toolchain (`build-essential` di Debian, `clang` di Termux).
- `main.node` gede banget (~50MB) → normal njir, itu whatsmeow + sqlite + goja distatic-link semua.

</details>

---

## Quickstart

Bot minimal yang beneran jalan:

```js
import { Container, makeClient } from "mywagua"

// 1. buka database sesi
const db = Container("sqlite3", "file:bot.db?_foreign_keys=on", "INFO")

// 2. ambil device (atau bikin baru kalo belum ada)
const device = db.GetFirstDevice()

// 3. bikin client
const conn = makeClient(device, {
    Logger: { Client: "INFO" },
    OsName: "Chrome"
})

conn.Connect()

// 4. pairing kalo belum login
if (!conn.IsLoggedIn()) {
    await conn.WaitForConnection(10)
    console.log("kode pairing:", conn.PairPhone("6281234567890"))
}

// 5. dengerin event
conn.Event(async ({ type, evt }) => {
    if (type === "*events.Connected") return console.log("connected njir finally")

    if (type === "*events.Message") {
        const m = await conn.simple(conn, evt)
        console.log(`[${m.pushname}] ${m.text}`)

        if (m.command === ".ping") {
            await conn.sendMessage(m.from, { text: "pong njir" }, { quoted: evt })
        }
    }
})
```

---

## Container — penyimpanan sesi

Ini `sqlstore.New`-nya whatsmeow, versi JS.

```js
Container(driver?, dsn?, logLevel?)
```

| Param | Default | Isi |
|---|---|---|
| `driver` | `"sqlite3"` | `"sqlite3"` atau `"pgx"` |
| `dsn` | `"file:mywagua.db?_foreign_keys=on"` | connection string |
| `logLevel` | `""` (silent) | `"DEBUG"` `"INFO"` `"WARN"` `"ERROR"` |

```js
Container()                                                          // sqlite, default path
Container("sqlite3", "file:bot.db?_foreign_keys=on", "INFO")
Container("pgx", "postgres://user:pass@localhost:5432/dbname", "DEBUG")
```

> `_foreign_keys=on` itu penting njir. Tanpa itu row sesi bisa nyangkut pas device dihapus.

### Method

```js
db.GetFirstDevice()   // device pertama, atau device baru kosong kalo DB masih kosong
db.GetAllDevices()    // array semua device
db.GetDevice(jid)     // cari by JID, null kalo gak ada
db.PutDevice()        // bikin device baru (buat nambah akun)
db.DeleteDevice(dev)  // hapus sesi — terima object device atau handle-nya
```

Semuanya balikin object yang sama bentuknya:

```js
const device = db.GetFirstDevice()
// { handle: '1', dbPath: 'bot.db', jid: '6281234567890:4@s.whatsapp.net' }
```

`handle` itu ID internal (string angka) yang nunjuk ke `*store.Device` di sisi Go. `jid` kosong
(`''`) berarti device-nya belum pernah login.

---

## makeClient — bikin socket

```js
const conn = makeClient(device, config)
```

`device` boleh object dari `db.*`, boleh handle string-nya langsung.

```js
makeClient(device, {
    Logger: {
        Client: "INFO",       // level log; kosongin buat silent total
        Color:  true,         // warna ANSI (cuma kalo output ke stdout)
        File:   "wa.log"      // tulis ke file, bukan stdout
    },
    OsName: "Chrome"          // nama device yang muncul di Linked Devices
})
```

Default: `OsName: "Chrome"`, platform type `CHROME`,
dan `AutomaticMessageRerequestFromPhone` nyala (jadi pesan "Waiting for this message" auto
di-request ulang dari HP njir enak).

### Banyak bot dalam satu proses

```js
const db      = Container("pgx", "postgres://user:pass@localhost/db")
const devices = db.GetAllDevices()

const bots = devices.map(d => {
    const conn = makeClient(d, { OsName: "Chrome" })
    conn.Connect()
    return conn
})

// atau pilih akun spesifik
const botA = makeClient(devices.find(d => d.jid.startsWith("62851")), config)
```

Tiap client punya event queue sendiri, jadi gak nyampur njir aman.

---

## Koneksi & login

```js
conn.Connect()                     // konek (non-blocking)
conn.ConnectContext()              // konek pake context
conn.Disconnect()                  // putus, sesi tetep ada
conn.IsConnected()                 // → boolean
conn.IsLoggedIn()                  // → boolean
conn.WaitForConnection(10)         // tunggu max 10 detik
conn.ResetConnection()             // reset socket paksa
conn.Logout()                      // unlink device, sesi HANGUS njir hati-hati
conn.PairPhone("6281234567890")    // → kode 8 digit buat "Link with phone number"
conn.Store()                       // → object store (ID, LID, push name, dll)
conn.SetPassive(true)              // mode pasif, gak ngirim receipt
```

Alur pairing pertama kali:

```js
conn.Connect()
await conn.WaitForConnection(15)
if (!conn.IsLoggedIn()) {
    const code = conn.PairPhone("6281234567890")
    console.log("masukin kode ini di HP:", code)
}
```

---

## Event

```js
const timer = conn.Event(async ({ type, evt }) => {
    // ...
})

clearInterval(timer)   // stop listener
```

Di baliknya: Go nampung semua event di queue, JS nge-poll tiap **100ms** terus nge-drain.
Jadi callback lu gak pernah dipanggil dari thread Go — aman njir, gak ada race.

Pesan tanpa isi (protocol message, sender-key doang, dll) otomatis di-skip biar handler lu
gak kebanjiran sampah.

### Event whatsmeow

| Type | Kapan |
|---|---|
| `*events.Connected` | Nyambung njir finally |
| `*events.Disconnected` | Putus hadeh |
| `*events.LoggedOut` | Kena unlink / banned |
| `*events.PairSuccess` | Pairing berhasil |
| `*events.QR` | Kode QR keluar |
| `*events.Message` | Pesan masuk |
| `*events.Receipt` | Read / delivery receipt |
| `*events.Presence` | Online / offline |
| `*events.ChatPresence` | Lagi ngetik / rekam voice note |
| `*events.CallOffer` | Ada telpon masuk (event mentah) |
| `*events.GroupInfo` | Grup diubah (nama, member, setting) |
| `*events.JoinedGroup` | Baru masuk grup |
| `*events.Contact` / `*events.PushName` | Kontak keupdate |

Semua tipe event whatsmeow lewat sini njir, `type`-nya persis nama Go-nya (`fmt.Sprintf("%T")`).

### Event voice call

| Type | `evt` |
|---|---|
| `meowcaller.IncomingCall` | `{ callId, peer, isVideo }` |
| `meowcaller.CallReady` | `{ callId, peer }` — media udah nyambung, siap play audio |
| `meowcaller.CallStateChange` | `{ callId, phase }` |
| `meowcaller.CallEnd` | `{ callId, reason }` |
| `meowcaller.AudioFrame` | `{ callId, pcm }` — float32 array, cuma kalo `receivePCM()` dipanggil |

---

## conn.simple() — parse pesan

Event `*events.Message` mentahnya berlapis-lapis njir pusing. `simple()` ngeratain jadi object
yang enak dipake:

```js
const m = await conn.simple(conn, evt)
```

> ⚠️ Perhatiin njir: `conn` dipassing dua kali. Fungsinya `swmeow(conn, m)` dan dia gak dibind,
> jadi arg pertama harus `conn`.

| Field | Isi |
|---|---|
| `m.id` | Message ID |
| `m.from` | JID chat (grup atau personal), sufiks device dibuang |
| `m.sender` | JID pengirim asli (`SenderAlt` diprioritasin, buat era LID) |
| `m.lid` | LID pengirim |
| `m.fromMe` | Pesan dari diri sendiri? |
| `m.isGroup` | Dari grup? |
| `m.pushname` | Nama tampilan pengirim |
| `m.key` | `{ remoteJid, id, fromMe, participant? }` — kompatibel gaya Baileys |
| `m.type` | Tipe konten, mis. `conversation`, `imageMessage` |
| `m.msg` | Isi pesan sesuai tipenya |
| `m.text` | Teks / caption / selectedId, udah dinormalisasi |
| `m.prefix` | Prefix command kalo ada (`! # % . /` `\`) |
| `m.command` | Kata pertama, lowercase — mis. `.ping` |
| `m.args` | Sisa kata jadi array |
| `m.q` | Sisa kata jadi satu string |
| `m.mentionedJid` | Array JID yang di-mention |
| `m.quoted` | Object pesan yang dibalas (lihat bawah) |
| `m.full` | Event mentahnya, kalo butuh sesuatu yang gak dipetakan |

Lapisan yang otomatis dibuka: `viewOnceMessageV2`, `documentWithCaptionMessage`,
`editedMessage`, `deviceSentMessage`.

### m.quoted

```js
if (m.quoted.id) {
    m.quoted.type          // tipe pesan yang dibalas
    m.quoted.text          // teks / caption-nya
    m.quoted.sender        // siapa yang nulis
    m.quoted.fromMe        // itu pesan kita sendiri?
    m.quoted.key           // key lengkap, siap dipake buat revoke/react
    m.quoted.mentionedJid
    m.quoted.groupMentions
    m.quoted.full          // raw
}
```

---

## Kirim pesan

```js
await conn.sendMessage(jid, content, options)
```

### Teks

```js
await conn.sendMessage(jid, { text: "Halo njir!" })

await conn.sendMessage(jid, {
    text: "Woy @628123456789 sini",
    mentions: ["628123456789@s.whatsapp.net"]
})

// balas pesan — passing event mentahnya, bukan hasil simple()
await conn.sendMessage(jid, { text: "Nih balesannya" }, { quoted: evt })
```

Kalo teks polos tanpa mention/quoted, dia dikirim sebagai `conversation` (bukan
`extendedTextMessage`) — biar sama kayak WhatsApp asli njir.

### Media

Tiga bentuk sumber, semuanya boleh:

```js
// dari URL
await conn.sendMessage(jid, { image: { url: "https://example.com/a.jpg" }, caption: "Foto" })

// dari file lokal
await conn.sendMessage(jid, { video: { url: "./clip.mp4" }, caption: "Video njir" })

// dari Buffer
await conn.sendMessage(jid, { image: readFileSync("./a.jpg") })

// dari base64
await conn.sendMessage(jid, { audio: { base64: "SUQzB..." } })
```

Tipe yang didukung: `image` · `video` · `audio` · `document` · `sticker`

```js
await conn.sendMessage(jid, {
    document: { url: "./laporan.pdf" },
    fileName: "Laporan Q3.pdf",
    mimetype: "application/pdf"
})

await conn.sendMessage(jid, { sticker: { url: "./stiker.webp" } })
```

Opsi konten yang bisa ditambahin: `caption` · `mimetype` · `fileName` · `mentions` ·
`contextInfo`

<details>
<summary><b>Yang terjadi di balik layar</b></summary>

1. **Buffer ≤ 5MB** → dikonvert ke base64 dan dikirim inline.
   **Buffer > 5MB** → ditulis dulu ke `assets/<random>` terus dibaca sebagai file, dihapus
   setelah upload. Jadi folder `assets/` harus ada dan writable njir.
2. **Thumbnail** — image & video digenerate thumbnail JPEG max 200×200 quality 50 (pake jimp).
   Video diambil frame pertamanya pake ffmpeg. `width`/`height` asli ikut dikirim.
3. **Optimasi download** — kalo medianya dari URL dan ukurannya ≤20MB, byte-nya dipake ulang
   buat upload, jadi gak download dua kali. Lebih dari itu cuma 512KB pertama yang diambil
   buat probe thumbnail.
4. **Mimetype** otomatis: `image/jpeg`, `video/mp4`, `audio/mpeg`, `image/webp`.
5. Thumbnail gagal (jimp/ffmpeg gak ada) → cuma di-log `thumbnail skipped:`, pesannya tetep
   kekirim.

</details>

### Edit & hapus

```js
// hapus (revoke)
await conn.sendMessage(jid, { delete: m.key })

// edit
await conn.sendMessage(jid, { edit: m.key, text: "Teks baru njir" })

// atau lewat helper
await conn.editMessage(jid, m.key, { text: "Teks baru" })
```

### Reaction, poll, revoke lewat builder

Method `Build*` cuma **bikin** object pesannya — masih harus di-relay:

```js
conn.relayMessage(jid, conn.BuildReaction(chat, sender, id, "👍"))
conn.relayMessage(jid, conn.BuildPollCreation("Makan apa njir?", ["Nasi", "Mie", "Gak makan"], 1))
conn.relayMessage(jid, conn.BuildEdit(chat, id, { conversation: "Diedit hadeh" }))
conn.relayMessage(jid, conn.BuildRevoke(chat, sender, id))
```

Lainnya:

```js
conn.RevokeMessage(chat, sender, id)
conn.BuildMessageKey(chat, sender, id)
conn.BuildUnavailableMessageRequest(chat, sender, id)
conn.GenerateMessageID()                 // format: 28 hex uppercase + "-FRM"
conn.DecryptPollVote(pollMsg, vote)
conn.DecryptReaction(reactionMsg)
conn.SendPeerMessage(message)
conn.ParseWebMessage(chatJID, webMsg)
```

### Presence & receipt

```js
conn.MarkRead([m.id], Date.now(), m.from, m.sender)   // centang biru

conn.SendChatPresence(jid, "composing", "")           // "sedang menulis..."
conn.SendChatPresence(jid, "composing", "audio")      // "sedang merekam..."
conn.SendChatPresence(jid, "paused", "")              // berhenti

conn.SendPresence("available")                        // online
conn.SendPresence("unavailable")                      // offline
conn.SubscribePresence(jid)                           // biar dapet event presence-nya
```

---

## Media manual

```js
conn.Upload({ File: "./a.jpg" }, "WhatsApp Image Keys")
conn.UploadNewsletter({ Url: "https://..." }, "WhatsApp Video Keys")
conn.DownloadAny(message)          // → Buffer
conn.SendMediaRetryReceipt(msgInfo, mediaKey)
conn.FetchStickerPack(id)
```

Sumber: `{ File }` `{ Url }` `{ Base64 }` `{ Byte }`.
Tipe key: `"WhatsApp Image Keys"` · `"WhatsApp Video Keys"` · `"WhatsApp Audio Keys"` ·
`"WhatsApp Document Keys"`.

---

## Grup

```js
// baca
conn.GetGroupInfo(jid)
conn.GetGroupInfoFromLink(code)
conn.GetGroupInfoFromInvite(inviter, jid, code, expiration)
conn.GetJoinedGroups()
conn.GetGroupInviteLink(jid, false)          // true = reset link lama
conn.GetGroupRequestParticipants(jid)

// bikin & keluar
conn.CreateGroup("Nama Grup", ["628xxx@s.whatsapp.net"])
conn.LeaveGroup(jid)
conn.JoinGroupWithLink(code)
conn.JoinGroupWithInvite(inviter, jid, code, expiration)

// member
conn.UpdateGroupParticipants(jid, jids, "add")       // add | remove | promote | demote
conn.UpdateGroupRequestParticipants(jid, jids, "approve")   // approve | reject

// setting
conn.SetGroupName(jid, "Nama Baru")
conn.SetGroupDescription(jid, "Deskripsi baru")
conn.SetGroupTopic(jid, prevID, newID, "Topic")
conn.SetGroupPhoto(jid, { File: "./pp.jpg" })
conn.SetGroupAnnounce(jid, true)             // true = cuma admin yang bisa kirim
conn.SetGroupLocked(jid, true)               // true = cuma admin yang bisa edit info
conn.SetGroupJoinApprovalMode(jid, true)
conn.SetGroupMemberAddMode(jid, "admin_add")
conn.SetDisappearingTimer(jid, 86400)        // detik; 0 = matiin
conn.SetDefaultDisappearingTimer(86400)

// community
conn.GetSubGroups(jid)
conn.GetLinkedGroupsParticipants(jid)
conn.LinkGroup(parent, child)
conn.UnlinkGroup(parent, child)
```

---

## Profil, kontak & privasi

```js
conn.GetUserInfo(["628xxx@s.whatsapp.net"])
conn.GetUserDevices(jids)
conn.IsOnWhatsApp(["628xxx"])                // cek nomor kedaftar WA apa kagak
conn.GetProfilePictureInfo(jid, { Preview: false })
conn.GetBusinessProfile(jid)
conn.SetStatusMessage("Lagi sibuk njir")

conn.GetContactQRLink(false)                 // true = revoke QR lama
conn.ResolveContactQRLink(code)
conn.ResolveBusinessMessageLink(code)

conn.GetPrivacySettings()
conn.TryFetchPrivacySettings(true)
conn.SetPrivacySetting("last_seen", "all")   // all | contacts | contact_blacklist | none
conn.GetStatusPrivacy()

conn.UpdateBlocklist(jid, "block")
conn.UpdateBlocklist(jid, "unblock")
conn.GetBlocklist()

conn.GetBotListV2()
conn.GetBotProfiles(jids)
conn.StoreLIDPNMapping(lid, pn)
conn.AcceptTOSNotice(stage, privacyActtoken)
```

---

## Newsletter (Channel)

```js
conn.CreateNewsletter("Nama", "Deskripsi", { File: "./pp.jpg" })
conn.GetNewsletterInfo(jid)
conn.GetNewsletterInfoWithInvite(key)
conn.GetSubscribedNewsletters()
conn.FollowNewsletter(jid)
conn.UnfollowNewsletter(jid)
conn.NewsletterToggleMute(jid, true)
conn.NewsletterSendReaction(jid, serverID, "👍", messageID)
conn.NewsletterMarkViewed(jid, [serverID])
conn.GetNewsletterMessages(jid, 50, 0)       // count, before
conn.GetNewsletterMessageUpdates(jid, 50)
conn.NewsletterSubscribeLiveUpdates(jid)
```

---

## Voice Call

Pake [meowcaller](https://github.com/purpshell/meowcaller) di baliknya — njir nama library-nya
aja udah kocak parah wkwk.

### Terima telpon

```js
conn.Event(async ({ type, evt }) => {
    if (type === "meowcaller.IncomingCall") {
        console.log("telpon masuk dari", evt.peer, evt.isVideo ? "(video)" : "(audio)")
        conn.answerCall(evt.callId)
        // atau: conn.rejectCall(evt.callId)
    }

    if (type === "meowcaller.CallReady") {
        conn.playAudio(evt.callId, "halo.mp3")      // putar audio ke penelpon
        conn.receiveAudio(evt.callId, "rekaman.wav") // rekam suara dia ke file
    }

    if (type === "meowcaller.CallEnd") {
        console.log("telpon kelar:", evt.reason)
    }
})
```

### Nelpon keluar

```js
const callId = await conn.placeCall("+6281234567890")   // ini async njir, Promise
conn.hangupCall(callId)
```

### Method

```js
conn.answerCall(callId)
conn.rejectCall(callId)
conn.hangupCall(callId)
conn.playAudio(callId, path)      // .mp3 .ogg .opus .wav — file video auto-konvert
conn.playVideo(callId, path)
conn.receiveAudio(callId, path)   // rekam ke WAV
conn.receivePCM(callId)           // stream float32 lewat event meowcaller.AudioFrame
```

**Format audio.** `.mp3` → MP3 decoder, `.ogg`/`.opus` → Opus, sisanya diperlakuin sebagai WAV.
File video (`.mp4` `.mkv` `.avi` `.mov` `.flv` `.webm` `.m4v` `.3gp`) langsung dikonvert otomatis
ke mp3 mono 16kHz via ffmpeg njir gak perlu ribet konvert manual wkwk serius dah. ffmpeg dicari di
`PATH` dulu, terus fallback ke path Termux / `/usr/bin` / `/usr/local/bin`.

**Proses PCM realtime:**

```js
if (type === "meowcaller.CallReady") conn.receivePCM(evt.callId)
if (type === "meowcaller.AudioFrame") {
    // evt.pcm = array float32, kirim ke STT / VAD / apa aja
}
```

> ⚠️ `playVideo` belum divalidasi meowcaller hadeh jangan ngarep jalan mulus.

---

## 🚪 Escape hatch — akses whatsmeow mentah

Ini bagian paling gokil njir. Di dalem ada [goja](https://github.com/dop251/goja) (JS interpreter
di Go) yang udah dikasih variabel `client` (`*whatsmeow.Client`) dan `ctx`. Jadi **method whatsmeow
apapun** bisa dipanggil walaupun belum dibungkus:

```js
// eval JS langsung di sisi Go
conn.run(`client.IsLoggedIn()`)
conn.run(`JSON.stringify(client.Store.ID)`)

// atau versi terstruktur — argumen otomatis di-serialize
conn.Call("GetGroupInfo", conn.ctx, "1234567890@g.us")
```

Cek method apa aja yang ada:

```js
console.log(conn.run(`JSON.stringify(Object.keys(client))`))
```

Bonus: semua method yang dibungkus punya signature asli pas di-`toString()`, jadi enak buat
introspeksi di REPL:

```js
conn.SetGroupName.toString()
// → function SetGroupName(jid, name) { [native code] }
```

---

## Utility

```js
conn.simple(conn, evt)             // parse pesan jadi object enak
conn.decodeJid("628xxx:4@s.whatsapp.net")   // → "628xxx@s.whatsapp.net"
conn.getDevice(messageID)          // → 'ios' | 'web' | 'android' | 'desktop' | 'unknown'
conn.ParseMention("@628xxx woy")   // → ["628xxx@s.whatsapp.net"]
conn.getContentType(content)       // cari key tipe konten
conn.generateWAMessageFromContent(jid, content, options)
conn.relayMessage(jid, message, options)
conn.Store()                       // store lengkap sebagai object JS
conn.GetStore()                    // versi JSON string
conn.RemoveEventHandlers()
conn.SetForceActiveDeliveryReceipts(true)
conn.SetMaxParallelRetryReceiptHandling(4)
conn.MarkNotDirty(name, timestamp)
```

---

## Catatan & batasan

- **Event polling 100ms.** Latency event maksimal ~100ms. Sengaja gitu biar aman dari race
  antara goroutine Go dan event loop Node njir.
- **Folder `assets/` harus ada.** Buffer >5MB ditulis sementara di situ.
- **`main.node` gak masuk git.** Ada di `.gitignore` — tiap environment build sendiri.
- **`conn.Call()` nge-log ke console.** Kalo berisik, pake `conn.run()` langsung.
- **JID vs LID.** WhatsApp lagi migrasi ke LID. `simple()` udah handle: `SenderAlt` dipake buat
  `sender`, `Sender` masuk ke `lid`.
- **Error dari Go dilempar jadi JS exception.** Bungkus pake `try/catch` njir jangan sok jago.

---

## Kontribusi

Nemu bug atau ada method whatsmeow yang belum dibungkus? Buka issue atau PR.
Pola nambah method di `conn.go`:

```go
reg("NamaMethod", "param1, param2", func(param1 string, param2 bool) any {
    res, err := Cli.NamaMethod(ctx, param1, param2)
    if err != nil { return Throw(env, err) }
    return Res(res)
})
```

`reg()` sekalian nyimpen string parameternya biar `toString()` di JS-nya bener njir.

---

## Credits

- [whatsmeow](https://github.com/tulir/whatsmeow) — protokol WhatsApp-nya
- [meowcaller](https://github.com/purpshell/meowcaller) — voice call
- [napi-go](https://sirherobrine23.com.br/Sirherobrine23/napi-go) — jembatan N-API
- [goja](https://github.com/dop251/goja) — JS runtime di Go
- [jimp](https://github.com/jimp-dev/jimp) — thumbnail

---

<div align="center">

**Made with 💩**

*tai lu bang*

<sub>readme by <img src="https://www.anthropic.com/favicon.ico" width="12" height="12"> claude</sub>

</div>
