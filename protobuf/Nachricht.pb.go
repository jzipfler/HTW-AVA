// Code generated by protoc-gen-go.
// source: Nachricht.proto
// DO NOT EDIT!

/*
Package protobuf is a generated protocol buffer package.

It is generated from these files:
	Nachricht.proto

It has these top-level messages:
	Nachricht
	MessageTwo
*/
package protobuf

import proto "code.google.com/p/goprotobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type Nachricht_NachrichtenTyp int32

const (
	Nachricht_KONTROLLNACHRICHT   Nachricht_NachrichtenTyp = 0
	Nachricht_ANWENDUNGSNACHRICHT Nachricht_NachrichtenTyp = 1
)

var Nachricht_NachrichtenTyp_name = map[int32]string{
	0: "KONTROLLNACHRICHT",
	1: "ANWENDUNGSNACHRICHT",
}
var Nachricht_NachrichtenTyp_value = map[string]int32{
	"KONTROLLNACHRICHT":   0,
	"ANWENDUNGSNACHRICHT": 1,
}

func (x Nachricht_NachrichtenTyp) Enum() *Nachricht_NachrichtenTyp {
	p := new(Nachricht_NachrichtenTyp)
	*p = x
	return p
}
func (x Nachricht_NachrichtenTyp) String() string {
	return proto.EnumName(Nachricht_NachrichtenTyp_name, int32(x))
}
func (x *Nachricht_NachrichtenTyp) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Nachricht_NachrichtenTyp_value, data, "Nachricht_NachrichtenTyp")
	if err != nil {
		return err
	}
	*x = Nachricht_NachrichtenTyp(value)
	return nil
}

type Nachricht_KontrollTyp int32

const (
	// Use this option if you want to assign multiple definitions for the same value
	// For example: "INITIALISIEREN = 0;" && "START = 0;".
	// option allow_alias = true;
	Nachricht_INITIALISIEREN Nachricht_KontrollTyp = 0
	Nachricht_BEENDEN        Nachricht_KontrollTyp = 1
)

var Nachricht_KontrollTyp_name = map[int32]string{
	0: "INITIALISIEREN",
	1: "BEENDEN",
}
var Nachricht_KontrollTyp_value = map[string]int32{
	"INITIALISIEREN": 0,
	"BEENDEN":        1,
}

func (x Nachricht_KontrollTyp) Enum() *Nachricht_KontrollTyp {
	p := new(Nachricht_KontrollTyp)
	*p = x
	return p
}
func (x Nachricht_KontrollTyp) String() string {
	return proto.EnumName(Nachricht_KontrollTyp_name, int32(x))
}
func (x *Nachricht_KontrollTyp) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Nachricht_KontrollTyp_value, data, "Nachricht_KontrollTyp")
	if err != nil {
		return err
	}
	*x = Nachricht_KontrollTyp(value)
	return nil
}

type MessageTwo_MessageType int32

const (
	MessageTwo_CONTROLMESSAGE     MessageTwo_MessageType = 0
	MessageTwo_APPLICATIONMESSAGE MessageTwo_MessageType = 1
)

var MessageTwo_MessageType_name = map[int32]string{
	0: "CONTROLMESSAGE",
	1: "APPLICATIONMESSAGE",
}
var MessageTwo_MessageType_value = map[string]int32{
	"CONTROLMESSAGE":     0,
	"APPLICATIONMESSAGE": 1,
}

func (x MessageTwo_MessageType) Enum() *MessageTwo_MessageType {
	p := new(MessageTwo_MessageType)
	*p = x
	return p
}
func (x MessageTwo_MessageType) String() string {
	return proto.EnumName(MessageTwo_MessageType_name, int32(x))
}
func (x *MessageTwo_MessageType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(MessageTwo_MessageType_value, data, "MessageTwo_MessageType")
	if err != nil {
		return err
	}
	*x = MessageTwo_MessageType(value)
	return nil
}

type MessageTwo_ControlType int32

const (
	// Use this option if you want to assign multiple definitions for the same value
	// For example: "INITIALIZE = 0;" && "START = 0;".
	// option allow_alias = true;
	MessageTwo_INITIALIZE MessageTwo_ControlType = 0
	MessageTwo_QUIT       MessageTwo_ControlType = 1
)

var MessageTwo_ControlType_name = map[int32]string{
	0: "INITIALIZE",
	1: "QUIT",
}
var MessageTwo_ControlType_value = map[string]int32{
	"INITIALIZE": 0,
	"QUIT":       1,
}

func (x MessageTwo_ControlType) Enum() *MessageTwo_ControlType {
	p := new(MessageTwo_ControlType)
	*p = x
	return p
}
func (x MessageTwo_ControlType) String() string {
	return proto.EnumName(MessageTwo_ControlType_name, int32(x))
}
func (x *MessageTwo_ControlType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(MessageTwo_ControlType_value, data, "MessageTwo_ControlType")
	if err != nil {
		return err
	}
	*x = MessageTwo_ControlType(value)
	return nil
}

type MessageTwo_NodeType int32

const (
	MessageTwo_COMPANY  MessageTwo_NodeType = 0
	MessageTwo_CUSTOMER MessageTwo_NodeType = 1
)

var MessageTwo_NodeType_name = map[int32]string{
	0: "COMPANY",
	1: "CUSTOMER",
}
var MessageTwo_NodeType_value = map[string]int32{
	"COMPANY":  0,
	"CUSTOMER": 1,
}

func (x MessageTwo_NodeType) Enum() *MessageTwo_NodeType {
	p := new(MessageTwo_NodeType)
	*p = x
	return p
}
func (x MessageTwo_NodeType) String() string {
	return proto.EnumName(MessageTwo_NodeType_name, int32(x))
}
func (x *MessageTwo_NodeType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(MessageTwo_NodeType_value, data, "MessageTwo_NodeType")
	if err != nil {
		return err
	}
	*x = MessageTwo_NodeType(value)
	return nil
}

// Nachrichtendefinition fuer Aufgabe 1.
type Nachricht struct {
	SourceIP          *string                   `protobuf:"bytes,1,req,name=sourceIP" json:"sourceIP,omitempty"`
	SourcePort        *int32                    `protobuf:"varint,2,req,name=sourcePort" json:"sourcePort,omitempty"`
	SourceID          *int32                    `protobuf:"varint,3,req,name=sourceID" json:"sourceID,omitempty"`
	NachrichtenTyp    *Nachricht_NachrichtenTyp `protobuf:"varint,4,req,name=nachrichtenTyp,enum=protobuf.Nachricht_NachrichtenTyp,def=1" json:"nachrichtenTyp,omitempty"`
	KontrollTyp       *Nachricht_KontrollTyp    `protobuf:"varint,5,opt,name=kontrollTyp,enum=protobuf.Nachricht_KontrollTyp,def=0" json:"kontrollTyp,omitempty"`
	NachrichtenInhalt *string                   `protobuf:"bytes,6,req,name=nachrichtenInhalt" json:"nachrichtenInhalt,omitempty"`
	ZeitStempel       *string                   `protobuf:"bytes,7,opt,name=zeitStempel" json:"zeitStempel,omitempty"`
	XXX_unrecognized  []byte                    `json:"-"`
}

func (m *Nachricht) Reset()         { *m = Nachricht{} }
func (m *Nachricht) String() string { return proto.CompactTextString(m) }
func (*Nachricht) ProtoMessage()    {}

const Default_Nachricht_NachrichtenTyp Nachricht_NachrichtenTyp = Nachricht_ANWENDUNGSNACHRICHT
const Default_Nachricht_KontrollTyp Nachricht_KontrollTyp = Nachricht_INITIALISIEREN

func (m *Nachricht) GetSourceIP() string {
	if m != nil && m.SourceIP != nil {
		return *m.SourceIP
	}
	return ""
}

func (m *Nachricht) GetSourcePort() int32 {
	if m != nil && m.SourcePort != nil {
		return *m.SourcePort
	}
	return 0
}

func (m *Nachricht) GetSourceID() int32 {
	if m != nil && m.SourceID != nil {
		return *m.SourceID
	}
	return 0
}

func (m *Nachricht) GetNachrichtenTyp() Nachricht_NachrichtenTyp {
	if m != nil && m.NachrichtenTyp != nil {
		return *m.NachrichtenTyp
	}
	return Default_Nachricht_NachrichtenTyp
}

func (m *Nachricht) GetKontrollTyp() Nachricht_KontrollTyp {
	if m != nil && m.KontrollTyp != nil {
		return *m.KontrollTyp
	}
	return Default_Nachricht_KontrollTyp
}

func (m *Nachricht) GetNachrichtenInhalt() string {
	if m != nil && m.NachrichtenInhalt != nil {
		return *m.NachrichtenInhalt
	}
	return ""
}

func (m *Nachricht) GetZeitStempel() string {
	if m != nil && m.ZeitStempel != nil {
		return *m.ZeitStempel
	}
	return ""
}

// Message definition for exercise 2.
type MessageTwo struct {
	SourceIP         *string                 `protobuf:"bytes,1,req,name=sourceIP" json:"sourceIP,omitempty"`
	SourcePort       *int32                  `protobuf:"varint,2,req,name=sourcePort" json:"sourcePort,omitempty"`
	SourceID         *int32                  `protobuf:"varint,3,req,name=sourceID" json:"sourceID,omitempty"`
	MessageType      *MessageTwo_MessageType `protobuf:"varint,4,req,name=messageType,enum=protobuf.MessageTwo_MessageType,def=1" json:"messageType,omitempty"`
	ControlType      *MessageTwo_ControlType `protobuf:"varint,5,opt,name=controlType,enum=protobuf.MessageTwo_ControlType,def=0" json:"controlType,omitempty"`
	NodeType         *MessageTwo_NodeType    `protobuf:"varint,6,req,name=nodeType,enum=protobuf.MessageTwo_NodeType,def=1" json:"nodeType,omitempty"`
	MessageContent   *string                 `protobuf:"bytes,7,req,name=messageContent" json:"messageContent,omitempty"`
	Timestamp        *string                 `protobuf:"bytes,8,opt,name=timestamp" json:"timestamp,omitempty"`
	XXX_unrecognized []byte                  `json:"-"`
}

func (m *MessageTwo) Reset()         { *m = MessageTwo{} }
func (m *MessageTwo) String() string { return proto.CompactTextString(m) }
func (*MessageTwo) ProtoMessage()    {}

const Default_MessageTwo_MessageType MessageTwo_MessageType = MessageTwo_APPLICATIONMESSAGE
const Default_MessageTwo_ControlType MessageTwo_ControlType = MessageTwo_INITIALIZE
const Default_MessageTwo_NodeType MessageTwo_NodeType = MessageTwo_CUSTOMER

func (m *MessageTwo) GetSourceIP() string {
	if m != nil && m.SourceIP != nil {
		return *m.SourceIP
	}
	return ""
}

func (m *MessageTwo) GetSourcePort() int32 {
	if m != nil && m.SourcePort != nil {
		return *m.SourcePort
	}
	return 0
}

func (m *MessageTwo) GetSourceID() int32 {
	if m != nil && m.SourceID != nil {
		return *m.SourceID
	}
	return 0
}

func (m *MessageTwo) GetMessageType() MessageTwo_MessageType {
	if m != nil && m.MessageType != nil {
		return *m.MessageType
	}
	return Default_MessageTwo_MessageType
}

func (m *MessageTwo) GetControlType() MessageTwo_ControlType {
	if m != nil && m.ControlType != nil {
		return *m.ControlType
	}
	return Default_MessageTwo_ControlType
}

func (m *MessageTwo) GetNodeType() MessageTwo_NodeType {
	if m != nil && m.NodeType != nil {
		return *m.NodeType
	}
	return Default_MessageTwo_NodeType
}

func (m *MessageTwo) GetMessageContent() string {
	if m != nil && m.MessageContent != nil {
		return *m.MessageContent
	}
	return ""
}

func (m *MessageTwo) GetTimestamp() string {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return ""
}

func init() {
	proto.RegisterEnum("protobuf.Nachricht_NachrichtenTyp", Nachricht_NachrichtenTyp_name, Nachricht_NachrichtenTyp_value)
	proto.RegisterEnum("protobuf.Nachricht_KontrollTyp", Nachricht_KontrollTyp_name, Nachricht_KontrollTyp_value)
	proto.RegisterEnum("protobuf.MessageTwo_MessageType", MessageTwo_MessageType_name, MessageTwo_MessageType_value)
	proto.RegisterEnum("protobuf.MessageTwo_ControlType", MessageTwo_ControlType_name, MessageTwo_ControlType_value)
	proto.RegisterEnum("protobuf.MessageTwo_NodeType", MessageTwo_NodeType_name, MessageTwo_NodeType_value)
}
