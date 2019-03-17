package helper

import (
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	FakeBigFileSize = (1 << 20) * 20 //20MB
	FakeBigFilePath = "./FakeBigFile.txt"

	FakeSmallFileSize = (1 << 20) * 2 //2MB
	FakeSmallFilePath = "./FakeSmallFile.txt"

	ConfigFile = "config.json"

	PageSize = 1 << 12 //4K
)

func GenerateFakefile(filepath string, fsize int) {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic("创建测试文件失败，失败信息为：" + err.Error())
	}
	defer f.Close()
	bytes := make([]byte, PageSize, PageSize) //以 4K 一次大小写文件。
	for i := 0; i < PageSize; i++ {
		bytes[i] = 'm' //全部填充 m
	}

	for i := PageSize; i <= fsize; i += PageSize {
		f.Write(bytes)
	}
}

func GenerateUniqKey() string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	randInt := seededRand.Int()
	return strconv.Itoa(randInt) + ".txt"
}
