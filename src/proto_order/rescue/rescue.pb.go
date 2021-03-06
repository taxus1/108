// Code generated by protoc-gen-go.
// source: rescue/rescue.proto
// DO NOT EDIT!

/*
Package proto_rescue is a generated protocol buffer package.

It is generated from these files:
	rescue/rescue.proto

It has these top-level messages:
	RepairRescueArgs
	DirectRescueArgs
	RescueCompleteArgs
	RescueGrabArgs
	RescueAcceptArgs
	RescueRejectArgs
*/
package proto_rescue

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

// RepairRescueArgs 维修救援参数
type RepairRescueArgs struct {
	// 下单司机ID，对应数据库user_id
	DriverID int64 `protobuf:"varint,1,opt,name=driverID" json:"driverID,omitempty"`
	// 操作订单时的坐标
	Location   *proto_order.Location    `protobuf:"bytes,2,opt,name=location" json:"location,omitempty"`
	OrderInfos []*proto_order.OrderInfo `protobuf:"bytes,3,rep,name=orderInfos" json:"orderInfos,omitempty"`
}

func (m *RepairRescueArgs) Reset()                    { *m = RepairRescueArgs{} }
func (m *RepairRescueArgs) String() string            { return proto.CompactTextString(m) }
func (*RepairRescueArgs) ProtoMessage()               {}
func (*RepairRescueArgs) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *RepairRescueArgs) GetLocation() *proto_order.Location {
	if m != nil {
		return m.Location
	}
	return nil
}

func (m *RepairRescueArgs) GetOrderInfos() []*proto_order.OrderInfo {
	if m != nil {
		return m.OrderInfos
	}
	return nil
}

// DirectRescueArgs 针对维修厂直接救援参数
type DirectRescueArgs struct {
	// 技工ID
	RepairID int64 `protobuf:"varint,1,opt,name=repairID" json:"repairID,omitempty"`
	// 技工和司机直线距离
	Distance float32 `protobuf:"fixed32,2,opt,name=distance" json:"distance,omitempty"`
	// 救援基本参数
	RescueArgs *RepairRescueArgs `protobuf:"bytes,3,opt,name=rescueArgs" json:"rescueArgs,omitempty"`
}

func (m *DirectRescueArgs) Reset()                    { *m = DirectRescueArgs{} }
func (m *DirectRescueArgs) String() string            { return proto.CompactTextString(m) }
func (*DirectRescueArgs) ProtoMessage()               {}
func (*DirectRescueArgs) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *DirectRescueArgs) GetRescueArgs() *RepairRescueArgs {
	if m != nil {
		return m.RescueArgs
	}
	return nil
}

// RescueCompleteArgs 完善救援订单参数
type RescueCompleteArgs struct {
	// 技工ID
	OrderID int64 `protobuf:"varint,1,opt,name=orderID" json:"orderID,omitempty"`
	// 技工和司机直线距离
	Distance float32  `protobuf:"fixed32,2,opt,name=distance" json:"distance,omitempty"`
	Image    []string `protobuf:"bytes,3,rep,name=image" json:"image,omitempty"`
	// 自定义设置的金额
	CustomerPrice *proto_order.CustomerPrice `protobuf:"bytes,4,opt,name=CustomerPrice" json:"CustomerPrice,omitempty"`
	// 操作订单时的坐标
	Location   *proto_order.Location    `protobuf:"bytes,5,opt,name=location" json:"location,omitempty"`
	OrderInfos []*proto_order.OrderInfo `protobuf:"bytes,6,rep,name=orderInfos" json:"orderInfos,omitempty"`
}

func (m *RescueCompleteArgs) Reset()                    { *m = RescueCompleteArgs{} }
func (m *RescueCompleteArgs) String() string            { return proto.CompactTextString(m) }
func (*RescueCompleteArgs) ProtoMessage()               {}
func (*RescueCompleteArgs) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *RescueCompleteArgs) GetCustomerPrice() *proto_order.CustomerPrice {
	if m != nil {
		return m.CustomerPrice
	}
	return nil
}

func (m *RescueCompleteArgs) GetLocation() *proto_order.Location {
	if m != nil {
		return m.Location
	}
	return nil
}

func (m *RescueCompleteArgs) GetOrderInfos() []*proto_order.OrderInfo {
	if m != nil {
		return m.OrderInfos
	}
	return nil
}

// RescueGrabArgs 救援订单抢单参数
type RescueGrabArgs struct {
	// 技工ID
	OrderID int64 `protobuf:"varint,1,opt,name=orderID" json:"orderID,omitempty"`
	// 技工和司机直线距离
	Distance float32 `protobuf:"fixed32,2,opt,name=distance" json:"distance,omitempty"`
	// 操作订单时的坐标
	Location *proto_order.Location `protobuf:"bytes,5,opt,name=location" json:"location,omitempty"`
}

func (m *RescueGrabArgs) Reset()                    { *m = RescueGrabArgs{} }
func (m *RescueGrabArgs) String() string            { return proto.CompactTextString(m) }
func (*RescueGrabArgs) ProtoMessage()               {}
func (*RescueGrabArgs) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *RescueGrabArgs) GetLocation() *proto_order.Location {
	if m != nil {
		return m.Location
	}
	return nil
}

// RescueAcceptArgs 救援订单接受单参数
type RescueAcceptArgs struct {
	AcceptArgs *RescueGrabArgs `protobuf:"bytes,1,opt,name=acceptArgs" json:"acceptArgs,omitempty"`
}

func (m *RescueAcceptArgs) Reset()                    { *m = RescueAcceptArgs{} }
func (m *RescueAcceptArgs) String() string            { return proto.CompactTextString(m) }
func (*RescueAcceptArgs) ProtoMessage()               {}
func (*RescueAcceptArgs) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *RescueAcceptArgs) GetAcceptArgs() *RescueGrabArgs {
	if m != nil {
		return m.AcceptArgs
	}
	return nil
}

// RescueRejectArgs 救援订单p2p拒绝接单参数
type RescueRejectArgs struct {
	RejectArgs *RescueGrabArgs `protobuf:"bytes,1,opt,name=rejectArgs" json:"rejectArgs,omitempty"`
}

func (m *RescueRejectArgs) Reset()                    { *m = RescueRejectArgs{} }
func (m *RescueRejectArgs) String() string            { return proto.CompactTextString(m) }
func (*RescueRejectArgs) ProtoMessage()               {}
func (*RescueRejectArgs) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *RescueRejectArgs) GetRejectArgs() *RescueGrabArgs {
	if m != nil {
		return m.RejectArgs
	}
	return nil
}

func init() {
	proto.RegisterType((*RepairRescueArgs)(nil), "proto.rescue.RepairRescueArgs")
	proto.RegisterType((*DirectRescueArgs)(nil), "proto.rescue.DirectRescueArgs")
	proto.RegisterType((*RescueCompleteArgs)(nil), "proto.rescue.RescueCompleteArgs")
	proto.RegisterType((*RescueGrabArgs)(nil), "proto.rescue.RescueGrabArgs")
	proto.RegisterType((*RescueAcceptArgs)(nil), "proto.rescue.RescueAcceptArgs")
	proto.RegisterType((*RescueRejectArgs)(nil), "proto.rescue.RescueRejectArgs")
}

func init() { proto.RegisterFile("rescue/rescue.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 347 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xa4, 0x92, 0xdd, 0x4a, 0xc3, 0x40,
	0x10, 0x85, 0x49, 0x63, 0x6b, 0x9d, 0xa8, 0x94, 0xf8, 0x43, 0x28, 0x22, 0x25, 0x57, 0xbd, 0x8a,
	0x58, 0xc1, 0x2b, 0x11, 0xa5, 0x05, 0x29, 0x08, 0x96, 0x7d, 0x01, 0x49, 0xb7, 0x63, 0x59, 0x69,
	0xbb, 0x61, 0xb2, 0x15, 0x7c, 0x05, 0xc1, 0x07, 0xf0, 0x6d, 0xc5, 0x59, 0xb3, 0x4d, 0x7a, 0x21,
	0x55, 0xaf, 0x96, 0x93, 0x39, 0x7b, 0xf6, 0xe3, 0x4c, 0xe0, 0x80, 0x30, 0x97, 0x4b, 0x3c, 0xb3,
	0x47, 0x92, 0x91, 0x36, 0x3a, 0xdc, 0xe5, 0x23, 0xb1, 0xdf, 0xda, 0x81, 0xa6, 0x09, 0x92, 0x1d,
	0xc5, 0x1f, 0x1e, 0xb4, 0x04, 0x66, 0xa9, 0x22, 0xc1, 0xd3, 0x5b, 0x9a, 0xe6, 0x61, 0x1b, 0x9a,
	0x13, 0x52, 0x2f, 0x48, 0xc3, 0x41, 0xe4, 0x75, 0xbc, 0xae, 0x2f, 0x9c, 0x0e, 0xcf, 0xa1, 0x39,
	0xd3, 0x32, 0x35, 0x4a, 0x2f, 0xa2, 0x5a, 0xc7, 0xeb, 0x06, 0xbd, 0x23, 0x1b, 0xf5, 0x68, 0x63,
	0xef, 0xbf, 0x87, 0xc2, 0xd9, 0xc2, 0x4b, 0x00, 0x9e, 0x0d, 0x17, 0x4f, 0x3a, 0x8f, 0xfc, 0x8e,
	0xdf, 0x0d, 0x7a, 0xc7, 0x95, 0x4b, 0x0f, 0xc5, 0x58, 0x94, 0x9c, 0xf1, 0x9b, 0x07, 0xad, 0x81,
	0x22, 0x94, 0xa6, 0xca, 0x46, 0xcc, 0xbb, 0x62, 0x2b, 0x34, 0x73, 0xab, 0xdc, 0xa4, 0x0b, 0x89,
	0xcc, 0x56, 0x13, 0x4e, 0x87, 0xd7, 0x00, 0xe4, 0x52, 0x22, 0x9f, 0xc9, 0x4f, 0x93, 0x72, 0x31,
	0xc9, 0x7a, 0x0f, 0xa2, 0x74, 0x23, 0x7e, 0xaf, 0x41, 0x68, 0x47, 0x7d, 0x3d, 0xcf, 0x66, 0x68,
	0x2c, 0x4e, 0x04, 0xdb, 0x96, 0xb8, 0xa0, 0x29, 0xe4, 0x8f, 0x30, 0x87, 0x50, 0x57, 0xf3, 0x74,
	0x8a, 0x5c, 0xc6, 0x8e, 0xb0, 0x22, 0xbc, 0x81, 0xbd, 0xfe, 0x32, 0x37, 0x7a, 0x8e, 0x34, 0x22,
	0x25, 0x31, 0xda, 0x62, 0xca, 0x76, 0xa5, 0xaa, 0x8a, 0x43, 0x54, 0x2f, 0x54, 0x96, 0x53, 0xff,
	0xcb, 0x72, 0x1a, 0x1b, 0x2f, 0xe7, 0x15, 0xf6, 0x6d, 0x1d, 0x77, 0x94, 0x8e, 0xff, 0x51, 0xc5,
	0xef, 0x91, 0xe3, 0xd1, 0xd7, 0x2f, 0xcb, 0x8b, 0x91, 0x12, 0x33, 0xc3, 0x8f, 0x5f, 0x01, 0xa4,
	0x4e, 0xf1, 0xfb, 0x41, 0xef, 0x64, 0x7d, 0xbd, 0x65, 0x5c, 0x51, 0xf2, 0xaf, 0x12, 0x05, 0x3e,
	0xa3, 0x74, 0x89, 0xe4, 0xd4, 0x66, 0x89, 0x2b, 0xff, 0xb8, 0xc1, 0xc6, 0x8b, 0xcf, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x1e, 0x1d, 0x38, 0x1d, 0x90, 0x03, 0x00, 0x00,
}
