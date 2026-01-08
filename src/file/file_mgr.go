package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// FileMgrは、ディスクへのブロックの読み書きを管理します。
// Fig. 3.15
//
// 書籍のコードとの差分：
// - RandomAccessFile → *os.File を使用
// - synchronized → sync.Mutex を使用
// - 例外 → error を返す（Goの慣習に従う）
type FileMgr struct {
	dbDirectory string
	blockSize   int
	isNew       bool
	openFiles   map[string]*os.File
	mu          sync.Mutex
}

// NewFileMgrは、指定されたディレクトリに対する新しいFileMgrを作成します。
func NewFileMgr(dbDirectory string, blockSize int) (*FileMgr, error) {
	// ディレクトリが存在するか確認
	info, err := os.Stat(dbDirectory)
	isNew := os.IsNotExist(err)

	switch {
	case isNew:
		// データベースが新規の場合、ディレクトリを作成
		if err := os.MkdirAll(dbDirectory, 0o755); err != nil {
			return nil, fmt.Errorf("cannot create directory %s: %w", dbDirectory, err)
		}
	case err != nil:
		return nil, fmt.Errorf("cannot access directory %s: %w", dbDirectory, err)
	case !info.IsDir():
		return nil, fmt.Errorf("%s is not a directory", dbDirectory)
	}

	// 残存する一時テーブルを削除
	entries, err := os.ReadDir(dbDirectory)
	if err != nil {
		return nil, fmt.Errorf("cannot read directory %s: %w", dbDirectory, err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "temp") {
			tempPath := filepath.Join(dbDirectory, entry.Name())
			if err := os.Remove(tempPath); err != nil {
				// 一時ファイルの削除エラーは無視
				continue
			}
		}
	}

	return &FileMgr{
		dbDirectory: dbDirectory,
		blockSize:   blockSize,
		isNew:       isNew,
		openFiles:   make(map[string]*os.File),
	}, nil
}

// Readは、指定されたブロックの内容をPageに読み込みます。
func (fm *FileMgr) Read(blk BlockID, p *Page) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	f, err := fm.getFile(blk.FileName)
	if err != nil {
		return fmt.Errorf("cannot read block %s: %w", blk, err)
	}

	offset := int64(blk.BlockNum) * int64(fm.blockSize)
	if _, err := f.Seek(offset, 0); err != nil {
		return fmt.Errorf("cannot read block %s: %w", blk, err)
	}

	if _, err := f.Read(p.Contents()); err != nil {
		return fmt.Errorf("cannot read block %s: %w", blk, err)
	}

	return nil
}

// Writeは、Pageの内容を指定されたブロックに書き込みます。
func (fm *FileMgr) Write(blk BlockID, p *Page) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	f, err := fm.getFile(blk.FileName)
	if err != nil {
		return fmt.Errorf("cannot write block %s: %w", blk, err)
	}

	offset := int64(blk.BlockNum) * int64(fm.blockSize)
	if _, err := f.Seek(offset, 0); err != nil {
		return fmt.Errorf("cannot write block %s: %w", blk, err)
	}

	if _, err := f.Write(p.Contents()); err != nil {
		return fmt.Errorf("cannot write block %s: %w", blk, err)
	}

	return nil
}

// Appendは、指定されたファイルに新しいブロックを追加し、そのBlockIDを返します。
func (fm *FileMgr) Append(filename string) (BlockID, error) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	newBlkNum, err := fm.length(filename)
	if err != nil {
		return BlockID{}, fmt.Errorf("cannot append to file %s: %w", filename, err)
	}

	blk := NewBlockID(filename, newBlkNum)
	b := make([]byte, fm.blockSize)

	f, err := fm.getFile(blk.FileName)
	if err != nil {
		return BlockID{}, fmt.Errorf("cannot append to file %s: %w", filename, err)
	}

	offset := int64(blk.BlockNum) * int64(fm.blockSize)
	if _, err := f.Seek(offset, 0); err != nil {
		return BlockID{}, fmt.Errorf("cannot append to file %s: %w", filename, err)
	}

	if _, err := f.Write(b); err != nil {
		return BlockID{}, fmt.Errorf("cannot append to file %s: %w", filename, err)
	}

	return blk, nil
}

// Lengthは、指定されたファイル内のブロック数を返します。
func (fm *FileMgr) Length(filename string) (int, error) {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	return fm.length(filename)
}

// lengthは、ロックなしの内部バージョンです。
func (fm *FileMgr) length(filename string) (int, error) {
	f, err := fm.getFile(filename)
	if err != nil {
		return 0, err
	}

	info, err := f.Stat()
	if err != nil {
		return 0, fmt.Errorf("cannot get file info for %s: %w", filename, err)
	}

	return int(info.Size()) / fm.blockSize, nil
}

// getFileは、指定されたファイル名のファイルハンドルを返します。
// ファイルがまだオープンされていない場合、オープンします。
// このメソッドはmutexを保持した状態で呼び出す必要があります。
func (fm *FileMgr) getFile(filename string) (*os.File, error) {
	f, ok := fm.openFiles[filename]
	if ok {
		return f, nil
	}

	path := filepath.Join(fm.dbDirectory, filename)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %s: %w", path, err)
	}

	fm.openFiles[filename] = f
	return f, nil
}

// IsNewは、データベースディレクトリが新規作成された場合にtrueを返します。
func (fm *FileMgr) IsNew() bool {
	return fm.isNew
}

// BlockSizeは、ブロックサイズを返します。
func (fm *FileMgr) BlockSize() int {
	return fm.blockSize
}
