// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gpx.proto

package main

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ActivityPB struct {
	Id                   uint64   `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Athlete              *Athlete `protobuf:"bytes,2,opt,name=athlete,proto3" json:"athlete,omitempty"`
	Map                  *Map     `protobuf:"bytes,3,opt,name=map,proto3" json:"map,omitempty"`
	Startdttm            string   `protobuf:"bytes,4,opt,name=startdttm,proto3" json:"startdttm,omitempty"`
	Type                 string   `protobuf:"bytes,5,opt,name=type,proto3" json:"type,omitempty"`
	Device               string   `protobuf:"bytes,6,opt,name=device,proto3" json:"device,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ActivityPB) Reset()         { *m = ActivityPB{} }
func (m *ActivityPB) String() string { return proto.CompactTextString(m) }
func (*ActivityPB) ProtoMessage()    {}
func (*ActivityPB) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a293e36f114101d, []int{0}
}

func (m *ActivityPB) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ActivityPB.Unmarshal(m, b)
}
func (m *ActivityPB) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ActivityPB.Marshal(b, m, deterministic)
}
func (m *ActivityPB) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ActivityPB.Merge(m, src)
}
func (m *ActivityPB) XXX_Size() int {
	return xxx_messageInfo_ActivityPB.Size(m)
}
func (m *ActivityPB) XXX_DiscardUnknown() {
	xxx_messageInfo_ActivityPB.DiscardUnknown(m)
}

var xxx_messageInfo_ActivityPB proto.InternalMessageInfo

func (m *ActivityPB) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ActivityPB) GetAthlete() *Athlete {
	if m != nil {
		return m.Athlete
	}
	return nil
}

func (m *ActivityPB) GetMap() *Map {
	if m != nil {
		return m.Map
	}
	return nil
}

func (m *ActivityPB) GetStartdttm() string {
	if m != nil {
		return m.Startdttm
	}
	return ""
}

func (m *ActivityPB) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *ActivityPB) GetDevice() string {
	if m != nil {
		return m.Device
	}
	return ""
}

type Athlete struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Athlete) Reset()         { *m = Athlete{} }
func (m *Athlete) String() string { return proto.CompactTextString(m) }
func (*Athlete) ProtoMessage()    {}
func (*Athlete) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a293e36f114101d, []int{1}
}

func (m *Athlete) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Athlete.Unmarshal(m, b)
}
func (m *Athlete) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Athlete.Marshal(b, m, deterministic)
}
func (m *Athlete) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Athlete.Merge(m, src)
}
func (m *Athlete) XXX_Size() int {
	return xxx_messageInfo_Athlete.Size(m)
}
func (m *Athlete) XXX_DiscardUnknown() {
	xxx_messageInfo_Athlete.DiscardUnknown(m)
}

var xxx_messageInfo_Athlete proto.InternalMessageInfo

func (m *Athlete) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type Map struct {
	Polyline             string   `protobuf:"bytes,1,opt,name=polyline,proto3" json:"polyline,omitempty"`
	Summarypolyline      string   `protobuf:"bytes,2,opt,name=summarypolyline,proto3" json:"summarypolyline,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Map) Reset()         { *m = Map{} }
func (m *Map) String() string { return proto.CompactTextString(m) }
func (*Map) ProtoMessage()    {}
func (*Map) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a293e36f114101d, []int{2}
}

func (m *Map) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Map.Unmarshal(m, b)
}
func (m *Map) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Map.Marshal(b, m, deterministic)
}
func (m *Map) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Map.Merge(m, src)
}
func (m *Map) XXX_Size() int {
	return xxx_messageInfo_Map.Size(m)
}
func (m *Map) XXX_DiscardUnknown() {
	xxx_messageInfo_Map.DiscardUnknown(m)
}

var xxx_messageInfo_Map proto.InternalMessageInfo

func (m *Map) GetPolyline() string {
	if m != nil {
		return m.Polyline
	}
	return ""
}

func (m *Map) GetSummarypolyline() string {
	if m != nil {
		return m.Summarypolyline
	}
	return ""
}

func init() {
	proto.RegisterType((*ActivityPB)(nil), "main.ActivityPB")
	proto.RegisterType((*Athlete)(nil), "main.Athlete")
	proto.RegisterType((*Map)(nil), "main.Map")
}

func init() {
	proto.RegisterFile("gpx.proto", fileDescriptor_6a293e36f114101d)
}

var fileDescriptor_6a293e36f114101d = []byte{
	// 229 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0xcd, 0x4a, 0x03, 0x31,
	0x14, 0x85, 0x49, 0x26, 0x4e, 0xcd, 0x15, 0x15, 0xee, 0x42, 0xe2, 0xcf, 0x62, 0x98, 0x8d, 0x59,
	0xcd, 0x42, 0x9f, 0xa0, 0xee, 0x8a, 0x14, 0x24, 0x6f, 0x10, 0x9b, 0xa0, 0x81, 0x49, 0x27, 0x4c,
	0xaf, 0xc5, 0xbc, 0x93, 0x0f, 0x29, 0x93, 0xa9, 0x15, 0xba, 0xcb, 0xf9, 0xce, 0x17, 0x72, 0x08,
	0xc8, 0x8f, 0xf4, 0xdd, 0xa5, 0x71, 0xa0, 0x01, 0x45, 0xb4, 0x61, 0xdb, 0xfe, 0x30, 0x80, 0xe5,
	0x86, 0xc2, 0x3e, 0x50, 0x7e, 0x7b, 0xc1, 0x2b, 0xe0, 0x2b, 0xa7, 0x58, 0xc3, 0xb4, 0x30, 0x7c,
	0xe5, 0xf0, 0x11, 0x16, 0x96, 0x3e, 0x7b, 0x4f, 0x5e, 0xf1, 0x86, 0xe9, 0x8b, 0xa7, 0xcb, 0x6e,
	0xba, 0xd6, 0x2d, 0x67, 0x68, 0xfe, 0x5a, 0xbc, 0x87, 0x2a, 0xda, 0xa4, 0xaa, 0x22, 0xc9, 0x59,
	0x5a, 0xdb, 0x64, 0x26, 0x8a, 0x0f, 0x20, 0x77, 0x64, 0x47, 0x72, 0x44, 0x51, 0x89, 0x86, 0x69,
	0x69, 0xfe, 0x01, 0x22, 0x08, 0xca, 0xc9, 0xab, 0xb3, 0x52, 0x94, 0x33, 0xde, 0x40, 0xed, 0xfc,
	0x3e, 0x6c, 0xbc, 0xaa, 0x0b, 0x3d, 0xa4, 0xf6, 0x16, 0x16, 0x87, 0xa7, 0xa7, 0xa9, 0x61, 0x9e,
	0x5a, 0x19, 0x1e, 0x5c, 0xfb, 0x0a, 0xd5, 0xda, 0x26, 0xbc, 0x83, 0xf3, 0x34, 0xf4, 0xb9, 0x0f,
	0x5b, 0x5f, 0x4a, 0x69, 0x8e, 0x19, 0x35, 0x5c, 0xef, 0xbe, 0x62, 0xb4, 0x63, 0x3e, 0x2a, 0xbc,
	0x28, 0xa7, 0xf8, 0xbd, 0x2e, 0x7f, 0xf4, 0xfc, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x7b, 0x12, 0xf2,
	0x34, 0x30, 0x01, 0x00, 0x00,
}
