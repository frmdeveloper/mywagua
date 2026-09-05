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
    "os/exec"
    "path/filepath"
    "regexp"
    "strings"
    "sync"
    "time"
    "sirherobrine23.com.br/Sirherobrine23/napi-go"
    waProto "go.mau.fi/whatsmeow/binary/proto"
    meowcaller "github.com/purpshell/meowcaller"
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
type Conn struct {
    C *whatsmeow.Client
}
func Sends(env napi.EnvType, Cli *whatsmeow.Client) map[string]any {
c := &Conn{C:Cli}
conn := make(map[string]any)
p := map[string]string{}
conn["_params"] = p
reg := func(name, params string, fn any) {
    conn[name] = fn
    p[name] = params
}

reg("GetStore", "", func() any {
    return ToJson(Cli.Store)
})

reg("BuildRevoke", "chat, sender, id", func(chat string, sender string, id string) any {
    Chat,err := types.ParseJID(chat)
    if err != nil { return Throw(env, err) }
    Sender,err := types.ParseJID(sender)
    if err != nil { return Throw(env, err) }
    res := Cli.BuildRevoke(Chat, Sender, types.MessageID(id))
    return Res(res)
})

conn["Connect"] = func() any {
    err := Cli.Connect()
    if err != nil { return Throw(env,err) }
    return nil
}; p["Connect"] = ""

conn["Disconnect"] = func() {
    Cli.Disconnect()
}; p["Disconnect"] = ""

conn["DownloadAny"] = func(msg any) any {
    var mes *waProto.Message
    ua := json.Unmarshal([]byte(ToJson(msg)), &mes)
    if ua != nil { return Throw(env,ua) }
    ok,err := Cli.DownloadAny(ctx, mes)
    if err != nil { return Throw(env,err) }
    buf,_ := napi.CopyBuffer(env, ok)
    return buf
}; p["DownloadAny"] = "message"

conn["FollowNewsletter"] = func(jid string) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    err = Cli.FollowNewsletter(ctx, Jid)
    if err != nil { return Throw(env, err) }
    return nil
}; p["FollowNewsletter"] = "jid"

conn["GetBlocklist"] = func() any {
    res,err := Cli.GetBlocklist(ctx)
    if err != nil { return Throw(env,err) }
    return Res(res)
}; p["GetBlocklist"] = ""

conn["GetBusinessProfile"] = func(jid string) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    res,err := Cli.GetBusinessProfile(ctx, Jid)
    if err != nil { return Throw(env,err) }
    return Res(res)
}; p["GetBusinessProfile"] = "jid"

conn["GetContactQRLink"] = func(revoke bool) any {
    res,err := Cli.GetContactQRLink(ctx, revoke)
    if err != nil { return Throw(env,err) }
    return res
}; p["GetContactQRLink"] = "jid"

conn["GetGroupInfo"] = func(jid string) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    res,err := Cli.GetGroupInfo(ctx, Jid)
    if err != nil { return Throw(env,err) }
    return Res(res)
}; p["GetGroupInfo"] = "jid"

conn["GetGroupInfoFromLink"] = func(code string) any {
    res,err := Cli.GetGroupInfoFromLink(ctx, code)
    if err != nil { return Throw(env,err) }
    return Res(res)
}; p["GetGroupInfoFromLink"] = "link"

conn["GetGroupInviteLink"] = func(jid string, reset bool) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    res,err := Cli.GetGroupInviteLink(ctx, Jid, reset)
    if err != nil { return Throw(env,err) }
    return Res(res)
}; p["GetGroupInviteLink"] = "jid, reset"

conn["GetGroupRequestParticipants"] = func(jid string) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    res,err := Cli.GetGroupRequestParticipants(ctx,Jid)
    if err != nil { return Throw(env,err) }
    return Res(res)
}; p["GetGroupRequestParticipants"] = "jid"

conn["GetJoinedGroups"] = func() any {
    res,err := Cli.GetJoinedGroups(ctx)
    if err != nil { return Throw(env,err) }
    return Res(res)
}; p["GetJoinedGroups"] = ""

conn["GetProfilePictureInfo"] = func(jid string, params any) any {
    Jid, _ := types.ParseJID(jid)
    var param *whatsmeow.GetProfilePictureParams
    ua := json.Unmarshal([]byte(ToJson(params)), &param)
    if ua != nil { return Throw(env,ua) }

    pp, err := Cli.GetProfilePictureInfo(ctx, Jid, param)
    if err != nil { return Throw(env,err) }
    return Res(pp)
}; p["GetProfilePictureInfo"] = "jid, options"

conn["GetUserInfo"] = func(jids any) any {
    var Jids []types.JID
    err := json.Unmarshal([]byte(ToJson(jids)), &Jids)
    if err != nil { return Throw(env,err) }
    res,err := Cli.GetUserInfo(ctx, Jids)
    if err != nil { return Throw(env,err) }
    return Res(res)
}; p["GetUserInfo"] = "jids"

conn["IsConnected"] = func() any {
    return Cli.IsConnected()
}; p["IsConnected"] = ""

conn["IsLoggedIn"] = func() any {
    return Cli.IsLoggedIn()
}; p["IsLoggedIn"] = ""

conn["IsOnWhatsApp"] = func(phones any) any {
    var phone []string
    err := json.Unmarshal([]byte(ToJson(phones)), &phone)
    if err != nil { return Throw(env,err) }

    ison, err := Cli.IsOnWhatsApp(ctx, phone)
    if err != nil { return Throw(env,err) }
    return Res(ison)
}; p["IsOnWhatsApp"] = "jids"

conn["JoinGroupWithLink"] = func(code string) any {
    jidGroup, err := Cli.JoinGroupWithLink(ctx, code)
    if err != nil { return Throw(env,err) }
    return Res(jidGroup)
}; p["JoinGroupWithLink"] = "link"

conn["LeaveGroup"] = func(jid string) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    err = Cli.LeaveGroup(ctx, Jid)
    if err != nil { return Throw(env, err) }
    return nil
}; p["LeaveGroup"] = "jid"

conn["Logout"] = func() any {
    err := Cli.Logout(ctx)
    if err != nil { return Throw(env, err) }
    return nil
}; p["Logout"] = ""

conn["PairPhone"] = func(nomor string) any {
    linkingCode, err := Cli.PairPhone(ctx, nomor, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
    if err != nil { return Throw(env,err) }
    return linkingCode
}; p["PairPhone"] = "phone"

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
}; p["SendMessage"] = "jid, message, options"

conn["SendPresence"] = func(state string) any {
    err := Cli.SendPresence(ctx, types.Presence(state))
    if err != nil { return Throw(env,err) }
    return ""
}; p["SendPresence"] = "state"

conn["SetGroupAnnounce"] = func(jid string, announce bool) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    err = Cli.SetGroupAnnounce(ctx, Jid, announce)
    if err != nil { return Throw(env, err) }
    return nil
}; p["SetGroupAnnounce"] = "jid, announce"

conn["SetGroupDescription"] = func(jid string, description string) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    err = Cli.SetGroupDescription(ctx, Jid, description)
    if err != nil { return Throw(env, err) }
    return nil
}; p["SetGroupDescription"] = "jid, description"

conn["SetGroupJoinApprovalMode"] = func(jid string, mode bool) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    err = Cli.SetGroupJoinApprovalMode(ctx, Jid, mode)
    if err != nil { return Throw(env, err) }
    return nil
}; p["SetGroupJoinApprovalMode"] = "jid, mode"

conn["SetGroupLocked"] = func(jid string, locked bool) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    err = Cli.SetGroupLocked(ctx, Jid, locked)
    if err != nil { return Throw(env, err) }
    return nil
}; p["SetGroupLocked"] = "jid, locked"

conn["SetGroupMemberAddMode"] = func(jid string, mode string) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    err = Cli.SetGroupMemberAddMode(ctx, Jid, types.GroupMemberAddMode(mode))
    if err != nil { return Throw(env, err) }
    return nil
}; p["SetGroupMemberAddMode"] = "jid, mode"

conn["SetGroupName"] = func(jid string, name string) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    err = Cli.SetGroupName(ctx, Jid, name)
    if err != nil { return Throw(env, err) }
    return nil
}; p["SetGroupName"] = "jid, name"

conn["Upload"] = func(args L, tipeM string) any {
    msg,err := c.WaUpload(args, whatsmeow.MediaType(tipeM), false)
    if err != nil { return Throw(env,err) }
    return Res(msg)
}; p["Upload"] = "args"

conn["UploadNewsletter"] = func(args L, tipeM string) any {
    msg,err := c.WaUpload(args, whatsmeow.MediaType(tipeM), true)
    if err != nil { return Throw(env,err) }
    return Res(msg)
}; p["UploadNewsletter"] = "jid, args"

conn["UpdateGroupParticipants"] = func(jid string, participantChanges any, action string) any {
    Jid,err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    
    var Jids []types.JID
    err = json.Unmarshal([]byte(ToJson(participantChanges)), &Jids)
    if err != nil { return Throw(env,err) }

    res,err := Cli.UpdateGroupParticipants(ctx, Jid, Jids, whatsmeow.ParticipantChange(action))
    if err != nil { return Throw(env,err) }
    return Res(res)
}; p["UpdateGroupParticipants"] = "jid, participants, action"

conn["ParseMention"] = func(text string) []string {
    return parseMention(text)
}; p["ParseMention"] = "text"

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
}; p["relayMessage"] = "jid, message, options"

conn["MarkRead"] = func(ids any, timestamp int64, chat string, sender string) any {
    var msgIDs []types.MessageID
    if err := json.Unmarshal([]byte(ToJson(ids)), &msgIDs); err != nil { return Throw(env, err) }
    Chat, err := types.ParseJID(chat)
    if err != nil { return Throw(env, err) }
    Sender, err := types.ParseJID(sender)
    if err != nil { return Throw(env, err) }
    if err := Cli.MarkRead(ctx, msgIDs, time.Unix(timestamp, 0), Chat, Sender); err != nil { return Throw(env, err) }
    return nil
}; p["MarkRead"] = "ids, timestamp, chat, sender"

conn["GenerateMessageID"] = func() any {
    return Cli.GenerateMessageID()
}; p["GenerateMessageID"] = ""

conn["RevokeMessage"] = func(chat string, sender string, id string) any {
    Chat, err := types.ParseJID(chat)
    if err != nil { return Throw(env, err) }
    Sender, err := types.ParseJID(sender)
    if err != nil { return Throw(env, err) }
    revoke := Cli.BuildRevoke(Chat, Sender, types.MessageID(id))
    res, err := Cli.SendMessage(ctx, Chat, revoke)
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["RevokeMessage"] = "chat, sender, id"

conn["BuildReaction"] = func(chat string, sender string, id string, reaction string) any {
    Chat, err := types.ParseJID(chat)
    if err != nil { return Throw(env, err) }
    Sender, err := types.ParseJID(sender)
    if err != nil { return Throw(env, err) }
    return Res(Cli.BuildReaction(Chat, Sender, types.MessageID(id), reaction))
}; p["BuildReaction"] = "chat, sender, id, reaction"

conn["BuildPollCreation"] = func(name string, options any, maxAnswers int) any {
    var opts []string
    if err := json.Unmarshal([]byte(ToJson(options)), &opts); err != nil { return Throw(env, err) }
    return Res(Cli.BuildPollCreation(name, opts, maxAnswers))
}; p["BuildPollCreation"] = "name, options, maxAnswers"

conn["DecryptPollVote"] = func(pollMsg any, vote any) any {
    var pm *events.Message
    if err := json.Unmarshal([]byte(ToJson(pollMsg)), &pm); err != nil { return Throw(env, err) }
    var pv *events.Message
    if err := json.Unmarshal([]byte(ToJson(vote)), &pv); err != nil { return Throw(env, err) }
    if pm == nil || pv == nil || pv.Message == nil || pv.Message.PollUpdateMessage == nil { return Throw(env, fmt.Errorf("invalid poll or vote message")) }
    res, err := Cli.DecryptPollVote(ctx, pv)
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["DecryptPollVote"] = "voteMsg"

conn["DecryptReaction"] = func(reactionMsg any) any {
    var msg *events.Message
    if err := json.Unmarshal([]byte(ToJson(reactionMsg)), &msg); err != nil { return Throw(env, err) }
    res, err := Cli.DecryptReaction(ctx, msg)
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["DecryptReaction"] = "reactionMsg"

conn["SendChatPresence"] = func(jid string, state string, media string) any {
    Jid, err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    if err := Cli.SendChatPresence(ctx, Jid, types.ChatPresence(state), types.ChatPresenceMedia(media)); err != nil { return Throw(env, err) }
    return nil
}; p["SendChatPresence"] = "jid, state, media"

conn["SubscribePresence"] = func(jid string) any {
    Jid, err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    if err := Cli.SubscribePresence(ctx, Jid); err != nil { return Throw(env, err) }
    return nil
}; p["SubscribePresence"] = "jid"

conn["SetStatusMessage"] = func(msg string) any {
    if err := Cli.SetStatusMessage(ctx, types.SetStatusInput{Text: &msg}); err != nil { return Throw(env, err) }
    return nil
}; p["SetStatusMessage"] = "msg"

conn["GetPrivacySettings"] = func() any {
    res := Cli.GetPrivacySettings(ctx)
    return Res(res)
}; p["GetPrivacySettings"] = ""

conn["TryFetchPrivacySettings"] = func(prefetch bool) any {
    res, err := Cli.TryFetchPrivacySettings(ctx, prefetch)
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["TryFetchPrivacySettings"] = "prefetch"

conn["SetPrivacySetting"] = func(name string, value string) any {
    if _, err := Cli.SetPrivacySetting(ctx, types.PrivacySettingType(name), types.PrivacySetting(value)); err != nil { return Throw(env, err) }
    return nil
}; p["SetPrivacySetting"] = "name, value"

conn["GetStatusPrivacy"] = func() any {
    res, err := Cli.GetStatusPrivacy(ctx)
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["GetStatusPrivacy"] = ""

conn["GetUserDevices"] = func(jids any) any {
    var Jids []types.JID
    if err := json.Unmarshal([]byte(ToJson(jids)), &Jids); err != nil { return Throw(env, err) }
    res, err := Cli.GetUserDevices(ctx, Jids)
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["GetUserDevices"] = "jids"

conn["GetUserDevicesContext"] = func(jids any) any {
    var Jids []types.JID
    if err := json.Unmarshal([]byte(ToJson(jids)), &Jids); err != nil { return Throw(env, err) }
    res, err := Cli.GetUserDevicesContext(ctx, Jids)
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["GetUserDevicesContext"] = "jids"

conn["UpdateBlocklist"] = func(jid string, action string) any {
    Jid, err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    res, err := Cli.UpdateBlocklist(ctx, Jid, events.BlocklistChangeAction(action))
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["UpdateBlocklist"] = "jid, action"

conn["ConnectContext"] = func() any {
    if err := Cli.Connect(); err != nil { return Throw(env, err) }
    return nil
}; p["ConnectContext"] = ""

conn["ResetConnection"] = func() {
    Cli.ResetConnection()
}; p["ResetConnection"] = ""

conn["AcceptTOSNotice"] = func(stage string, privacyActtoken string) any {
    if err := Cli.AcceptTOSNotice(ctx, stage, privacyActtoken); err != nil { return Throw(env, err) }
    return nil
}; p["AcceptTOSNotice"] = "stage, privacyActtoken"

conn["CreateGroup"] = func(name string, participants any) any {
    var Jids []types.JID
    if err := json.Unmarshal([]byte(ToJson(participants)), &Jids); err != nil { return Throw(env, err) }
    res, err := Cli.CreateGroup(ctx, whatsmeow.ReqCreateGroup{Name: name, Participants: Jids})
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["CreateGroup"] = "name, participants"

conn["SetGroupPhoto"] = func(jid string, args L) any {
    Jid, err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    data, err := GetByte(args)
    if err != nil { return Throw(env, err) }
    id, err := Cli.SetGroupPhoto(ctx, Jid, data.Byte)
    if err != nil { return Throw(env, err) }
    return id
}; p["SetGroupPhoto"] = "jid, args"

conn["SetGroupTopic"] = func(jid string, prevID string, newID string, topic string) any {
    Jid, err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    if err := Cli.SetGroupTopic(ctx, Jid, prevID, newID, topic); err != nil { return Throw(env, err) }
    return nil
}; p["SetGroupTopic"] = "jid, prevID, newID, topic"

conn["SetDisappearingTimer"] = func(jid string, seconds int64) any {
    Jid, err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    if err := Cli.SetDisappearingTimer(ctx, Jid, time.Duration(seconds)*time.Second, time.Now()); err != nil { return Throw(env, err) }
    return nil
}; p["SetDisappearingTimer"] = "jid, seconds"

conn["SetDefaultDisappearingTimer"] = func(seconds int64) any {
    if err := Cli.SetDefaultDisappearingTimer(ctx, time.Duration(seconds)*time.Second); err != nil { return Throw(env, err) }
    return nil
}; p["SetDefaultDisappearingTimer"] = "seconds"

conn["JoinGroupWithInvite"] = func(inviter string, jid string, code string, expiration int64) any {
    Inviter, err := types.ParseJID(inviter)
    if err != nil { return Throw(env, err) }
    Jid, err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    if err := Cli.JoinGroupWithInvite(ctx, Jid, Inviter, code, expiration); err != nil { return Throw(env, err) }
    return nil
}; p["JoinGroupWithInvite"] = "inviter, jid, code, expiration"

conn["GetGroupInfoFromInvite"] = func(inviter string, jid string, code string, expiration int64) any {
    Inviter, err := types.ParseJID(inviter)
    if err != nil { return Throw(env, err) }
    Jid, err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    res, err := Cli.GetGroupInfoFromInvite(ctx, Jid, Inviter, code, expiration)
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["GetGroupInfoFromInvite"] = "inviter, jid, code, expiration"

conn["UpdateGroupRequestParticipants"] = func(jid string, participants any, action string) any {
    Jid, err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    var Jids []types.JID
    if err := json.Unmarshal([]byte(ToJson(participants)), &Jids); err != nil { return Throw(env, err) }
    res, err := Cli.UpdateGroupRequestParticipants(ctx, Jid, Jids, whatsmeow.ParticipantRequestChange(action))
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["UpdateGroupRequestParticipants"] = "jid, participants, action"

conn["GetSubGroups"] = func(jid string) any {
    Jid, err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    res, err := Cli.GetSubGroups(ctx, Jid)
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["GetSubGroups"] = "jid"

conn["GetLinkedGroupsParticipants"] = func(jid string) any {
    Jid, err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    res, err := Cli.GetLinkedGroupsParticipants(ctx, Jid)
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["GetLinkedGroupsParticipants"] = "jid"

conn["LinkGroup"] = func(parent string, child string) any {
    Parent, err := types.ParseJID(parent)
    if err != nil { return Throw(env, err) }
    Child, err := types.ParseJID(child)
    if err != nil { return Throw(env, err) }
    if err := Cli.LinkGroup(ctx, Parent, Child); err != nil { return Throw(env, err) }
    return nil
}; p["LinkGroup"] = "parent, child"

conn["UnlinkGroup"] = func(parent string, child string) any {
    Parent, err := types.ParseJID(parent)
    if err != nil { return Throw(env, err) }
    Child, err := types.ParseJID(child)
    if err != nil { return Throw(env, err) }
    if err := Cli.UnlinkGroup(ctx, Parent, Child); err != nil { return Throw(env, err) }
    return nil
}; p["UnlinkGroup"] = "parent, child"

conn["GetNewsletterInfo"] = func(jid string) any {
    Jid, err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    res, err := Cli.GetNewsletterInfo(ctx, Jid)
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["GetNewsletterInfo"] = "jid"

conn["GetNewsletterInfoWithInvite"] = func(key string) any {
    res, err := Cli.GetNewsletterInfoWithInvite(ctx, key)
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["GetNewsletterInfoWithInvite"] = "key"

conn["GetSubscribedNewsletters"] = func() any {
    res, err := Cli.GetSubscribedNewsletters(ctx)
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["GetSubscribedNewsletters"] = ""

conn["GetNewsletterMessages"] = func(jid string, count int, before int) any {
    Jid, err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    res, err := Cli.GetNewsletterMessages(ctx, Jid, &whatsmeow.GetNewsletterMessagesParams{
        Count:  count,
        Before: types.MessageServerID(before),
    })
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["GetNewsletterMessages"] = "jid, count, before"

conn["GetNewsletterMessageUpdates"] = func(jid string, count int) any {
    Jid, err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    res, err := Cli.GetNewsletterMessageUpdates(ctx, Jid, &whatsmeow.GetNewsletterUpdatesParams{
        Count: count,
    })
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["GetNewsletterMessageUpdates"] = "jid, count"

conn["UnfollowNewsletter"] = func(jid string) any {
    Jid, err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    if err := Cli.UnfollowNewsletter(ctx, Jid); err != nil { return Throw(env, err) }
    return nil
}; p["UnfollowNewsletter"] = "jid"

conn["NewsletterToggleMute"] = func(jid string, mute bool) any {
    Jid, err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    if err := Cli.NewsletterToggleMute(ctx, Jid, mute); err != nil { return Throw(env, err) }
    return nil
}; p["NewsletterToggleMute"] = "jid, mute"

conn["NewsletterSendReaction"] = func(jid string, serverID int, reaction string, messageID string) any {
    Jid, err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    if err := Cli.NewsletterSendReaction(ctx, Jid, types.MessageServerID(serverID), reaction, types.MessageID(messageID)); err != nil { return Throw(env, err) }
    return nil
}; p["NewsletterSendReaction"] = "jid, serverID, reaction, messageID"

conn["NewsletterMarkViewed"] = func(jid string, serverIDs any) any {
    Jid, err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    var ids []types.MessageServerID
    if err := json.Unmarshal([]byte(ToJson(serverIDs)), &ids); err != nil { return Throw(env, err) }
    if err := Cli.NewsletterMarkViewed(ctx, Jid, ids); err != nil { return Throw(env, err) }
    return nil
}; p["NewsletterMarkViewed"] = "jid, serverIDs"

conn["CreateNewsletter"] = func(name string, description string, pictureArgs L) any {
    params := whatsmeow.CreateNewsletterParams{Name: name, Description: description}
    if pictureArgs.Byte != nil || pictureArgs.File != "" || pictureArgs.Url != "" || pictureArgs.Base64 != "" {
        data, err := GetByte(pictureArgs)
        if err != nil { return Throw(env, err) }
        params.Picture = data.Byte
    }
    res, err := Cli.CreateNewsletter(ctx, params)
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["CreateNewsletter"] = "name, description, pictureArgs"

conn["ResolveBusinessMessageLink"] = func(code string) any {
    res, err := Cli.ResolveBusinessMessageLink(ctx, code)
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["ResolveBusinessMessageLink"] = "code"

conn["ResolveContactQRLink"] = func(code string) any {
    res, err := Cli.ResolveContactQRLink(ctx, code)
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["ResolveContactQRLink"] = "code"

conn["GetBotListV2"] = func() any {
    res, err := Cli.GetBotListV2(ctx)
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["GetBotListV2"] = ""

conn["GetBotProfiles"] = func(jids any) any {
    var infos []types.BotListInfo
    if err := json.Unmarshal([]byte(ToJson(jids)), &infos); err != nil { return Throw(env, err) }
    res, err := Cli.GetBotProfiles(ctx, infos)
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["GetBotProfiles"] = "jids"

conn["BuildEdit"] = func(chat string, id string, newContent any) any {
    Chat, err := types.ParseJID(chat)
    if err != nil { return Throw(env, err) }
    var msg *waProto.Message
    if err := json.Unmarshal([]byte(ToJson(newContent)), &msg); err != nil { return Throw(env, err) }
    return Res(Cli.BuildEdit(Chat, types.MessageID(id), msg))
}; p["BuildEdit"] = "chat, id, newContent"

conn["BuildMessageKey"] = func(chat string, sender string, id string) any {
    Chat, err := types.ParseJID(chat)
    if err != nil { return Throw(env, err) }
    Sender, err := types.ParseJID(sender)
    if err != nil { return Throw(env, err) }
    return Res(Cli.BuildMessageKey(Chat, Sender, types.MessageID(id)))
}; p["BuildMessageKey"] = "chat, sender, id"

conn["BuildUnavailableMessageRequest"] = func(chat string, sender string, id string) any {
    Chat, err := types.ParseJID(chat)
    if err != nil { return Throw(env, err) }
    Sender, err := types.ParseJID(sender)
    if err != nil { return Throw(env, err) }
    return Res(Cli.BuildUnavailableMessageRequest(Chat, Sender, id))
}; p["BuildUnavailableMessageRequest"] = "chat, sender, id"

conn["RejectCall"] = func(callID string, callerJID string) any {
    Caller, err := types.ParseJID(callerJID)
    if err != nil { return Throw(env, err) }
    if err := Cli.RejectCall(ctx, Caller, callID); err != nil { return Throw(env, err) }
    return nil
}; p["RejectCall"] = "callID, callerJID"

conn["RemoveEventHandler"] = func(id uint32) any {
    return Cli.RemoveEventHandler(id)
}; p["RemoveEventHandler"] = "id"

conn["RemoveEventHandlers"] = func() {
    Cli.RemoveEventHandlers()
}; p["RemoveEventHandlers"] = ""

conn["WaitForConnection"] = func(seconds int64) any {
    return Cli.WaitForConnection(time.Duration(seconds) * time.Second)
}; p["WaitForConnection"] = "seconds"

conn["SetPassive"] = func(passive bool) any {
    if err := Cli.SetPassive(ctx, passive); err != nil { return Throw(env, err) }
    return nil
}; p["SetPassive"] = "passive"

conn["SetForceActiveDeliveryReceipts"] = func(active bool) {
    Cli.SetForceActiveDeliveryReceipts(active)
}; p["SetForceActiveDeliveryReceipts"] = "active"

conn["SetMaxParallelRetryReceiptHandling"] = func(max int64) {
    Cli.SetMaxParallelRetryReceiptHandling(max)
}; p["SetMaxParallelRetryReceiptHandling"] = "max"

conn["FetchStickerPack"] = func(id string) any {
    res, err := Cli.FetchStickerPack(ctx, id)
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["FetchStickerPack"] = "id"

conn["NewsletterSubscribeLiveUpdates"] = func(jid string) any {
    Jid, err := types.ParseJID(jid)
    if err != nil { return Throw(env, err) }
    dur, err := Cli.NewsletterSubscribeLiveUpdates(ctx, Jid)
    if err != nil { return Throw(env, err) }
    return dur.Seconds()
}; p["NewsletterSubscribeLiveUpdates"] = "jid"

conn["ParseWebMessage"] = func(chatJID string, webMsg any) any {
    Chat, err := types.ParseJID(chatJID)
    if err != nil { return Throw(env, err) }
    var wm *waProto.WebMessageInfo
    if err := json.Unmarshal([]byte(ToJson(webMsg)), &wm); err != nil { return Throw(env, err) }
    res, err := Cli.ParseWebMessage(Chat, wm)
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["ParseWebMessage"] = "chatJID, webMsg"

conn["SendPeerMessage"] = func(message any) any {
    var msg *waProto.Message
    if err := json.Unmarshal([]byte(ToJson(message)), &msg); err != nil { return Throw(env, err) }
    res, err := Cli.SendPeerMessage(ctx, msg)
    if err != nil { return Throw(env, err) }
    return Res(res)
}; p["SendPeerMessage"] = "message"

conn["SendMediaRetryReceipt"] = func(msgInfo any, mediaKey string) any {
    var info types.MessageInfo
    if err := json.Unmarshal([]byte(ToJson(msgInfo)), &info); err != nil { return Throw(env, err) }
    key, err := base64.StdEncoding.DecodeString(mediaKey)
    if err != nil { return Throw(env, err) }
    if err := Cli.SendMediaRetryReceipt(ctx, &info, key); err != nil { return Throw(env, err) }
    return nil
}; p["SendMediaRetryReceipt"] = "msgInfo, mediaKey"

conn["MarkNotDirty"] = func(name string, timestamp int64) any {
    if err := Cli.MarkNotDirty(ctx, name, time.Unix(timestamp, 0)); err != nil { return Throw(env, err) }
    return nil
}; p["MarkNotDirty"] = "name, timestamp"

conn["StoreLIDPNMapping"] = func(lid string, pn string) {
    Lid, err := types.ParseJID(lid)
    if err != nil { return }
    Pn, err := types.ParseJID(pn)
    if err != nil { return }
    Cli.StoreLIDPNMapping(ctx, Lid, Pn)
}; p["StoreLIDPNMapping"] = "lid, pn"

return conn
}

func (c *Conn) GenerateMessageID() types.MessageID {
    id := make([]byte, 14)
    _, err := rand.Read(id)
    if err != nil {
        panic(err)
    }
    return strings.ToUpper(hex.EncodeToString(id)) + "0FRM"
}
var polaMention = regexp.MustCompile(`@([0-9]{5,20}|0)(?:@(s\.whatsapp\.net|lid|c\.us))?`)

func parseMention(text string) []string {
    res := []string{}
    sudah := map[string]bool{}
    for _, match := range polaMention.FindAllStringSubmatch(text, -1) {
        server := "s.whatsapp.net"
        if match[2] == "lid" {
            server = "lid"
        }
        jid := match[1] + "@" + server
        if sudah[jid] {
            continue
        }
        sudah[jid] = true
        res = append(res, jid)
    }
    return res
}

func (c *Conn) ParseMention(text string) []string {
    return parseMention(text)
}
func (c *Conn) WaUpload(args L, tipeM whatsmeow.MediaType, newsletter bool) (J, error) {
    dow, err := GetByte(args)
    if err != nil { return nil, err }
    var uploaded whatsmeow.UploadResponse
    var uperr error
    if newsletter {
        uploaded, err = c.C.UploadNewsletter(context.Background(), dow.Byte, tipeM)
    } else {
        uploaded, err = c.C.Upload(context.Background(), dow.Byte, tipeM)
    }
    if uperr != nil { return nil, uperr }
    return J{
        "URL": uploaded.URL,
        "directPath": uploaded.DirectPath,
        "handle": uploaded.Handle,
        "objectID": uploaded.ObjectID,
        "mediaKey": uploaded.MediaKey,
        "fileEncSHA256": uploaded.FileEncSHA256,
        "fileSHA256": uploaded.FileSHA256,
        "fileLength": uploaded.FileLength,
    }, nil
}
func (c *Conn) RelayMessage(jid string, message *waProto.Message, a L) (*events.Message, error) {
    Jid, _ := types.ParseJID(jid)
    if a.Edit != "" {
        message = c.C.BuildEdit(Jid, a.Edit, message)
    }
    send, err := c.C.SendMessage(context.Background(), Jid, message, whatsmeow.SendRequestExtra{ID:c.GenerateMessageID()})
    return &events.Message{
        Info: types.MessageInfo{
            ID: send.ID,
            ServerID: send.ServerID,
            Timestamp: send.Timestamp,
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

func SetupCaller(env napi.EnvType, caller *meowcaller.Client, push func(any), conn map[string]any, p map[string]string) {
    var callMap sync.Map

    registerCB := func(call *meowcaller.Call) {
        call.OnReady(func() {
            push(map[string]any{
                "type": "meowcaller.CallReady",
                "evt":  map[string]any{
                    "callId": call.ID(),
                    "peer":   call.Peer().String(),
                },
            })
        })
        call.OnEnd(func(reason string) {
            push(map[string]any{
                "type": "meowcaller.CallEnd",
                "evt":  map[string]any{
                    "callId": call.ID(),
                    "reason": reason,
                },
            })
            callMap.Delete(call.ID())
        })
        call.OnStateChange(func(phase meowcaller.CallPhase) {
            push(map[string]any{
                "type": "meowcaller.CallStateChange",
                "evt":  map[string]any{
                    "callId": call.ID(),
                    "phase":  int(phase),
                },
            })
        })
    }

    caller.OnIncomingCall(func(call *meowcaller.Call) {
        callMap.Store(call.ID(), call)
        registerCB(call)
        push(map[string]any{
            "type": "meowcaller.IncomingCall",
            "evt":  map[string]any{
                "callId":  call.ID(),
                "peer":    call.Peer().String(),
                "isVideo": call.IsVideo(),
            },
        })
    })

    getCall := func(callId string) (*meowcaller.Call, error) {
        v, ok := callMap.Load(callId)
        if !ok {
            return nil, fmt.Errorf("unknown call: %s", callId)
        }
        return v.(*meowcaller.Call), nil
    }

    conn["answerCall"] = func(callId string) any {
        call, err := getCall(callId)
        if err != nil { return Throw(env, err) }
        if err := call.Answer(); err != nil { return Throw(env, err) }
        return nil
    }; p["answerCall"] = "callId"

    conn["rejectCall"] = func(callId string) any {
        call, err := getCall(callId)
        if err != nil { return Throw(env, err) }
        if err := call.Reject(); err != nil { return Throw(env, err) }
        return nil
    }; p["rejectCall"] = "callId"

    conn["hangupCall"] = func(callId string) any {
        call, err := getCall(callId)
        if err != nil { return Throw(env, err) }
        if err := call.Hangup(); err != nil { return Throw(env, err) }
        return nil
    }; p["hangupCall"] = "callId"

    conn["playAudio"] = func(callId, path string) any {
        call, err := getCall(callId)
        if err != nil { return Throw(env, err) }

        videoExts := map[string]bool{
            ".mp4": true, ".mkv": true, ".avi": true, ".mov": true,
            ".flv": true, ".webm": true, ".m4v": true, ".3gp": true,
        }
        actualPath := path
        if videoExts[strings.ToLower(filepath.Ext(path))] {
            ffmpeg, lookErr := exec.LookPath("ffmpeg")
            if lookErr != nil {
                for _, candidate := range []string{
                    "/data/data/com.termux/files/usr/bin/ffmpeg",
                    "/usr/bin/ffmpeg",
                    "/usr/local/bin/ffmpeg",
                } {
                    if _, e := os.Stat(candidate); e == nil {
                        ffmpeg = candidate
                        lookErr = nil
                        break
                    }
                }
            }
            if lookErr != nil { return Throw(env, fmt.Errorf("ffmpeg not found: %v", lookErr)) }
            tmp := path + "_converted.mp3"
            out, ffErr := exec.Command(ffmpeg, "-y", "-i", path, "-vn", "-ar", "16000", "-ac", "1", "-b:a", "64k", tmp).CombinedOutput()
            if ffErr != nil { return Throw(env, fmt.Errorf("ffmpeg: %v\n%s", ffErr, string(out))) }
            defer os.Remove(tmp)
            actualPath = tmp
        }

        var src meowcaller.AudioSource
        switch strings.ToLower(filepath.Ext(actualPath)) {
        case ".mp3":
            src, err = meowcaller.MP3File(actualPath)
        case ".ogg", ".opus":
            src, err = meowcaller.OpusFile(actualPath)
        default:
            src, err = meowcaller.WAVFile(actualPath)
        }
        if err != nil { return Throw(env, err) }
        call.Play(src)
        return nil
    }; p["playAudio"] = "callId, path"

    conn["receiveAudio"] = func(callId, path string) any {
        call, err := getCall(callId)
        if err != nil { return Throw(env, err) }
        sink, err := meowcaller.WAVRecorder(path)
        if err != nil { return Throw(env, err) }
        call.Receive(sink)
        return nil
    }; p["receiveAudio"] = "callId, path"

    conn["receivePCM"] = func(callId string) any {
        call, err := getCall(callId)
        if err != nil { return Throw(env, err) }
        call.Receive(meowcaller.SinkFunc(func(frame []float32) {
            push(map[string]any{
                "type": "meowcaller.AudioFrame",
                "evt": map[string]any{
                    "callId": callId,
                    "pcm":    frame,
                },
            })
        }))
        return nil
    }; p["receivePCM"] = "callId"

    conn["placeCall"] = func(target string) any {
        var (
            outCall *meowcaller.Call
            callErr error
        )
        worker, err := napi.CreateAsyncWorker(env,
            func(env napi.EnvType) {
                outCall, callErr = caller.Call(context.Background(), target)
            },
            func(env napi.EnvType, Resolve, Reject func(napi.ValueType)) {
                if callErr != nil {
                    s, _ := napi.CreateString(env, callErr.Error())
                    Reject(s)
                    return
                }
                callMap.Store(outCall.ID(), outCall)
                registerCB(outCall)
                s, _ := napi.CreateString(env, outCall.ID())
                Resolve(s)
            },
        )
        if err != nil { return Throw(env, err) }
        return worker
    }; p["placeCall"] = "target"

    conn["playVideo"] = func(callId string, path string) any {
        call, err := getCall(callId)
        if err != nil { return Throw(env, err) }
        au, err := os.ReadFile(path)
        if err != nil { return Throw(env, err) }
        if err := call.SendVideo(au); err != nil { return Throw(env, err) }
        return nil
    }; p["playVideo"] = "callId, path"
}