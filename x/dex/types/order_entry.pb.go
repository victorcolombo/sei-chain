// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dex/order_entry.proto

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

type OrderEntry struct {
	Price             github_com_cosmos_cosmos_sdk_types.Dec   `protobuf:"bytes,1,opt,name=price,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"price" yaml:"price"`
	Quantity          github_com_cosmos_cosmos_sdk_types.Dec   `protobuf:"bytes,2,opt,name=quantity,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"quantity" yaml:"quantity"`
	AllocationCreator []string                                 `protobuf:"bytes,3,rep,name=allocationCreator,proto3" json:"allocation_creator"`
	Allocation        []github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,4,rep,name=allocation,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"allocation" yaml:"allocation"`
	PriceDenom        Denom                                    `protobuf:"varint,5,opt,name=priceDenom,proto3,enum=seiprotocol.seichain.dex.Denom" json:"price_denom"`
	AssetDenom        Denom                                    `protobuf:"varint,6,opt,name=assetDenom,proto3,enum=seiprotocol.seichain.dex.Denom" json:"asset_denom"`
}

func (m *OrderEntry) Reset()         { *m = OrderEntry{} }
func (m *OrderEntry) String() string { return proto.CompactTextString(m) }
func (*OrderEntry) ProtoMessage()    {}
func (*OrderEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_25878922effe12c2, []int{0}
}
func (m *OrderEntry) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderEntry.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OrderEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderEntry.Merge(m, src)
}
func (m *OrderEntry) XXX_Size() int {
	return m.Size()
}
func (m *OrderEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderEntry.DiscardUnknown(m)
}

var xxx_messageInfo_OrderEntry proto.InternalMessageInfo

func (m *OrderEntry) GetAllocationCreator() []string {
	if m != nil {
		return m.AllocationCreator
	}
	return nil
}

func (m *OrderEntry) GetPriceDenom() Denom {
	if m != nil {
		return m.PriceDenom
	}
	return Denom_SEI
}

func (m *OrderEntry) GetAssetDenom() Denom {
	if m != nil {
		return m.AssetDenom
	}
	return Denom_SEI
}

func init() {
	proto.RegisterType((*OrderEntry)(nil), "seiprotocol.seichain.dex.OrderEntry")
}

func init() { proto.RegisterFile("dex/order_entry.proto", fileDescriptor_25878922effe12c2) }

var fileDescriptor_25878922effe12c2 = []byte{
	// 392 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0xcf, 0x6a, 0xea, 0x40,
	0x18, 0xc5, 0x93, 0xeb, 0x1f, 0xae, 0x73, 0x2f, 0x15, 0x43, 0x5b, 0x82, 0x8b, 0x8c, 0x64, 0x51,
	0xdc, 0x98, 0x40, 0xbb, 0xeb, 0x32, 0xb5, 0xb8, 0x2a, 0x6d, 0xb3, 0x2c, 0x14, 0x89, 0x93, 0x41,
	0x87, 0x26, 0x19, 0xcd, 0x8c, 0x60, 0xde, 0xa2, 0xef, 0xd2, 0x97, 0x70, 0xe9, 0xb2, 0x74, 0x31,
	0x14, 0xdd, 0xb9, 0xf4, 0x09, 0x4a, 0x66, 0xaa, 0x09, 0x94, 0x2e, 0x5c, 0x39, 0x73, 0xce, 0x77,
	0x7e, 0x67, 0xf8, 0x0c, 0x38, 0x0b, 0xf1, 0xc2, 0xa5, 0x69, 0x88, 0xd3, 0x21, 0x4e, 0x78, 0x9a,
	0x39, 0xd3, 0x94, 0x72, 0x6a, 0x98, 0x0c, 0x13, 0x79, 0x42, 0x34, 0x72, 0x18, 0x26, 0x68, 0x12,
	0x90, 0xc4, 0x09, 0xf1, 0xa2, 0x7d, 0x3a, 0xa6, 0x63, 0x2a, 0x2d, 0x37, 0x3f, 0xa9, 0xf9, 0x76,
	0x33, 0xc7, 0xe0, 0x64, 0x1e, 0x33, 0x25, 0xd8, 0x6f, 0x55, 0x00, 0xee, 0x73, 0xec, 0x6d, 0x4e,
	0x35, 0x9e, 0x41, 0x6d, 0x9a, 0x12, 0x84, 0x4d, 0xbd, 0xa3, 0x77, 0x1b, 0xde, 0x60, 0x29, 0xa0,
	0xf6, 0x21, 0xe0, 0xc5, 0x98, 0xf0, 0xc9, 0x7c, 0xe4, 0x20, 0x1a, 0xbb, 0x88, 0xb2, 0x98, 0xb2,
	0xef, 0x9f, 0x1e, 0x0b, 0x5f, 0x5c, 0x9e, 0x4d, 0x31, 0x73, 0xfa, 0x18, 0x6d, 0x05, 0x54, 0xf1,
	0x9d, 0x80, 0xff, 0xb3, 0x20, 0x8e, 0xae, 0x6d, 0x79, 0xb5, 0x7d, 0x25, 0x1b, 0x04, 0xfc, 0x9d,
	0xcd, 0x83, 0x84, 0x13, 0x9e, 0x99, 0x7f, 0x64, 0xc3, 0xdd, 0xd1, 0x0d, 0x07, 0xc2, 0x4e, 0xc0,
	0xa6, 0x2a, 0xd9, 0x2b, 0xb6, 0x7f, 0x30, 0x8d, 0x3e, 0x68, 0x05, 0x51, 0x44, 0x51, 0xc0, 0x09,
	0x4d, 0x6e, 0x52, 0x1c, 0x70, 0x9a, 0x9a, 0x95, 0x4e, 0xa5, 0xdb, 0xf0, 0xce, 0xb7, 0x02, 0x1a,
	0x85, 0x39, 0x44, 0xca, 0xf5, 0x7f, 0x06, 0x8c, 0x19, 0x00, 0x85, 0x68, 0x56, 0x65, 0xfc, 0xf1,
	0xe8, 0x27, 0x97, 0x18, 0x3b, 0x01, 0x5b, 0xea, 0xd1, 0x85, 0x66, 0xfb, 0xa5, 0x01, 0xe3, 0x01,
	0x00, 0xb9, 0xac, 0x3e, 0x4e, 0x68, 0x6c, 0xd6, 0x3a, 0x7a, 0xf7, 0xe4, 0x12, 0x3a, 0xbf, 0xfd,
	0xcf, 0x8e, 0x1c, 0xf3, 0x9a, 0x5b, 0x01, 0xff, 0xc9, 0xd8, 0x30, 0xcc, 0x05, 0xbf, 0xc4, 0xc8,
	0x89, 0x01, 0x63, 0x98, 0x2b, 0x62, 0xfd, 0x08, 0xa2, 0x8c, 0xed, 0x89, 0x05, 0xc3, 0x1b, 0x2c,
	0xd7, 0x96, 0xbe, 0x5a, 0x5b, 0xfa, 0xe7, 0xda, 0xd2, 0x5f, 0x37, 0x96, 0xb6, 0xda, 0x58, 0xda,
	0xfb, 0xc6, 0xd2, 0x9e, 0x7a, 0xa5, 0xa5, 0x30, 0x4c, 0x7a, 0xfb, 0x0a, 0x79, 0x91, 0x1d, 0xee,
	0xc2, 0xcd, 0x3f, 0x42, 0xb9, 0x9f, 0x51, 0x5d, 0xfa, 0x57, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff,
	0x26, 0x8d, 0x13, 0xd0, 0xdf, 0x02, 0x00, 0x00,
}

func (m *OrderEntry) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderEntry) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OrderEntry) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.AssetDenom != 0 {
		i = encodeVarintOrderEntry(dAtA, i, uint64(m.AssetDenom))
		i--
		dAtA[i] = 0x30
	}
	if m.PriceDenom != 0 {
		i = encodeVarintOrderEntry(dAtA, i, uint64(m.PriceDenom))
		i--
		dAtA[i] = 0x28
	}
	if len(m.Allocation) > 0 {
		for iNdEx := len(m.Allocation) - 1; iNdEx >= 0; iNdEx-- {
			{
				size := m.Allocation[iNdEx].Size()
				i -= size
				if _, err := m.Allocation[iNdEx].MarshalTo(dAtA[i:]); err != nil {
					return 0, err
				}
				i = encodeVarintOrderEntry(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.AllocationCreator) > 0 {
		for iNdEx := len(m.AllocationCreator) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.AllocationCreator[iNdEx])
			copy(dAtA[i:], m.AllocationCreator[iNdEx])
			i = encodeVarintOrderEntry(dAtA, i, uint64(len(m.AllocationCreator[iNdEx])))
			i--
			dAtA[i] = 0x1a
		}
	}
	{
		size := m.Quantity.Size()
		i -= size
		if _, err := m.Quantity.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintOrderEntry(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.Price.Size()
		i -= size
		if _, err := m.Price.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintOrderEntry(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintOrderEntry(dAtA []byte, offset int, v uint64) int {
	offset -= sovOrderEntry(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *OrderEntry) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Price.Size()
	n += 1 + l + sovOrderEntry(uint64(l))
	l = m.Quantity.Size()
	n += 1 + l + sovOrderEntry(uint64(l))
	if len(m.AllocationCreator) > 0 {
		for _, s := range m.AllocationCreator {
			l = len(s)
			n += 1 + l + sovOrderEntry(uint64(l))
		}
	}
	if len(m.Allocation) > 0 {
		for _, e := range m.Allocation {
			l = e.Size()
			n += 1 + l + sovOrderEntry(uint64(l))
		}
	}
	if m.PriceDenom != 0 {
		n += 1 + sovOrderEntry(uint64(m.PriceDenom))
	}
	if m.AssetDenom != 0 {
		n += 1 + sovOrderEntry(uint64(m.AssetDenom))
	}
	return n
}

func sovOrderEntry(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozOrderEntry(x uint64) (n int) {
	return sovOrderEntry(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *OrderEntry) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowOrderEntry
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
			return fmt.Errorf("proto: OrderEntry: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OrderEntry: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Price", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrderEntry
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
				return ErrInvalidLengthOrderEntry
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOrderEntry
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Price.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Quantity", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrderEntry
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
				return ErrInvalidLengthOrderEntry
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOrderEntry
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Quantity.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AllocationCreator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrderEntry
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
				return ErrInvalidLengthOrderEntry
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOrderEntry
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AllocationCreator = append(m.AllocationCreator, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Allocation", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrderEntry
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
				return ErrInvalidLengthOrderEntry
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOrderEntry
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_cosmos_cosmos_sdk_types.Dec
			m.Allocation = append(m.Allocation, v)
			if err := m.Allocation[len(m.Allocation)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PriceDenom", wireType)
			}
			m.PriceDenom = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrderEntry
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PriceDenom |= Denom(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AssetDenom", wireType)
			}
			m.AssetDenom = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrderEntry
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AssetDenom |= Denom(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipOrderEntry(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthOrderEntry
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
func skipOrderEntry(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowOrderEntry
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
					return 0, ErrIntOverflowOrderEntry
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
					return 0, ErrIntOverflowOrderEntry
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
				return 0, ErrInvalidLengthOrderEntry
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupOrderEntry
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthOrderEntry
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthOrderEntry        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowOrderEntry          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupOrderEntry = fmt.Errorf("proto: unexpected end of group")
)
