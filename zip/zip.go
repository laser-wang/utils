package zip

import (
	//	"archive/zip"
	//	"fmt"
	"io"
	//	"log"
	"bytes"
	"compress/zlib"
	"encoding/base64"
)

func ZipStr(src string) (bool, string) {
	var in bytes.Buffer
	b := []byte(src)
	w := zlib.NewWriter(&in)
	w.Write(b)
	w.Close()

	encodeString := base64.StdEncoding.EncodeToString(in.Bytes())
	return true, encodeString

	return true, in.String()
}

func UnZipStr(zipStr string) (bool, string) {
	var out bytes.Buffer
	var in bytes.Buffer

	decodeBytes, _ := base64.StdEncoding.DecodeString(zipStr)
	_, errWt := in.WriteString(string(decodeBytes))
	if errWt != nil {
		return false, ""
	}

	r, _ := zlib.NewReader(&in)
	io.Copy(&out, r)
	return true, out.String()
}

func ZipFile(srcFile string, zipFile string) bool {
	//// 创建一个缓冲区用来保存压缩文件内容
	//    buf := new(bytes.Buffer)

	//    // 创建一个压缩文档
	//    w := zip.NewWriter(buf)

	//    // 将文件加入压缩文档
	//    var files = []struct {
	//        Name, Body string
	//    }{
	//        {"readme.txt", "This archive contains some text files."},
	//        {"gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
	//        {"todo.txt", "Get animal handling licence.\nWrite more examples."},
	//    }
	//    for _, file := range files {
	//        f, err := w.Create(file.Name)
	//        if err != nil {
	//            log.Fatal(err)
	//        }
	//        _, err = f.Write([]byte(file.Body))
	//        if err != nil {
	//            log.Fatal(err)
	//        }
	//    }

	//    // 关闭压缩文档
	//    err := w.Close()
	//    if err != nil {
	//        log.Fatal(err)
	//    }

	//    // 将压缩文档内容写入文件
	//    f, err := os.OpenFile("file.zip", os.O_CREATE|os.O_WRONLY, 0666)
	//    if err != nil {
	//        log.Fatal(err)
	//    }
	//    buf.WriteTo(f)
	return true
}

// UnZipFile
func UnZipFile(zipFile string, unzipFile string) bool {
	//	// 打开一个zip格式文件
	//	r, err := zip.OpenReader("file.zip")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	defer r.Close()

	//	// 迭代压缩文件中的文件，打印出文件中的内容
	//	for _, f := range r.File {
	//		fmt.Printf("文件名 %s:\n", f.Name)
	//		rc, err := f.Open()
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//		_, err = io.CopyN(os.Stdout, rc, int64(f.UncompressedSize64))
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//		rc.Close()
	//		fmt.Println()
	//	}

	return true
}
