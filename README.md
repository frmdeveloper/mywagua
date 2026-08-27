<div align="center">
<img src="assets/bukankahinimy.jpg" width="30%">

# 🫪 My WA Gua 🫪

### *Whatsmeow tapi Node.js, njir parah*

<sub>`whatsmeow` (Go) → N-API → JavaScript · sqlite / postgres · voice call · goja escape hatch</sub>

</div>

---

## Jadi gini ceritanya njir

[whatsmeow](https://github.com/tulir/whatsmeow) itu library WhatsApp paling waras yang pernah ada.
Masalahnya? **Dia Go.** Sementara isi dunia bot WA tuh Node.js semua. Hadeh, kocak parah.

Mau rewrite protokolnya sendiri? Kagak sanggup gua wkwk. Jadi ya udah, whatsmeow-nya gua bungkus
jadi **native Node addon** lewat N-API. Hasilnya lu nulis JavaScript santuy, tapi yang capek-capek
ngurus enkripsi Signal, websocket, retry, upload media — Go semua. Enak kan njir.

```
┌──────────────────────────────────────────────────┐
│  bot lu (JavaScript)  ← lu main di sini          │
│    index.js  · message.js  · simple.js           │
├──────────────────────────────────────────────────┤
│  main.node   ← jembatan N-API (napi-go)          │
├──────────────────────────────────────────────────┤
│  whatsmeow · meowcaller · goja  (Go, yang capek) │
├──────────────────────────────────────────────────┤
│  sqlite3  /  postgres (pgx)                      │
└──────────────────────────────────────────────────┘
```

### Kenapa lu harusnya pake ini

| | |
|---|---|
| ⚡ **Kenceng anjir** | Crypto & socket digarap Go, bukan JS. Event queue-nya lock-protected, aman. |
| 🔀 **Multi-akun** | Satu proses, satu database, bot sebanyak yang lu mau. `GetAllDevices()` udah kelar. |
| 📞 **Voice call** | Beneran bisa nerima/nelpon + muter audio, pake [meowcaller](https://github.com/purpshell/meowcaller). |
| 🖼️ **Thumbnail otomatis** | Foto & video digenerate thumbnail-nya sendiri, lu gak usah mikir. |
| 🚪 **Ada pintu darurat** | Method whatsmeow belum dibungkus? Sikat langsung lewat `conn.run()` njir. |
| 🗄️ **SQLite / Postgres** | Sesi disimpen di DB. Bukan folder JSON yang dikit-dikit korup terus hadeh. |

---

## Bahan-bahan dulu njir

| | Wajib? | Buat apaan |
|---|---|---|
| **Node.js 20+** | wajib | ESM, `Readable.fromWeb`, fetch bawaan |
| **Go 1.26+** | wajib | `main.node` dicompile pas install |
| **ffmpeg** | boleh skip | thumbnail video + konvert video→audio buat voice call |
| **jimp** | boleh skip | generate JPEG thumbnail foto/video |

> Gak ada ffmpeg/jimp tetep jalan kok njir, cuma pesan medianya jadi gak ada thumbnail-nya. Santuy.

---

## Gas install

```bash
npm install mywagua@github:frmdeveloper/mywagua
```

`postinstall` bakal jalanin `go get -u && go build -buildmode=c-shared -o main.node .` — jadi install
pertama emang lama njir sabar dulu, dia compile whatsmeow dari nol. Ngopi aja dulu.

Mau build ulang manual:

```bash
npm run build
```

<details>
<summary><b>Install-nya error? Buka nih</b></summary>

- `go: command not found` → Go-nya belum keinstall, atau gak masuk `PATH` njir.
- `gcc failed` → butuh C toolchain. `build-essential` di Debian, `clang` di Termux.
- `main.node` gede banget ~50MB → **normal njir santuy**, itu whatsmeow + sqlite + goja
  distatic-link jadi satu.

</details>

---

## Langsung gas — bot minimal

Ini beneran jalan, copas aja:

```js
import { Container, makeClient } from "mywagua"

// 1. buka database sesi
const db = Container("sqlite3", "file:bot.db?_foreign_keys=on", "INFO")

// 2. ambil device — kalo DB masih kosong, dia bikin baru sendiri
const device = db.GetFirstDevice()

// 3. bikin client-nya
const conn = makeClient(device, {
    Logger: { Client: "INFO" },
    OsName: "Chrome"
})

conn.Connect()

// 4. belum login? pairing dulu njir
if (!conn.IsLoggedIn()) {
    await conn.WaitForConnection(10)
    console.log("nih kode pairing-nya:", conn.PairPhone("6281234567890"))
}

// 5. dengerin event
conn.Event(async ({ type, evt }) => {
    if (type === "*events.Connected") return console.log("nyambung njir finally wkwk")

    if (type === "*events.Message") {
        const m = await conn.simple(conn, evt)
        console.log(`[${m.pushname}] ${m.text}`)

        if (m.command === ".ping") {
            await conn.sendMessage(m.from, { text: "pong njir" }, { quoted: evt })
        }
    }
})
```

Udah. Segitu doang buat bot yang bisa bales njir wkwk.

---

## Container — tempat nyimpen sesi

Ini `sqlstore.New`-nya whatsmeow, tapi versi JS biar lu gak pusing.

```js
Container(driver?, dsn?, logLevel?)
```

| Param | Kalo dikosongin | Isinya apa |
|---|---|---|
| `driver` | `"sqlite3"` | `"sqlite3"` atau `"pgx"` |
| `dsn` | `"file:mywagua.db?_foreign_keys=on"` | connection string |
| `logLevel` | `""` = diem total | `"DEBUG"` `"INFO"` `"WARN"` `"ERROR"` |

```js
Container()                                                          // sqlite, path default
Container("sqlite3", "file:bot.db?_foreign_keys=on", "INFO")
Container("pgx", "postgres://user:pass@localhost:5432/dbname", "DEBUG")
```

> `_foreign_keys=on` jangan dilupain njir. Tanpa itu row sesi bisa nyangkut pas device dihapus,
> terus lu bingung sendiri kenapa DB-nya makin gendut hadeh.

### Method-methodnya

```js
db.GetFirstDevice()   // device pertama. DB kosong? dia bikinin yang baru
db.GetAllDevices()    // array semua device
db.GetDevice(jid)     // cari pake JID, null kalo gak nemu
db.PutDevice()        // bikin device baru — buat nambah akun
db.DeleteDevice(dev)  // hapus sesi. terima object device atau handle-nya, bebas
```

Semuanya balikin bentuk yang sama:

```js
const device = db.GetFirstDevice()
// { handle: '1', dbPath: 'bot.db', jid: '6281234567890:4@s.whatsapp.net' }
```

`handle` itu ID internal (string angka) yang nunjuk ke `*store.Device` di sisi Go — anggep aja
tiket parkir njir. `jid` kosong (`''`) artinya device-nya belum pernah login.

---

## makeClient — bikin socket-nya

```js
const conn = makeClient(device, config)
```

`device` boleh object dari `db.*`, boleh handle string-nya doang. Bebas.

```js
makeClient(device, {
    Logger: {
        Client: "INFO",       // level log. kosongin kalo mau diem total
        Color:  true,         // warna ANSI, cuma ngaruh kalo output ke stdout
        File:   "wa.log"      // tulis ke file, bukan nyampah di terminal
    },
    OsName: "Chrome"          // nama device yang nongol di Linked Devices HP lu
})
```

Default bawaan: `OsName: "Chrome"`, platform type `CHROME`, dan
`AutomaticMessageRerequestFromPhone` udah **nyala** — jadi pesan "Waiting for this message" auto
di-request ulang dari HP njir enak banget gak usah mikir.

### Banyak bot satu proses, gaskeun

```js
const db      = Container("pgx", "postgres://user:pass@localhost/db")
const devices = db.GetAllDevices()

const bots = devices.map(d => {
    const conn = makeClient(d, { OsName: "Chrome" })
    conn.Connect()
    return conn
})

// mau pilih akun tertentu? gampang
const botA = makeClient(devices.find(d => d.jid.startsWith("62851")), config)
```

Tiap client punya event queue masing-masing, jadi gak bakal nyampur njir aman.

---

## Konek & login

```js
conn.Connect()                     // konek, non-blocking
conn.ConnectContext()              // konek pake context
conn.Disconnect()                  // putus, tapi sesi masih aman
conn.IsConnected()                 // → boolean
conn.IsLoggedIn()                  // → boolean
conn.WaitForConnection(10)         // tungguin max 10 detik
conn.ResetConnection()             // reset socket paksa kalo nyangkut
conn.Logout()                      // unlink device — SESI HANGUS njir hati-hati banget
conn.PairPhone("6281234567890")    // → kode 8 digit buat "Link with phone number"
conn.Store()                       // → object store (ID, LID, push name, dll)
conn.SetPassive(true)              // mode kalem, gak ngirim receipt
```

Alur pairing pertama kali:

```js
conn.Connect()
await conn.WaitForConnection(15)
if (!conn.IsLoggedIn()) {
    const code = conn.PairPhone("6281234567890")
    console.log("buruan masukin kode ini di HP:", code)
}
```

---

## Event

```js
const timer = conn.Event(async ({ type, evt }) => {
    // ...
})

clearInterval(timer)   // udahan, stop
```

Di baliknya: Go nampung semua event di queue, JS-nya nge-poll tiap **100ms** terus dikuras.
Jadi callback lu **gak pernah** dipanggil dari thread Go — aman njir, zero race condition.

Pesan kosongan (protocol message, sender-key doang, dll) otomatis di-skip biar handler lu gak
kebanjiran sampah wkwk.

### Event whatsmeow

| Type | Kapan nongol |
|---|---|
| `*events.Connected` | Nyambung njir finally |
| `*events.Disconnected` | Putus hadeh |
| `*events.LoggedOut` | Kena unlink / banned. RIP |
| `*events.PairSuccess` | Pairing sukses, gacor |
| `*events.QR` | Kode QR keluar |
| `*events.Message` | Ada pesan masuk |
| `*events.Receipt` | Read / delivery receipt |
| `*events.Presence` | Online / offline |
| `*events.ChatPresence` | Lagi ngetik / lagi rekam VN |
| `*events.CallOffer` | Ada telpon masuk (versi mentah) |
| `*events.GroupInfo` | Grup diobok-obok (nama, member, setting) |
| `*events.JoinedGroup` | Baru dimasukin grup |
| `*events.Contact` / `*events.PushName` | Kontak keupdate |

Semua tipe event whatsmeow lewat sini njir, `type`-nya persis nama Go-nya
(soalnya diambil dari `fmt.Sprintf("%T")` wkwk males mikir).

### Event voice call

| Type | isi `evt` |
|---|---|
| `meowcaller.IncomingCall` | `{ callId, peer, isVideo }` |
| `meowcaller.CallReady` | `{ callId, peer }` — media nyambung, udah siap muter audio |
| `meowcaller.CallStateChange` | `{ callId, phase }` |
| `meowcaller.CallEnd` | `{ callId, reason }` |
| `meowcaller.AudioFrame` | `{ callId, pcm }` — array float32, cuma keluar kalo `receivePCM()` dipanggil |

---

## conn.simple() — biar gak pusing baca pesan

Event `*events.Message` mentahnya berlapis-lapis kayak bawang njir bikin nangis. `simple()` ngeratain
jadi satu object yang enak dipake:

```js
const m = await conn.simple(conn, evt)
```

> ⚠️ **Perhatiin njir**: `conn` diketik dua kali, itu bukan typo. Fungsi aslinya `swmeow(conn, m)`
> dan dia gak dibind, jadi argumen pertama wajib `conn`. Iya gua tau agak kocak wkwk.

| Field | Isinya |
|---|---|
| `m.id` | Message ID |
| `m.from` | JID chat-nya (grup atau personal), sufiks device udah dibuang |
| `m.sender` | JID pengirim asli (`SenderAlt` diutamain, buat jaman LID) |
| `m.lid` | LID pengirim |
| `m.fromMe` | Pesan dari diri sendiri? |
| `m.isGroup` | Dari grup? |
| `m.pushname` | Nama tampilan pengirim |
| `m.key` | `{ remoteJid, id, fromMe, participant? }` — gayanya mirip Baileys biar lu gak kaget |
| `m.type` | Tipe kontennya, mis. `conversation`, `imageMessage` |
| `m.msg` | Isi pesan sesuai tipenya |
| `m.text` | Teks / caption / selectedId, udah dirapiin |
| `m.prefix` | Prefix command kalo ada (`! # % . /` `\`) |
| `m.command` | Kata pertama, lowercase — mis. `.ping` |
| `m.args` | Sisa katanya jadi array |
| `m.q` | Sisa katanya jadi satu string |
| `m.mentionedJid` | Array JID yang di-tag |
| `m.quoted` | Object pesan yang dibales (lihat bawah) |
| `m.full` | Event mentahnya, kalo lu butuh sesuatu yang gak dipetakan |

Lapisan yang otomatis dibongkar: `viewOnceMessageV2`, `documentWithCaptionMessage`, `editedMessage`,
`deviceSentMessage`. Jadi pesan sekali-lihat pun kebaca njir wkwk.

### m.quoted

```js
if (m.quoted.id) {
    m.quoted.type          // tipe pesan yang dibales
    m.quoted.text          // teks / caption-nya
    m.quoted.sender        // siapa yang nulis
    m.quoted.fromMe        // itu pesan kita sendiri?
    m.quoted.key           // key lengkap, siap dipake buat revoke/react
    m.quoted.mentionedJid
    m.quoted.groupMentions
    m.quoted.full          // mentahan
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
    text: "Woy @628123456789 sini bre",
    mentions: ["628123456789@s.whatsapp.net"]
})

// bales pesan — passing event MENTAHNYA njir, bukan hasil simple()
await conn.sendMessage(jid, { text: "Nih balesannya" }, { quoted: evt })
```

Kalo teks polos tanpa mention/quoted, dia dikirim sebagai `conversation` (bukan
`extendedTextMessage`) — biar persis kayak WA asli njir, gak keliatan bot wkwk.

### Media

Tiga cara ngasih sumbernya, bebas mau yang mana:

```js
// dari URL
await conn.sendMessage(jid, { image: { url: "https://example.com/a.jpg" }, caption: "Nih foto" })

// dari file lokal
await conn.sendMessage(jid, { video: { url: "./clip.mp4" }, caption: "Video njir" })

// dari Buffer
await conn.sendMessage(jid, { image: readFileSync("./a.jpg") })

// dari base64
await conn.sendMessage(jid, { audio: { base64: "SUQzB..." } })
```

Yang didukung: `image` · `video` · `audio` · `document` · `sticker`

```js
await conn.sendMessage(jid, {
    document: { url: "./laporan.pdf" },
    fileName: "Laporan Q3.pdf",
    mimetype: "application/pdf"
})

await conn.sendMessage(jid, { sticker: { url: "./stiker.webp" } })
```

Bumbu tambahan yang bisa lu tempel: `caption` · `mimetype` · `fileName` · `mentions` ·
`contextInfo`

<details>
<summary><b>Yang sebenernya kejadian di balik layar (buka kalo penasaran)</b></summary>

1. **Buffer ≤ 5MB** → dikonvert ke base64, dikirim inline.
   **Buffer > 5MB** → ditulis dulu ke `assets/<random>`, dibaca sebagai file, terus dihapus abis
   upload. Makanya folder `assets/` **wajib ada** dan writable njir.
2. **Thumbnail** — foto & video digenerate thumbnail JPEG max 200×200 quality 50 (pake jimp).
   Video diambil frame pertamanya pake ffmpeg. `width`/`height` aslinya ikut dikirim.
3. **Hemat kuota** — kalo medianya dari URL dan ukurannya ≤20MB, byte-nya dipake ulang buat upload,
   jadi gak download dua kali njir pinter kan wkwk. Lebih dari itu cuma 512KB pertama yang diambil
   buat ngintip thumbnail.
4. **Mimetype** otomatis: `image/jpeg`, `video/mp4`, `audio/mpeg`, `image/webp`.
5. Thumbnail gagal (jimp/ffmpeg gak ada) → cuma di-log `thumbnail skipped:` doang, pesannya **tetep
   kekirim** njir santuy gak error.

</details>

### Edit & hapus

```js
// hapus (revoke)
await conn.sendMessage(jid, { delete: m.key })

// edit
await conn.sendMessage(jid, { edit: m.key, text: "Teks baru njir, yang tadi typo wkwk" })

// atau pake helper
await conn.editMessage(jid, m.key, { text: "Teks baru" })
```

### Reaction, polling, revoke pake builder

Method `Build*` itu cuma **ngerakit** object pesannya doang njir, masih harus di-relay:

```js
conn.relayMessage(jid, conn.BuildReaction(chat, sender, id, "👍"))
conn.relayMessage(jid, conn.BuildPollCreation("Makan apa njir?", ["Nasi", "Mie", "Puasa"], 1))
conn.relayMessage(jid, conn.BuildEdit(chat, id, { conversation: "Diedit hadeh" }))
conn.relayMessage(jid, conn.BuildRevoke(chat, sender, id))
```

Yang lain-lain:

```js
conn.RevokeMessage(chat, sender, id)
conn.BuildMessageKey(chat, sender, id)
conn.BuildUnavailableMessageRequest(chat, sender, id)
conn.GenerateMessageID()                 // formatnya: 28 hex uppercase + "-FRM"
conn.DecryptPollVote(pollMsg, vote)
conn.DecryptReaction(reactionMsg)
conn.SendPeerMessage(message)
conn.ParseWebMessage(chatJID, webMsg)
```

### Presence & receipt — biar keliatan manusia

```js
conn.MarkRead([m.id], Date.now(), m.from, m.sender)   // centang biru

conn.SendChatPresence(jid, "composing", "")           // "sedang menulis..."
conn.SendChatPresence(jid, "composing", "audio")      // "sedang merekam..."
conn.SendChatPresence(jid, "paused", "")              // udah, berhenti

conn.SendPresence("available")                        // online
conn.SendPresence("unavailable")                      // ngilang
conn.SubscribePresence(jid)                           // biar dapet event presence dia
```

---

## Media manual (kalo mau ngoprek sendiri)

```js
conn.Upload({ File: "./a.jpg" }, "WhatsApp Image Keys")
conn.UploadNewsletter({ Url: "https://..." }, "WhatsApp Video Keys")
conn.DownloadAny(message)          // → Buffer
conn.SendMediaRetryReceipt(msgInfo, mediaKey)
conn.FetchStickerPack(id)
```

Sumbernya: `{ File }` `{ Url }` `{ Base64 }` `{ Byte }`.
Tipe key-nya: `"WhatsApp Image Keys"` · `"WhatsApp Video Keys"` · `"WhatsApp Audio Keys"` ·
`"WhatsApp Document Keys"`. Sticker pake key **Image** ya njir, jangan bingung.

---

## Grup

```js
// baca-baca
conn.GetGroupInfo(jid)
conn.GetGroupInfoFromLink(code)
conn.GetGroupInfoFromInvite(inviter, jid, code, expiration)
conn.GetJoinedGroups()
conn.GetGroupInviteLink(jid, false)          // true = reset link lama
conn.GetGroupRequestParticipants(jid)

// bikin & cabut
conn.CreateGroup("Nama Grup", ["628xxx@s.whatsapp.net"])
conn.LeaveGroup(jid)
conn.JoinGroupWithLink(code)
conn.JoinGroupWithInvite(inviter, jid, code, expiration)

// urusan member
conn.UpdateGroupParticipants(jid, jids, "add")             // add | remove | promote | demote
conn.UpdateGroupRequestParticipants(jid, jids, "approve")  // approve | reject

// setelan grup
conn.SetGroupName(jid, "Nama Baru")
conn.SetGroupDescription(jid, "Deskripsi baru")
conn.SetGroupTopic(jid, prevID, newID, "Topic")
conn.SetGroupPhoto(jid, { File: "./pp.jpg" })
conn.SetGroupAnnounce(jid, true)             // true = cuma admin yang boleh ngomong
conn.SetGroupLocked(jid, true)               // true = cuma admin yang boleh edit info
conn.SetGroupJoinApprovalMode(jid, true)
conn.SetGroupMemberAddMode(jid, "admin_add")
conn.SetDisappearingTimer(jid, 86400)        // detik. 0 = matiin
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
conn.IsOnWhatsApp(["628xxx"])                // cek nomornya kedaftar WA apa kagak
conn.GetProfilePictureInfo(jid, { Preview: false })
conn.GetBusinessProfile(jid)
conn.SetStatusMessage("Lagi sibuk njir jangan chat")

conn.GetContactQRLink(false)                 // true = revoke QR lama
conn.ResolveContactQRLink(code)
conn.ResolveBusinessMessageLink(code)

conn.GetPrivacySettings()
conn.TryFetchPrivacySettings(true)
conn.SetPrivacySetting("last_seen", "all")   // all | contacts | contact_blacklist | none
conn.GetStatusPrivacy()

conn.UpdateBlocklist(jid, "block")           // bye
conn.UpdateBlocklist(jid, "unblock")         // yaudah maafin
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

Pake [meowcaller](https://github.com/purpshell/meowcaller) di baliknya — njir nama library-nya aja
udah kocak parah wkwk.

### Nerima telpon

```js
conn.Event(async ({ type, evt }) => {
    if (type === "meowcaller.IncomingCall") {
        console.log("ada telpon njir dari", evt.peer, evt.isVideo ? "(video)" : "(audio)")
        conn.answerCall(evt.callId)
        // males? conn.rejectCall(evt.callId)
    }

    if (type === "meowcaller.CallReady") {
        conn.playAudio(evt.callId, "halo.mp3")       // muterin audio ke penelpon
        conn.receiveAudio(evt.callId, "rekaman.wav") // rekam suara dia ke file wkwk
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

### Method-methodnya

```js
conn.answerCall(callId)
conn.rejectCall(callId)
conn.hangupCall(callId)
conn.playAudio(callId, path)      // .mp3 .ogg .opus .wav — file video auto dikonvert
conn.playVideo(callId, path)
conn.receiveAudio(callId, path)   // rekam ke WAV
conn.receivePCM(callId)           // streaming float32 lewat event meowcaller.AudioFrame

conn.RejectCall(callID, callerJID)  // versi whatsmeow, buat *events.CallOffer mentah
```

**Soal format audio njir.** `.mp3` → decoder MP3, `.ogg`/`.opus` → Opus, sisanya dianggep WAV.
File video (`.mp4` `.mkv` `.avi` `.mov` `.flv` `.webm` `.m4v` `.3gp`) langsung dikonvert otomatis ke
mp3 mono 16kHz via ffmpeg — **gak perlu ribet konvert manual** njir wkwk serius dah. ffmpeg dicari di
`PATH` dulu, gak nemu ya fallback ke path Termux / `/usr/bin` / `/usr/local/bin`.

**Mau proses PCM realtime? Bisa:**

```js
if (type === "meowcaller.CallReady") conn.receivePCM(evt.callId)
if (type === "meowcaller.AudioFrame") {
    // evt.pcm = array float32. lempar ke STT, VAD, terserah lu mau diapain
}
```

> ⚠️ `playVideo` belum divalidasi meowcaller hadeh, jangan ngarep jalan mulus.

---

## 🚪 Pintu darurat — sikat whatsmeow mentah

Nah ini bagian paling gokil njir. Di dalem ada [goja](https://github.com/dop251/goja) (interpreter JS
yang ditulis pake Go) dan dia udah dikasih variabel `client` (`*whatsmeow.Client`) plus `ctx`. Jadi
**method whatsmeow APAPUN** bisa lu panggil, walaupun gua belum bungkus:

```js
// eval JS langsung di sisi Go, gila kan wkwk
conn.run(`client.IsLoggedIn()`)
conn.run(`JSON.stringify(client.Store.ID)`)

// atau versi rapi — argumennya auto di-serialize
conn.Call("GetGroupInfo", conn.ctx, "1234567890@g.us")
```

Mau ngintip ada method apa aja?

```js
console.log(conn.run(`JSON.stringify(Object.keys(client))`))
```

Bonus njir: semua method yang udah dibungkus nyimpen signature aslinya, jadi enak buat ngintip di
REPL:

```js
conn.SetGroupName.toString()
// → function SetGroupName(jid, name) { [native code] }
```

---

## Printilan lain

```js
conn.simple(conn, evt)             // parse pesan jadi object enak
conn.decodeJid("628xxx:4@s.whatsapp.net")   // → "628xxx@s.whatsapp.net"
conn.getDevice(messageID)          // → 'ios' | 'web' | 'android' | 'desktop' | 'unknown'
conn.ParseMention("@628xxx woy")   // → ["628xxx@s.whatsapp.net"]
conn.getContentType(content)       // nyari key tipe kontennya
conn.generateWAMessageFromContent(jid, content, options)
conn.relayMessage(jid, message, options)
conn.Store()                       // store lengkap jadi object JS
conn.GetStore()                    // versi JSON string
conn.RemoveEventHandler(id)
conn.RemoveEventHandlers()
conn.SetForceActiveDeliveryReceipts(true)
conn.SetMaxParallelRetryReceiptHandling(4)
conn.MarkNotDirty(name, timestamp)
```

---

## Yang perlu lu tau (jangan sampe kaget)

- **Polling 100ms.** Delay event maksimal ~100ms. Sengaja gitu njir, biar gak ada rebutan antara
  goroutine Go sama event loop Node. Mending delay dikit daripada crash wkwk.
- **Folder `assets/` wajib ada.** Buffer >5MB nginep sementara di situ.
- **`main.node` gak masuk git.** Udah di `.gitignore` — tiap environment build sendiri-sendiri.
- **`conn.Call()` nyampah ke console.** Kalo berisik, pake `conn.run()` aja langsung.
- **JID vs LID.** WhatsApp lagi pindahan ke LID njir. `simple()` udah nanganin: `SenderAlt` dipake
  buat `sender`, `Sender` dilempar ke `lid`.
- **Error dari Go dilempar jadi JS exception.** Bungkus pake `try/catch` njir, jangan sok jago wkwk.

---

## Mau nyumbang kode?

Nemu bug, atau ada method whatsmeow yang belum kebungkus? Buka issue atau langsung PR aja.
Polanya gampang, di `conn.go`:

```go
reg("NamaMethod", "param1, param2", func(param1 string, param2 bool) any {
    res, err := Cli.NamaMethod(ctx, param1, param2)
    if err != nil { return Throw(env, err) }
    return Res(res)
})
```

`reg()` sekalian nyimpen string parameternya, biar `toString()` di JS-nya bener njir jangan males.

---

## Makasih buat

- [whatsmeow](https://github.com/tulir/whatsmeow) — yang ngurus protokol WhatsApp-nya
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
