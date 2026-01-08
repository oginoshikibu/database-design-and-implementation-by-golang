package file

import "encoding/binary"

const (
	// ASCII文字のみをサポート（UTF-8のASCII範囲と互換）
	BytesPerChar = 1
	INT_SIZE     = 4
)

// Pageは、ディスクブロックの内容を保持するための構造体です。
// Fig. 3.14
//
// 書籍のコードとの差分：
// - ByteBuffer → []byte + encoding/binary パッケージを使用
// - CHARSET (US_ASCII) → BytesPerChar=1 として定義（GoにはCharset型がないため）
// - コンストラクタのオーバーロードができないため、NewPage / NewPageFromBytes に分割
type Page struct {
	byteBuffer []byte
}

// NewPageは、指定されたブロックサイズで新しいページを作成します。
func NewPage(blockSize int) *Page {
	return &Page{byteBuffer: make([]byte, blockSize)}
}

// NewPageFromBytesは、既存のバイトスライスから新しいページを作成します。
func NewPageFromBytes(b []byte) *Page {
	return &Page{byteBuffer: b}
}

// GetIntは、指定されたオフセットから整数を読み取ります。
func (p *Page) GetInt(offset int) int {
	return int(binary.BigEndian.Uint32(p.byteBuffer[offset : offset+INT_SIZE]))
}

// SetIntは、指定されたオフセットに整数を書き込みます。
func (p *Page) SetInt(offset, val int) {
	binary.BigEndian.PutUint32(p.byteBuffer[offset:offset+INT_SIZE], uint32(val))
}

// GetBytesは、指定されたオフセットからバイトスライスを読み取ります。
// データは [長さ(4バイト)][データ本体] の形式で格納されています。
func (p *Page) GetBytes(offset int) []byte {
	length := int(binary.BigEndian.Uint32(p.byteBuffer[offset:]))
	return p.byteBuffer[offset+INT_SIZE : offset+INT_SIZE+length]
}

// SetBytesは、指定されたオフセットにバイトスライスを書き込みます。
// データは [長さ(4バイト)][データ本体] の形式で格納されます。
func (p *Page) SetBytes(offset int, b []byte) {
	binary.BigEndian.PutUint32(p.byteBuffer[offset:], uint32(len(b)))
	copy(p.byteBuffer[offset+INT_SIZE:], b)
}

// GetStringは、指定されたオフセットから文字列を読み取ります。
func (p *Page) GetString(offset int) string {
	b := p.GetBytes(offset)
	return string(b)
}

// SetStringは、指定されたオフセットに文字列を書き込みます。
func (p *Page) SetString(offset int, s string) {
	p.SetBytes(offset, []byte(s))
}

// MaxLengthは、指定された長さの文字列を格納するために必要な最大バイト数を返します。
// 文字列の長さを格納するためにINT_SIZE(4バイト)が必要です。
func MaxLength(strlen int) int {
	return INT_SIZE + strlen*BytesPerChar
}

// Contentsは、底層のバイトバッファを返します。
func (p *Page) Contents() []byte {
	return p.byteBuffer
}
