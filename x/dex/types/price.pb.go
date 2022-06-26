// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dex/price.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type Price struct {
	SnapshotTimestampInSeconds uint64                                 `protobuf:"varint,1,opt,name=snapshotTimestampInSeconds,proto3" json:"snapshot_timestamp_in_seconds"`
	Price                      github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=price,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"price" yaml:"price"`
	Pair                       *Pair                                  `protobuf:"bytes,3,opt,name=pair,proto3" json:"pair"`
}

func (m *Price) Reset()         { *m = Price{} }
func (m *Price) String() string { return proto.CompactTextString(m) }
func (*Price) ProtoMessage()    {}
func (*Price) Descriptor() ([]byte, []int) {
	return fileDescriptor_bd5d1c9d490efb8c, []int{0}
}
func (m *Price) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Price) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Price.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Price) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Price.Merge(m, src)
}
func (m *Price) XXX_Size() int {
	return m.Size()
}
func (m *Price) XXX_DiscardUnknown() {
	xxx_messageInfo_Price.DiscardUnknown(m)
}

var xxx_messageInfo_Price proto.InternalMessageInfo

func (m *Price) GetSnapshotTimestampInSeconds() uint64 {
	if m != nil {
		return m.SnapshotTimestampInSeconds
	}
	return 0
}

func (m *Price) GetPair() *Pair {
	if m != nil {
		return m.Pair
	}
	return nil
}

func init() {
	proto.RegisterType((*Price)(nil), "seiprotocol.seichain.dex.Price")
}

func init() { proto.RegisterFile("dex/price.proto", fileDescriptor_bd5d1c9d490efb8c) }

var fileDescriptor_bd5d1c9d490efb8c = []byte{
	// 316 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4f, 0x49, 0xad, 0xd0,
	0x2f, 0x28, 0xca, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x92, 0x28, 0x4e, 0xcd,
	0x04, 0xb3, 0x92, 0xf3, 0x73, 0xf4, 0x8a, 0x53, 0x33, 0x93, 0x33, 0x12, 0x33, 0xf3, 0xf4, 0x52,
	0x52, 0x2b, 0xa4, 0x44, 0xd2, 0xf3, 0xd3, 0xf3, 0xc1, 0x52, 0xfa, 0x20, 0x16, 0x44, 0xbd, 0x14,
	0x1f, 0xd8, 0x80, 0xc4, 0xcc, 0x22, 0x08, 0x5f, 0xa9, 0x9d, 0x89, 0x8b, 0x35, 0x00, 0x64, 0x9e,
	0x50, 0x22, 0x97, 0x54, 0x71, 0x5e, 0x62, 0x41, 0x71, 0x46, 0x7e, 0x49, 0x48, 0x66, 0x6e, 0x6a,
	0x71, 0x49, 0x62, 0x6e, 0x81, 0x67, 0x5e, 0x70, 0x6a, 0x72, 0x7e, 0x5e, 0x4a, 0xb1, 0x04, 0xa3,
	0x02, 0xa3, 0x06, 0x8b, 0x93, 0xe2, 0xab, 0x7b, 0xf2, 0xb2, 0x30, 0x55, 0xf1, 0x25, 0x30, 0x65,
	0xf1, 0x99, 0x79, 0xf1, 0xc5, 0x10, 0x85, 0x41, 0x78, 0x0c, 0x11, 0x8a, 0xe5, 0x62, 0x05, 0xbb,
	0x5d, 0x82, 0x49, 0x81, 0x51, 0x83, 0xd3, 0xc9, 0xfd, 0xc4, 0x3d, 0x79, 0x86, 0x5b, 0xf7, 0xe4,
	0xd5, 0xd2, 0x33, 0x4b, 0x32, 0x4a, 0x93, 0xf4, 0x92, 0xf3, 0x73, 0xf5, 0x93, 0xf3, 0x8b, 0x73,
	0xf3, 0x8b, 0xa1, 0x94, 0x6e, 0x71, 0x4a, 0xb6, 0x7e, 0x49, 0x65, 0x41, 0x6a, 0xb1, 0x9e, 0x4b,
	0x6a, 0xf2, 0xab, 0x7b, 0xf2, 0x10, 0xed, 0x9f, 0xee, 0xc9, 0xf3, 0x54, 0x26, 0xe6, 0xe6, 0x58,
	0x29, 0x81, 0xb9, 0x4a, 0x41, 0x10, 0x61, 0x21, 0x1b, 0x2e, 0x16, 0x90, 0xcf, 0x24, 0x98, 0x15,
	0x18, 0x35, 0xb8, 0x8d, 0xe4, 0xf4, 0x70, 0x05, 0x8d, 0x5e, 0x40, 0x62, 0x66, 0x91, 0x13, 0xc7,
	0xab, 0x7b, 0xf2, 0x60, 0xf5, 0x41, 0x60, 0xd2, 0xc9, 0xfd, 0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f,
	0xe4, 0x18, 0x1f, 0x3c, 0x92, 0x63, 0x9c, 0xf0, 0x58, 0x8e, 0xe1, 0xc2, 0x63, 0x39, 0x86, 0x1b,
	0x8f, 0xe5, 0x18, 0xa2, 0x74, 0x91, 0xdc, 0x57, 0x9c, 0x9a, 0xa9, 0x0b, 0x33, 0x14, 0xcc, 0x01,
	0x9b, 0xaa, 0x5f, 0xa1, 0x0f, 0x0a, 0x57, 0xb0, 0x53, 0x93, 0xd8, 0xc0, 0xf2, 0xc6, 0x80, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x1a, 0xe3, 0x94, 0xe7, 0xac, 0x01, 0x00, 0x00,
}

func (m *Price) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Price) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Price) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pair != nil {
		{
			size, err := m.Pair.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPrice(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	{
		size := m.Price.Size()
		i -= size
		if _, err := m.Price.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPrice(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.SnapshotTimestampInSeconds != 0 {
		i = encodeVarintPrice(dAtA, i, uint64(m.SnapshotTimestampInSeconds))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintPrice(dAtA []byte, offset int, v uint64) int {
	offset -= sovPrice(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Price) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SnapshotTimestampInSeconds != 0 {
		n += 1 + sovPrice(uint64(m.SnapshotTimestampInSeconds))
	}
	l = m.Price.Size()
	n += 1 + l + sovPrice(uint64(l))
	if m.Pair != nil {
		l = m.Pair.Size()
		n += 1 + l + sovPrice(uint64(l))
	}
	return n
}

func sovPrice(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPrice(x uint64) (n int) {
	return sovPrice(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Price) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPrice
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Price: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Price: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SnapshotTimestampInSeconds", wireType)
			}
			m.SnapshotTimestampInSeconds = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPrice
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SnapshotTimestampInSeconds |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Price", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPrice
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPrice
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPrice
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Price.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pair", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPrice
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPrice
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPrice
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pair == nil {
				m.Pair = &Pair{}
			}
			if err := m.Pair.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPrice(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPrice
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipPrice(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPrice
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowPrice
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowPrice
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthPrice
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupPrice
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthPrice
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthPrice        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPrice          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupPrice = fmt.Errorf("proto: unexpected end of group")
)
