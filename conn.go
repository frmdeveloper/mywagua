package main
import (
    "context"
    "crypto/rand"
    "encoding/base64"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "google.golang.org/protobuf/proto"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types"
    "go.mau.fi/whatsmeow/types/events"
    "io/ioutil"
    "net/http"
    "os"
    "regexp"
    "strings"
    "sirherobrine23.com.br/Sirherobrine23/napi-go"
    waProto "go.mau.fi/whatsmeow/binary/proto"
)

type L struct {
    Edit string
    Caption string
    Mentions []string
    ParseMention bool
    Quoted *events.Message
    ContextInfo *waProto.ContextInfo
    Url string
    Base64 string
    File string
    Byte []byte
    Text string
    String string
    Mimetype *string
    Type whatsmeow.MediaType
}
var MediaType = J{
    "WhatsApp Image Keys": "ImageMessage",
    "WhatsApp Video Keys": "VideoMessage",
    "WhatsApp Audio Keys": "AudioMessage",
    "WhatsApp Document Keys": "DocumentMessage",
}
func ToJson(ani any) string {
    jsons,_ := json.Marshal(ani)
    return string(jsons)
}
func Res(a any) interface{} {
    var jsonData interface{}
    ua := json.Unmarshal([]byte(ToJson(a)), &jsonData)
    if ua != nil {
        return nil
    }
    return jsonData
}
var conn = make(map[string]any)
type Conn struct {
    C *whatsmeow.Client
}
func Sends(env napi.EnvType, Cli *whatsmeow.Client) map[string]any {
c := &Conn{C:Cli}
conn["GetStore"] = func() any {
    return ToJson(Cli.Store)
}
conn["Connect"] = func() any {
    err := Cli.Connect()
    if err != nil { return Throw(env,err) }
    return ""
}
conn["Disconnect"] = func() {
    Cli.Disconnect()
}
conn["DownloadAny"] = func(msg any) any {
    var mes *waProto.Message
    ua := json.Unmarshal([]byte(ToJson(msg)), &mes)
    if ua != nil { return Throw(env,ua) }
    ok,err := Cli.DownloadAny(ctx, mes)
    if err != nil { return Throw(env,err) }
    buf,_ := napi.CopyBuffer(env, ok)
    return buf
}
conn["FollowNewsletter"] = func(jid string) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    err = Cli.FollowNewsletter(ctx, Jid)
    if err != nil { return Throw(env, err) }
    return nil
}
conn["GetBlocklist"] = func() any {
    res,err := Cli.GetBlocklist(ctx)
    if err != nil { return Throw(env,err) }
    return Res(res)
}
conn["GetBusinessProfile"] = func(jid string) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    res,err := Cli.GetBusinessProfile(ctx, Jid)
    if err != nil { return Throw(env,err) }
    return Res(res)
}
conn["GetContactQRLink"] = func(revoke bool) any {
    res,err := Cli.GetContactQRLink(ctx, revoke)
    if err != nil { return Throw(env,err) }
    return res
}
conn["GetGroupInfo"] = func(jid string) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    res,err := Cli.GetGroupInfo(ctx, Jid)
    if err != nil { return Throw(env,err) }
    return Res(res)
}
conn["GetGroupRequestParticipants"] = func(jid string) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    res,err := Cli.GetGroupRequestParticipants(ctx,Jid)
    if err != nil { return Throw(env,err) }
    return Res(res)
}
conn["GetJoinedGroups"] = func() any {
    res,err := Cli.GetJoinedGroups(ctx)
    if err != nil { return Throw(env,err) }
    return Res(res)
}
conn["GetProfilePictureInfo"] = func(jid string, params any) any {
    Jid, _ := types.ParseJID(jid)
    var param *whatsmeow.GetProfilePictureParams
    ua := json.Unmarshal([]byte(ToJson(params)), &param)
    if ua != nil { return Throw(env,ua) }

    pp, err := Cli.GetProfilePictureInfo(ctx, Jid, param)
    if err != nil { return Throw(env,err) }
    return Res(pp)
}
conn["GetUserInfo"] = func(jids any) any {
    var Jids []types.JID
    err := json.Unmarshal([]byte(ToJson(jids)), &Jids)
    if err != nil { return Throw(env,err) }
    res,err := Cli.GetUserInfo(ctx, Jids)
    if err != nil { return Throw(env,err) }
    return Res(res)
}
conn["IsConnected"] = func() any {
    return Cli.IsConnected()
}
conn["IsLoggedIn"] = func() any {
    return Cli.IsLoggedIn()
}
conn["IsOnWhatsApp"] = func(phones any) any {
    var phone []string
    err := json.Unmarshal([]byte(ToJson(phones)), &phone)
    if err != nil { return Throw(env,err) }

    ison, err := Cli.IsOnWhatsApp(ctx, phone)
    if err != nil { return Throw(env,err) }
    return Res(ison)
}
conn["JoinGroupWithLink"] = func(code string) any {
    jidGroup, err := Cli.JoinGroupWithLink(ctx, code)
    if err != nil { return Throw(env,err) }
    return Res(jidGroup)
}
conn["LeaveGroup"] = func(jid string) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    err = Cli.LeaveGroup(ctx, Jid)
    if err != nil { return Throw(env, err) }
    return nil
}
conn["Logout"] = func() any {
    err := Cli.Logout(ctx)
    if err != nil { return Throw(env, err) }
    return nil
}
conn["PairPhone"] = func(nomor string) any {
    linkingCode, err := Cli.PairPhone(ctx, nomor, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
    if err != nil { return Throw(env,err) }
    return linkingCode
}
conn["SendMessage"] = func(to string, message any, sendRequestExtra any) any {
    To, err := types.ParseJID(to)
    if err != nil { return Throw(env, err) }

    var mes *waProto.Message
    err = json.Unmarshal([]byte(ToJson(message)), &mes)
    if err != nil { return Throw(env,err) }
    
    var extra whatsmeow.SendRequestExtra
    err = json.Unmarshal([]byte(ToJson(sendRequestExtra)), &extra)
    if err != nil { return Throw(env,err) }

    resp, err := Cli.SendMessage(ctx, To, mes, extra)
    if err != nil { return Throw(env,err) }
    return Res(resp)
}
conn["SendPresence"] = func(state string) any {
    err := Cli.SendPresence(ctx, types.Presence(state))
    if err != nil { return Throw(env,err) }
    return ""
}
conn["SetGroupAnnounce"] = func(jid string, announce bool) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    err = Cli.SetGroupAnnounce(ctx, Jid, announce)
    if err != nil { return Throw(env, err) }
    return nil
}
conn["SetGroupDescription"] = func(jid string, description string) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    err = Cli.SetGroupDescription(ctx, Jid, description)
    if err != nil { return Throw(env, err) }
    return nil
}
conn["SetGroupJoinApprovalMode"] = func(jid string, mode bool) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    err = Cli.SetGroupJoinApprovalMode(ctx, Jid, mode)
    if err != nil { return Throw(env, err) }
    return nil
}
conn["SetGroupLocked"] = func(jid string, locked bool) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    err = Cli.SetGroupLocked(ctx, Jid, locked)
    if err != nil { return Throw(env, err) }
    return nil
}
conn["Upload"] = func(args L, tipeM string) any {
    msg,err := c.WaUpload(args, whatsmeow.MediaType(tipeM))
    if err != nil { return Throw(env,err) }
    return Res(msg)
}
conn["UpdateGroupParticipants"] = func(jid string, participantChanges any, action string) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    
    var Jids []types.JID
    err = json.Unmarshal([]byte(ToJson(participantChanges)), &Jids)
    if err != nil { return Throw(env,err) }

    res,err := Cli.UpdateGroupParticipants(ctx, Jid, Jids, whatsmeow.ParticipantChange(action))
    if err != nil { return Throw(env,err) }
    return Res(res)
}
conn["ParseMention"] = func(text string) []string {
    res := []string{}
    matches := regexp.MustCompile("@([0-9]{5,16}|0)").FindAllStringSubmatch(text, -1)
    for _, match := range matches {
        res = append(res, match[1]+"@s.whatsapp.net")
    }
    return res
}
conn["relayMessage"] = func(jid string, message any, a any) any {
    var mes *waProto.Message
    ua := json.Unmarshal([]byte(ToJson(message)), &mes)
    if ua != nil { return Throw(env,ua) }

    var b L
    er := json.Unmarshal([]byte(ToJson(a)), &b)
    if er != nil { return Throw(env,er) }

    ok,err := c.RelayMessage(jid, mes, b)
    if err != nil { return Throw(env,err) }
    return Res(ok)
}
conn["sendText"] = func(jid string, text any, a any) any {
    var b L
    er := json.Unmarshal([]byte(ToJson(a)), &b)
    if er != nil { return Throw(env,er) }

    ok,err := c.SendText(jid, fmt.Sprintf("%s",text) ,b)
    if err != nil { return Throw(env,err) }
    return Res(ok)
}
conn["sendMedia"] = func(jid string, a any) any {
    var b L
    er := json.Unmarshal([]byte(ToJson(a)), &b)
    if er != nil { return Throw(env,er) }

    ok,err := c.SendMedia(jid,b)
    if err != nil { return Throw(env,err) }
    return Res(ok)
}
conn["sendImage"] = func(jid string, a any) any {
    var b L
    er := json.Unmarshal([]byte(ToJson(a)), &b)
    if er != nil { return Throw(env,er) }
    b.Type = whatsmeow.MediaType("WhatsApp Image Keys")
    ok,err := c.SendMedia(jid,b)
    if err != nil { return Throw(env,err) }
    return Res(ok)
}
conn["sendVideo"] = func(jid string, a any) any {
    var b L
    er := json.Unmarshal([]byte(ToJson(a)), &b)
    if er != nil { return Throw(env,er) }
    b.Type = whatsmeow.MediaType("WhatsApp Video Keys")
    ok,err := c.SendMedia(jid,b)
    if err != nil { return Throw(env,err) }
    return Res(ok)
}
conn["sendSticker"] = func(jid string, a any) any {
    var b L
    er := json.Unmarshal([]byte(ToJson(a)), &b)
    if er != nil { return Throw(env,er) }
    b.Type = whatsmeow.MediaType("WhatsApp Sticker Keys")
    ok,err := c.SendMedia(jid,b)
    if err != nil { return Throw(env,err) }
    return Res(ok)
}
conn["sendAudio"] = func(jid string, a any) any {
    var b L
    er := json.Unmarshal([]byte(ToJson(a)), &b)
    if er != nil { return Throw(env,er) }
    b.Type = whatsmeow.MediaType("WhatsApp Audio Keys")
    ok,err := c.SendMedia(jid,b)
    if err != nil { return Throw(env,err) }
    return Res(ok)
}
return conn
}

func (c *Conn) GenerateMessageID() types.MessageID {
    id := make([]byte, 14)
    _, err := rand.Read(id)
    if err != nil {
        panic(err)
    }
    return strings.ToUpper(hex.EncodeToString(id)) + "-FRM"
}
func (c *Conn) ParseMention(text string) []string {
    res := []string{}
    matches := regexp.MustCompile("@([0-9]{5,16}|0)").FindAllStringSubmatch(text, -1)
    for _, match := range matches {
        res = append(res, match[1]+"@s.whatsapp.net")
    }
    return res
}
func (c *Conn) WaUpload(args L, tipeM whatsmeow.MediaType) (*waProto.ImageMessage, error) {
    dow, err := GetByte(args)
    if err != nil { return nil, err }
    uploaded, err := c.C.Upload(context.Background(), dow.Byte, tipeM)
    if err != nil { return nil, err }
    var mtype *string
    if (tipeM == "WhatsApp Image Keys") { mtype = proto.String("image/jpeg") }
    if (tipeM == "WhatsApp Video Keys") { mtype = proto.String("video/mp4") }
    if (tipeM == "WhatsApp Audio Keys") { mtype = proto.String("audio/mpeg") }
    if (args.Mimetype != nil) { mtype = args.Mimetype }
    return &waProto.ImageMessage{
        URL:           proto.String(uploaded.URL),
        DirectPath:    proto.String(uploaded.DirectPath),
        MediaKey:      uploaded.MediaKey,
        Mimetype:      mtype,
        FileEncSHA256: uploaded.FileEncSHA256,
        FileSHA256:    uploaded.FileSHA256,
        FileLength:    proto.Uint64(uint64(dow.Length)),
    }, nil
}
func (c *Conn) RelayMessage(jid string, message *waProto.Message, a L) (*events.Message, error) {
    Jid, _ := types.ParseJID(jid)
    if a.Edit != "" {
        message = c.C.BuildEdit(Jid, a.Edit, message)
    }
    send, err := c.C.SendMessage(context.Background(), Jid, message, whatsmeow.SendRequestExtra{ID:c.C.GenerateMessageID()})
    return &events.Message{
        Info: types.MessageInfo{
            ID: send.ID,
            MessageSource: types.MessageSource{
                Chat: Jid,
                Sender: *c.C.Store.ID,
                IsFromMe: true,
                IsGroup: Jid.Server == types.GroupServer,
            },
        },
        Message: message,
    }, err
}
func (c *Conn) quoted(a L) *waProto.ContextInfo {
    var kontek = &waProto.ContextInfo{}
    if a.ContextInfo != nil {
        kontek = a.ContextInfo
    }
    if a.Quoted != nil {
        kontek.StanzaID = &a.Quoted.Info.ID
        kontek.Participant = proto.String(a.Quoted.Info.Sender.String())
        kontek.QuotedMessage = a.Quoted.Message
    }
    return kontek
}
func (c *Conn) SendText(jid string, text string, a L) (*events.Message, error) {
    var mentionedjid []string
    if a.ParseMention {
        mentionedjid = c.ParseMention(text)
    } else {
        mentionedjid = a.Mentions
    }
    co := c.quoted(a)
    co.MentionedJID = mentionedjid
    return c.RelayMessage(jid, &waProto.Message{
        ExtendedTextMessage: &waProto.ExtendedTextMessage{
            Text: &text,
            ContextInfo: co,
        },
    }, a)
}
func (c *Conn) SendMedia(jid string, a L) (*events.Message, error) {
    var mentionedjid []string
    if a.ParseMention {
        mentionedjid = c.ParseMention(a.Caption)
    } else {
        mentionedjid = a.Mentions
    }
    var typenya string
    if a.Type == "WhatsApp Sticker Keys" {
        typenya = "StickerMessage"
        a.Type = "WhatsApp Image Keys"
        a.Mimetype = proto.String("image/webp")
    } else {
        typenya = MediaType[fmt.Sprintf("%s",a.Type)].(string)
    }
    up, err := c.WaUpload(a, a.Type)
    if err != nil { return nil, err }
    co := c.quoted(a)
    co.MentionedJID = mentionedjid
    medias := J{
      typenya: J{
        "URL": up.URL,
        "DirectPath": up.DirectPath,
        "MediaKey": up.MediaKey,
        "Mimetype": up.Mimetype,
        "FileEncSHA256": up.FileEncSHA256,
        "FileSHA256": up.FileSHA256,
        "FileLength": up.FileLength,
        "Caption": &a.Caption,
        "ContextInfo": co,
      },
    }
    var msgg *waProto.Message
    json.Unmarshal([]byte(ToJson(medias)), &msgg)
    return c.RelayMessage(jid, msgg, a)
}

func Atob(base string) ([]byte) {
    b,_ := base64.StdEncoding.DecodeString(base)
    return b
}
func Btoa(buffer []byte) string {
  return base64.StdEncoding.EncodeToString(buffer)
}
type Getbyte struct {
    Byte []byte
    Mimetype string
    Length int
}
func GetByte(args L) (*Getbyte, error) {
    if args.Byte != nil {
        return &Getbyte{
            Byte: args.Byte, 
            Mimetype: http.DetectContentType(args.Byte),
            Length: len(args.Byte),
        }, nil
    }
    if args.File != "" {
        bacaf,erbacaf := os.ReadFile(args.File)
        return &Getbyte{
            Byte: bacaf, 
            Mimetype: http.DetectContentType(bacaf),
            Length: len(bacaf),
        }, erbacaf
    }
    if args.Url != "" {
        res, err := http.Get(args.Url)
        if err != nil { return nil, err }
        defer res.Body.Close()
        rio,errio := ioutil.ReadAll(res.Body)
        return &Getbyte{
            Byte: rio,
            Mimetype: http.DetectContentType(rio),
            Length: len(rio),
        }, errio
    }
    if args.Base64 != "" {
        rtob := Atob(args.Base64)
        if rtob == nil {
            return nil,fmt.Errorf("error base64")
        }
        return &Getbyte{
            Byte: rtob,
            Mimetype: http.DetectContentType(rtob),
            Length: len(rtob),
        }, nil
    }
    if args.Text != "" {
        tobyte := []byte(args.String)
        return &Getbyte{
            Byte: tobyte, 
            Mimetype: http.DetectContentType(tobyte),
            Length: len(tobyte),
        }, nil
    }
    return nil,nil
}