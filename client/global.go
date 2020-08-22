package client

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/wfjsw/MiraiGo/binary"
	devinfo "github.com/wfjsw/MiraiGo/client/pb"
	"github.com/wfjsw/MiraiGo/client/pb/msg"
	"github.com/wfjsw/MiraiGo/message"
	"github.com/wfjsw/MiraiGo/utils"
	"google.golang.org/protobuf/proto"
)

type DeviceInfo struct {
	Display      []byte
	Product      []byte
	Device       []byte
	Board        []byte
	Brand        []byte
	Model        []byte
	Bootloader   []byte
	FingerPrint  []byte
	BootId       []byte
	ProcVersion  []byte
	BaseBand     []byte
	SimInfo      []byte
	OSType       []byte
	MacAddress   []byte
	IpAddress    []byte
	WifiBSSID    []byte
	WifiSSID     []byte
	IMSIMd5      []byte
	IMEI         string
	AndroidId    []byte
	APN          []byte
	Guid         []byte
	TgtgtKey     []byte
	Version      *Version
	VendorName   string
	VendorOSName string
}

type Version struct {
	Incremental []byte
	Release     []byte
	CodeName    []byte
	Sdk         uint32
}

type DeviceInfoFile struct {
	Display      string       `json:"display"`
	Product      string       `json:"product"`
	Device       string       `json:"device"`
	Board        string       `json:"board"`
	Model        string       `json:"model"`
	Bootloader   string       `json:"bootloader"`
	FingerPrint  string       `json:"finger_print"`
	BootId       string       `json:"boot_id"`
	ProcVersion  string       `json:"proc_version"`
	SimInfo      string       `json:"sim_info"`
	MacAddress   string       `json:"mac_address"`
	WifiBSSID    string       `json:"wifi_bssid"`
	WifiSSID     string       `json:"wifi_ssid"`
	IMEI         string       `json:"imei"`
	AndroidId    string       `json:"android_id"`
	APN          string       `json:"apn"`
	Version      *VersionFile `json:"version"`
	VendorName   string       `json:"vendor_name"`
	VendorOSName string       `json:"vendor_os_name"`
}

type VersionFile struct {
	Incremental []byte `json:"incremental"`
	Release     []byte `json:"release"`
	CodeName    []byte `json:"codename"`
	Sdk         uint32 `json:"sdk"`
}

type groupMessageBuilder struct {
	MessageSeq    int32
	MessageCount  int32
	MessageSlices []*msg.Message
}

// default
var SystemDeviceInfo = &DeviceInfo{
	// Display:     []byte("MIRAI.123456.001"),
	Display:     []byte("OPPO R9sk"),
	Product:     []byte("OPPO R9sk"),
	Device:      []byte("OPPO R9sk"),
	Board:       []byte("msm8953"),
	Brand:       []byte("OPPO"),
	Model:       []byte("OPPO R9sk (R9sk)"),
	Bootloader:  []byte("unknown"),
	FingerPrint: []byte("OPPO/R9sk/R9sk:6.0.1/MMB29M/1234567890:user/release-keys"),
	BootId:      []byte("cb886ae2-00b6-4d68-a230-787f111d12c7"),
	ProcVersion: []byte("Linux version 3.18.24-perf-cb886ae2 (eng.root.20170603.124902)"),
	BaseBand:    []byte{},
	SimInfo:     []byte("CMCC"),
	OSType:      []byte("android"),
	MacAddress:  []byte("00:50:56:C0:00:08"),
	IpAddress:   []byte{192, 168, 1, 101}, // 10.0.1.3
	WifiBSSID:   []byte("00:50:56:C0:00:08"),
	WifiSSID:    []byte("<unknown ssid>"),
	IMEI:        "468356291846738",
	AndroidId:   []byte("123456789abcdef"),
	APN:         []byte("wifi"),
	Version: &Version{
		Incremental: []byte("5891938"),
		Release:     []byte("10"),
		CodeName:    []byte("REL"),
		Sdk:         29,
	},
	VendorName:   "OnePlus",
	VendorOSName: "ONEPLUS A5000_23_17",
}

const IMEI_BASE_DIGITS_COUNT int = 14

var EmptyBytes = []byte{}
var NumberRange = "0123456789"

func init() {
	r := make([]byte, 16)
	rand.Read(r)
	t := md5.Sum(r)
	SystemDeviceInfo.IMSIMd5 = t[:]
	SystemDeviceInfo.GenNewGuid()
	SystemDeviceInfo.GenNewTgtgtKey()
}

func GenIMEI() string {
	sum := 0 // the control sum of digits
	var final strings.Builder

	randSrc := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSrc)

	for i := 0; i < IMEI_BASE_DIGITS_COUNT; i++ { // generating all the base digits
		toAdd := randGen.Intn(10)
		if (i+1)%2 == 0 { // special proc for every 2nd one
			toAdd *= 2
			if toAdd >= 10 {
				toAdd = (toAdd % 10) + 1
			}
		}
		sum += toAdd
		final.WriteString(fmt.Sprintf("%d", toAdd)) // and even printing them here!
	}
	var ctrlDigit int = (sum * 9) % 10 // calculating the control digit
	final.WriteString(fmt.Sprintf("%d", ctrlDigit))
	return final.String()
}

func GenRandomDevice() {
	r := make([]byte, 16)
	rand.Read(r)
	// SystemDeviceInfo.Display = []byte("MIRAI." + utils.RandomStringRange(6, NumberRange) + ".001")
	SystemDeviceInfo.Display = []byte("OPPO R9sk - " + utils.RandomStringRange(6, NumberRange))
	SystemDeviceInfo.FingerPrint = []byte("OPPO/R9sk/R9sk:6.0.1/MMB29M/" + utils.RandomStringRange(10, NumberRange) + ":user/release-keys")
	SystemDeviceInfo.BootId = []byte(binary.GenUUID(r))
	SystemDeviceInfo.ProcVersion = []byte("Linux version 3.18.24-perf-" + utils.RandomString(8) + " (eng.root.20170603.124902)")
	rand.Read(r)
	t := md5.Sum(r)
	SystemDeviceInfo.IMSIMd5 = t[:]
	SystemDeviceInfo.IMEI = GenIMEI()
	SystemDeviceInfo.AndroidId = []byte(utils.RandomString(15))
	SystemDeviceInfo.GenNewGuid()
	SystemDeviceInfo.GenNewTgtgtKey()
}

func (info *DeviceInfo) ToJson() []byte {
	v := &VersionFile{
		Incremental: string(info.Version.Incremental),
		Release:     string(info.Version.Release),
		CodeName:    string(info.Version.CodeName),
		Sdk:         info.Version.Sdk,
	}
	f := &DeviceInfoFile{
		Display:      string(info.Display),
		Product:      string(info.Product),
		Device:       string(info.Device),
		Board:        string(info.Board),
		Model:        string(info.Model),
		Bootloader:   string(info.Bootloader),
		FingerPrint:  string(info.FingerPrint),
		BootId:       string(info.BootId),
		ProcVersion:  string(info.ProcVersion),
		SimInfo:      string(info.SimInfo),
		MacAddress:   string(info.MacAddress),
		WifiBSSID:    string(info.WifiBSSID),
		WifiSSID:     string(info.WifiSSID),
		IMEI:         string(info.IMEI),
		AndroidId:    string(info.IMEI),
		APN:          string(info.APN),
		Version:      v,
		VendorName:   info.VendorName,
		VendorOSName: info.VendorOSName,
	}
	d, _ := json.Marshal(f)
	return d
}

func (info *DeviceInfo) ReadJson(d []byte) error {
	var f DeviceInfoFile
	if err := json.Unmarshal(d, &f); err != nil {
		return err
	}
	info.Display = []byte(f.Display)
	if f.Product != "" {
		info.Product = []byte(f.Product)
		info.Device = []byte(f.Device)
		info.Board = []byte(f.Board)
		info.Model = []byte(f.Model)
	}
	if f.Bootloader != "" {
		info.Bootloader = []byte(f.Bootloader)
	}
	info.FingerPrint = []byte(f.FingerPrint)
	info.BootId = []byte(f.BootId)
	info.ProcVersion = []byte(f.ProcVersion)
	if f.SimInfo != "" {
		info.SimInfo = []byte(f.SimInfo)
	}
	if f.MacAddress != "" {
		info.MacAddress = []byte(f.MacAddress)
	}
	if f.WifiBSSID != "" {
		info.WifiBSSID = []byte(f.WifiBSSID)
	}
	if f.WifiSSID != "" {
		info.WifiSSID = []byte(f.WifiSSID)
	}
	info.IMEI = f.IMEI
	info.AndroidId = []byte(f.AndroidId)
	if f.APN != "" {
		info.APN = []byte(f.APN)
	}
	if f.Version != nil {
		info.Version.Incremental = []byte(f.Version.Incremental)
		info.Version.Release = []byte(f.Version.Release)
		info.Version.CodeName = []byte(f.Version.CodeName)
		info.Version.Sdk = f.Version.Sdk
	}
	if f.VendorName != "" {
		info.VendorName = f.VendorName
	}
	if f.VendorOSName != "" {
		info.VendorOSName = f.VendorOSName
	}
	SystemDeviceInfo.GenNewGuid()
	SystemDeviceInfo.GenNewTgtgtKey()
	return nil
}

func (info *DeviceInfo) GenNewGuid() {
	t := md5.Sum(append(info.AndroidId, info.MacAddress...))
	info.Guid = t[:]
}

func (info *DeviceInfo) GenNewTgtgtKey() {
	r := make([]byte, 16)
	rand.Read(r)
	t := md5.Sum(append(r, info.Guid...))
	info.TgtgtKey = t[:]
}

func (info *DeviceInfo) GenDeviceInfoData() []byte {
	m := &devinfo.DeviceInfo{
		Bootloader:   string(info.Bootloader),
		ProcVersion:  string(info.ProcVersion),
		Codename:     string(info.Version.CodeName),
		Incremental:  string(info.Version.Incremental),
		Fingerprint:  string(info.FingerPrint),
		BootId:       string(info.BootId),
		AndroidId:    string(info.AndroidId),
		BaseBand:     string(info.BaseBand),
		InnerVersion: string(info.Version.Incremental),
	}
	data, err := proto.Marshal(m)
	if err != nil {
		panic(err)
	}
	return data
}

func (c *QQClient) parsePrivateMessage(msg *msg.Message) *message.PrivateMessage {
	friend := c.FindFriend(msg.Head.FromUin)
	if friend == nil {
		return nil
	}
	ret := &message.PrivateMessage{
		Id:     msg.Head.MsgSeq,
		Target: c.Uin,
		Time:   msg.Head.MsgTime,
		Sender: &message.Sender{
			Uin:      friend.Uin,
			Nickname: friend.Nickname,
		},
		Elements: message.ParseMessageElems(msg.Body.RichText.Elems),
	}
	if msg.Body.RichText.Attr != nil {
		ret.InternalId = msg.Body.RichText.Attr.Random
	}
	return ret
}

func (c *QQClient) parseTempMessage(msg *msg.Message) *message.TempMessage {
	group := c.FindGroupByUin(msg.Head.C2CTmpMsgHead.GroupUin)
	mem := group.FindMember(msg.Head.FromUin)
	return &message.TempMessage{
		Id:        msg.Head.MsgSeq,
		GroupCode: group.Code,
		GroupName: group.Name,
		Sender: &message.Sender{
			Uin:      mem.Uin,
			Nickname: mem.Nickname,
			CardName: mem.CardName,
		},
		Elements: message.ParseMessageElems(msg.Body.RichText.Elems),
	}
}

func (c *QQClient) parseGroupMessage(m *msg.Message) *message.GroupMessage {
	group := c.FindGroup(m.Head.GroupInfo.GroupCode)
	if group == nil {
		return nil
	}
	var anonInfo *msg.AnonymousGroupMessage
	for _, e := range m.Body.RichText.Elems {
		if e.AnonGroupMsg != nil {
			anonInfo = e.AnonGroupMsg
		}
	}
	var sender *message.Sender
	if anonInfo != nil {
		sender = &message.Sender{
			Uin:      80000000,
			Nickname: string(anonInfo.AnonNick),
			IsFriend: false,
		}
	} else {
		mem := group.FindMember(m.Head.FromUin)
		if mem == nil {
			return nil
		}
		sender = &message.Sender{
			Uin:      mem.Uin,
			Nickname: mem.Nickname,
			CardName: mem.CardName,
			IsFriend: c.FindFriend(mem.Uin) != nil,
		}
	}
	g := &message.GroupMessage{
		Id:        m.Head.MsgSeq,
		GroupCode: group.Code,
		GroupName: string(m.Head.GroupInfo.GroupName),
		Sender:    sender,
		Time:      m.Head.MsgTime,
		Elements:  message.ParseMessageElems(m.Body.RichText.Elems),
		//OriginalElements: m.Body.RichText.Elems,
	}
	if m.Body.RichText.Ptt != nil {
		g.Elements = []message.IMessageElement{
			&message.VoiceElement{
				Name: m.Body.RichText.Ptt.FileName,
				Md5:  m.Body.RichText.Ptt.FileMd5,
				Size: m.Body.RichText.Ptt.FileSize,
				Url:  "http://grouptalk.c2c.qq.com" + string(m.Body.RichText.Ptt.DownPara),
			},
		}
	}
	if m.Body.RichText.Attr != nil {
		g.InternalId = m.Body.RichText.Attr.Random
	}
	return g
}

func (b *groupMessageBuilder) build() *msg.Message {
	sort.Slice(b.MessageSlices, func(i, j int) bool {
		return b.MessageSlices[i].Content.PkgIndex < b.MessageSlices[i].Content.PkgIndex
	})
	base := b.MessageSlices[0]
	for _, m := range b.MessageSlices[1:] {
		base.Body.RichText.Elems = append(base.Body.RichText.Elems, m.Body.RichText.Elems...)
	}
	return base
}

func packRequestDataV3(data []byte) (r []byte) {
	r = append([]byte{0x0A}, data...)
	r = append(r, 0x0B)
	return
}

func genForwardTemplate(resId, preview, title, brief, source, summary string, ts int64) *message.SendingMessage {
	template := fmt.Sprintf(`<?xml version='1.0' encoding='UTF-8'?><msg serviceID="35" templateID="1" action="viewMultiMsg" brief="%s" m_resid="%s" m_fileName="%d" tSum="3" sourceMsgId="0" url="" flag="3" adverSign="0" multiMsgFlag="0"><item layout="1"><title color="#000000" size="34">%s</title> %s<hr></hr><summary size="26" color="#808080">%s</summary></item><source name="%s"></source></msg>`,
		brief, resId, ts, title, preview, summary, source,
	)
	return &message.SendingMessage{Elements: []message.IMessageElement{
		&message.ServiceElement{
			Id:      35,
			Content: template,
			ResId:   resId,
			SubType: "Forward",
		},
	}}
}

func genLongTemplate(resId, brief string, ts int64) *message.SendingMessage {
	limited := func() string {
		if len(brief) > 30 {
			return brief[:30] + "…"
		}
		return brief
	}()
	template := fmt.Sprintf(`<?xml version='1.0' encoding='UTF-8' standalone='yes' ?><msg serviceID="35" templateID="1" action="viewMultiMsg" brief="%s" m_resid="%s" m_fileName="%d" sourceMsgId="0" url="" flag="3" adverSign="0" multiMsgFlag="1"> <item layout="1"> <title>%s</title> <hr hidden="false" style="0"/> <summary>点击查看完整消息</summary> </item> <source name="聊天记录" icon="" action="" appid="-1"/> </msg>`,
		limited, resId, ts, limited,
	)
	return &message.SendingMessage{Elements: []message.IMessageElement{
		&message.ServiceElement{
			Id:      35,
			Content: template,
			ResId:   resId,
			SubType: "Long",
		},
	}}
}
