package sofia

import (
   "errors"
   "io"
   "log/slog"
)

// ISO/IEC 14496-12
//
//   aligned(8) class ProtectionSchemeInfoBox(fmt) extends Box('sinf') {
//      OriginalFormatBox(fmt) original_format;
//      SchemeTypeBox scheme_type_box; // optional
//      SchemeInformationBox info; // optional
//   }
type ProtectionSchemeInfo struct {
   BoxHeader      BoxHeader
   Boxes          []Box
   OriginalFormat OriginalFormat
}

func (p *ProtectionSchemeInfo) read(r io.Reader) error {
   for {
      var head BoxHeader
      err := head.read(r)
      if err == io.EOF {
         return nil
      } else if err != nil {
         return err
      }
      box_type := head.GetType()
      r := head.payload(r)
      slog.Debug("BoxHeader", "Type", box_type)
      switch box_type {
      case "schi", // Roku
         "schm": // Roku
         b := Box{BoxHeader: head}
         err := b.read(r)
         if err != nil {
            return err
         }
         p.Boxes = append(p.Boxes, b)
      case "frma":
         p.OriginalFormat.BoxHeader = head
         err := p.OriginalFormat.read(r)
         if err != nil {
            return err
         }
      default:
         return errors.New("ProtectionSchemeInfo.Decode")
      }
   }
}

func (p ProtectionSchemeInfo) write(w io.Writer) error {
   err := p.BoxHeader.write(w)
   if err != nil {
      return err
   }
   for _, b := range p.Boxes {
      err := b.write(w)
      if err != nil {
         return err
      }
   }
   return p.OriginalFormat.write(w)
}
