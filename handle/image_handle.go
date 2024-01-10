package handle

import (
	"bufio"
	"fmt"
	"gin/go-poem/bean"
	"gin/go-poem/db"
	"github.com/golang/freetype"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

const contentReStr = ".*?[，。、！：？；]"
const bracketStr= `.*?\(.*?[^)]*\)?`
//const bracketStr= `.*?\([^)]*\)?`
//const bracketStr= `.*?[\([^)]*\)]`
//const bracketStr= ".*?[\( ,\)]"

var (
	dpi               = float64(72)
	fontFile          = "../static/font/kaiti.TTF"
	qrcodeFile        = "../static/img/gzh.jpeg"
	hinting           = "none"
	size              = float64(44)
	width             = 443
	height            = 959
	spacing           = float64(1.5)
	leftSpace         = 100  //距离左侧距离
	topSpace          = 50  //距离顶部距离
	titleLineSize     = 8  //标题字符长度
	wonb              = false
	pageSize          = 15
	calcSize          = float64(13)
	contentRes        = regexp.MustCompile(contentReStr)
	bracketRes        = regexp.MustCompile(bracketStr)
	imgDir            = "../static/img/"   //todo:: 这里换成自己的项目路径
	imgFileName       = imgDir + "%s.png"
	imgFileNameQrCode = imgDir + "qr/%s.jpg"
	imgFileNameSave   = "%s.jpg"
)

//计算图片
func calcImage(poem db.Poem) (imgBean bean.ImageBean, err error) {
	imgBean.Height = height
	imgBean.Spacing = spacing
	imgBean.Size = size
	imgBean.LeftSpace = leftSpace
	imgBean.TopSpace = topSpace

fmt.Println("poem.Content:",poem.Content)
	content := bracketRes.FindAllString(strings.Replace(poem.Content, "(", "--", -1), -1)
	content = bracketRes.FindAllString(strings.Replace(poem.Content, ")", " == ", -1), -1)
	fmt.Println(content)

	content = contentRes.FindAllString(strings.Replace(poem.Content, " ", "", -1), -1)
	fmt.Println(content)
	os.Exit(22)

	resultContent := make([]string, 0)
	if utf8.RuneCountInString(poem.Title) > titleLineSize {
		resultTitle := SubStringTitle(poem.Title)
		resultContent = append(resultContent, resultTitle...)
	} else {
		resultContent = append(resultContent, poem.Title)
	}

	resultContent = append(resultContent, poem.Author+" "+ "["+poem.Dynasty+"]")
	resultContent = append(resultContent, "")
	resultContent = append(resultContent, content...)
	imgBean.Lines = len(resultContent)
	for index := range resultContent {
		if utf8.RuneCountInString(resultContent[index]) > imgBean.MaxLen {
			imgBean.MaxLen = utf8.RuneCountInString(resultContent[index])
		}
	}

	if imgBean.MaxLen > 7 {
		imgBean.LeftSpace = 50
	}
	if imgBean.MaxLen > 10 {
		imgBean.Size = 32
	}
	if imgBean.Lines > pageSize {
		lensF := float64(imgBean.Lines) / calcSize
		if imgBean.Size < 36 {
			lensF = float64(imgBean.Lines) / float64(pageSize)
		}
		imgBean.Height = int(lensF*float64(height)) + 80 // 图片高度
	}
	if imgBean.Lines < 9 {
		imgBean.TopSpace = 200
	}

	imgBean.Content = resultContent
	fmt.Println("imgBean.MaxLen=", imgBean.MaxLen)
	fmt.Println("imgBean.Height=", imgBean.Height)
	fmt.Println("imgBean.Size=", imgBean.Size)
	fmt.Println("imgBean.LeftSpace=", imgBean.LeftSpace)
	fmt.Println("imgBean.TopSpace=", imgBean.TopSpace)
	fmt.Println("imgBean.Lines=", imgBean.Lines)
	fmt.Println("imgBean.Spacing=", imgBean.Spacing)
	fmt.Println("imgBean.Content=", imgBean.Content)
	return imgBean, nil
}

// CreateShiImage 创建诗图片
func CreateShiImage(poem db.Poem) {
	//1.读取字体文件
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		log.Println(err)
		return
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	//2.计算图片
	imgBean, _ := calcImage(poem)
	resultContent := imgBean.Content

	// 3.初始化内容
	fg, bg := image.Black, image.NewUniform(color.RGBA{189, 153, 95, 0xff}) //文字黑色，背景色
	if wonb {
		fg, bg = image.White, image.Black
	}

	rgba := image.NewRGBA(image.Rect(0, 0, width, imgBean.Height))
	draw.Draw(rgba, rgba.Bounds(), bg, image.Point{}, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(f)
	c.SetFontSize(imgBean.Size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	switch hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// 绘制文本
	pt := freetype.Pt(imgBean.LeftSpace, imgBean.TopSpace+int(c.PointToFixed(imgBean.Size)>>6))
	for _, s := range resultContent {
		_, err = c.DrawString(s, pt)
		if err != nil {
			log.Println(err)
			return
		}
		pt.Y += c.PointToFixed(imgBean.Size * imgBean.Spacing)
	}

	// 将该RGBA映像保存到磁盘，文件名格式：作者_标题
	fileName := poem.Author + "_" + strings.Replace(poem.Title, "/", "", -1)

	//err = os.MkdirAll(imgDir, 755)
	//if err != nil {
	//	log.Println("创建文件夹失败")
	//}

	outFile, err := os.Create(fmt.Sprintf(imgFileName, fileName))
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	//生成带有二维码的图片
	CreateQrCodeImg(poem, fmt.Sprintf(imgFileName, fileName), fmt.Sprintf(imgFileNameQrCode, fileName+"_qr"))
	fmt.Println("Wrote out.png OK.")
}

// CreateQrCodeImg 创建二维码诗图片
func CreateQrCodeImg(poem db.Poem, fileName string, qrFileName string) {
	imgb, _ := os.Open(fileName)
	pngPic, _ := png.Decode(imgb)

	wmb, _ := os.Open(qrcodeFile)
	watermark, _ := jpeg.Decode(wmb)
	defer wmb.Close()

	offset := image.Pt(pngPic.Bounds().Dx()-watermark.Bounds().Dx()-100, pngPic.Bounds().Dy()-watermark.Bounds().Dy()-45)
	b := pngPic.Bounds()
	m := image.NewNRGBA(b)

	draw.Draw(m, b, pngPic, image.Point{}, draw.Src)
	draw.Draw(m, watermark.Bounds().Add(offset), watermark, image.ZP, draw.Over)

	imgw, _ := os.Create(qrFileName)
	jpeg.Encode(imgw, m, &jpeg.Options{100})
	defer imgw.Close()
}

func SubStringTitle(title string) (titles []string) {
	titleSize := utf8.RuneCountInString(title)
	num := titleSize / titleLineSize
	//fmt.Println(num)
	titles = make([]string, 0)
	titleRune := []rune(title)
	for i := 0; i < num; i++ {
		titles = append(titles, string(titleRune[i*titleLineSize:(i+1)*titleLineSize]))
	}
	if titleSize > titleLineSize*num {
		titles = append(titles, string(titleRune[num*titleLineSize:]))
	}
	return titles
}
