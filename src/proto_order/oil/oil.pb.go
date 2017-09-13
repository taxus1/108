// Code generated by protoc-gen-go.
// source: oil/oil.proto
// DO NOT EDIT!

/*
Package proto_oil is a generated protocol buffer package.

It is generated from these files:
	oil/oil.proto

It has these top-level messages:
	OilCreateArgs
*/
package proto_oil

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import proto_order "ocenter/src/proto_order"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// OilCreateArgs 创建加油订单参数
type OilCreateArgs struct {
	// 下单基本
	RepairID int64 `protobuf:"varint,1,opt,name=repairID" json:"repairID,omitempty"`
	// 操作订单时的坐标
	Location *proto_order.Location `protobuf:"bytes,2,opt,name=location" json:"location,omitempty"`
	// 油枪编号
	OilGun    string                 `protobuf:"bytes,3,opt,name=oilGun" json:"oilGun,omitempty"`
	OrderInfo *proto_order.OrderInfo `protobuf:"bytes,4,opt,name=orderInfo" json:"orderInfo,omitempty"`
}

func (m *OilCreateArgs) Reset()                    { *m = OilCreateArgs{} }
func (m *OilCreateArgs) String() string            { return proto.CompactTextString(m) }
func (*OilCreateArgs) ProtoMessage()               {}
func (*OilCreateArgs) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *OilCreateArgs) GetLocation() *proto_order.Location {
	if m != nil {
		return m.Location
	}
	return nil
}

func (m *OilCreateArgs) GetOrderInfo() *proto_order.OrderInfo {
	if m != nil {
		return m.OrderInfo
	}
	return nil
}

func init() {
	proto.RegisterType((*OilCreateArgs)(nil), "proto.oil.OilCreateArgs")
}

func init() { proto.RegisterFile("oil/oil.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 173 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0xcf, 0xcc, 0xd1,
	0xcf, 0xcf, 0xcc, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x04, 0x53, 0x7a, 0xf9, 0x99,
	0x39, 0x52, 0xdc, 0xf9, 0x45, 0x29, 0xa9, 0x45, 0x10, 0x71, 0xa5, 0x35, 0x8c, 0x5c, 0xbc, 0xfe,
	0x99, 0x39, 0xce, 0x45, 0xa9, 0x89, 0x25, 0xa9, 0x8e, 0x45, 0xe9, 0xc5, 0x42, 0x52, 0x5c, 0x1c,
	0x45, 0xa9, 0x05, 0x89, 0x99, 0x45, 0x9e, 0x2e, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0xcc, 0x41, 0x70,
	0xbe, 0x90, 0x21, 0x17, 0x47, 0x4e, 0x7e, 0x72, 0x62, 0x49, 0x66, 0x7e, 0x9e, 0x04, 0x93, 0x02,
	0xa3, 0x06, 0xb7, 0x91, 0x28, 0xc4, 0x9c, 0x78, 0x88, 0x99, 0x3e, 0x50, 0xc9, 0x20, 0xb8, 0x32,
	0x21, 0x31, 0x2e, 0xb6, 0xfc, 0xcc, 0x1c, 0xf7, 0xd2, 0x3c, 0x09, 0x66, 0x05, 0x46, 0x0d, 0xce,
	0x20, 0x28, 0x4f, 0xc8, 0x84, 0x8b, 0x13, 0xac, 0xc7, 0x33, 0x2f, 0x2d, 0x5f, 0x82, 0x05, 0x6c,
	0x96, 0x18, 0x8a, 0x59, 0xfe, 0x30, 0xd9, 0x20, 0x84, 0xc2, 0x24, 0x36, 0xb0, 0x0a, 0x63, 0x40,
	0x00, 0x00, 0x00, 0xff, 0xff, 0xb4, 0x11, 0x4e, 0x88, 0xde, 0x00, 0x00, 0x00,
}