package file

import (
	"fmt"
)

// BlockIDは、ファイル内の特定のブロックを識別するための構造体です。
// Fig. 3.13
//
// 書籍のコードとの差分：
// - getterメソッドを削除し、フィールドを公開する
// - Equals() や HashCode() の実装を削除
//
// Goでは、比較可能なフィールドを持つ構造体は
// そのままマップのキーとして利用でき、== で比較できる
// そのため、Equals() や HashCode() は不要
//
// 例：
//
//	blocks := make(map[BlockID]SomeValue)
//	blocks[blockID] = value
//
//	if blockID1 == blockID2 { ... }
type BlockID struct {
	FileName string
	BlockNum int
}

func NewBlockID(filename string, blockNum int) BlockID {
	return BlockID{FileName: filename, BlockNum: blockNum}
}

func (b BlockID) String() string {
	return fmt.Sprintf("[file %s, block %d]", b.FileName, b.BlockNum)
}
