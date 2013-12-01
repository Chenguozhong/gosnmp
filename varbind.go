package gosnmp

import (
	"fmt"
	"net"
)

type Varbind interface {
	//encodeValue returns the number of bytes written to the encoder
	encodeValue(encoder *berEncoder) (int, error)
	decodeValue(decoder *berDecoder, valueLength int) error
	getOid() ObjectIdentifier
	setOid(oid ObjectIdentifier)
}

func (encoder *berEncoder) encodeVarbind(vb Varbind) (int, error) {
	header := encoder.newHeader(SEQUENCE)
	oidLen, err := encoder.encodeObjectIdentifier(vb.getOid())
	if err != nil {
		return 0, err
	}
	valLen, err := vb.encodeValue(encoder)
	if err != nil {
		return 0, err
	}
	_, marshalledLen := header.setContentLength(oidLen + valLen)
	return marshalledLen, nil
}

type baseVarbind struct {
	oid ObjectIdentifier
}

func (vb *baseVarbind) getOid() ObjectIdentifier {
	return vb.oid
}

func (vb *baseVarbind) setOid(oid ObjectIdentifier) {
	vb.oid = oid
}

type IntegerVarbind struct { // type 0x02
	baseVarbind
	val int32
}

func NewIntegerVarbind(oid ObjectIdentifier, val int32) *IntegerVarbind {
	vb := new(IntegerVarbind)
	vb.oid = oid
	return vb
}

func (vb *IntegerVarbind) encodeValue(encoder *berEncoder) (int, error) {
	return encoder.encodeInteger(int64(vb.val)), nil
}

func (vb *IntegerVarbind) decodeValue(decoder *berDecoder, valueLength int) (err error) {
	vb.val, err = decoder.decodeInt32(valueLength)
	return
}

// BitStringVarbind is the structure to use when you want an ASN.1 BIT STRING type. A bit string is padded up to the nearest byte in memory
// and the number of valid bits is recorded. Padding bits will be zero
type BitStringVarbind struct { // type 0x03
	baseVarbind
	val *BitString
}

func NewBitStringVarbind(oid ObjectIdentifier, val *BitString) *BitStringVarbind {
	vb := new(BitStringVarbind)
	vb.oid = oid
	return vb
}

func (vb *BitStringVarbind) encodeValue(encoder *berEncoder) (int, error) {
	return encoder.encodeBitString(vb.val), nil
}

func (vb *BitStringVarbind) decodeValue(decoder *berDecoder, valueLength int) (err error) {
	vb.val, err = decoder.decodeBitString(valueLength)
	return
}

type OctetStringVarbind struct { // type 0x04
	baseVarbind
	val []byte
}

func NewOctetStringVarbind(oid ObjectIdentifier, val []byte) *OctetStringVarbind {
	vb := new(OctetStringVarbind)
	vb.oid = oid
	return vb
}

func (vb *OctetStringVarbind) encodeValue(encoder *berEncoder) (int, error) {
	return encoder.encodeOctetString(vb.val), nil
}

func (vb *OctetStringVarbind) decodeValue(decoder *berDecoder, valueLength int) (err error) {
	vb.val = make([]byte, valueLength)
	if numRead, err := decoder.Read(vb.val); err != nil || numRead != valueLength {
		return fmt.Errorf("Couldn't decode octet string of length %d. Number of bytes read from stream: %d, err: %s", valueLength, numRead, err)
	}
	return
}

type NullVarbind struct { // type 0x05
	baseVarbind
}

func NewNullVarbind(oid ObjectIdentifier) *NullVarbind {
	vb := new(NullVarbind)
	vb.oid = oid
	return vb
}

func (vb *NullVarbind) encodeValue(encoder *berEncoder) (int, error) {
	return encoder.encodeNull(), nil
}

func (vb *NullVarbind) decodeValue(decoder *berDecoder, valueLength int) (err error) {
	if valueLength != 0 {
		return fmt.Errorf("Non-zero value length found for NULL varbind: %d", valueLength)
	}
	return
}

type ObjectIdentifierVarbind struct { // type 0x06
	baseVarbind
	val ObjectIdentifier
}

func NewObjectIdentifierVarbind(oid ObjectIdentifier, val ObjectIdentifier) *ObjectIdentifierVarbind {
	vb := new(ObjectIdentifierVarbind)
	vb.oid = oid
	return vb
}

func (vb *ObjectIdentifierVarbind) encodeValue(encoder *berEncoder) (int, error) {
	return encoder.encodeObjectIdentifier(vb.val)
}

func (vb *ObjectIdentifierVarbind) decodeValue(decoder *berDecoder, valueLength int) (err error) {
	vb.val, err = decoder.decodeObjectIdentifier(valueLength)
	return
}

type IPv4AddressVarbind struct { // type 0x40
	baseVarbind
	val net.IP
}

func NewIPv4AddressVarbind(oid ObjectIdentifier, val net.IP) *IPv4AddressVarbind {
	vb := new(IPv4AddressVarbind)
	vb.oid = oid

	return vb
}

func (vb *IPv4AddressVarbind) encodeValue(encoder *berEncoder) (int, error) {
	return encoder.encodeIPv4Address(vb.val)
}

func (vb *IPv4AddressVarbind) decodeValue(decoder *berDecoder, valueLength int) (err error) {
	vb.val, err = decoder.decodeIPv4Address(valueLength)
	return
}

type Counter32Varbind struct { // type 0x41
	baseVarbind
	val uint32
}

func NewCounter32Varbind(oid ObjectIdentifier) *Counter32Varbind {
	vb := new(Counter32Varbind)
	vb.oid = oid
	return vb
}

type Gauge32Varbind struct { // type 0x42
	baseVarbind
	val uint32
}

func NewGauge32Varbind(oid ObjectIdentifier) *Gauge32Varbind {
	vb := new(Gauge32Varbind)
	vb.oid = oid
	return vb
}

type TimeTicksVarbind struct { // type 0x43
	baseVarbind
	val uint32
}

func NewTimeTicksVarbind(oid ObjectIdentifier) *TimeTicksVarbind {
	vb := new(TimeTicksVarbind)
	vb.oid = oid
	return vb
}

type OpaqueVarbind struct { // type 0x44
	baseVarbind
	val []byte
}

func NewOpaqueVarbind(oid ObjectIdentifier) *OpaqueVarbind {
	vb := new(OpaqueVarbind)
	vb.oid = oid
	return vb
}

type NsapAddressVarbind struct { // type 0x45
	baseVarbind
	val [6]byte
}

func NewNsapAddressVarbind(oid ObjectIdentifier) *NsapAddressVarbind {
	vb := new(NsapAddressVarbind)
	vb.oid = oid
	return vb
}

type Counter64Varbind struct { // type 0x46
	baseVarbind
	val uint64
}

func NewCounter64Varbind(oid ObjectIdentifier) *Counter64Varbind {
	vb := new(Counter64Varbind)
	vb.oid = oid
	return vb
}

type Uint32Varbind struct { // type 0x47
	baseVarbind
	val uint32
}

func NewUint32Varbind(oid ObjectIdentifier) *Uint32Varbind {
	vb := new(Uint32Varbind)
	vb.oid = oid
	return vb
}

func decodeVarbind(decoder *berDecoder) (varbind Varbind, err error) {
	varbindHeaderType, varbindLength, err := decoder.decodeHeader()
	if err != nil {
		return nil, fmt.Errorf("Unable to decode varbind header - err: %s", err)
	}
	startDecoderLen := decoder.Len()
	if varbindHeaderType != SEQUENCE {
		return nil, fmt.Errorf("Invalid varbind header type 0x%x - not 0x%x", varbindHeaderType, SEQUENCE)
	}
	oid, err := decoder.decodeObjectIdentifierWithHeader()
	if err != nil {
		return nil, fmt.Errorf("Failed to decode object identifier - err: %s", err)
	}
	valueType, value, err := decoder.decodeValue()
	if err != nil {
		return nil, fmt.Errorf("Unable to decode value header - err: %s", err)
	}
	switch valueType {
	case INTEGER:
		varbind = NewIntegerVarbind(oid, value.(int32))
	case BIT_STRING:
		varbind = NewBitStringVarbind(oid, value.(*BitString))
	case OCTET_STRING:
		varbind = NewOctetStringVarbind(oid, value.(OctectString))
	case NULL:
		varbind = NewNullVarbind(oid)
	case OBJECT_IDENTIFIER:
		varbind = NewObjectIdentifierVarbind(oid, value.(ObjectIdentifier))
	case IP_ADDRESS:
		varbind = NewIPv4AddressVarbind(oid, value.(net.IP))
	// case COUNTER_32:
	// 	varbind = NewCounter32Varbind(oid)
	// case GAUGE_32:
	// 	varbind = NewGauge32Varbind(oid)
	// case TIME_TICKS:
	// 	varbind = NewTimeTicksVarbind(oid)
	// case OPAQUE:
	// 	varbind = NewOpaqueVarbind(oid)
	// case COUNTER_64:
	// 	varbind = NewCounter64Varbind(oid)
	// case UINT_32:
	// 	varbind = NewUint32Varbind(oid)
	default:
		return nil, fmt.Errorf("Unknown value type 0x%x", valueType)
	}
	if startDecoderLen-decoder.Len() != varbindLength {
		return nil, fmt.Errorf("Decoding varbind consumed too many bytes. Expected: %d, actual: %d", varbindLength, startDecoderLen-decoder.Len())
	}
	return
}
