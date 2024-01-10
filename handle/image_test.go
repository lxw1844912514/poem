package handle

import (
	"gin/go-poem/db"
	"testing"
)

func Test_CreateImg(t *testing.T) {
	//poems, err := db.QueryPoems("author", "王安石")
	poems, err := db.QueryPoems("author", "李白")
	if err != nil {
		return
	}


	for index := range poems {

		CreateShiImage(poems[index])
	}
}
