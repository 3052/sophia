package stream

import (
   "154.pages.dev/sofia"
   "bytes"
   "encoding/hex"
   "fmt"
   "io"
   "os"
)

func decode_sidx(data []byte, sidx, moof uint32) ([][2]uint32, error) {
   var f sofia.File
   if err := f.Decode(bytes.NewReader(data[sidx:moof])); err != nil {
      return nil, err
   }
   for _, ref := range f.SegmentIndex.References {
      fmt.Println(ref[0], ref.Referenced_Size())
   }
   return f.SegmentIndex.ByteRanges(moof), nil
}

func segment_base() error {
   key, err := hex.DecodeString("dee726e9015a608a3db559a6b9a9c034")
   if err != nil {
      return err
   }
   file, err := os.Create("dec.mp4")
   if err != nil {
      return err
   }
   defer file.Close()
   var (
      sidx uint32 = 1530
      moof uint32 = 16178
   )
   data, err := os.ReadFile("enc.mp4")
   if err != nil {
      return err
   }
   if err := encode_init(file, bytes.NewReader(data[:sidx])); err != nil {
      return err
   }
   file.Write(data[sidx:moof])
   byte_ranges, err := decode_sidx(data, sidx, moof)
   if err != nil {
      return err
   }
   for _, r := range byte_ranges {
      segment := data[r[0]:r[1]+1]
      err := encode_segment(file, bytes.NewReader(segment), key)
      if err != nil {
         return err
      }
   }
   return nil
}

func encode_segment(dst io.Writer, src io.Reader, key []byte) error {
   var f sofia.File
   if err := f.Decode(src); err != nil {
      return err
   }
   for i, data := range f.MediaData.Data {
      sample := f.MovieFragment.TrackFragment.SampleEncryption.Samples[i]
      err := sample.Decrypt_CENC(data, key)
      if err != nil {
         return err
      }
   }
   return f.Encode(dst)
}

func encode_init(dst io.Writer, src io.Reader) error {
   var f sofia.File
   if err := f.Decode(src); err != nil {
      return err
   }
   for _, b := range f.Movie.Boxes {
      if b.Header.BoxType() == "pssh" {
         copy(b.Header.Type[:], "free") // Firefox
      }
   }
   sd := &f.Movie.Track.Media.MediaInformation.SampleTable.SampleDescription
   if as := sd.AudioSample; as != nil {
      copy(as.ProtectionScheme.Header.Type[:], "free") // Firefox
      copy(
         as.Entry.Header.Type[:],
         as.ProtectionScheme.OriginalFormat.DataFormat[:],
      ) // Firefox
   }
   if vs := sd.VisualSample; vs != nil {
      copy(vs.ProtectionScheme.Header.Type[:], "free") // Firefox
      copy(
         vs.Entry.Header.Type[:],
         vs.ProtectionScheme.OriginalFormat.DataFormat[:],
      ) // Firefox
   }
   return f.Encode(dst)
}
