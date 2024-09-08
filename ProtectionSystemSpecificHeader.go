package sofia

import (
	"154.pages.dev/sofia/box"
	"encoding/binary"
	"io"
)

// dashif.org/identifiers/content_protection
func (p ProtectionSystemSpecificHeader) Widevine() bool {
	return p.SystemId.String() == "edef8ba979d64acea3c827dcd51d21ed"
}

// ISO/IEC 23001-7
//
//	aligned(8) class ProtectionSystemSpecificHeaderBox extends FullBox(
//	   'pssh', version, flags=0,
//	) {
//	   unsigned int(8)[16] SystemID;
//	   if (version > 0) {
//	      unsigned int(32) KID_count;
//	      {
//	         unsigned int(8)[16] KID;
//	      } [KID_count];
//	   }
//	   unsigned int(32) DataSize;
//	   unsigned int(8)[DataSize] Data;
//	}
type ProtectionSystemSpecificHeader struct {
	BoxHeader     box.Header
	FullBoxHeader box.FullBoxHeader
	SystemId      box.UUID
	KidCount      uint32
	KID           []box.UUID
	DataSize      uint32
	Data          []uint8
}

func (p *ProtectionSystemSpecificHeader) Read(r io.Reader) error {
	err := p.FullBoxHeader.Read(r)
	if err != nil {
		return err
	}
	_, err = io.ReadFull(r, p.SystemId[:])
	if err != nil {
		return err
	}
	if p.FullBoxHeader.Version > 0 {
		err := binary.Read(r, binary.BigEndian, &p.KidCount)
		if err != nil {
			return err
		}
		p.KID = make([]box.UUID, p.KidCount)
		err = binary.Read(r, binary.BigEndian, p.KID)
		if err != nil {
			return err
		}
	}
	err = binary.Read(r, binary.BigEndian, &p.DataSize)
	if err != nil {
		return err
	}
	p.Data = make([]uint8, p.DataSize)
	_, err = io.ReadFull(r, p.Data)
	if err != nil {
		return err
	}
	return nil
}

func (p ProtectionSystemSpecificHeader) Write(w io.Writer) error {
	err := p.BoxHeader.Write(w)
	if err != nil {
		return err
	}
	err = p.FullBoxHeader.Write(w)
	if err != nil {
		return err
	}
	_, err = w.Write(p.SystemId[:])
	if err != nil {
		return err
	}
	if p.FullBoxHeader.Version > 0 {
		err := binary.Write(w, binary.BigEndian, p.KidCount)
		if err != nil {
			return err
		}
		err = binary.Write(w, binary.BigEndian, p.KID)
		if err != nil {
			return err
		}
	}
	err = binary.Write(w, binary.BigEndian, p.DataSize)
	if err != nil {
		return err
	}
	_, err = w.Write(p.Data)
	if err != nil {
		return err
	}
	return nil
}
